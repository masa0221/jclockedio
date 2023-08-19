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

	// Set wait for timeout
	// See: https://www.seleniumqref.com/api/ruby/time_set/Ruby_implicit_wait.html
	page.SetImplicitWait(10000)

	return &AgoutiBrowser{driver: driver, page: page}, nil
}

func (ab *AgoutiBrowser) Close() {
	ab.driver.Stop()
}

func (ab *AgoutiBrowser) Open(url string) error {
	err := ab.page.Navigate(url)
	if err != nil {
		return fmt.Errorf("Failed to Open to page. url: %v", url)
	}

	return nil
}

func (ab *AgoutiBrowser) Submit(postData map[string]string, submitBtnClass string) error {
	for elementID, value := range postData {
		if err := ab.fillElementByID(elementID, value); err != nil {
			return fmt.Errorf("Failed to fill the specified data in elements. elementID: %v", elementID)
		}
	}

	// submit
	if err := ab.page.FindByClass(submitBtnClass).Submit(); err != nil {
		return fmt.Errorf("Failed to submit to the page with the specified data. specified class: %v", submitBtnClass)
	}

	return nil
}

func (ab *AgoutiBrowser) ClickElementByID(id string) error {
	if err := ab.page.FindByID(id).Click(); err != nil {
		return fmt.Errorf("Failed to click the element with the specified ID: %v", id)
	}

	return nil
}

func (ab *AgoutiBrowser) GetElementValueByID(id string) (string, error) {
	element := ab.page.FindByID(id)
	elementText, err := element.Text()
	if err != nil {
		return "", fmt.Errorf("Failed to fetch the element text of the specified ID: %v", id)
	}

	return elementText, nil
}

func (ab *AgoutiBrowser) fillElementByID(id, value string) error {
	if err := ab.page.FindByID(id).Fill(value); err != nil {
		return fmt.Errorf("Failed to fill the value in the specified ID: %v", id)
	}

	return nil
}
