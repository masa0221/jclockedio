package jobcan_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/masa0221/jclockedio/internal/jobcan"
)

func TestAdit(t *testing.T) {
	wantUserEmail := "test@example.com"
	wantUserPassword := "dummy-password"
	wantBeforeWorkingStatus := "Not attending work"
	wantAfterWorkingStatus := "Working"
	wantClock := "12:23:34"

	jobcanLoginPageHtml := `
    <html><body>
      <form action="/users/sign_in" method="post">
        <input type="text" name="user[email]" id="user_email" />
        <input type="text" name="user[client_code]" id="user_client_code" />
        <input type="password" name="user[password]" id="user_password" />
        <input type="submit" name="commit" class="form__login" />
      </form>
    </body></html>`

	jobcanAditPageHtml := fmt.Sprintf(`
    <html>
      <script>
        document.getElementById('adit-button').click = document.getElementById('working_status').innerHTML= '%v';
      </script>
    <body>
      <h3 id="working_status">%v</h3>
      <div id="clock">%v</div>
      <div id="adit-button">
        <span id="adit-button-push">PUSH</span>
      </div>
    </body></html>`, wantAfterWorkingStatus, wantBeforeWorkingStatus, wantClock)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/jbcoauth/login" {
			fmt.Fprintln(w, jobcanLoginPageHtml)
		}
		if r.Method == "POST" && r.URL.Path == "/users/sign_in" {
			if err := r.ParseForm(); err != nil {
				t.Errorf("ParseForm() err: %v", err)
				return
			}
			gotUserEmail := r.PostFormValue("user[email]")
			if gotUserEmail != wantUserEmail {
				t.Errorf("want %v, got %v", wantUserEmail, gotUserEmail)
			}
			gotUserPassword := r.PostFormValue("user[password]")
			if gotUserPassword != wantUserPassword {
				t.Errorf("want %v, got %v", wantUserPassword, gotUserPassword)
			}
		}

		if r.Method == "GET" && r.URL.Path == "/employee" {
			fmt.Fprintln(w, jobcanAditPageHtml)
		}
	})
	ts := httptest.NewServer(h)
	defer ts.Close()

	jobcanClient := jobcan.New(wantUserEmail, wantUserPassword)
	jobcanClient.BaseUrl = ts.URL
	jobcanClient.Verbose = true
	// TODO: remove this
	jobcanClient.NoAdit = true
	aditResult := jobcanClient.Adit()
	gotBeforeWorkingStatus := aditResult.BeforeWorkingStatus
	gotAfterWorkingStatus := aditResult.AfterWorkingStatus
	gotClock := aditResult.Clock

	if gotBeforeWorkingStatus != wantBeforeWorkingStatus {
		t.Errorf("want %v, got %v", wantBeforeWorkingStatus, gotBeforeWorkingStatus)
	}
	if gotAfterWorkingStatus != wantAfterWorkingStatus {
		t.Errorf("want %v, got %v", wantAfterWorkingStatus, gotAfterWorkingStatus)
	}
	if gotClock != wantClock {
		t.Errorf("want %v, got %v", wantClock, gotClock)
	}
}
