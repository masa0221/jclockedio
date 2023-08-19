package jobcan

import (
	"fmt"
	"time"
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
		BaseUrl:     "",
		NoAdit:      false,
	}
}

func (jc *DefaultJobcanClient) Login() error {
	// Open Login page
	err := jc.browser.Open(jc.getLoginUrl())
	if err != nil {
		return fmt.Errorf("Failed to open the Login page: %v", jc.getLoginUrl())
	}

	// Post email and password
	postData := map[string]string{
		"user_email":    jc.credentials.Email,
		"user_password": jc.credentials.Password,
	}
	submitBtnClass := "form__login"
	if err := jc.browser.Submit(postData, submitBtnClass); err != nil {
		return fmt.Errorf("Failed to submit to login page")
	}

	// Wait for rendering
	time.Sleep(1 * time.Second)

	return nil
}

func (jc *DefaultJobcanClient) Adit() (*AditResult, error) {
	clock, err := jc.browser.GetElementValueByID("clock")
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the clock")
	}
	beforeStatus, err := jc.fetchWorkingStatus()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the status before clocking in/out")
	}

	if !jc.NoAdit {
		if err := jc.browser.ClickElementByID("adit-button-push"); err != nil {
			return nil, fmt.Errorf("Failed to clocked in or out! (Failed to click adit button)")
		}
	}

	// Wait for rendering
	time.Sleep(1 * time.Second)
	afterStatus, err := jc.fetchAfterStatus(beforeStatus, 5)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the status after clocking in/out")
	}

	return &AditResult{
		BeforeWorkingStatus: beforeStatus,
		AfterWorkingStatus:  afterStatus,
		Clock:               clock,
	}, nil
}

func (jc *DefaultJobcanClient) getLoginUrl() string {
	return fmt.Sprintf("%s/jbcoauth/login", jc.BaseUrl)
}

func (jc *DefaultJobcanClient) fetchWorkingStatus() (string, error) {
	return jc.browser.GetElementValueByID("working_status")
}

func (jc *DefaultJobcanClient) fetchAfterStatus(beforeStatus string, retry int) (string, error) {
	afterStatus, err := jc.fetchWorkingStatus()
	if err != nil {
		return "", fmt.Errorf("Failed to fetch the afterStatus. beforeStatus: %v, retry: %v", beforeStatus, retry)
	}

	if beforeStatus != afterStatus || retry <= 0 {
		return afterStatus, nil
	}
	time.Sleep(3 * time.Second)

	return jc.fetchAfterStatus(beforeStatus, retry-1)
}
