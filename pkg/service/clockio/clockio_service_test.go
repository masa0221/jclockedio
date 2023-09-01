package clockio_test

import (
	"errors"
	"testing"

	"github.com/masa0221/jclockedio/pkg/client/jobcan"
	"github.com/masa0221/jclockedio/pkg/service/clockio"
)

func TestAdit(t *testing.T) {
	mockJobcanClient := &MockJobcanClient{}
	mockNotificationService := &MockNotificationService{t: t}
	mockConfig := &clockio.Config{
		NotifyEnabled:         true,
		ClockedIOResultFormat: "{{.clock}}: {{.beforeStatus}} -> {{.afterStatus}}",
	}
	service := clockio.NewClockIOService(mockJobcanClient, mockNotificationService, mockConfig)

	result, err := service.Adit()
	if err != nil {
		t.Errorf("Error in Adit: %v", err)
	}

	// Clockのテスト
	expectedClock := "12:00"
	if result.Clock != expectedClock {
		t.Errorf("Expected clock %s but got %s", expectedClock, result.Clock)
	}

	// BeforeWorkingStatusのテスト
	expectedBeforeStatus := "IN"
	if result.BeforeWorkingStatus != expectedBeforeStatus {
		t.Errorf("Expected before status %s but got %s", expectedBeforeStatus, result.BeforeWorkingStatus)
	}

	// AfterWorkingStatusのテスト
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

type MockNotificationService struct {
	t *testing.T
}

func (m *MockNotificationService) Notify(message string) error {
	expected := "12:00: IN -> OUT"
	if message != expected {
		m.t.Errorf("Expected message is \"%s\" but got \"%s\"", expected, message)
	}

	return nil
}

func TestNoAdit(t *testing.T) {
	mockJobcanClient := &MockJobcanClient{}
	mockNotificationService := &MockFailNotificationService{}
	mockConfig := &clockio.Config{NotifyEnabled: false}
	service := clockio.NewClockIOService(mockJobcanClient, mockNotificationService, mockConfig)

	if _, err := service.Adit(); err != nil {
		t.Error("Expected to execute Notify process, but it was actually executed")
	}
}

type MockFailNotificationService struct {
}

func (m *MockFailNotificationService) Notify(message string) error {
	return errors.New("Notify is executed")
}
