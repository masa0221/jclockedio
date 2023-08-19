package jobcan_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/masa0221/jclockedio/pkg/client/jobcan"
	"github.com/masa0221/jclockedio/pkg/client/jobcan/browser"
)

func TestJobcanClient(t *testing.T) {
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
	<body>
	  <h3 id="working_status">%v</h3>
	  <div id="clock">%v</div>
	  <div id="adit-button"><span id="adit-button-push">PUSH</span></div>
	</body>
	<script>
	  document.getElementById('adit-button').onclick = function() {
	    setTimeout(function() {
	      document.getElementById('working_status').innerHTML= '%v';
	    }, 10000);
	  };
	</script>
	</html>`, wantBeforeWorkingStatus, wantClock, wantAfterWorkingStatus)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/jbcoauth/login" {
			fmt.Fprintln(w, jobcanLoginPageHtml)
		}
		if r.Method == "POST" && r.URL.Path == "/users/sign_in" {
			if err := r.ParseForm(); err != nil {
				t.Errorf("ParseForm() err: %v", err)
				return
			}
			assertEquals(t, r.PostFormValue("user[email]"), wantUserEmail, "check post data(email)")
			assertEquals(t, r.PostFormValue("user[password]"), wantUserPassword, "check post data(password)")

			w.Header().Set("Content-Type", "text/html")
			w.Header().Set("location", "/employee")
			w.WriteHeader(http.StatusMovedPermanently)
		}

		if r.URL.Path == "/employee" {
			fmt.Fprintln(w, jobcanAditPageHtml)
		}
	})

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	browser, err := browser.NewAgoutiBrowser()
	defer browser.Close()

	if err != nil {
		t.Errorf("Error in NewAgoutiBrowser: %v", err)
	}
	client := jobcan.NewJobcanClient(browser, &jobcan.JobcanCredentials{Email: wantUserEmail, Password: wantUserPassword})
	client.BaseUrl = testServer.URL

	err = client.Login()
	if err != nil {
		t.Errorf("Error in Login: %v", err)
	}

	// Do not adit
	client.NoAdit = true
	aditResult, err := client.Adit()
	if err != nil {
		t.Errorf("Error in Adit: %v", err)
	}
	assertEquals(t, aditResult.BeforeWorkingStatus, wantBeforeWorkingStatus, "check before working status(do not adit)")
	assertEquals(t, aditResult.AfterWorkingStatus, wantBeforeWorkingStatus, "check after working status(do not adit)")
	assertEquals(t, aditResult.Clock, wantClock, "check clock(do not adit)")

	// // Do adit
	aditResult, err = client.Adit()
	if err != nil {
		t.Errorf("Error in Adit: %v", err)
	}
	assertEquals(t, aditResult.BeforeWorkingStatus, wantBeforeWorkingStatus, "check before working status(do not adit)")
	assertEquals(t, aditResult.AfterWorkingStatus, wantAfterWorkingStatus, "check after working status(do not adit)")
	assertEquals(t, aditResult.Clock, wantClock, "check clock(do not adit)")
	assertEquals(t, "", "", "check clock(do not adit)")
}

func assertEquals(t *testing.T, got interface{}, want interface{}, description string) {
	if got != want {
		t.Errorf("FAIL: %v. want %v, got %v", description, want, got)
	} else {
		t.Logf("PASS: %v", description)
	}
}
