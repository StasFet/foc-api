package internal_test

import (
	internal "foc_api/internal"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

//aDIeu

func getTestPerformance() *internal.Performance {
	testTime1 := time.Now()
	testTime2 := testTime1.Add(time.Hour)

	testPerformance := &internal.Performance{
		ItemName:  "Test ItemName",
		GenreName: "Test GenreName",
		GroupName: "Test GroupName",
		Location:  "Test Location",
		StartTime: testTime1,
		EndTime:   testTime2,
	}

	return testPerformance
}

func getTestPerformer() *internal.Performer {
	return &internal.Performer{
		Name:  "TestFirstName, TestLastName",
		Email: "test@test.com",
	}
}

func TestCreatePerformance(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	testPerformance := getTestPerformance()

	result, err := dbw.CreatePerformance(testPerformance)
	require.NoError(t, err, "CreatePerformance() has failed: %v", err)

	// main testable thing about createperformance is that it returns an id
	assert.IsType(t, 0, result.Id, "Performance ID not assigned correctly")
	assert.NotEqual(t, result.Id, 0, "Performance ID assigned as 0")
}

func TestCreatePerformer(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	testPerformer := getTestPerformer()

	result, err := dbw.CreatePerformer(testPerformer)
	require.NoError(t, err, "CreatePerformer() has failed: %v", err)

	assert.IsType(t, 0, result.Id, "Performer ID not assigned correctly")
	assert.NotEqual(t, result.Id, 0, "Performer ID assigned as 0")
}

func TestGetPerformance(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	expected := getTestPerformance()

	performance, err := dbw.CreatePerformance(expected)
	require.NoError(t, err, "Error creating performance while testing getPerformance")

	expected = performance

	actual, err := dbw.GetPerformanceById(expected.Id)
	require.NoError(t, err, "GetPerformanceById() has failed: %v", err)

	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.ItemName, actual.ItemName)
	assert.Equal(t, expected.GenreName, actual.GenreName)
	assert.Equal(t, expected.GroupName, actual.GroupName)
	assert.Equal(t, expected.Location, actual.Location)

	assert.True(t, expected.StartTime.Equal(actual.StartTime), "StartTime not equal")
	assert.True(t, expected.EndTime.Equal(actual.EndTime), "EndTime not equal")

	assert.Equal(t, expected.Duration, actual.Duration)
}

func TestGetPerformer(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	expected := getTestPerformer()

	expected, err := dbw.CreatePerformer(expected)
	require.NoError(t, err, "CreatePerformer() has failed: %v", err)

	actual, err := dbw.GetPerformerById(expected.Id)
	require.NoError(t, err, "GerPerformerById() has failed: %v", err)

	assert.Equal(t, expected, actual, "Expected and Actual performers not equal")
}
