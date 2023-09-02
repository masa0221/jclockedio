package jobcan

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type JobcanClient interface {
	Login() error
	Adit() (*AditResult, error)
}

type DefaultJobcanClient struct {
	browser     Browser
	credentials *JobcanCredentials
	BaseUrl     string
	NoAdit      bool
}

type JobcanCredentials struct {
	Email    string
	Password string
}

type AditResult struct {
	BeforeWorkingStatus string
	AfterWorkingStatus  string
	Clock               string
}

type Browser interface {
	Open(url string) error
	Close()
	Submit(postData map[string]string, submitBtnClass string) error
	GetElementValueByID(id string) (string, error)
	ClickElementByID(id string) error
}

func NewJobcanClient(b Browser, credentials *JobcanCredentials) *DefaultJobcanClient {
	return &DefaultJobcanClient{
		browser:     b,
		credentials: credentials,
		BaseUrl:     "https://ssl.jobcan.jp",
		NoAdit:      false,
	}
}

func (jc *DefaultJobcanClient) Login() error {
	// Open Login page
	log.Debug("Starting the login process")
	err := jc.browser.Open(jc.getLoginUrl())
	if err != nil {
		logMsg := fmt.Sprintf("failed to open the login page: %v", jc.getLoginUrl())
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	// Post email and password
	log.Debug("Submitting email and password.")
	postData := map[string]string{
		"user_email":    jc.credentials.Email,
		"user_password": jc.credentials.Password,
	}
	submitBtnClass := "form__login"
	if err := jc.browser.Submit(postData, submitBtnClass); err != nil {
		logMsg := "failed to submit to login page"
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	// Wait for rendering
	log.Debug("Waiting for page to render.")
	time.Sleep(1 * time.Second)

	log.Debug("Successfully logged in.")
	return nil
}

func (jc *DefaultJobcanClient) Adit() (*AditResult, error) {
	log.Debug("Starting the Adit process.")

	clock, err := jc.browser.GetElementValueByID("clock")
	if err != nil {
		logMsg := "failed to fetch the clock"
		log.Error(logMsg)
		return nil, errors.New(logMsg)
	}

	beforeStatus, err := jc.fetchWorkingStatus()
	if err != nil {
		logMsg := "failed to fetch the status before clocking in/out"
		log.Error(logMsg)
		return nil, errors.New(logMsg)
	}

	if !jc.NoAdit {
		log.Debug("Click the Adit button")
		if err := jc.browser.ClickElementByID("adit-button-push"); err != nil {
			logMsg := "failed to clock in or out (failed to click adit button)"
			log.Error(logMsg)
			return nil, errors.New(logMsg)
		}
	}

	// Wait for rendering
	log.Debug("Waiting for page to render.")
	time.Sleep(1 * time.Second)
	afterStatus, err := jc.fetchAfterStatus(beforeStatus, 5)
	if err != nil {
		logMsg := "failed to fetch the status after clocking in/out"
		log.Error(logMsg)
		return nil, errors.New(logMsg)
	}

	log.Debug("Successfully completed the Adit process.")
	return &AditResult{
		BeforeWorkingStatus: beforeStatus,
		AfterWorkingStatus:  afterStatus,
		Clock:               clock,
	}, nil
}

func (jc *DefaultJobcanClient) getLoginUrl() string {
	log.Debugf("Login url is: %v", jc.BaseUrl)
	return fmt.Sprintf("%s/jbcoauth/login", jc.BaseUrl)
}

func (jc *DefaultJobcanClient) fetchWorkingStatus() (string, error) {
	log.Debug("Fetching the working status.")
	return jc.browser.GetElementValueByID("working_status")
}

func (jc *DefaultJobcanClient) fetchAfterStatus(beforeStatus string, retry int) (string, error) {
	log.Debug("Fetching the after working status.")
	afterStatus, err := jc.fetchWorkingStatus()
	if err != nil {
		logMsg := fmt.Sprintf("Failed to fetch the afterStatus. beforeStatus: %v, retry: %v", beforeStatus, retry)
		log.Error(logMsg)
		return "", errors.New(logMsg)
	}

	if beforeStatus != afterStatus || retry <= 0 {
		return afterStatus, nil
	}
	time.Sleep(2 * time.Second)

	return jc.fetchAfterStatus(beforeStatus, retry-1)
}
