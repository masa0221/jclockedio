package clockio

type ClockIOService interface {
	ClockIn() error
	ClockOut() error
}
