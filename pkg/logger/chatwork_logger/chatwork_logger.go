package chatwork_logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ChatworkLogger struct {
	BaseUrl  string
	apiToken string
	config   *Config
}

type Config struct {
	ToRoomId string
	Unread   bool
}

type postMessageResult struct {
	MessageId string `json:"message_id"`
}

func NewChatworkLogger(token string, config *Config) *ChatworkLogger {
	return &ChatworkLogger{
		BaseUrl:  "https://api.chatwork.com/v2",
		config:   config,
		apiToken: token,
	}
}

func (cl *ChatworkLogger) Name() string {
	return "Chatwork"
}

func (cl *ChatworkLogger) Log(message string) error {
	log.Debug("Starting to send a message to Chatwork")

	endpoint := fmt.Sprintf("%v/rooms/%v/messages", cl.BaseUrl, cl.config.ToRoomId)
	self_unread := 0
	if cl.config.Unread {
		self_unread = 1
	}
	payload := strings.NewReader(fmt.Sprintf("self_unread=%v&body=%v", self_unread, url.QueryEscape(message)))

	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		logMsg := fmt.Sprintf("Failed to create request: %s", err)
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	req.Header.Add("Aclept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-ChatWorkToken", cl.apiToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logMsg := fmt.Sprintf("Failed to send request to Chatwork: %s", err)
		log.Error(logMsg)
		return errors.New(logMsg)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logMsg := fmt.Sprintf("Failed to read Chatwork API response body: %s", err)
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	if res.StatusCode >= 400 {
		logMsg := fmt.Sprintf("Received bad status code %d from Chatwork: %s", res.StatusCode, body)
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	result := postMessageResult{}
	if err = json.Unmarshal([]byte(string(body)), &result); err != nil {
		logMsg := fmt.Sprintf("Failed to unmarshal Chatwork response body: %s", err)
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	log.Debugf("Message sent suceessfully to Chatwork on %v", cl.generateChatworkMessageUrl(result.MessageId))
	return nil
}

func (cl *ChatworkLogger) generateChatworkMessageUrl(messageId string) string {
	url := fmt.Sprintf("https://www.chatwork.com/#!rid%v-%v", cl.config.ToRoomId, messageId)
	log.Debugf("Generated Chatwork message URL: %s", url)
	return url
}
