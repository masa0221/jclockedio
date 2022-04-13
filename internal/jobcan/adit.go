package jobcan

import (
	"fmt"
	"log"
	"time"

	"github.com/sclevine/agouti"
)

type jobcanClient struct {
	email    string
	password string
	NoAdit   bool
	Verbose  bool
	BaseUrl  string
}

type aditResult struct {
	BeforeWorkingStatus string
	AfterWorkingStatus  string
	Clock               string
}

func New(email string, password string) *jobcanClient {
	client := jobcanClient{}
	client.BaseUrl = "https://ssl.jobcan.jp"
	client.email = email
	client.password = password

	return &client
}

func (c *jobcanClient) outputVerboseMessage(message string) {
	if c.Verbose {
		log.Println("[jobcan]", message)
	}
}

func (c *jobcanClient) generateLoginUrl() string {
	return fmt.Sprintf("%s/jbcoauth/login", c.BaseUrl)
}

func (c *jobcanClient) Adit() aditResult {
	webBrowser := c.openWebBrowser()
	defer webBrowser.closeWebBrowser()
	webBrowser.Verbose = c.Verbose

	webBrowser.login(c.generateLoginUrl(), c.email, c.password)
	// Wait for rendering
	time.Sleep(1 * time.Second)

	aditResult := aditResult{}
	aditResult.BeforeWorkingStatus = webBrowser.fetchWorkingStatus()

	c.outputVerboseMessage(fmt.Sprintf("Adit process. noAdit: %v", c.NoAdit))
	if c.NoAdit {
		c.outputVerboseMessage("The adit was no execute")
	} else {
		webBrowser.adit()
	}

	// Wait for rendering
	time.Sleep(1 * time.Second)
	aditResult.Clock = webBrowser.fetchClock()
	aditResult.AfterWorkingStatus = webBrowser.fetchWorkingStatus()

	return aditResult
}

type webBrowser struct {
	driver  *agouti.WebDriver
	page    *agouti.Page
	Verbose bool
}

func (b *webBrowser) outputVerboseMessage(message string) {
	if b.Verbose {
		log.Println("[jobcan] ChromeDriver:", message)
	}
}

func (c *jobcanClient) openWebBrowser() *webBrowser {
	browser := webBrowser{}
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
	c.outputVerboseMessage(fmt.Sprintf("Start chrome driver. err: %v", err))
	if err != nil {
		driver.Stop()
		log.Fatalf("Failed to start driver:%v", err)
	}

	page, err := driver.NewPage()
	c.outputVerboseMessage(fmt.Sprintf("Create new page on chrome driver. err: %v", err))
	if err != nil {
		driver.Stop()
		log.Fatalf("Failed to open page:%v", err)
	}

	browser.driver = driver
	browser.page = page
	// Wait for driver
	// See: https://www.seleniumqref.com/api/ruby/time_set/Ruby_implicit_wait.html
	page.SetImplicitWait(10000)

	return &browser
}

func (b *webBrowser) closeWebBrowser() {
	b.driver.Stop()
}

func (b *webBrowser) login(url string, email string, password string) {
	err := b.page.Navigate(url)
	b.outputVerboseMessage(fmt.Sprintf("Open login page. url: %v err: %v", url, err))
	if err != nil {
		b.driver.Stop()
		log.Fatalf("Failed to navigate at Login page:%v", err)
	}
	// Input login form
	identityElement := b.page.FindByID("user_email")
	passwordElement := b.page.FindByID("user_password")
	identityElement.Fill(email)
	passwordElement.Fill(password)

	// submit
	if err := b.page.FindByClass("form__login").Submit(); err != nil {
		b.driver.Stop()
		log.Fatalf("Failed to login: %v", err)
	}

	b.outputVerboseMessage(fmt.Sprintf("Execute login process. err: %v", err))
}

func (b *webBrowser) fetchWorkingStatus() string {
	return b.fetchElementTextById("working_status")
}

func (b *webBrowser) fetchClock() string {
	return b.fetchElementTextById("clock")
}

func (b *webBrowser) fetchElementTextById(id string) string {
	b.outputVerboseMessage(fmt.Sprintf("Fetch %v text", id))
	element := b.page.FindByID(id)
	elementText, err := element.Text()
	if err != nil {
		b.driver.Stop()
		log.Fatalf("Failed to fetch %s: %v", id, err)
	}

	return elementText
}

func (b *webBrowser) adit() {
	b.outputVerboseMessage("Execute the adit process...")

	aditButtonElement := b.page.FindByID("adit-button-push")
	if err := aditButtonElement.Click(); err != nil {
		b.driver.Stop()
		log.Fatalf("Failed to clocked in or out! (Failed to click adit button): %v", err)
	}
}
