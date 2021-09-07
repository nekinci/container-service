package specification

import (
	"gopkg.in/yaml.v2"
	"strconv"
)

type Specification struct {
	Version int 	`yaml:"version"`
	Name string 	`yaml:"name"`
	Port int		`yaml:"port"`
	Image string	`yaml:"image"`
	Username string	`yaml:"username"`
	Password string	`yaml:"password"`
	Timeout int64	`yaml:"timeout"`
}

func NewApplication(stream []byte) (*Specification, error){
	app := Specification{}
	err := yaml.Unmarshal(stream, &app)
	if err != nil {
		return nil, err
	}
	app.Timeout = 3000
	return &app, nil
}

func (s *Specification) GetPort() string {
	return strconv.Itoa(s.Port)
}