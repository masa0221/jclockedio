package chatwork

type ChatworkClient interface {
	SendMessage(roomID string, message string) error
}
