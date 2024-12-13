package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService  struct {
	db   * sql.DB
}

func  NewUserService(db  *sql.DB)  *UserService{
	return  &UserService{db: db}
}


type User struct  {
	Id   string   				`json:"id"`
	Email  string      			`json:"email"`
    Username  string  			`json:"username"`
    Password  string 			`json:"password,omitempty"`
    ConfirmPassword  string  	`json:"confirmPassword,omitempty"`
	CreatedAt time.Time 		`json:"created_at"`

}

func validateUserInput(user User) error {
	if len(user.Username) < 3 || len(user.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("invalid email format")
	}

	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}
func (s *UserService) isUserExists(username, email string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM users 
		WHERE username = $1 OR email = $2
	`

	var count int
	err := s.db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}



func (s *UserService) CreateUser(user *User) error {
    // Validate user input
    if err := validateUserInput(*user); err != nil {
        return fmt.Errorf("invalid user input: %w", err)
    }

    exists, err := s.isUserExists(user.Username, user.Email)
    if err != nil {
        return fmt.Errorf("error checking user existence: %w", err)
    }
    if exists {
        return errors.New("user or email already exists")
    }

    userID := uuid.New().String()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

    // Prepare SQL query
    query := `
        INSERT INTO users (id, username, email, password, created_at) 
        VALUES ($1, $2, $3, $4, $5)
    `

    // Execute database insert
    _, err = s.db.Exec(query, 
        userID, 
        user.Username, 
        user.Email, 
        hashedPassword, 
        time.Now(),
    )
    if err != nil {
        return fmt.Errorf("failed to create user in database: %w", err)
    }

    return nil
}

func (s *UserService) FindByUsername(username string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at 
		FROM users 
		WHERE username = $1
	`

	var user User
	err := s.db.QueryRow(query, username).Scan(
		&user.Id, 
		&user.Username, 
		&user.Email, 
		&user.Password, 
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, err
	}

	return &user, nil
}


func (s *UserService) GetProfileUser(userId string) (*User, error) {
    // Prepare the query to select user profile information
    query := `
        SELECT id, username, email, created_at
        FROM users
        WHERE id = $1
    `

    // Create a User variable to scan the results into
    var user User

    // Use QueryRow to retrieve a single row
    err := s.db.QueryRow(query, userId).Scan(
        &user.Id, 
        &user.Username, 
        &user.Email, 
        &user.CreatedAt,
    )

    // Handle potential errors
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    if err != nil {
        log.Printf("Error finding user: %v", err)
        return nil, fmt.Errorf("failed to retrieve user: %w", err)
    }

    // Return the user pointer and nil error
    return &user, nil
}
