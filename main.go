package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type Health struct {
	Status string `json:"status"`
}

func health(c *gin.Context) {
	c.IndentedJSON(200, Health{Status: "ok"})
}

func main() {
	apiRouter := gin.Default()

	// Healthcheck Endpoint
	apiRouter.GET("/health", health)

	// API Endpoints
	apiRouter.GET("/link", getLinks)
	apiRouter.GET("/link/:uuid", getLinkByUUID)
	apiRouter.POST("/link", addLink)

	apiRouter.GET("/tracking/", getTrackingInfo)

	apiRouter.GET("/target/:payload", trackUrl)

	apiRouterStartErr := apiRouter.Run(`localhost:8080`)
	if apiRouterStartErr != nil {
		fmt.Printf("%+v", apiRouterStartErr)
		os.Exit(1)
	}
}
