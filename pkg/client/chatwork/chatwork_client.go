package chatwork

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
	endpoint := fmt.Sprintf("%v/rooms/%v/messages", cc.BaseUrl, cc.sendMessageConfig.ToRoomId)
	self_unread := 0
	if cc.sendMessageConfig.Unread {
		self_unread = 1
	}
	payload := strings.NewReader(fmt.Sprintf("self_unread=%v&body=%v", self_unread, url.QueryEscape(message)))

	req, err := http.NewRequest("POST", endpoint, payload)
	if err != nil {
		return fmt.Errorf("Failed to create request: %s", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-ChatWorkToken", cc.apiToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to request to Chatwork: %s", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body of ChatworkAPI: %s", err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("Bad response status code %d :%v", res.StatusCode, string(body))
	}

	result := postMessageResult{}
	err = json.Unmarshal([]byte(string(body)), &result)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal from response body. err: %s", err)
	}

	return nil
}

func (cc *ChatworkClient) generateChatworkMessageUrl(messageId string) string {
	return fmt.Sprintf("https://www.chatwork.com/#!rid%v-%v", cc.sendMessageConfig.ToRoomId, messageId)
}
