package config

import "time"

type Rest struct {
	Server               string
	ContextPath          string
	ReadTimeout          int
	WriteTimeout         int
	IdleTimeout          int
	GraceShutdownTimeout int
}

func (rest Rest) ReadTimeoutDuration() time.Duration {
	return time.Duration(rest.ReadTimeout) * time.Second
}

func (rest Rest) WriteTimeoutDuration() time.Duration {
	return time.Duration(rest.WriteTimeout) * time.Second
}

func (rest Rest) IdleTimeoutDuration() time.Duration {
	return time.Duration(rest.IdleTimeout) * time.Second
}

func (rest Rest) GraceShutdownTimeoutDuration() time.Duration {
	return time.Duration(rest.GraceShutdownTimeout) * time.Second
}
