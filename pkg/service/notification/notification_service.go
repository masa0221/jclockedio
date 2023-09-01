package notification

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type NotificationClient interface {
	SendMessage(message string) error
}

type NotificationService struct {
	chatworkClient NotificationClient
}

func NewNotificationService(notificationClient NotificationClient) *NotificationService {
	return &NotificationService{
		chatworkClient: notificationClient,
	}
}

func (ns *NotificationService) Notify(message string) error {
	log.Debug("Starting the notification process")

	if err := ns.chatworkClient.SendMessage(message); err != nil {
		logMsg := "failed to send message to Chatwork"
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	log.Debug("Notification process completed successfully")
	return nil
}
