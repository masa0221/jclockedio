package jobcan

import (
	"fmt"
	"log"
	"time"
)

type JobcanClient interface {
	Login() error
	Adit(noAdit bool) (*AditResult, error)
}

type DefaultJobcanClient struct {
	browser     Browser
	credentials *JobcanCredentials
	BaseUrl     string
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
	Submit(postData map[string]string, submitBtnClass string) error
	GetElementValueByID(id string) (string, error)
	ClickElementByID(id string) error
	Close()
}

func NewJobcanClient(b Browser, credentials *JobcanCredentials) *DefaultJobcanClient {
	return &DefaultJobcanClient{
		browser:     b,
		credentials: credentials,
		BaseUrl:     "",
	}
}

func (jc *DefaultJobcanClient) Login() error {
	// Open Login page
	err := jc.browser.Open(jc.getLoginUrl())
	if err != nil {
		log.Fatalf("Failed to open the Login page:%v", err)
	}

	// Post email and password
	postData := map[string]string{
		"user_email":    jc.credentials.Email,
		"user_password": jc.credentials.Password,
	}
	submitBtnClass := "form__login"
	if err := jc.browser.Submit(postData, submitBtnClass); err != nil {
		return fmt.Errorf("Failed to submit to login page: %v", err)
	}

	// Wait for rendering
	time.Sleep(1 * time.Second)

	return nil
}

func (jc *DefaultJobcanClient) Adit(noAdit bool) (*AditResult, error) {
	defer jc.browser.Close()

	clock, err := jc.browser.GetElementValueByID("clock")
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the clock: %v", err)
	}
	beforeStatus, err := jc.fetchWorkingStatus()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the status before clocking in/out: %v", err)
	}

	if !noAdit {
		if err := jc.browser.ClickElementByID("adit-button-push"); err != nil {
			return nil, fmt.Errorf("Failed to clocked in or out! (Failed to click adit button): %v", err)
		}
	}

	// Wait for rendering
	time.Sleep(1 * time.Second)
	afterStatus, err := jc.fetchAfterStatus(beforeStatus, 5)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the status after clocking in/out: %v", err)
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
		return "", fmt.Errorf("Failed to fetch the afterStatus: %v", err)
	}

	if beforeStatus != afterStatus || retry <= 0 {
		return afterStatus, nil
	}
	time.Sleep(3 * time.Second)

	return jc.fetchAfterStatus(beforeStatus, retry-1)
}
