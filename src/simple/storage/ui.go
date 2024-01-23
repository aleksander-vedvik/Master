package storage

import (
	"bytes"
	"fmt"
	"net/http"
)

func ui(addr string) {
	requestURL := fmt.Sprintf("http://localhost:8080/register/%s", addr)
	_, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("could not connect to ui: %s\n", err)
	}
}

func sendStatus(addr, method, val string, msgs int) {
	payload := fmt.Sprintf(`{"address": "%s", "messages": "%v", "method": "%s", "value": "%s"}`, addr, msgs, method, val)
	jsonBody := []byte(payload)
	bodyReader := bytes.NewReader(jsonBody)

	requestURL := "http://localhost:8080/node"
	http.NewRequest(http.MethodPost, requestURL, bodyReader)
}
