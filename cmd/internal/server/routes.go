package server

import (
	"fmt"
	"golang/cmd/internal/auth/middleware"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() {
    log.Println("Registering routes...")

	s.router.Use(middleware.JWTMiddleware())

    // Public routes
    s.router.GET("/", s.HelloWorldHandler)
    s.router.GET("/health", s.healthHandler)
    s.router.GET("/websocket", s.websocketHandler)

    // Auth routes
    // PASTIKAN MENGGUNAKAN POINTER RECEIVER YANG BENAR
    s.router.POST("/auth/register", func(c *gin.Context) {
        log.Println("Register route called")
        s.deps.AuthModule.AuthController.Register(c)
    })

    s.router.POST("/auth/login", func(c *gin.Context) {
        log.Println("Login route called")
        s.deps.AuthModule.AuthController.Login(c)
    })
	s.router.GET("/auth/profile", func(ctx *gin.Context) {
		s.deps.AuthModule.AuthController.GetProfile(ctx)
	})


    // Log routes untuk verifikasi
    routes := s.router.Routes()
    log.Printf("Registered Routes (%d):", len(routes))
    for _, route := range routes {
        log.Printf("Method: %s, Path: %s", route.Method, route.Path)
    }
}
func (s *Server) HelloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

func (s *Server) healthHandler(c *gin.Context) {
	health := s.db.Health()
	c.JSON(http.StatusOK, health)
}

func (s *Server) websocketHandler(c *gin.Context) {
	// Upgrade koneksi ke websocket
	socket, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open websocket",
		})
		return
	}
	defer socket.Close(websocket.StatusGoingAway, "Server closing websocket")

	ctx := c.Request.Context()
	socketCtx := socket.CloseRead(ctx)

	// Kirim pesan berkala
	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		if err := socket.Write(socketCtx, websocket.MessageText, []byte(payload)); err != nil {
			c.Error(err)
			break
		}
		time.Sleep(2 * time.Second)
	}
}

