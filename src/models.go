package main

import (
	"database/sql"
	"time"
)

type Performance struct {
	Id        int       `json:"id"`
	ItemName  string    `json:"itemName"`
	GenreName string    `json:"genreName"`
	GroupName string    `json:"groupName"`
	Location  string    `json:"location"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  int       `json:"duration"`
}

type Performer struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"pwHash"`
}

// just a little wrapper so we can make actions methodic rather than functional
type DBWrapper struct {
	db *sql.DB
}

func CreateDBWrapper(db *sql.DB) *DBWrapper {
	return &DBWrapper{db}
}

func (dbw *DBWrapper) CreatePerformance(p *Performance) (*Performance, error) {
	dbQuery := `
		INSERT INTO performances (itemName, genreName, groupName, location, startTime, endTime, duration)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	err := dbw.db.QueryRow(dbQuery, p.ItemName, p.GenreName, p.GroupName, p.Location, p.StartTime, p.EndTime, p.Duration).
		Scan(&p.Id)

	if err != nil {
		return nil, err
	}

	return p, nil
}
