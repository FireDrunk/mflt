package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type trackerInfo struct {
	Payload trackingPayload `json:"payload"`

	Time      time.Time `json:"timestamp"`
	UserAgent string    `json:"userAgent"`
	SourceIP  string    `json:"sourceIP"`
}

var trackingStorage = make(map[uuid.UUID]trackerInfo)

func getTrackingInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, trackingStorage)
}

func trackUrl(c *gin.Context) {
	now := time.Now()

	encodedPayload := c.Param("payload")

	trackingPayloadJson, decodeErr := base64.StdEncoding.DecodeString(encodedPayload)
	if decodeErr != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Create new trackingPayload struct
	parsedTrackingPayload := trackingPayload{}

	// Unmarshal JSON into the struct
	unmarshalErr := json.Unmarshal(trackingPayloadJson, &parsedTrackingPayload)
	if unmarshalErr != nil {
		return
	}

	parsedUUID, err := uuid.Parse(parsedTrackingPayload.Id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	originalTrackingPayload, ok := memoryStorage[parsedUUID]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	newTrackerInfo := trackerInfo{
		Payload: originalTrackingPayload,

		Time:      now,
		UserAgent: c.Request.UserAgent(),
		SourceIP:  c.ClientIP(),
	}

	trackingStorage[parsedUUID] = newTrackerInfo

	fmt.Printf("Tracked URL: %v\n", newTrackerInfo.Payload.Id)

	c.Redirect(http.StatusFound, originalTrackingPayload.OriginalLink)
}
