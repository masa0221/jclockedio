package clockio_test

import (
	"errors"
	"testing"

	"github.com/masa0221/jclockedio/pkg/client/jobcan"
	"github.com/masa0221/jclockedio/pkg/service/clockio"
)

func TestAdit(t *testing.T) {
	mockJobcanClient := &MockJobcanClient{}
	mockLoggingService := &MockLoggingService{t: t}
	mockConfig := &clockio.Config{
		LoggingEnabled:        true,
		ClockedIOResultFormat: "{{.clock}}: {{.beforeStatus}} -> {{.afterStatus}}",
	}
	service := clockio.NewClockIOService(mockJobcanClient, mockLoggingService, mockConfig)

	result, err := service.Adit()
	if err != nil {
		t.Errorf("Error in Adit: %v", err)
	}

	// testing Clock
	expectedClock := "12:00"
	if result.Clock != expectedClock {
		t.Errorf("Expected clock %s but got %s", expectedClock, result.Clock)
	}

	// testing BeforeWorkingStatus
	expectedBeforeStatus := "IN"
	if result.BeforeWorkingStatus != expectedBeforeStatus {
		t.Errorf("Expected before status %s but got %s", expectedBeforeStatus, result.BeforeWorkingStatus)
	}

	// testing AfterWorkingStatus
	expectedAfterStatus := "OUT"
	if result.AfterWorkingStatus != expectedAfterStatus {
		t.Errorf("Expected after status %s but got %s", expectedAfterStatus, result.AfterWorkingStatus)
	}
}

type MockJobcanClient struct{}

func (m *MockJobcanClient) Login() error {
	return nil
}

func (m *MockJobcanClient) Adit() (*jobcan.AditResult, error) {
	return &jobcan.AditResult{
		BeforeWorkingStatus: "IN",
		AfterWorkingStatus:  "OUT",
		Clock:               "12:00",
	}, nil
}

type MockLoggingService struct {
	t *testing.T
}

func (m *MockLoggingService) Broadcast(message string) error {
	expected := "12:00: IN -> OUT"
	if message != expected {
		m.t.Errorf("Expected message is \"%s\" but got \"%s\"", expected, message)
	}

	return nil
}

func TestNoAdit(t *testing.T) {
	mockJobcanClient := &MockJobcanClient{}
	mockLoggingService := &MockFailLoggingService{}
	mockConfig := &clockio.Config{LoggingEnabled: false}
	service := clockio.NewClockIOService(mockJobcanClient, mockLoggingService, mockConfig)

	if _, err := service.Adit(); err != nil {
		t.Error("Expected to execute Log process, but it was actually executed")
	}
}

type MockFailLoggingService struct {
}

func (m *MockFailLoggingService) Broadcast(message string) error {
	return errors.New("Log is executed")
}
