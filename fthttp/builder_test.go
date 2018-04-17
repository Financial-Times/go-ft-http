package fthttp

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewClientBuilder(t *testing.T) {
	client := NewClientBuilder().
		WithTimeout(8 * time.Second).
		WithSysInfo(testPlatform, testSystemCode).
		WithLogging(logrus.New()).
		Build()

	assert.NotNil(t, client)
}
