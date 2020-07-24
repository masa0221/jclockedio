package jobcan

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sclevine/agouti"
)

type Env struct {
	Username string `required:"true"`
	Password string `required:"true"`
	Debug    bool   `default:false`
}

type Result struct {
	ClockTime    string
	BeforeStatus string
	AfterStatus  string
}

func ClockedInOut() Result {
	// 環境変数
	var env Env
	err := envconfig.Process("jclockedio_jobcan", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	// ブラウザはChromeを指定して起動
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--disable-gpu",
			"--disable-extensions",
			"--no-sandbox",
		}),
	)
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Fatalf("Failed to open page:%v", err)
	}
	// 読み込み待ち時間を設定 https://www.seleniumqref.com/api/ruby/time_set/Ruby_implicit_wait.html
	page.SetImplicitWait(2000)

	// ログインページに遷移
	if err := page.Navigate("https://ssl.jobcan.jp/jbcoauth/login"); err != nil {
		log.Fatalf("Failed to navigate at Login page:%v", err)
	}
	// ID, Passの要素を取得し、値を設定
	identityElement := page.FindByID("user_email")
	passwordElement := page.FindByID("user_password")
	identityElement.Fill(env.Username)
	passwordElement.Fill(env.Password)
	// ログイン
	if err := page.FindByClass("form__login").Submit(); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	// JsでHTMLが書き換わるまで待機
	time.Sleep(1 * time.Second)

	// 勤務状況取得(出勤中|未出勤)
	beforeStatusElement := page.FindByID("working_status")
	beforeStatus, err := beforeStatusElement.Text()
	if err != nil {
		log.Fatalf("Failed to fetch before status: %v", err)
	}

	// 打刻の時刻を取得
	clockElement := page.FindByID("clock")
	clockTime, err := clockElement.Text()
	if err != nil {
		log.Fatalf("Failed to fetch clock time: %v", err)
	}

	// 打刻
	if env.Debug {
		log.Println("DEBUG: It's not clocked in/out because debug mode")
	} else {
		log.Println("Clocked in/out!")
		// aditButtonElement := page.FindByID("adit-button-push")
		// if err := aditButtonElement.Click(); err != nil {
		// 	log.Fatalf("Failed to clocked in or out! (Failed to click adit button): %v", err)
		// }
	}

	// JsでHTMLが書き換わるまで待機
	time.Sleep(1 * time.Second)

	// 勤務状況取得(出勤中|未出勤)
	afterStatusElement := page.FindByID("working_status")
	afterStatus, err := afterStatusElement.Text()
	if err != nil {
		log.Fatalf("Failed to fetch after status: %v", err)
	}

	return Result{
		ClockTime:    clockTime,
		BeforeStatus: beforeStatus,
		AfterStatus:  afterStatus,
	}
}
