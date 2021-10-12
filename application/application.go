package application

import "github.com/nekinci/paas/specification"

type Info struct {
	Id   string
	Name string
}
type Status int8
type LogType int8

const (
	READY   Status = 0
	WAITING Status = 1
	RUNNING Status = 2
	STOPPED Status = 3
	PAUSED  Status = 4
	ZOMBIE  Status = 5
)

var (
	reservedNames = []string{"frontend", "", "www"}
)

const (
	REMOVE LogType = 1
	INFO   LogType = 2
)

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
