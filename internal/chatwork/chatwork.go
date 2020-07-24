package chatwork

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	APIToken string `required:"true" split_words:"true"`
	RoomID   string `required:"true" split_words:"true"`
	Debug    bool   `default:false`
}

func Send(message string) {
	// 環境変数
	var env Env
	err := envconfig.Process("jclockedio_chatwork", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	// チャット送信
	apiUrl := "https://api.chatwork.com"
	resource := fmt.Sprintf("/v2/rooms/%s/messages", env.RoomID)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)

	data := url.Values{}
	data.Set("body", message)
	log.Println(data.Encode())

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("X-ChatWorkToken", env.APIToken)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if env.Debug {
		log.Println("DEBUG: It's not send to Chatwork because debug mode")
	} else {
		log.Println("Send to Chatwork")

		resp, _ := client.Do(r)

		defer resp.Body.Close()
		contents, _ := ioutil.ReadAll(resp.Body)

		log.Println(resp.Status)
		log.Printf("result: %s\n", contents)
	}
	log.Println(urlStr)
}
