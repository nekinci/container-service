package application

import (
	"github.com/nekinci/paas/specification"
	"io"
	"time"
)

type Info struct {
	Id        string
	Name      string
	UserEmail string
	StartTime string
	Status    string
	Image     string
}
type Status int8

const (
	READY   Status = 0
	WAITING Status = 1
	RUNNING Status = 2
	STOPPED Status = 3
	PAUSED  Status = 4
	ZOMBIE  Status = 5
)

var (
	reservedNames = []string{"frontend", "", "www", "api"}
)

type ProcessPipe struct {
	Stdin  *io.WriteCloser
	Stdout *io.ReadCloser
	Stderr *io.ReadCloser
}

// Application is a interface for abstraction of running applications.
type Application interface {
	// Run runs application on the system.
	Run() error

	// Kill kills the current application.
	// It returns killed container id or error if it occurred.
	Kill() (*string, error)

	// RemoveFromFileSystem removes current application on the file system.
	// It returns error if it occurred.
	RemoveFromFileSystem() error

	// RemoveApplication removes current application on the application runtime if it required.
	// It returns error if it occurred.
	RemoveApplication() error

	// GetPort returns bounded port.
	GetPort() string

	// GetProtocol returns application's protocol.
	GetProtocol() string

	// GetAddress returns tcp or udp address.
	GetAddress() string

	// GetType returns application type.
	GetType() string

	// GetSpecification returns its own specification.
	GetSpecification() *specification.Specification

	// GetApplicationInfo returns application info.
	GetApplicationInfo() Info

	// SetStatus sets application's status.
	SetStatus(status Status)

	// GetStatus gives application's status.
	GetStatus() Status

	// GetLogs returns application's logs.
	GetLogs() []Log

	// LogStream streams log when its emitted.
	LogStream(handlerFunc LogHandler)

	// AddNewLog adds new log to application logs.
	AddNewLog(log Log)

	// ListenLogs listens logs from app inside container.
	ListenLogs()

	// OpenTerminal open connection to docker container shell.
	// It returns a writer and reader.
	OpenTerminal() (processPipe *ProcessPipe, cancel func() error, err error)

	// GetCacheTime returns cache time of application.
	GetCacheTime() time.Duration
}

// NewApplication returns new application by given specification.
func NewApplication(spec specification.Specification) Application {

	if isReservedName(spec.Name) {
		return nil
	}

	if spec.Type == "docker" {
		return NewContainer(spec)
	}
	return nil
}

func isReservedName(key string) bool {
	for _, k := range reservedNames {
		if key == k {
			return true
		}
	}

	return false
}

func (s Status) String() string {

	if s == RUNNING {
		return "RUNNING"
	} else if s == WAITING {
		return "WAITING"
	} else if s == STOPPED {
		return "STOPPED"
	} else if s == ZOMBIE {
		return "ZOMBIE"
	} else if s == PAUSED {
		return "PAUSED"
	} else if s == READY {
		return "READY"
	}

	return ""
}
