package fthttp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testPlatform   = "C64"
	testSystemCode = "an-awesome-service-as-usual"
)

func TestNewClientTimeoutSetting(t *testing.T) {
	client := NewClient(time.Second, testPlatform, testSystemCode)
	assert.Equal(t, time.Second, client.Timeout)
}

func TestNewClientWithDefaultTimeoutSetting(t *testing.T) {
	client := NewClientWithDefaultTimeout(testPlatform, testSystemCode)
	assert.Equal(t, defaultClientTimeout, client.Timeout)

}
