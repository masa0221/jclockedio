package jobcan_test

import (
	"fmt"
	"testing"

	"github.com/masa0221/jclockedio/pkg/client/jobcan"
)

func TestLogin(t *testing.T) {
	mockBrowser := &MockBrowser{}
	credentials := &jobcan.JobcanCredentials{
		Email:    "test@example.com",
		Password: "password",
	}
	client := jobcan.NewJobcanClient(mockBrowser, credentials)

	err := client.Login()
	if err != nil {
		t.Errorf("Error in Login: %v", err)
	}
}

func TestAdit(t *testing.T) {
	mockBrowser := &MockBrowser{}
	credentials := &jobcan.JobcanCredentials{
		Email:    "test@example.com",
		Password: "password",
	}
	client := jobcan.NewJobcanClient(mockBrowser, credentials)

	result, err := client.Adit()
	if err != nil {
		t.Errorf("Error in Adit: %v", err)
	}

	expectedClock := "12:00"
	if result.Clock != expectedClock {
		t.Errorf("Expected clock %s but got %s", expectedClock, result.Clock)
	}
}

type MockBrowser struct{}

func (m *MockBrowser) Open(url string) error {
	return nil
}
func (m *MockBrowser) Close() {}
func (m *MockBrowser) Submit(postData map[string]string, submitBtnClass string) error {
	return nil
}
func (m *MockBrowser) GetElementValueByID(id string) (string, error) {
	return "12:00", nil
}
func (m *MockBrowser) ClickElementByID(id string) error {
	return nil
}

func TestAditWithNoAdit(t *testing.T) {
	mockBrowser := &MockBrowserWithNoAdit{}
	credentials := &jobcan.JobcanCredentials{
		Email:    "test@example.com",
		Password: "password",
	}
	client := jobcan.NewJobcanClient(mockBrowser, credentials)
	client.NoAdit = true

	result, err := client.Adit()
	if err != nil {
		t.Errorf("Error in Adit with noAdit true: %v", err)
	}

	expectedClock := "12:00"
	if result.Clock != expectedClock {
		t.Errorf("Expected clock %s but got %s", expectedClock, result.Clock)
	}
}

type MockBrowserWithNoAdit struct{}

func (m *MockBrowserWithNoAdit) Open(url string) error {
	return nil
}
func (m *MockBrowserWithNoAdit) Close() {}
func (m *MockBrowserWithNoAdit) Submit(postData map[string]string, submitBtnClass string) error {
	return nil
}
func (m *MockBrowserWithNoAdit) GetElementValueByID(id string) (string, error) {
	return "12:00", nil
}
func (m *MockBrowserWithNoAdit) ClickElementByID(id string) error {
	return fmt.Errorf("ClickElementByID should not be called when noAdit is true")
}
