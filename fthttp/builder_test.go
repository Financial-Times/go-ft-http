package fthttp

import (
	"testing"
	"time"

	"github.com/Financial-Times/go-logger/v2"

	"github.com/stretchr/testify/assert"
)

func TestNewClientBuilder(t *testing.T) {
	log := logger.NewUPPLogger(testSystemCode, "info")
	client := NewClientBuilder().
		WithTimeout(8*time.Second).
		WithSysInfo(testPlatform, testSystemCode).
		WithLogging(log).
		Build()

	assert.NotNil(t, client)
}
