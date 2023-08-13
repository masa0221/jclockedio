package notification

type NotificationService interface {
	Notify(message string) error
}
