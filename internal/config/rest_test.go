package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRest_ReadTimeoutDuration(t *testing.T) {
	rest := Rest{ReadTimeout: 10}
	expected := time.Duration(10) * time.Second

	assert.Equal(t, expected, rest.ReadTimeoutDuration())
}

func TestRest_WriteTimeoutDuration(t *testing.T) {
	rest := Rest{WriteTimeout: 20}
	expected := time.Duration(20) * time.Second

	assert.Equal(t, expected, rest.WriteTimeoutDuration())
}

func TestRest_IdleTimeoutDuration(t *testing.T) {
	rest := Rest{IdleTimeout: 30}
	expected := time.Duration(30) * time.Second

	assert.Equal(t, expected, rest.IdleTimeoutDuration())
}

func TestRest_GraceShutdownTimeoutDuration(t *testing.T) {
	rest := Rest{GraceShutdownTimeout: 40}
	expected := time.Duration(40) * time.Second

	assert.Equal(t, expected, rest.GraceShutdownTimeoutDuration())
}
