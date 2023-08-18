package notification_test

import (
	"fmt"
	"testing"

	"github.com/masa0221/jclockedio/pkg/service/notification"
)

func TestNotify(t *testing.T) {
	mockClient := &MockClient{}
	config := &notification.NotificationConfig{
		NotifyEnabled:         true,
		ClockedIOResultFormat: "{{.clock}}: {{.beforeStatus}} -> {{.afterStatus}}",
	}
	service := notification.NewNotificationService(config, mockClient)

	err := service.Notify("12:00", "IN", "OUT")
	if err != nil {
		t.Errorf("Error notifying: %v", err)
	}
}

type MockClient struct{}

func (m *MockClient) SendMessage(message string) error {
	expected := "12:00: IN -> OUT"
	if message != expected {
		return fmt.Errorf("Expected %s but got %s", expected, message)
	}
	return nil
}
