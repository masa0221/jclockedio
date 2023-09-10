package chatwork_logger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	logger "github.com/masa0221/jclockedio/pkg/logger/chatwork_logger"
)

func TestLog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message_id": "12345"}`))
	}))
	defer server.Close()

	config := &logger.Config{
		ToRoomId: "1001",
		Unread:   true,
	}
	client := logger.NewChatworkLogger("apiToken", config)
	client.BaseUrl = server.URL

	if err := client.Log("Hello, logger!"); err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
}

func TestLog_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	config := &logger.Config{
		ToRoomId: "1001",
		Unread:   true,
	}
	client := logger.NewChatworkLogger("apiToken", config)
	client.BaseUrl = server.URL

	if err := client.Log("Hello, logger!"); err == nil {
		t.Fatalf("Expected an error, but got nil")
	}
}
