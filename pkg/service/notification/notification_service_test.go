package notification_test

import (
	"testing"

	"github.com/masa0221/jclockedio/pkg/service/notification"
)

func TestNotify(t *testing.T) {
	mockClient := &MockClient{t: t}
	service := notification.NewNotificationService(mockClient)

	service.Notify("12:00: IN -> OUT")
}

type MockClient struct {
	t *testing.T
}

func (m *MockClient) SendMessage(message string) error {
	expected := "12:00: IN -> OUT"
	if message != expected {
		m.t.Errorf("Expected \"%s\" but got \"%s\"", expected, message)
	}
	return nil
}
