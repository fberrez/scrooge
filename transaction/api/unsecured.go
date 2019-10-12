package api

import "github.com/gin-gonic/gin"

type (
	// Health represents the API health status.
	Health struct {
		// Status is the API health status.
		Status string `json:"status" description:"API health status."`
	}
)

// health returns the API health status.
func (a *API) health(c *gin.Context) (*Health, error) {
	return &Health{
		Status: "OK",
	}, nil
}
