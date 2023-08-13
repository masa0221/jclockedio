package jobcan

import (
	"fmt"
	"time"
)

type JobcanClient struct {
	browser Browser
	baseUrl string
	noAdit  bool
}

type AditResult struct {
	BeforeWorkingStatus string
	AfterWorkingStatus  string
	Clock               string
}

type Browser interface {
	Close()
	Submit(url string, credentials map[string]string, submitBtnClass string) error
	WaitForRender(time.Duration)
	GetElementValueByID(id string) (string, error)
	ClickElementByID(id string) error
}

func NewJobcanClient(b Browser, url string, noAdit bool) *JobcanClient {
	return &JobcanClient{
		browser: b,
		baseUrl: url,
		noAdit:  noAdit,
	}
}

func (jc *JobcanClient) Login(email string, password string) error {
	credentials := map[string]string{
		"user_email":    email,
		"user_password": password,
	}
	submitBtnClass := "form__login"

	if err := jc.browser.Submit(jc.getLoginUrl(), credentials, submitBtnClass); err != nil {
		return fmt.Errorf("Failed to create browser: %v", err)
	}
	defer jc.browser.Close()

	// Wait for rendering
	time.Sleep(1 * time.Second)

	return nil
}

func (jc *JobcanClient) Adit() (*AditResult, error) {
	clock, err := jc.browser.GetElementValueByID("clock")
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the clock: %v", err)
	}
	beforeStatus, err := jc.fetchWorkingStatus()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the status before clocking in/out: %v", err)
	}

	if !jc.noAdit {
		if err := jc.browser.ClickElementByID("adit-button-push"); err != nil {
			jc.browser.Close()

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

func (jc *JobcanClient) getLoginUrl() string {
	return fmt.Sprintf("%s/jbcoauth/login", jc.baseUrl)
}

func (jc *JobcanClient) fetchWorkingStatus() (string, error) {
	return jc.browser.GetElementValueByID("working_status")
}

func (jc *JobcanClient) fetchAfterStatus(beforeStatus string, retry int) (string, error) {
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
