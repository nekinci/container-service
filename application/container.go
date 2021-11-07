package application

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/nekinci/paas/specification"
	"github.com/nekinci/paas/util"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// TODO: handle errors

type Container struct {
	Id            string
	Specification *specification.Specification
	StartTime     time.Time
	Status        Status
	BindingPort   *int32
	IPV4          string
	Logs          []Log
	Removable     bool
	CacheTime     time.Duration // The CacheTime describes remove time of application after expiry container.
	LogHandlers   []LogHandler
}

func (container *Container) Run() error {
	port := container.Specification.GetPort()
	cmd := exec.Command("docker", "run", "-d", "-p", port, container.Specification.Image)
	cmd.Stderr = os.Stderr
	executedContainer, err := cmd.Output()

	if err != nil {
		log.Printf("An error occurred while creating container: %v", err)
		container.AddNewLog(FormatString("An error occurred while creating container: %v", err).ToInfoLog())
		return err
	}

	container.StartTime = time.Now()
	containerId := string(executedContainer)
	container.AddNewLog(FormatString("Container created with id: %s", containerId).ToInfoLog())
	container.Id = containerId
	return nil
}

func (container *Container) Kill() (*string, error) {
	cmd := exec.Command("docker", "kill", container.Id[:6])
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()

	if err != nil {
		container.AddNewLog(FormatString("An error occurred while killing container: %v", err).ToErrorLog())
		return nil, err
	}

	res := container.Id[:6]
	if string(out[:6]) == container.Id[:6] {
		container.AddNewLog(FormatString("Container killed as expected!").ToInfoLog())
		return &res, nil
	}

	container.AddNewLog(FormatString("An error occurred while killing container: %v", err).ToErrorLog())
	return &res, err
}

func (container *Container) RemoveFromFileSystem() error {
	var err error = nil
	cmd := exec.Command("docker", "rmi", "-f", container.Specification.Image)

	errorBuff := make([]byte, 512)

	done := make(chan bool)
	go func() {
		_, err = cmd.Stderr.Write(errorBuff)
		if err != nil {
			log.Printf("%v", err)
		}
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
		Id:        container.Id,
		Name:      container.Specification.Name,
		UserEmail: container.Specification.Email,
		StartTime: container.StartTime.Format(time.RFC3339),
		Status:    container.Status.String(),
		Image:     container.Specification.Image,
	}
}

func (container *Container) SetStatus(status Status) {
	container.Status = status
}

func (container *Container) GetStatus() Status {
	return container.Status
}

func (container *Container) GetLogs() []Log {
	return container.Logs
}

func (container *Container) LogStream(handler LogHandler) {
	container.LogHandlers = append(container.LogHandlers, handler)
}

func (container *Container) AddNewLog(log Log) {
	container.Logs = append(container.Logs, log)
	for _, handler := range container.LogHandlers {
		handler(log)
	}
}

func (container *Container) ListenLogs() {
	r, w := io.Pipe()
	defer w.Close()
	defer r.Close()
	id := container.Id[:6]
	cmd := exec.Command("docker", "logs", "-f", id)
	cmd.Stdout = w
	cmd.Stderr = cmd.Stdout
	cmd.Start()
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			log := scanner.Text()
			container.AddNewLog(FormatString("%v", log).ToInfoLog())
		}
	}()

	cmd.Wait()

	fmt.Printf("CIKTI.....")

}

func (container *Container) OpenTerminal() (*ProcessPipe, func() error, error) {
	cmd := exec.Command("docker", "exec", "-i", container.Id[:6], "/bin/bash")

	container.AddNewLog(FormatString("Terminal session created.\n").ToInfoLog())

	stdin, stdinErr := cmd.StdinPipe()
	if stdinErr != nil {
		container.AddNewLog(FormatString("%v", stdinErr).ToErrorLog())
		return nil, nil, fmt.Errorf("Stdin pipe err: %v\n", stdinErr)
	}

	stdout, stdoutErr := cmd.StdoutPipe()
	if stdoutErr != nil {
		container.AddNewLog(FormatString("%v", stdoutErr).ToErrorLog())
		return nil, nil, fmt.Errorf("Stdout pipe err: %v\n", stdoutErr)
	}

	stderr, stderrPipe := cmd.StderrPipe()
	if stderrPipe != nil {
		container.AddNewLog(FormatString("%v", stderrPipe).ToErrorLog())
		return nil, nil, fmt.Errorf("Stderr pipe err: %v\n", stderrPipe)
	}

	err := cmd.Start()
	if err != nil {
		container.AddNewLog(FormatString("Terminal not opening, %v", err).ToErrorLog())
	}

	go func() {
		waitErr := cmd.Wait()
		if waitErr != nil {
			container.AddNewLog(FormatString("Process not forked as expected! %v", waitErr).ToErrorLog())
		}
	}()

	// cancel kills child process
	cancel := func() error {
		log.Println("Process killing...")
		err := cmd.Process.Kill()
		if err != nil {
			container.AddNewLog(FormatString("An error occurred while killing terminal process, %v", err).ToErrorLog())
		}
		return err
	}

	return &ProcessPipe{
		Stdin:  &stdin,
		Stdout: &stdout,
		Stderr: &stderr,
	}, cancel, nil
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
	pp := strings.Replace(addr[1], "\n", "", 1)
	port, err := strconv.Atoi(pp)
	if err != nil {
		panic(err)
	}

	p := int32(port)
	return addr[0], &p
}

func (container *Container) GetProtocol() string {
	return "tcp"
}

func (container *Container) GetCacheTime() time.Duration {
	return container.CacheTime
}

func (container *Container) IsRemovable() bool {
	return container.Removable
}

func NewContainer(specification2 specification.Specification, removable bool) *Container {
	return &Container{
		Id:            "",
		Specification: &specification2,
		StartTime:     time.Time{},
		Status:        WAITING,
		BindingPort:   nil,
		IPV4:          "",
		Logs:          []Log{},
		Removable:     removable,
		CacheTime:     time.Minute * 30,
	}
}
