package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestApiHandler(t *testing.T) {
	host := "localhost"
	port := "7333"

	url := fmt.Sprintf("http://%s:%s/apis", host, port)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close() // Ensure the response body is closed

	// Read the response body
	_, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	want := "200 OK"
	if res.Status != want {
		t.Errorf("got %q want %q", res.Status, want)
	}

}
