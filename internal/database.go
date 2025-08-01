package internal

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite", filepath)

	// checks for errors within sql.Open()
	if err != nil {
		log.Fatal(err)
	}

	// checks for errors when pinging the db
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables(db)
	return db
}

func createTables(db *sql.DB) {
	createPerformancesString := `
		CREATE TABLE IF NOT EXISTS performances (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			itemName TEXT NOT NULL,
			genreName TEXT NOT NULL,
			groupName TEXT NOT NULL,
			location TEXT NOT NULL,
			startTime DATETIME,
			endTime DATETIME
		);
	`

	createPerformersString := `
		CREATE TABLE IF NOT EXISTS performers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		);
	`

	createJunctionString := `
		CREATE TABLE IF NOT EXISTS junction (
			performer_id INTEGER NOT NULL,
			performance_id INTEGER NOT NULL,
			PRIMARY KEY (performer_id, performance_id),
			FOREIGN KEY (performer_id) REFERENCES performers(id),
			FOREIGN KEY (performance_id) REFERENCES performances(id)
		);
	`

	// creates performances table
	_, err := db.Exec(createPerformancesString)
	if err != nil {
		log.Fatal(err)
	}

	// creates performers table
	_, err = db.Exec(createPerformersString)
	if err != nil {
		log.Fatal(err)
	}

	// creates junction table that stores pairs of performers and performances
	_, err = db.Exec(createJunctionString)
	if err != nil {
		log.Fatal(err)
	}
}
