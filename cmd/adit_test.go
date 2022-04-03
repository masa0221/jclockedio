package cmd

import (
	"testing"
)

func TestGenerateOutputMessage(t *testing.T) {
	messageFormat := "clock: {{ .clock }}, {{ .beforeStatus }} -> {{ .afterStatus }}"
	gotMessage := generateOutputMessage(messageFormat, "12:23:34", "Not attending work", "Working")
	wantMessage := "clock: 12:23:34, Not attending work -> Working"
	if gotMessage != wantMessage {
		t.Errorf("got %s, want %s", gotMessage, wantMessage)
	}
}
