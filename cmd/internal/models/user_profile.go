package models

import (
	"database/sql"
	"log"
)


type  ProfileUSer struct  {}


func (t * ProfileUSer)  CreateProfileUserTable(db *sql.DB) {
    query := `
        CREATE TABLE IF NOT EXISTS profiles (
            id VARCHAR(36) PRIMARY KEY,
            userId VARCHAR(36) UNIQUE NOT NULL,
            fullName VARCHAR(200),
            bio TEXT,
            phoneNumber VARCHAR(20),
            address TEXT,
            avatarUrl VARCHAR(255),
            
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        )
    `
    _, err := db.Exec(query)
    if err != nil {
        log.Fatalf("Error creating profile table: %v", err)
    }
    log.Println("Profile table created successfully")
}