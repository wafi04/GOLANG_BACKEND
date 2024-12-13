package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang/cmd/internal/auth"
	"golang/cmd/internal/database"
	"golang/cmd/internal/server"
	"golang/cmd/internal/user"
)

func gracefulShutdown(srv *http.Server, db *database.Database, done chan bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, 
		syscall.SIGINT,  
		syscall.SIGTERM, 
	)

	<-sigChan

	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Tutup koneksi database
	if err := db.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	log.Println("Server and database shutdown completed")
	done <- true
}

func main() {
	// Inisialisasi logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Inisialisasi database
	db := database.New()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Database close error: %v", err)
		}
	}()

	

	// Inisialisasi dependencies
	deps := server.Dependencies{
		DB: db,
		AuthModule: server.AuthDependencies{
			UserService:     user.NewUserService(db.DB),
			AuthService:     nil, // Akan diisi di bawah
			AuthController:  nil, // Akan diisi di bawah
		},
	}

	// Inisialisasi auth service
	deps.AuthModule.AuthService = auth.NewAuthService(
		deps.AuthModule.UserService, 
		db,
	)

	// Inisialisasi auth controller
	deps.AuthModule.AuthController = auth.NewAuthController(
		deps.AuthModule.AuthService,
		deps.AuthModule.UserService,
	)

	// Buat server dengan dependency injection
	srv := server.NewServer(deps)

	// Dapatkan HTTP Server
	httpServer := srv.HttpServer()

	// Channel untuk sinkronisasi shutdown
	done := make(chan bool, 1)

	// Jalankan goroutine untuk graceful shutdown
	go gracefulShutdown(httpServer, db, done)

	// Log informasi server
	log.Printf("Starting server on %s", httpServer.Addr)

	// Jalankan server
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	// Tunggu proses shutdown selesai
	<-done
	log.Println("Server shutdown completed")
}