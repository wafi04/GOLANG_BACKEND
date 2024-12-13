package models

import (
	"database/sql"
	"log"
)

type  VenuesTable  struct{}


func (t *VenuesTable)  CreateVenuesTable(db *sql.DB){
	query := `
		CREATE TABLE IF NOT EXISTS venues (
			id VARCHAR(36) PRIMARY KEY,
			name VARCHAR(200) NOT NULL UNIQUE,
			description TEXT,
			address VARCHAR(300),
			category VARCHAR(100),
			capacity INTEGER,
			price_per_hour DECIMAL(10,2),
			userId  VARCHAR(36)  NOT NULL,
			
			// Lokasi
			longitude VARCHAR(50),
			latitude VARCHAR(50),
			
			// Fasilitas
			has_parking BOOLEAN DEFAULT false,
			has_wifi BOOLEAN DEFAULT false,
			has_ac BOOLEAN DEFAULT false,
			
			// Status
			is_available BOOLEAN DEFAULT true,
			max_booking_duration INTEGER,
				
			// Waktu operasional
			open_time TIME,
			close_time TIME,
			
			// Audit trail
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			
			// Tambahan untuk booking
			min_booking_hours INTEGER DEFAULT 1,
			max_booking_hours INTEGER DEFAULT 12

			CONSTRAINT unique_venue_contact UNIQUE (name)
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE,

		)
	`
	
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}