package auth

import (
	"golang/cmd/internal/database"
	"golang/cmd/internal/user"
)

type AuthModule struct {
	Controller *AuthController
	Service    *AuthService
}

func NewAuthModule(userService *user.UserService, db *database.Database) *AuthModule {
	authService := NewAuthService(userService, db)
	authController := NewAuthController(authService,userService)

	return &AuthModule{
		Controller: authController,
		Service:    authService,
	}
}