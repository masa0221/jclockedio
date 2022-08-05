package chatwork_test

import (
	"fmt"
	"github.com/masa0221/jclockedio/internal/chatwork"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessage(t *testing.T) {
	wantRoomId := "12345"
	wantMessage := "[12:23:23] before -> after"
	wantApiToken := "qwert12345"
	wantPath := "/v2/rooms/12345/messages"
	wantMessageId := "0987654321"
	wantContentType := "application/x-www-form-urlencoded"

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedMethod := r.Method
		if r.Method != "POST" {
			t.Errorf("Unsupported http method. %v", requestedMethod)
		}
		gotPath := r.URL.Path
		if gotPath != wantPath {
			t.Errorf("want %v, got %v", wantPath, gotPath)
		}

		gotApiToken := r.Header.Get("X-ChatWorkToken")
		if gotApiToken != wantApiToken {
			t.Errorf("want %v, got %v", wantApiToken, gotApiToken)
		}

		gotContentType := r.Header.Get("Content-Type")
		if gotContentType != wantContentType {
			t.Errorf("want %v, got %v", wantContentType, gotContentType)
		}

		gotMessage := r.PostFormValue("body")
		if gotMessage != wantMessage {
			t.Errorf("want %v, got %v", wantMessage, gotMessage)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("ParseForm() err: %v", err)
			return
		}
		fmt.Fprintln(w, fmt.Sprintf("{\"message_id\": \"%s\"}", wantMessageId))
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	chatworkClient := chatwork.New(wantApiToken)
	// chatworkClient.Verbose = true
	chatworkClient.BaseUrl = ts.URL
	gotMessageId, err := chatworkClient.SendMessage(wantMessage, wantRoomId)
	if err != nil {
		t.Errorf("error occurred. err: %v", err)
	}

	if gotMessageId != wantMessageId {
		t.Errorf("want %v, got %v", wantMessageId, gotMessageId)
	}
}
