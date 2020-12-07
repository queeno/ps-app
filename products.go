package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const ProductsURLEnvVar = "PRODUCTS_URL"

type DataJson struct {
	Data []map[string]interface{} `json:"data"`
}

func HandleProducts(w http.ResponseWriter, _ *http.Request) {
	productsUrl := os.Getenv(ProductsURLEnvVar)
	if _, err := url.ParseRequestURI(productsUrl); err != nil {
		http.Error(w, errors.Wrap(err, "url error").Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.Get(productsUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("request to get products returned %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &DataJson{}
	if err := json.Unmarshal(body, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.New("index.gohtml").ParseFiles("index.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pd := struct {
		Products []map[string]interface{}
	}{
		Products: data.Data,
	}

	if err := t.Execute(w, pd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
