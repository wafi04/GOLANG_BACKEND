// internal/server/server_test.go
package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorldHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Buat router
	router := gin.New()
	server := &Server{} 
	router.GET("/", server.HelloWorldHandler)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Hello World")
}