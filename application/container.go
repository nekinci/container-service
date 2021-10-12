package application

import (
	"errors"
	"fmt"
	"github.com/nekinci/paas/specification"
	"github.com/nekinci/paas/util"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// TODO: handle errors

type Log struct {
	logType       LogType
	log           string
	logTime       time.Time
	shouldBeShown bool
}

type Container struct {
	Id            string
	Specification *specification.Specification
	StartTime     time.Time
	Status        Status
	BindingPort   *int32
	IPV4          string
	Logs          []Log
	IsRemovable   bool
	CacheTime     time.Duration // The CacheTime describes remove time of application after expiry container.
}

func (container *Container) Run() error {
	port := container.Specification.GetPort()
	cmd := exec.Command("docker", "run", "-d", "-p", port, container.Specification.Image)
	cmd.Stderr = os.Stderr
	executedContainer, err := cmd.Output()

	if err != nil {
		log.Printf("An error occurred while creating container: %v", err)
		container.Logs = append(container.Logs, NewLog(fmt.Sprintf("An error occurred while creating container: %v", err), INFO))
		return err
	}

	container.StartTime = time.Now()
	containerId := string(executedContainer)
	container.Logs = append(container.Logs, NewLog(fmt.Sprintf("Container created with id: %s", containerId), INFO))
	container.Id = containerId
	return nil
}

func (container *Container) Kill() (*string, error) {
	cmd := exec.Command("docker", "kill", container.Id[:6])
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()

	if err != nil {
		container.Logs = append(container.Logs, NewLog(fmt.Sprintf("An error occurred while killing container: %v", err), INFO))
		return nil, err
	}

	res := container.Id[:6]
	if string(out) == container.Id[:6] {
		container.Logs = append(container.Logs, NewLog(fmt.Sprintf("Container killed as expected!"), INFO))
		return &res, nil
	}

	return &res, err
}

func (container *Container) RemoveFromFileSystem() error {
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

	if err == nil && util.IsEmpty(errorBuff) {
		err = errors.New(string(errorBuff))
	}

	return err
}

func (container *Container) RemoveApplication() error {
	var err error = nil
	cmd := exec.Command("docker", "rm", "-f", container.Id[:6])

	errorBuff := make([]byte, 512)

	done := make(chan bool)
	go func() {
		cmd.Stderr.Write(errorBuff)
		done <- true
	}()

	out, err := cmd.Output()

	if string(out) != container.Id[:6] {
		err = errors.New("Removed container id mismatch with container id: " + container.Id[:6])
	}

	<-done

	if err == nil && !util.IsEmpty(errorBuff) {
		err = errors.New(string(errorBuff))
	}

	container.Status = STOPPED

	return err
}

func (container *Container) GetPort() string {
	if container.BindingPort == nil {
		container.IPV4, container.BindingPort = GetBindingAddress(*container)
	}
	p := *container.BindingPort
	return strconv.Itoa(int(p))
}

func (container *Container) GetAddress() string {
	if container.IPV4 == "" {
		container.IPV4, container.BindingPort = GetBindingAddress(*container)
	}

	return container.IPV4
}

func (container *Container) GetType() string {
	return container.Specification.Type
}

func (container *Container) GetSpecification() *specification.Specification {
	return container.Specification
}

func (container *Container) GetApplicationInfo() Info {
	return Info{
		Id:   container.Id,
		Name: container.Specification.Name,
	}
}

func (container *Container) SetStatus(status Status) {
	container.Status = status
}

func (container *Container) GetStatus() Status {
	return container.Status
}

// Returns the port and address that bound is to the container.
func GetBindingAddress(container Container) (string, *int32) {
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

	p := int32(port)
	return addr[0], &p
}

func (container *Container) GetProtocol() string {
	return "tcp"
}

func NewContainer(specification2 specification.Specification) *Container {
	return &Container{
		Id:            "",
		Specification: &specification2,
		StartTime:     time.Time{},
		Status:        WAITING,
		BindingPort:   nil,
		IPV4:          "",
		Logs:          []Log{},
		IsRemovable:   true,
		CacheTime:     time.Minute * 30,
	}
}

func NewLog(log string, logType LogType) Log {
	return Log{
		logType:       logType,
		log:           log,
		logTime:       time.Now(),
		shouldBeShown: true,
	}
}

func NewSecretLog(log string, logType LogType) Log {
	return Log{
		logType:       logType,
		log:           log,
		logTime:       time.Now(),
		shouldBeShown: false,
	}
}
