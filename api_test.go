package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestNewLink(t *testing.T) {
	r := SetUpRouter()

	r.GET("/link", getLinks)
	r.GET("/link/:uuid", getLinkByUUID)
	r.POST("/link", addLink)

	newLinkRequest := linkRequest{
		Link: "https://www.google.com",
	}

	jsonPayload, parsingError := json.Marshal(newLinkRequest)
	if parsingError != nil {
		t.Error(parsingError.Error())
	}

	w := httptest.NewRecorder()

	// Create a new Tracking Link
	createRequest, _ := http.NewRequest("POST", "/link", bytes.NewBuffer(jsonPayload))
	r.ServeHTTP(w, createRequest)

	// Check if creating the Tracking Link succeeded
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse the response to reuse the ID
	parsedTrackingPayload := trackingPayload{}
	parseError := json.Unmarshal(w.Body.Bytes(), &parsedTrackingPayload)
	if parseError != nil {
		t.Error("Error unmarshalling response")
	}

	// Reset the recorder
	w = httptest.NewRecorder()

	// Check if we can retrieve the stored information with the generated UUID
	getRequest, _ := http.NewRequest("GET", fmt.Sprintf("/link/%s", parsedTrackingPayload.Id), nil)
	r.ServeHTTP(w, getRequest)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response to reuse the ID
	parsedRetrievedTrackingPayload := trackingPayload{}
	retrieveParseError := json.Unmarshal(w.Body.Bytes(), &parsedRetrievedTrackingPayload)
	if retrieveParseError != nil {
		t.Error("Error unmarshalling response")
	}
	assert.Equal(t, parsedTrackingPayload, parsedRetrievedTrackingPayload)
}
