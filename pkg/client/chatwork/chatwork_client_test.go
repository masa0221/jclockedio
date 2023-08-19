package chatwork_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/masa0221/jclockedio/pkg/client/chatwork"
)

func TestSendMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message_id": "12345"}`))
	}))
	defer server.Close()

	config := &chatwork.ChatworkSendMessageConfig{
		ToRoomId: "1001",
		Unread:   true,
	}
	client := chatwork.NewChatworkClient("apiToken", config)
	client.BaseUrl = server.URL // サーバのURLをクライアントにセット

	if err := client.SendMessage("Hello, Chatwork!"); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
}

func TestSendMessage_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	config := &chatwork.ChatworkSendMessageConfig{
		ToRoomId: "1001",
		Unread:   true,
	}
	client := chatwork.NewChatworkClient("apiToken", config)
	client.BaseUrl = server.URL

	if err := client.SendMessage("Hello, Chatwork!"); err == nil {
		t.Fatalf("Expected an error, but got nil")
	}
}
