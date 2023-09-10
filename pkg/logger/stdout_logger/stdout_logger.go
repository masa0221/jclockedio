package stdout_logger

import (
	log "github.com/sirupsen/logrus"
)

type StdoutLogger struct {
}

func NewStdoutLogger() *StdoutLogger {
	return &StdoutLogger{}
}

func (cl *StdoutLogger) Name() string {
	return "Stdout"
}

func (cl *StdoutLogger) Log(message string) error {
	log.Infof(message)

	return nil
}
