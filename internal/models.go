package internal

import (
	"database/sql"
	"errors"
	"time"
)

// TODO: add slices of performers to each performance
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

// TODO: add slices of performances to each performer
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
		ORDER BY id ASC
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

// Returns all the performances associated with a particular performer
func (dbw *DBWrapper) GetPerformancesByPerformerId(performerId int) ([]*Performance, error) {
	dbQuery := `
		SELECT id, itemName, genreName, groupName, location, startTime, endTime
		FROM performances AS p
		JOIN junction AS j ON p.id = j.performance_id
		WHERE j.performer_id = ?
		ORDER BY p.id ASC
	`

	rows, err := dbw.db.Query(dbQuery, performerId)
	if err != nil {
		return nil, err
	}

	performances := []*Performance{}
	for rows.Next() {
		// scans the id of each performance
		p := &Performance{}
		err := rows.Scan(&p.Id, &p.ItemName, &p.GenreName, &p.GroupName, &p.Location, &p.StartTime, &p.EndTime)
		if err != nil {
			return nil, err
		}
		performances = append(performances, p)
	}

	return performances, nil
}

// Returns all the performers associated with a particular performance
func (dbw *DBWrapper) GetPerformersByPerformanceId(performanceId int) ([]*Performer, error) {
	dbQuery := `
		SELECT id, name, email
		FROM performers AS p
		JOIN junction AS j ON p.id = j.performer_id
		WHERE j.performance_id = ?
		ORDER BY p.id ASC
	`

	rows, err := dbw.db.Query(dbQuery, performanceId)
	if err != nil {
		return nil, err
	}

	performers := []*Performer{}
	for rows.Next() {
		p := &Performer{}
		err := rows.Scan(&p.Id, &p.Name, &p.Email)
		if err != nil {
			return nil, err
		}

		performers = append(performers, p)
	}

	return performers, nil
}

// Return the performance with the given id
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

// Return the performer with the given id
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

// Return all the performances that match a certain query, NOT to be exposed to api endpoint
func (dbw *DBWrapper) GetPerformancesUsingQuery(query string) ([]*Performance, error) {
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

// Get all the performers that match a given query. NOT to be exposed to any api endpoints.
func (dbw *DBWrapper) GetPerformerUsingQuery(query string) ([]*Performer, error) {
	rows, err := dbw.db.Query(query)
	if err != nil {
		return nil, err
	}

	performers := []*Performer{}
	for rows.Next() {
		p := &Performer{}
		err := rows.Scan(&p.Id, &p.Name, &p.Email)
		if err != nil {
			return nil, err
		}
		performers = append(performers, p)
	}

	return performers, nil
}

// TODO: change from just delete to archive
// Deletes the performance with the given id
func (dbw *DBWrapper) DeletePerformanceById(id int) error {
	dbQuery := `
		DELETE FROM performances WHERE id = ?;
	`
	_, err := dbw.db.Exec(dbQuery, id)
	if err != nil {
		return err
	}
	return nil
}

// TODO: change from just delete to archive
// Deletes the performer with the given id
func (dbw *DBWrapper) DeletePerformerById(id int) error {
	dbQuery := `
		DELETE FROM performers WHERE id = ?;
	`
	_, err := dbw.db.Exec(dbQuery, id)
	if err != nil {
		return err
	}
	return nil
}

// Updates the performance with the given id to have the details of the given performance
func (dbw *DBWrapper) UpdatePerformanceById(id int, p *Performance) (*Performance, error) {
	dbQuery := `
		UPDATE performances
		SET itemName = ?, genreName = ?, groupName = ?, location = ?, startTime = ?, endTime = ?
		WHERE id = ?
	`

	result, err := dbw.db.Exec(dbQuery, p.ItemName, p.GenreName, p.GroupName, p.Location, p.StartTime, p.EndTime)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	// error if no matching rows were found and updated
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return p, nil
}

// Updates the performer with the given id to have the details of the given performer
func (dbw *DBWrapper) UpdatePerformerById(id int, p *Performer) (*Performer, error) {
	dbQuery := `
		UPDATE performers
		SET name = ?, email = ?
		WHERE id = ?
	`

	result, err := dbw.db.Exec(dbQuery, p.Name, p.Email)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	// error if no matching rows were found and updated
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return p, nil
}

// creates a performer:performance relationship
func (dbw *DBWrapper) CreateJunction(performerId, performanceId int) error {
	dbQuery := `
		INSERT INTO junction (performer_id, performance_id)
		VALUES (?, ?)
	`

	row := dbw.db.QueryRow(dbQuery, performerId, performanceId)
	if row == nil {
		return errors.New("Error creating junction")
	}
	return nil
}

// deletes the performerId:performanceId pair
func (dbw *DBWrapper) DeleteJunction(performerId, performanceId int) error {
	dbQuery := `
		DELETE FROM junction WHERE performer_id = ? AND performance_id = ?;
	`
	_, err := dbw.db.Exec(dbQuery)
	if err != nil {
		return err
	}
	return nil
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
