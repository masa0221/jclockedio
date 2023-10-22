package logging

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Logger interface {
	Name() string
	Log(message string) error
}

type LoggingService struct {
	loggers []Logger
}

func NewLoggingService(loggers ...Logger) *LoggingService {
	return &LoggingService{
		loggers: loggers,
	}
}

func (ls *LoggingService) Broadcast(message string) error {
	log.Debug("Starting the logging process")

	errChan := make(chan error, len(ls.loggers))

	for _, logger := range ls.loggers {
		go func(l Logger) {
			if err := l.Log(message); err != nil {
				errMsg := fmt.Sprintf("failed to log on %s", l.Name())
				log.Error(errMsg)
				errChan <- errors.New(errMsg)
			} else {
				errChan <- nil
			}
		}(logger)
	}

	for i := 0; i < len(ls.loggers); i++ {
		err := <-errChan
		if err != nil {
			return err
		}
	}

	log.Debug("logging process completed successfully")
	return nil
}
