package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, StandardResponse{
		Status:  http.StatusCreated,
		Message: message,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, StandardResponse{
		Status:  http.StatusBadRequest,
		Message: message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, StandardResponse{
		Status:  http.StatusUnauthorized,
		Message: message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, StandardResponse{
		Status:  http.StatusNotFound,
		Message: message,
	})
}

func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, StandardResponse{
		Status:  http.StatusInternalServerError,
		Message: message,
	})
}