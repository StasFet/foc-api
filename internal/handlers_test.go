package internal_test

import (
	"bytes"
	"encoding/json"
	internal "foc_api/internal"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePerformanceEndpoint(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	dbw := internal.CreateDBWrapper(db)
	api := internal.NewAPI(dbw)

	testPerformance := getTestPerformance()
	testPerformanceJson, _ := json.Marshal(testPerformance)

	r := httptest.NewRequest("POST", "/performers", bytes.NewBuffer(testPerformanceJson))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// act
	api.CreateNewPerformance(w, r)

	// assert
	assert.Equal(t, w.Code, http.StatusCreated, "CreateNewPerformance() handler returned status %v", w.Code)

	var createdPerformance *internal.Performance
	err := json.NewDecoder(w.Body).Decode(&createdPerformance)
	require.NoError(t, err, "Error decoding response from endpoint: %v", err)

	createdPerformance, err = dbw.GetPerformanceById(createdPerformance.Id)
	require.NoError(t, err, "GetPerformanceById() has failed while checking creation of performance: %v", err)

	assert.Equal(t, testPerformance.ItemName, createdPerformance.ItemName, "ItemName fields are not equal")
	assert.Equal(t, testPerformance.GenreName, createdPerformance.GenreName, "GenreName fields are not equal")
	assert.Equal(t, testPerformance.GroupName, createdPerformance.GroupName, "GroupName fields are not equal")
	assert.Equal(t, testPerformance.Location, createdPerformance.Location, "Location fields are not equal")
	assert.True(t, testPerformance.StartTime.Equal(createdPerformance.StartTime), "Start times not equal!")
	assert.True(t, testPerformance.EndTime.Equal(createdPerformance.EndTime), "End times not equal!")
}
