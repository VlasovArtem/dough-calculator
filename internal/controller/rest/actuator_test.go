package rest

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	handler := NewActuatorHandler()

	req := httptest.NewRequest("GET", "/health", nil)
	resp := httptest.NewRecorder()

	handler.Health().ServeHTTP(resp, req)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"status":"UP"`)
}
