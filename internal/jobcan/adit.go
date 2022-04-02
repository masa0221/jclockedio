package jobcan

import (
	"fmt"
	"log"

	"github.com/sclevine/agouti"
)

type jobcanClient struct {
	email    string
	password string
	Host     string
	NoAdit   bool
	Verbose  bool
}

type aditResult struct {
	BeforeWorkingStatus string
	AfterWorkingStatus  string
	Clock               string
}

func New(email string, password string) *jobcanClient {
	client := jobcanClient{}
	// TODO: Remove localhost and uncomment below
	// client.Host = "ssl.jobcan.jp"
	client.Host = "localhost"
	client.email = email
	client.password = password

	return &client
}

func (c *jobcanClient) outputVerboseMessage(message string) {
	if c.Verbose {
		log.Println("[jobcan]", message)
	}
}

func (c *jobcanClient) Adit() aditResult {
	webBrowser := c.openWebBrowser()
	defer webBrowser.closeWebBrowser()
	webBrowser.Verbose = c.Verbose

	loginUrl := fmt.Sprintf("https://%s/jbcoauth/login", c.Host)
	webBrowser.login(loginUrl, c.email, c.password)

	return webBrowser.adit(c.NoAdit)
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
		log.Fatalf("Failed to start driver:%v", err)
	}

	page, err := driver.NewPage()
	c.outputVerboseMessage(fmt.Sprintf("Create new page on chrome driver. err: %v", err))
	if err != nil {
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
		log.Fatalf("Failed to navigate at Login page:%v", err)
	}

	b.outputVerboseMessage(fmt.Sprintf("Execute login process. err: %v", err))
}

func (b *webBrowser) adit(noAdit bool) aditResult {
	b.outputVerboseMessage(fmt.Sprintf("Adit process. noAdit: %v", noAdit))

	aditResult := aditResult{}
	if noAdit {
		b.outputVerboseMessage("The adit was no execute")
		aditResult.BeforeWorkingStatus = "Not attending work"
		aditResult.AfterWorkingStatus = "Working"
		aditResult.Clock = "12:23:34"
	} else {
		b.outputVerboseMessage("Execute the adit process...")
		aditResult.BeforeWorkingStatus = "Not attending work"
		aditResult.AfterWorkingStatus = "Working"
		aditResult.Clock = "12:23:34"
	}

	return aditResult
}
