package logging_test

import (
	"testing"

	"github.com/masa0221/jclockedio/pkg/service/logging"
)

func TestBroadcast(t *testing.T) {
	mockLogger := &MockLogger{t: t}
	service := logging.NewLoggingService(mockLogger, mockLogger)

	service.Broadcast("12:00: IN -> OUT")
}

type MockLogger struct {
	t *testing.T
}

func (ml *MockLogger) Name() string {
	return "Mock Logger"
}

func (ml *MockLogger) Log(message string) error {
	expected := "12:00: IN -> OUT"
	if message != expected {
		ml.t.Errorf("Expected \"%s\" but got \"%s\"", expected, message)
	}
	return nil
}
