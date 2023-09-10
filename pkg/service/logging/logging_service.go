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

	for _, logger := range ls.loggers {
		if err := logger.Log(message); err != nil {
			logMsg := fmt.Sprintf("failed to log on %s", logger.Name())
			log.Error(logMsg)
			return errors.New(logMsg)
		}

	}

	log.Debug("logging process completed successfully")
	return nil
}
