package auth

import (
	"errors"
	"golang/cmd/internal/auth/dto"
	"golang/cmd/internal/auth/middleware"
	"golang/cmd/internal/database"
	"golang/cmd/internal/user"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService  struct {
	userService   *user.UserService
	db     *database.Database
}

func  NewAuthService (userService *user.UserService,db  *database.Database)  * AuthService{
	return  &AuthService{
		userService: userService,
		db:  db,
	}
}

func (auth *AuthService) Register(dtos dto.RegisterDTO)  error {

	if dtos.Password  !=  dtos.ConfirmPassword  {
		return errors.New("Password do not match")
	}

	newUser :=  user.User{
		Username: dtos.Username,
		Email:  dtos.Email,
		Password: dtos.Password,
	}

	err  :=  auth.userService.CreateUser(&newUser)
	if err != nil {
		log.Fatal("error",err)
		return err
	}
	return nil
}

func (auth *AuthService) Login(dtos dto.LoginDTO) (string, error) {
	user, err := auth.userService.FindByUsername(dtos.Username)
	if err != nil {
		return "", errors.New("user not found")
	}

	log.Println(dtos.Password)
	log.Println(user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dtos.Password))
	if err != nil {
		log.Printf("Password comparison error: %v\n", err)
		return "", errors.New("invalid credentials")
	}

	token, err := middleware.GenerateToken(user.Id, user.Email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}