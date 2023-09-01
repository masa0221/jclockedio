package clockio

import (
	"github.com/masa0221/jclockedio/pkg/client/jobcan"
	log "github.com/sirupsen/logrus"
)

type ClockIOService struct {
	jobcanClient        jobcan.JobcanClient
	notificationService NotificationService
}

type Result struct {
	BeforeWorkingStatus string
	AfterWorkingStatus  string
	Clock               string
}

type NotificationService interface {
	Notify(clock string, beforeStatus string, afterStatus string) error
}

func NewClockIOService(jobcanClient jobcan.JobcanClient, ns NotificationService) *ClockIOService {
	return &ClockIOService{
		jobcanClient:        jobcanClient,
		notificationService: ns,
	}
}

func (cs *ClockIOService) Adit() (*Result, error) {
	// Login
	log.Debug("Log in to Jobcan")
	err := cs.jobcanClient.Login()
	if err != nil {
		log.Errorf("Failed to log in to Jobcan. err: %v", err)
		return nil, err
	}
	log.Debug("Successfully logged in to Jobcan")

	// clock in / out
	log.Debug("Start clock in/out process")
	result, err := cs.jobcanClient.Adit()
	if err != nil {
		log.Errorf("Failed to adit. err: %v", err)
		return nil, err
	}
	log.Debugf("Successfully clocked in/out. Status: %v", result.AfterWorkingStatus)

	// notify
	log.Debug("Start notification process")
	err = cs.notificationService.Notify(result.Clock, result.BeforeWorkingStatus, result.AfterWorkingStatus)
	if err != nil {
		log.Errorf("Failed to notify the messaging services. err: %v", err)
		return nil, err
	}
	log.Debug("Notification sent successfully")

	return &Result{
		BeforeWorkingStatus: result.BeforeWorkingStatus,
		AfterWorkingStatus:  result.AfterWorkingStatus,
		Clock:               result.Clock,
	}, nil
}
