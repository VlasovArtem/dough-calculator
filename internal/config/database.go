package config

import "time"

type Database struct {
	Uri               string
	ConnectionTimeout time.Duration
}
