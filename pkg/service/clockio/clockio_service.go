package clockio

import (
	"fmt"

	"github.com/masa0221/jclockedio/pkg/client/jobcan"
)

type ClockIOService struct {
	jobcanClient        jobcan.JobcanClient
	notificationService NotificationService
	noAdit              bool
	notifyFormat        string
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
		noAdit:              false,
	}
}

func (cs *ClockIOService) Adit() (*Result, error) {
	// Login
	err := cs.jobcanClient.Login()
	if err != nil {
		return nil, fmt.Errorf("Failed to log in to Jobcan: %v", err)
	}

	// clock in / out
	result, err := cs.jobcanClient.Adit(cs.noAdit)
	if err != nil {
		return nil, fmt.Errorf("Failed to clock in/out: %v", err)
	}

	// notify
	err = cs.notificationService.Notify(result.Clock, result.BeforeWorkingStatus, result.AfterWorkingStatus)
	if err != nil {
		return nil, fmt.Errorf("Failed to notify: %v", err)
	}

	return &Result{
		BeforeWorkingStatus: result.BeforeWorkingStatus,
		AfterWorkingStatus:  result.AfterWorkingStatus,
		Clock:               result.Clock,
	}, nil
}
