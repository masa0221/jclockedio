package clockio

import (
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/masa0221/jclockedio/pkg/client/jobcan"
	log "github.com/sirupsen/logrus"
)

type ClockIOService struct {
	jobcanClient   jobcan.JobcanClient
	loggingService LoggingService
	config         *Config
}

type Config struct {
	LoggingEnabled        bool
	ClockedIOResultFormat string
}

type Result struct {
	BeforeWorkingStatus string
	AfterWorkingStatus  string
	Clock               string
}

type LoggingService interface {
	Broadcast(message string) error
}

func NewClockIOService(jobcanClient jobcan.JobcanClient, ns LoggingService, config *Config) *ClockIOService {
	return &ClockIOService{
		jobcanClient:   jobcanClient,
		loggingService: ns,
		config:         config,
	}
}

func (cs *ClockIOService) Adit() (*Result, error) {
	// Login
	log.Debug("Log in to Jobcan")
	err := cs.jobcanClient.Login()
	if err != nil {
		logMsg := fmt.Sprintf("Failed to log in to Jobcan. err: %v", err)
		log.Error(logMsg)
		return nil, errors.New(logMsg)
	}
	log.Debug("Successfully logged in to Jobcan")

	// clock in / out
	log.Debug("Start clock in/out process")
	result, err := cs.jobcanClient.Adit()
	if err != nil {
		logMsg := fmt.Sprintf("failed to adit. err: %v", err)
		log.Error(logMsg)
		return nil, errors.New(logMsg)
	}
	log.Debugf("Successfully clocked in/out. Status: %v", result.AfterWorkingStatus)

	// generate message from result
	message, err := cs.generateOutputMessage(result.Clock, result.BeforeWorkingStatus, result.AfterWorkingStatus)
	if err != nil {
		logMsg := fmt.Sprintf("failed to generate message. err: %v", err)
		log.Error(logMsg)
		return nil, errors.New(logMsg)
	}
	log.Debugf("Successfully generate message: %v", message)

	// logging
	log.Debug("Start notification process")
	if cs.config.LoggingEnabled {
		err = cs.loggingService.Broadcast(message)
		if err != nil {
			logMsg := fmt.Sprintf("failed to log. err: %v", err)
			log.Error(logMsg)
			return nil, errors.New(logMsg)
		}
		log.Debug("Logging sent successfully")
	} else {
		log.Debug("Logging is disabled in the configuration")
	}

	return &Result{
		BeforeWorkingStatus: result.BeforeWorkingStatus,
		AfterWorkingStatus:  result.AfterWorkingStatus,
		Clock:               result.Clock,
	}, nil
}

func (cs *ClockIOService) generateOutputMessage(clock string, beforeStatus string, afterStatus string) (string, error) {
	log.Debug("Generating output message format")

	assignData := map[string]interface{}{
		"clock":        clock,
		"beforeStatus": beforeStatus,
		"afterStatus":  afterStatus,
	}

	tpl, err := template.New("notify_message").Parse(cs.config.ClockedIOResultFormat)
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
