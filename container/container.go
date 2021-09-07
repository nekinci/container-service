package container

import (
	"errors"
	"github.com/nekinci/paas/specification"
	"github.com/nekinci/paas/util"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Status int8


const (
	READY   Status = 0
	WAITING Status = 1
	RUNNING Status = 2
	STOPPED Status = 3
	PAUSED  Status = 4
	ORPHAN  Status = 5
)

type Log string

type RemoveLog struct {
	Logs	[]Log
	RetryCount	int
	Mutex	sync.Mutex

}

type Container struct {
	Id            string
	Specification *specification.Specification
	StartTime     time.Time
	Status        Status
	BindingPort   *int
	IPV4          string
	Context       *Context
	Logs		  []Log
	RemoveLogs	  *RemoveLog
	IsRemovable	  bool
	CacheTime	  time.Duration // Garbage collector collects the orphaned or removed containers after cache time expires.
}

// Creates the new container.
// It requires to specification which loaded from user and the garbagecollector should be passed to container
func NewContainer(specification*specification.Specification, context* Context) Container {
	return Container{
		Id:            "",
		Specification: specification,
		StartTime:     time.Time{},
		Status:        WAITING, // 0 Ready, 1 Waiting, 2 Running, 3 Stopping, 4 Paused // TODO Refactor it.
		Context:       context,
		Logs: 		   []Log{},
		RemoveLogs: &RemoveLog{
			Logs:       []Log{},
			RetryCount: 0,
		},
		IsRemovable: true,
	}
}

// Run container if queue is empty
func (container*Container) Run() error {
	container.Context.Acquire(container)

	port := container.Specification.GetPort()
	cmd := exec.Command("docker", "run", "-d", "-p", port, container.Specification.Image)
	cmd.Stderr = os.Stderr
	executedContainer, err := cmd.Output()

	if err != nil {
		log.Fatalf("An error occurred while creating container: %v", err)
		container.Context.Release(container)
		return err
	}

	container.StartTime = time.Now()
	containerId := string(executedContainer)

	container.Id = containerId
	container.Status = RUNNING
	container.IPV4, container.BindingPort = GetBindingAddress(container)

	go container.ScheduleKill()
	return nil
}

// Kill the running container
func (container*Container) Kill() (*string, error) {
	cmd := exec.Command("docker", "kill", container.Id[:6])
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	if string(out) == container.Id[:6]{
		container.Status = STOPPED
		return nil, nil
	}

	res := container.Id[:6]
	return &res, err
}

// Kills the container when its expires.
// container: Current Container
func (container *Container) ScheduleKill() {
	timer := time.NewTimer(2 * time.Minute)
	done := make(chan bool)
	go func() {
		<-timer.C
		done <- true
	}()
	<-done

	log.Printf("Container killing: %s\n", container.Id[:6])

	cid, err := container.Kill()
	if err != nil{
		log.Printf("An error occurred while killing container: %s\n", container.Id[:6])
		container.Status = ORPHAN
		return
	}

	if *cid == container.Id[:6]{
		container.Status = STOPPED
	}

	go container.Context.Release(container)

	log.Printf("Container killed: %s; Start Time: %s; End Time: %s\n",
		container.Id[:6],
		container.StartTime.Format("02-Jan-2006 15:04:05"),
		time.Now().Format("02-Jan-2006 15:04:05"))

}

func (container *Container) Remove() error {

	var err error = nil
	cmd := exec.Command("docker", "rmi", "-f", container.Specification.Name)

	errorBuff := make([]byte, 512)

	done := make(chan bool)
	go func() {
		cmd.Stderr.Write(errorBuff)
		done <- true
	}()

	_, err = cmd.Output()

	<-done

	if err == nil && util.IsEmpty(errorBuff){
		err = errors.New(string(errorBuff))
	}

	return err
}

// Returns the port and address that bound is to the container.
func GetBindingAddress(container *Container) (string, *int) {
	cmd := exec.Command("docker", "port", container.Id[:6], container.Specification.GetPort())
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()

	if err != nil {
		log.Printf("Error: %v", err)
		panic(err)
	}


	addr := strings.Split(strings.TrimSpace(string(output)), ":") // Only for ipv4
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		panic(err)
	}

	return addr[0], &port
}


