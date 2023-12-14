//go:generate mockgen -source=actuator.go -destination mocks/actuator.go -package mocks

package domain

import (
	"net/http"
)

type ActuatorHandler interface {
	Health() http.HandlerFunc
}
