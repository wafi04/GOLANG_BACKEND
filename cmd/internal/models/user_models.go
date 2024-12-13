package models

import (
	"database/sql"
	"log"
)

type UserTable struct{}

func (t * UserTable)  CreateTableUser(db *sql.DB) {
    query  := `
        CREATE TABLE IF NOT EXISTS users ( 
            id VARCHAR(36) PRIMARY KEY,
            username VARCHAR(100) UNIQUE NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            password VARCHAR(100) NOT NULL,
            role VARCHAR(50) DEFAULT 'user',
            status VARCHAR(20) DEFAULT 'active',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `
    _, err := db.Exec(query)
    if err != nil {
        log.Fatalf("Error creating user table: %v", err) 
    }

    log.Println("User table created successfully")
}

