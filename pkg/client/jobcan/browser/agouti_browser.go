package browser

import (
	"errors"
	"fmt"

	"github.com/sclevine/agouti"
	log "github.com/sirupsen/logrus"
)

type AgoutiBrowser struct {
	driver *agouti.WebDriver
	page   *agouti.Page
}

func NewAgoutiBrowser() (*AgoutiBrowser, error) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--disable-gpu",
			"--disable-extensions",
			"--disable-dev-shm-usage",
			"--no-sandbox",
		}),
	)

	log.Debug("Starting Chrome driver...")
	err := driver.Start()
	if err != nil {
		log.Error("Failed to start Chrome driver: ", err)
		driver.Stop()
		return nil, err
	}

	log.Debug("Creating new page...")
	page, err := driver.NewPage()

	if err != nil {
		log.Error("Failed to create new page: ", err)
		driver.Stop()
		return nil, err
	}

	// Set wait for timeout
	// See: https://www.seleniumqref.com/api/ruby/time_set/Ruby_implicit_wait.html
	log.Debug("Setting implicit wait timeout...")
	page.SetImplicitWait(10000)

	return &AgoutiBrowser{driver: driver, page: page}, nil
}

func (ab *AgoutiBrowser) Close() {
	log.Debug("Stopping Chrome driver...")
	ab.driver.Stop()
}

func (ab *AgoutiBrowser) Open(url string) error {
	log.Debugf("Navigating to URL: %s", url)

	err := ab.page.Navigate(url)
	if err != nil {
		logMsg := fmt.Sprintf("failed to navigate to URL %s: %v", url, err)
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	log.Debugf("Successfully navigated to URL: %s", url)
	return nil
}

func (ab *AgoutiBrowser) Submit(postData map[string]string, submitBtnClass string) error {
	log.Debugf("Attempting to fill form and submit")

	for elementID, value := range postData {
		if err := ab.fillElementByID(elementID, value); err != nil {
			logMsg := fmt.Sprintf("failed to fill data in element. elementID: %s", elementID)
			log.Error(logMsg)
			return errors.New(logMsg)
		}
	}

	if err := ab.page.FindByClass(submitBtnClass).Submit(); err != nil {
		logMsg := fmt.Sprintf("failed to submit the form with specified class: %s", submitBtnClass)
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	log.Debugf("Form submitted successfully")
	return nil
}

func (ab *AgoutiBrowser) ClickElementByID(id string) error {
	log.Debugf("Attempting to click element by ID: %s", id)

	if err := ab.page.FindByID(id).Click(); err != nil {
		logMsg := fmt.Sprintf("failed to click element with specified ID: %s", id)
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	log.Debugf("Successfully clicked element by ID: %s", id)
	return nil
}

func (ab *AgoutiBrowser) GetElementValueByID(id string) (string, error) {
	log.Debugf("Fetching element value by ID: %s", id)

	element := ab.page.FindByID(id)
	elementText, err := element.Text()
	if err != nil {
		logMsg := fmt.Sprintf("failed to fetch element text with specified ID: %s", id)
		log.Error(logMsg)
		return "", errors.New(logMsg)
	}

	log.Debugf("Successfully fetched value: %s", elementText)
	return elementText, nil
}

func (ab *AgoutiBrowser) fillElementByID(id, value string) error {
	log.Debugf("Filling element by ID: %s", id)

	if err := ab.page.FindByID(id).Fill(value); err != nil {
		logMsg := fmt.Sprintf("failed to fill value in element with specified ID: %s", id)
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	log.Debugf("Successfully filled element by ID: %s", id)
	return nil
}
