package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleSuccessResponse(c *gin.Context, data interface{}) {
	response := gin.H{
		"message":     "",
		"content":     data,
		"status_code": http.StatusOK,
		"error":       "",
	}

	jsonBytes, err := json.Marshal(response)

	if err != nil {
		http.Error(c.Writer, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w := c.Writer

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonBytes)))
	w.Header().Set("Connection", "close")
	w.Header().Del("Content-Encoding")
	w.WriteHeader(http.StatusOK)

	w.Write(jsonBytes)
}

func HandleErrorResponse(c *gin.Context, err error) {
	var resp interface{}
	c.JSON(http.StatusBadRequest, gin.H{
		"error":       err.Error(),
		"content":     resp,
		"status_code": http.StatusBadRequest,
		"message":     err.Error(),
	})
}
