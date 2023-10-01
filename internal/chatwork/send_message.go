package chatwork

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type chatworkClient struct {
	apiToken string
	BaseUrl  string
	Verbose  bool
}

func New(ApiToken string) chatworkClient {
	client := chatworkClient{}
	client.BaseUrl = "https://api.chatwork.com"
	client.apiToken = ApiToken

	return client
}

type postMessageResult struct {
	MessageId string `json:"message_id"`
}

func (c *chatworkClient) generateEndpointUrl() string {
	return fmt.Sprintf("%s/v2", c.BaseUrl)
}

func (c *chatworkClient) outputVerboseMessage(message string) {
	if c.Verbose {
		log.Println("[chatwork]", message)
	}
}

func (c *chatworkClient) SendMessage(message string, toRoomId string) (string, error) {
	endpoint := fmt.Sprintf("%v/rooms/%v/messages", c.generateEndpointUrl(), toRoomId)
	payload := strings.NewReader(fmt.Sprintf("self_unread=0&body=%v", url.QueryEscape(message)))

	req, err := http.NewRequest("POST", endpoint, payload)
	c.outputVerboseMessage(fmt.Sprintf("Create NewRequest. URL: %v err: %v", endpoint, err))
	if err != nil {
		return "", fmt.Errorf("Failed to create request: %s", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-ChatWorkToken", c.apiToken)

	res, err := http.DefaultClient.Do(req)
	c.outputVerboseMessage(fmt.Sprintf("Requested chatwork. HTTP status code is %v. err: %v", res.StatusCode, err))
	if err != nil {
		return "", fmt.Errorf("Failed to request to Chatwork: %s", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	c.outputVerboseMessage(fmt.Sprintf("Fetched from response body. Body: %v err: %v", string(body), err))
	if err != nil {
		return "", fmt.Errorf("Failed to read response body of ChatworkAPI: %s", err)
	}
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("Bad response status code %d :%v", res.StatusCode, string(body))
	}

	result := postMessageResult{}
	err = json.Unmarshal([]byte(string(body)), &result)
	c.outputVerboseMessage(fmt.Sprintf("Unmarshal to JSON from response body. result: %v err: %v", result, err))
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal from response body. err: %s", err)
	}

	c.outputVerboseMessage(fmt.Sprintf("Successed! Chatwork message URL is %v", generateChatworkMessageUrl(result.MessageId, toRoomId)))
	return result.MessageId, nil
}

func generateChatworkMessageUrl(chatworkId string, roomId string) string {
	return fmt.Sprintf("https://www.chatwork.com/#!rid%v-%v", roomId, chatworkId)
}