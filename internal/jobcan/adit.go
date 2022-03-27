package jobcan

import (
	"fmt"
	"log"

	"github.com/sclevine/agouti"
)

type jobcanClient struct {
	email    string
	password string
	NoAdit   bool
	Verbose  bool
	Debug    bool
}

func New(email string, password string) *jobcanClient {
	client := jobcanClient{}
	client.email = email
	client.password = password

	return &client
}

func (c *jobcanClient) outputVerboseMessage(message string) {
	if c.Verbose {
		log.Println("[jobcan]", message)
	}
}

func (c *jobcanClient) Adit() {
	webBrowser := c.openWebBrowser()
	defer webBrowser.closeWebBrowser()
	webBrowser.Verbose = c.Verbose

	var loginUrl string
	if c.Debug {
		loginUrl = "http://127.0.0.1:8080"
	} else {
		loginUrl = "http://127.0.0.1:1111"
		// loginUrl = "https://ssl.jobcan.jp/jbcoauth/login"
	}
	webBrowser.login(loginUrl, c.email, c.password)
	webBrowser.adit(c.NoAdit)
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

func (b *webBrowser) adit(noAdit bool) {
	b.outputVerboseMessage(fmt.Sprintf("Adit process. noAdit: %v", noAdit))
	if noAdit {
		b.outputVerboseMessage("The adit was no execute")
	} else {
		b.outputVerboseMessage("Execute the adit process...")
	}
}
