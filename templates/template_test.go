package templates

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestThatTheLoginPageCanBeRendered(t *testing.T) {
	w := httptest.NewRecorder()
	RenderLogin(w, LoginModel{GoogleAuthClientID: "the_client_id"})
	body, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Errorf("failed to read body: %v", err)
	}
	if !strings.Contains(string(body), "the_client_id") {
		t.Errorf("expected 'the_client_id', but didn't find it: %v", string(body))
	}
}
