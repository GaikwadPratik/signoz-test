package entity

import "time"

type WebserverTimeouts struct {
	StartWait time.Duration `yaml:"startwait"`
	Graceful  time.Duration `yaml:"graceful"`
	Write     time.Duration `yaml:"write"`
	Read      time.Duration `yaml:"read"`
	Idle      time.Duration `yaml:"idle"`
}

type Webserver struct {
	Port     int               `yaml:"port"`
	Timeouts WebserverTimeouts `yaml:"timeouts"`
}

type AppConfig struct {
	Webserver Webserver `yaml:"webserver"`
}
