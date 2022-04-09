package cmd

import (
	"testing"
)

func TestGenerateOutputMessage(t *testing.T) {
	tests := []struct {
		name          string
		messageFormat string
		want          string
	}{
		{name: "Use all variables", messageFormat: "clock: {{ .clock }}, {{ .beforeStatus }} -> {{ .afterStatus }}", want: "clock: 12:23:34, Not attending work -> Working"},

		// {{ if eq .afterStatus "Working" }}I'm working now!{{ else if eq .afterStatus "Not attending work" }}I'm done today.See you tomorrow.{{ else }}Opps! Problem occured.{{ end }} at {{ .clock }}
		{name: "Use if syntax", messageFormat: "{{ if eq .afterStatus \"Working\" }}I'm working now!{{ else if eq .afterStatus \"Not attending work\" }}I'm done today.See you tomorrow.{{ else }}Opps! Problem occured.{{ end }} at {{ .clock }}", want: "I'm working now! at 12:23:34"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				t.Log("cleanup!")
			})
			defer t.Log("defer!")

			if got := generateOutputMessage(tt.messageFormat, "12:23:34", "Not attending work", "Working"); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
