package it

import (
	"testing"
	"net/http"
	"encoding/json"
)

func requestURL(key string) string {
	return "http://localhost:8080/v1/proxy/" + key
}

func TestGettingData(t *testing.T) {
	response,err := http.Get(requestURL("it-test-1"))

	defer response.Body.Close()
	
	if err != nil {
		t.Errorf("Got a request error %v", err)
	}

	if response.StatusCode != 200 {
		t.Errorf("Expected 200 status code instead got %d", response.StatusCode)
	}

	testData := make(map[string]int)
	testData["data"] = 1234
	
	actualData := map[string]int{}
	json.NewDecoder(response.Body).Decode(&actualData)

	if actualData["data"] != testData["data"] {
		t.Errorf("Data received: %v did not match expected %v data", actualData, testData)
	}
}

func TestDataNotFound(t *testing.T) {
	response,_ := http.Get(requestURL("no data"))	

	if response.StatusCode != 404 {
		t.Errorf("Expected 200 status code instead got %d", response.StatusCode)
	}

}