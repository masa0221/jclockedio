package notification

import (
	"fmt"
	"html/template"
	"strings"
)

type NotificationClient interface {
	SendMessage(message string) error
}

type NotificationService struct {
	config         *NotificationConfig
	chatworkClient NotificationClient
}

type NotificationConfig struct {
	NotifyEnabled         bool
	ClockedIOResultFormat string
}

func NewNotificationService(config *NotificationConfig, notificationClient NotificationClient) *NotificationService {
	return &NotificationService{
		config:         config,
		chatworkClient: notificationClient,
	}
}

func (ns *NotificationService) Notify(clock string, beforeStatus string, afterStatus string) error {
	if ns.config.NotifyEnabled {
		if ns.config.ClockedIOResultFormat != "" {
			outputMessage, err := ns.generateOutputMessage(clock, beforeStatus, afterStatus)
			if err != nil {
				return fmt.Errorf("Failed to notify to Chatwork")
			}

			err = ns.chatworkClient.SendMessage(outputMessage)
			if err != nil {
				return fmt.Errorf("Failed to notify to Chatwork")
			}
		}
	}

	return nil
}

func (ns *NotificationService) generateOutputMessage(clock string, beforeStatus string, afterStatus string) (string, error) {
	assignData := map[string]interface{}{
		"clock":        clock,
		"beforeStatus": beforeStatus,
		"afterStatus":  afterStatus,
	}

	tpl, err := template.New("notify_message").Parse(ns.config.ClockedIOResultFormat)
	if err != nil {
		return "", fmt.Errorf("Failed to parse output format")
	}
	writer := new(strings.Builder)
	if err := tpl.Execute(writer, assignData); err != nil {
		return "", fmt.Errorf("Failed to generate output format")
	}

	return writer.String(), nil
}
