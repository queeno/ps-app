package main_test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	app "ps-app"
	"testing"
)

func createFakeHealthzRequest() *http.Response {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	app.HandleHealthz(w, req)
	return w.Result()
}

func TestHealthz(t *testing.T) {
	resp := createFakeHealthzRequest()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(body), "OK")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
