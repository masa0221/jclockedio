package browser

import (
	"fmt"

	"github.com/sclevine/agouti"
)

type AgoutiBrowser struct {
	driver *agouti.WebDriver
	page   *agouti.Page
}

func NewAgoutiBrowser() (*AgoutiBrowser, error) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",              // Runs Chrome in headless mode
			"--disable-gpu",           // Disable gpu acceleration
			"--disable-extensions",    // Disable Chrome extensions
			"--disable-dev-shm-usage", // Change shared memory to /tmp from /dev/shm
			"--no-sandbox",
		}),
	)

	err := driver.Start()
	if err != nil {
		driver.Stop()
		return nil, err
	}

	page, err := driver.NewPage()

	if err != nil {
		driver.Stop()
		return nil, err
	}

	// Wait for driver
	// See: https://www.seleniumqref.com/api/ruby/time_set/Ruby_implicit_wait.html
	page.SetImplicitWait(10000)

	return &AgoutiBrowser{driver: driver, page: page}, nil
}

func (ab *AgoutiBrowser) Close() {
	ab.driver.Stop()
}

func (ab *AgoutiBrowser) Submit(url string, postData map[string]string, submitBtnClass string) error {
	err := ab.page.Navigate(url)
	if err != nil {
		return fmt.Errorf("Failed to navigate to Login page: %v", err)
	}

	for elementID, value := range postData {
		ab.fillElementByID(elementID, value)
	}

	// submit
	if err := ab.page.FindByClass(submitBtnClass).Submit(); err != nil {
		ab.driver.Stop()
		return fmt.Errorf("Failed to login: %v", err)
	}

	return nil
}

func (ab *AgoutiBrowser) ClickElementByID(id string) error {
	return ab.page.FindByID(id).Click()
}

func (ab *AgoutiBrowser) GetElementValueByID(id string) (string, error) {
	element := ab.page.FindByID(id)
	return element.Text()
}

func (ab *AgoutiBrowser) fillElementByID(id, value string) error {
	return ab.page.FindByID(id).Fill(value)
}
