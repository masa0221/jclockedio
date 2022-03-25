package chatwork

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type chatworkClient struct {
	baseEndpoint string
	apiToken     string
}

func New(ApiToken string) chatworkClient {
	client := chatworkClient{}
	client.baseEndpoint = "https://api.chatwork.com/v2"
	client.apiToken = ApiToken

	return client
}

type postMessageResult struct {
	MessageId string `json:"message_id"`
}

func (c *chatworkClient) SendMessage(message string, toRoomId string) (string, error) {
	url := fmt.Sprintf("%v/rooms/%v/messages?body=%v", c.baseEndpoint, toRoomId, message)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create request: %s", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-ChatWorkToken", c.apiToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Failed to request to Chatwork: %s", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body of ChatworkAPI: %s", err)
	}
	fmt.Println(string(body))
	fmt.Println(res.StatusCode >= 400)
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("Bad response status code %d :%v", res.StatusCode, string(body))
	}

	result := postMessageResult{}
	json.Unmarshal([]byte(string(body)), &result)

	return result.MessageId, nil
}
