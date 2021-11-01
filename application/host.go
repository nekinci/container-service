package application

// Host is a port interface that relation between proxy and internal application.
type Host interface {
	// GetProtocol returns application's protocol. (e.g. tcp, udp, unix)
	GetProtocol() string

	// GetAddress returns application's address. (e.g. ipv4, ipv6)
	GetAddress() string

	// GetPort returns application's port.
	GetPort() string
}

type EmbeddedApplication struct {
	name     string
	protocol string
	address  string
	port     string
}

func (ea EmbeddedApplication) GetProtocol() string {
	return ea.protocol
}

func (ea EmbeddedApplication) GetAddress() string {
	return ea.address
}

func (ea EmbeddedApplication) GetPort() string {
	return ea.port
}

func NewEmbeddedApplication(name string, protocol string, address string, port string) EmbeddedApplication {
	return EmbeddedApplication{
		name:     name,
		protocol: protocol,
		address:  address,
		port:     port,
	}
}

func NewEmbeddedTcpApplication(name string, port string) EmbeddedApplication {
	return NewEmbeddedApplication(name, "tcp", "0.0.0.0", port)
}
