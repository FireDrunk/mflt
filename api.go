package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strings"
)

var TRACKING_DOMAIN = "http://localhost:8080"
var TRACKING_PATH = "/target/"

var memoryStorage = make(map[uuid.UUID]trackingPayload)

type trackingPayload struct {
	Id           string `json:"id"`
	TrackingLink string `json:"tracking_link"`
	OriginalLink string `json:"link"`
}

type linkRequest struct {
	Link string   `json:"link"`
	Tags []string `json:"tags"`
}

func addLink(c *gin.Context) {
	var newLinkRequest linkRequest
	err := c.BindJSON(&newLinkRequest)
	if err != nil {
		return
	}

	var generatedUuid = uuid.New()

	trackingLink, err := generateTrackingLink(generatedUuid, newLinkRequest.Link)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	newTrackingPayload := trackingPayload{
		Id:           generatedUuid.String(),
		OriginalLink: newLinkRequest.Link,
		TrackingLink: trackingLink,
	}

	memoryStorage[generatedUuid] = newTrackingPayload

	c.IndentedJSON(http.StatusCreated, newTrackingPayload)
}

func getLinks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, memoryStorage)
}

func getLinkByUUID(c *gin.Context) {
	uuidString := c.Param("uuid")
	var parsedUuid, err = uuid.Parse(uuidString)

	if err != nil {
		fmt.Printf("Error parsing uuid string:: %v", uuidString)
		c.Status(http.StatusBadRequest)
		return
	}

	linkResponse, ok := memoryStorage[parsedUuid]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	c.IndentedJSON(http.StatusOK, linkResponse)
}

func generateTrackingLink(trackingUuid uuid.UUID, link string) (string, error) {
	// Generate a new trackingPayload struct
	newTrackingPayload := trackingPayload{
		Id:           trackingUuid.String(),
		OriginalLink: link,
	}

	// Convert the trackingPayload to JSON
	newTrackingPayloadJson, err := json.Marshal(newTrackingPayload)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error marshalling payload to json: %v", newTrackingPayloadJson))
	}

	// Base64 encode the trackingPayload
	encodedString := base64.StdEncoding.EncodeToString(newTrackingPayloadJson)

	// Generate the Tracking URL
	parsed, err := url.Parse(TRACKING_DOMAIN)
	if err != nil {
		fmt.Printf("An invalid TRACKING_DOMAIN was specified")
		panic(err)
	}
	parsed.Path = fmt.Sprintf("%s/%s", strings.TrimRight(TRACKING_PATH, "/"), encodedString)

	return parsed.String(), nil
}
