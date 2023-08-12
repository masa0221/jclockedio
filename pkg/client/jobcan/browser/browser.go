type Browser interface {
	Open()
	Close()
	Login(url, email, password string)
	ClickElementByID(id string) error
	FillElementByID(id, value string) error
	WaitForRender(time.Duration)
	GetElementValueByID(id string) (string, error)
}
