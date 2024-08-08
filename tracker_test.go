package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTracker(t *testing.T) {
	r := SetUpRouter()

	r.POST("/link", addLink)
	r.GET("/target/:payload", trackUrl)
	r.GET("/tracking", getTrackingInfo)

	// Payload
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

	// Parse the response to reuse the ID
	parsedTrackingPayload := trackingPayload{}
	parseError := json.Unmarshal(w.Body.Bytes(), &parsedTrackingPayload)
	if parseError != nil {
		t.Error("Error unmarshalling response")
	}

	// Reset the recorder
	w = httptest.NewRecorder()

	// Follow the Tracking Link
	trackingRequest, _ := http.NewRequest("GET", parsedTrackingPayload.TrackingLink, nil)
	r.ServeHTTP(w, trackingRequest)

	// Expect a redirect
	assert.Equal(t, http.StatusFound, w.Code)
}
