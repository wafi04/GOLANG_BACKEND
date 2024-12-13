package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"golang/cmd/internal/auth"
	"golang/cmd/internal/database"
	"golang/cmd/internal/user"
)

// Struktur dependencies yang dapat dikonfigurasi
type AuthDependencies struct {
	UserService     *user.UserService
	AuthService     *auth.AuthService
	AuthController  *auth.AuthController
}

type Dependencies struct {
	DB          *database.Database
	AuthModule  AuthDependencies

}

type Server struct {
	port        int
	db          *database.Database
	router      *gin.Engine
	deps        Dependencies
}

func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Printf("Incoming Request: %s %s", c.Request.Method, c.Request.URL.Path)
        c.Next()
    }
}

func NewServer(deps Dependencies) *Server {
    port, err := strconv.Atoi(os.Getenv("PORT"))
    if err != nil {
        port = 8080
    }

    // Inisialisasi router Gin
    router := gin.New()
    
	// added logger for debugging
	router.Use(gin.Logger())
    router.Use(gin.Recovery())

    // Buat instance Server
    srv := &Server{
        port:   port,
        db:     deps.DB,
        router: router,
        deps:   deps,
    }

	// call allroutes 
	srv.RegisterRoutes()

    return srv
}

// HTTP Server dengan konfigurasi timeout
func (s *Server) HttpServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

