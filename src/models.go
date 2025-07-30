package main

import (
	"database/sql"
	"time"
)

/*
	Checklist for each entity:
		X Create
		X GetAll
		O GetById
		O Update
		O Delete
*/

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
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// just a little wrapper so we can make actions methodic rather than functional
type DBWrapper struct {
	db *sql.DB
}

func CreateDBWrapper(db *sql.DB) *DBWrapper {
	return &DBWrapper{db}
}

// creates a performance and puts it into the db
func (dbw *DBWrapper) CreatePerformance(p *Performance) (*Performance, error) {
	dbQuery := `
		INSERT INTO performances (itemName, genreName, groupName, location, startTime, endTime)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id
	`
	// the arguments after dbQuery get formatted into the ?s in the VALUES. this is an anti-injection measure
	err := dbw.db.QueryRow(dbQuery, p.ItemName, p.GenreName, p.GroupName, p.Location, p.StartTime, p.EndTime).
		Scan(&p.Id)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// creates a performer and puts it into the db
func (dbw *DBWrapper) CreatePerformer(p *Performer) (*Performer, error) {
	dbQuery := `
		INSERT INTO performers (name, email)
		VALUES (?, ?)
		RETURNING id
	`

	err := dbw.db.QueryRow(dbQuery, p.Name, p.Email).
		Scan(&p.Id)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// returns a slice with all the performances in the db
func (dbw *DBWrapper) GetAllPerformances() ([]*Performance, error) {
	dbQuery := `
		SELECT *
		FROM performances
		ORDER BY id startTime ASC
	`

	rows, err := dbw.db.Query(dbQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// seed performances
	var performances []*Performance
	for rows.Next() {
		p, err := getNextPerformance(rows)
		if err != nil {
			return nil, err
		}
		performances = append(performances, p)
	}

	return performances, nil
}

// returns a slice with all the performers in the db
func (dbw *DBWrapper) GetAllPerformers() ([]*Performer, error) {
	dbQuery := `
		SELECT *
		FROM performers
		ORDER BY id ASC
	`

	rows, err := dbw.db.Query(dbQuery)
	if err != nil {
		return nil, err
	}

	var performers []*Performer
	for rows.Next() {
		p, err := getNextPerformer(rows)
		if err != nil {
			return nil, err
		}
		performers = append(performers, p)
	}

	return performers, nil
}

// return the performance with the given id
func (dbw *DBWrapper) GetPerformanceById(id int) (*Performance, error) {
	dbQuery := `
		SELECT *
		FROM performances
		WHERE id = ?
	`

	p := &Performance{}
	err := dbw.db.QueryRow(dbQuery, id).
		Scan(&p.Id, &p.ItemName, &p.GenreName, &p.GroupName, &p.Location, &p.StartTime, &p.EndTime)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (dbw *DBWrapper) GetPerformerById(id int) (*Performer, error) {
	dbQuery := `
		SELECT *
		FROM performers
		WHERE id = ?
	`

	p := &Performer{}
	err := dbw.db.QueryRow(dbQuery, id).
		Scan(&p.Id, &p.Email, &p.Name)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// return all the performances that match a certain query
func (dbw *DBWrapper) GetRecordsUsingQuery(query string, returnType interface{}) ([]*Performance, error) {
	rows, err := dbw.db.Query(query)
	if err != nil {
		return nil, err
	}

	performances := []*Performance{}
	for rows.Next() {
		p := &Performance{}
		err := rows.Scan(&p.Id, &p.ItemName, &p.GenreName, &p.GroupName, &p.Location, &p.StartTime, &p.EndTime)
		if err != nil {
			return nil, err
		}
		performances = append(performances, p)
	}

	return performances, nil
}

/*


*	Utility Stuff


 */

// gets the head of rows and returns it as a Performance
func getNextPerformance(rows *sql.Rows) (*Performance, error) {
	p := &Performance{}
	err := rows.Scan(&p.Id, &p.ItemName, &p.GenreName, &p.GroupName, &p.Location, &p.StartTime, &p.EndTime)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// gets the head of rows and returns it as a Performer
func getNextPerformer(rows *sql.Rows) (*Performer, error) {
	p := &Performer{}
	err := rows.Scan(&p.Id, &p.Name, &p.Email)
	if err != nil {
		return nil, err
	}
	return p, nil
}
