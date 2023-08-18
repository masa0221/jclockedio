package clockio_test

import (
	"testing"

	"github.com/masa0221/jclockedio/pkg/client/jobcan"
	"github.com/masa0221/jclockedio/pkg/service/clockio"
)

func TestAdit(t *testing.T) {
	mockJobcanClient := &MockJobcanClient{}
	mockNotificationService := &MockNotificationService{}
	service := clockio.NewClockIOService(mockJobcanClient, mockNotificationService)

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

func (m *MockJobcanClient) Adit(noAdit bool) (*jobcan.AditResult, error) {
	return &jobcan.AditResult{
		BeforeWorkingStatus: "IN",
		AfterWorkingStatus:  "OUT",
		Clock:               "12:00",
	}, nil
}

type MockNotificationService struct{}

func (m *MockNotificationService) Notify(clock string, beforeStatus string, afterStatus string) error {
	return nil
}
