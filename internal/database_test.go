package internal_test

import (
	"database/sql"
	internal "foc_api/internal"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// sets up a mock database in memory
func setUpTestDB(t *testing.T) *sql.DB {
	db, err := internal.InitDB(":memory:")
	require.NoError(t, err)
	return db
}

func TestInitDB(t *testing.T) {
	result, err := internal.InitDB("../database/db.sqlite")

	if err != nil {
		t.Errorf("InitDB() failed: %v", err)
	}
	defer result.Close()

	if result.Stats().OpenConnections < 1 {
		t.Errorf("No open connections to database")
	}
}
