package models

import "database/sql"

type TableCreator  interface  {
	CreateTable(db *sql.DB)  error
}


