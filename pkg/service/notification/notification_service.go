package notification

import (
	"errors"
	"html/template"
	"strings"

	log "github.com/sirupsen/logrus"
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
	log.Debug("Starting the notification process")

	if ns.config.NotifyEnabled {
		if ns.config.ClockedIOResultFormat != "" {
			outputMessage, err := ns.generateOutputMessage(clock, beforeStatus, afterStatus)
			if err != nil {
				logMsg := "failed to generate output message for Chatwork"
				log.Error(logMsg)
				return errors.New(logMsg)
			}

			if err = ns.chatworkClient.SendMessage(outputMessage); err != nil {
				logMsg := "failed to send message to Chatwork"
				log.Error(logMsg)
				return errors.New(logMsg)
			}
		}
	} else {
		log.Debug("Notification is disabled in the configuration")
	}

	log.Debug("Notification process completed successfully")
	return nil
}

func (ns *NotificationService) generateOutputMessage(clock string, beforeStatus string, afterStatus string) (string, error) {
	log.Debug("Generating output message format")

	assignData := map[string]interface{}{
		"clock":        clock,
		"beforeStatus": beforeStatus,
		"afterStatus":  afterStatus,
	}

	tpl, err := template.New("notify_message").Parse(ns.config.ClockedIOResultFormat)
	if err != nil {
		logMsg := "failed to parse output format template"
		log.Error(logMsg)
		return "", errors.New(logMsg)
	}

	writer := new(strings.Builder)
	if err = tpl.Execute(writer, assignData); err != nil {
		logMsg := "failed to execute template for output format"
		log.Error(logMsg)
		return "", errors.New(logMsg)
	}

	outputMessage := writer.String()
	log.Debugf("Generated output message: %s", outputMessage)

	return outputMessage, nil
}
