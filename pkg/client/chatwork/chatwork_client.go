package chatwork

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

type ChatworkClient struct {
	BaseUrl           string
	apiToken          string
	sendMessageConfig *ChatworkSendMessageConfig
}

type ChatworkSendMessageConfig struct {
	ToRoomId string
	Unread   bool
}

type postMessageResult struct {
	MessageId string `json:"message_id"`
}

func NewChatworkClient(token string, sendMessageConfig *ChatworkSendMessageConfig) *ChatworkClient {
	return &ChatworkClient{
		BaseUrl:           "https://api.chatwork.com/v2",
		sendMessageConfig: sendMessageConfig,
		apiToken:          token,
	}
}

func (cc *ChatworkClient) SendMessage(message string) error {
	log.Debug("Starting to send a message to Chatwork")

	endpoint := fmt.Sprintf("%v/rooms/%v/messages", cc.BaseUrl, cc.sendMessageConfig.ToRoomId)
	self_unread := 0
	if cc.sendMessageConfig.Unread {
		self_unread = 1
	}
	payload := strings.NewReader(fmt.Sprintf("self_unread=%v&body=%v", self_unread, url.QueryEscape(message)))

	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		logMsg := fmt.Sprintf("Failed to create request: %s", err)
		log.Error(logMsg)
		return errors.New(logMsg)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-ChatWorkToken", cc.apiToken)

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

	log.Debug("Message sent successfully to Chatwork")
	return nil
}

func (cc *ChatworkClient) generateChatworkMessageUrl(messageId string) string {
	url := fmt.Sprintf("https://www.chatwork.com/#!rid%v-%v", cc.sendMessageConfig.ToRoomId, messageId)
	log.Debugf("Generated Chatwork message URL: %s", url)
	return url
}
