package main_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	app "ps-app"
	"testing"
)

const testData = "{\"data\":[{\"id\":1,\"name\":\"product\"}]}"

func TestMain(m *testing.M) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := io.WriteString(w, testData); err != nil {
			fmt.Printf("FATAL: %s\n", err.Error())
			os.Exit(1)
		}
	}))
	defer server.Close()

	if err := os.Setenv(app.ProductsURLEnvVar, server.URL); err != nil {
		fmt.Printf("FATAL: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func createFakeProductsRequest() *http.Response {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	app.HandleProducts(w, req)
	return w.Result()
}

func TestRetrieveProducts(t *testing.T) {
	resp := createFakeProductsRequest()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, string(body), "<th>name</th>")
	assert.Contains(t, string(body), "<td>product</td>")
	assert.NotContains(t, string(body), "<td>2</td>")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUnsetProductsURL(t *testing.T) {
	oldProductUrl := os.Getenv(app.ProductsURLEnvVar)
	defer func() {
		if err := os.Setenv(oldProductUrl, app.ProductsURLEnvVar); err != nil {
			t.Fatal(err)
		}
	}()

	if err := os.Setenv(app.ProductsURLEnvVar, ""); err != nil {
		t.Fatal(err)
	}

	resp := createFakeProductsRequest()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, string(body), "url error")
}
