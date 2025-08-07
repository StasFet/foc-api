package internal_test

import (
	internal "foc_api/internal"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

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
		Name:  "FirstName, LastName",
		Email: "test@test.com",
	}
}

func getTestPerformances(amount int) []*internal.Performance {
	template := getTestPerformance()

	res := []*internal.Performance{}
	for i := range amount {
		res = append(res, &internal.Performance{ItemName: template.ItemName + strconv.Itoa(i)})
	}

	return res
}

func getTestPerformers(amount int) []*internal.Performer {
	template := getTestPerformer()

	res := []*internal.Performer{}
	for i := range amount {
		res = append(res, &internal.Performer{Name: template.Name + strconv.Itoa(i), Email: template.Email + strconv.Itoa(i)})
	}

	return res
}

func TestCreatePerformance(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	testPerformance := getTestPerformance()

	// act
	result, err := dbw.CreatePerformance(testPerformance)

	// assert
	require.NoError(t, err, "CreatePerformance() has failed: %v", err)

	// main testable thing about createperformance is that it returns an id
	assert.IsType(t, 0, result.Id, "Performance ID not assigned correctly")
	assert.NotEqual(t, result.Id, 0, "Performance ID assigned as 0")
}

func TestCreatePerformer(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	testPerformer := getTestPerformer()

	// act
	result, err := dbw.CreatePerformer(testPerformer)

	//assert
	require.NoError(t, err, "CreatePerformer() has failed: %v", err)

	assert.IsType(t, 0, result.Id, "Performer ID not assigned correctly")
	assert.NotEqual(t, result.Id, 0, "Performer ID assigned as 0")
}

func TestCreateJunction(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	testPerformer := getTestPerformer()
	testPerformance := getTestPerformance()

	testPerformer, err := dbw.CreatePerformer(testPerformer)
	require.NoError(t, err, "CreatePerformer failed: %v", err)

	expected, err := dbw.CreatePerformance(testPerformance)
	require.NoError(t, err, "CreatePerformance() failed: %v", err)

	// act
	err = dbw.CreateJunction(testPerformer.Id, expected.Id)
	require.NoError(t, err, "CreateJunction() failed: %v", err)
}

func TestGetPerformance(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	expected := getTestPerformance()

	performance, err := dbw.CreatePerformance(expected)
	require.NoError(t, err, "Error creating performance while testing getPerformance")

	expected = performance

	// act
	actual, err := dbw.GetPerformanceById(expected.Id)

	// assert
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
	//arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	expected := getTestPerformer()

	expected, err := dbw.CreatePerformer(expected)
	require.NoError(t, err, "CreatePerformer() has failed: %v", err)

	// act
	actual, err := dbw.GetPerformerById(expected.Id)

	// assert
	require.NoError(t, err, "GerPerformerById() has failed: %v", err)
	assert.Equal(t, expected, actual, "Expected and Actual performers not equal")
}

func TestGetAllPerformances(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	expected := getTestPerformances(10)
	for i, val := range expected {
		p, err := dbw.CreatePerformance(val)
		require.NoError(t, err, "Error creating performance: %v", err)
		expected[i] = p
	}

	// act
	actual, err := dbw.GetAllPerformances()
	require.NoError(t, err, "GetAllPerformances() failed: %v", err)

	// assert

	for i := range expected {
		assert.Equal(t, expected[i], actual[i], "Expected performances and actual values ")
	}
}

func TestGetAllPerformers(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	expected := getTestPerformers(10)
	for i, val := range expected {
		p, err := dbw.CreatePerformer(val)
		require.NoError(t, err, "CreatePerformer() failed %v", err)
		expected[i] = p
	}

	// act
	actual, err := dbw.GetAllPerformers()
	require.NoError(t, err, "GetAllPerformers() failed: %v", err)

	// assert
	assert.True(t, (len(expected) == len(actual)), "Number of returns not equal")
	for i := range expected {
		assert.Equal(t, expected[i], actual[i], "Expected and ")
	}
}

func TestGetPerformancesByPerformerId(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	testPerformer := getTestPerformer()
	expected := getTestPerformances(3)

	testPerformer, err := dbw.CreatePerformer(testPerformer)
	require.NoError(t, err, "CreatePerformer() failed: %v", err)

	for i, val := range expected {
		expected[i], err = dbw.CreatePerformance(val)
		require.NoError(t, err, "CreatePerformance() failed: %v", err)
	}

	for _, val := range expected {
		err := dbw.CreateJunction(testPerformer.Id, val.Id)
		require.NoError(t, err, "CreateJunction() failed: %v", err)
	}

	// act
	actual, err := dbw.GetPerformancesByPerformerId(testPerformer.Id)
	require.NoError(t, err, "GetPerformancesByPerformerId() failed: %v", err)

	// assert
	assert.True(t, (len(expected) == len(actual)), "Number of returns not equal")
	for i, val := range actual {
		assert.Equal(t, val.Id, expected[i].Id, "Expected performances not equal to actual")
		assert.Equal(t, val.ItemName, expected[i].ItemName, "Expected performances not equal to actual")
		assert.Equal(t, val.GroupName, expected[i].GroupName, "Expected performances not equal to actual")
		assert.Equal(t, val.GenreName, expected[i].GenreName, "Expected performances not equal to actual")
		assert.Equal(t, val.Location, expected[i].Location, "Expected performances not equal to actual")
		assert.Equal(t, val.StartTime, expected[i].StartTime, "Expected performances not equal to actual")
		assert.Equal(t, val.EndTime, expected[i].EndTime, "Expected performances not equal to actual")
		assert.Equal(t, val.Duration, expected[i].Duration, "Expected performances not equal to actual")
	}
}

func TestGetPerformersByPerformanceId(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	testPerformance := getTestPerformance()
	expected := getTestPerformers(3)

	testPerformance, err := dbw.CreatePerformance(testPerformance)
	require.NoError(t, err, "CreatePerformance() failed: %v", err)

	for i, val := range expected {
		expected[i], err = dbw.CreatePerformer(val)
		require.NoError(t, err, "CreatePerformer() failed: %v", err)
	}

	for _, val := range expected {
		err := dbw.CreateJunction(val.Id, testPerformance.Id)
		require.NoError(t, err, "CreateJunction() failed: %v", err)
	}

	// act
	actual, err := dbw.GetPerformersByPerformanceId(testPerformance.Id)
	require.NoError(t, err, "GetPerformersByPerformanceId failed: %v", err)

	// assert
	assert.True(t, (len(expected) == len(actual)), "Number of returns not equal")
	for i, val := range actual {
		assert.Equal(t, val.Id, expected[i].Id, "Expected performers not equal to actual")
		assert.Equal(t, val.Name, expected[i].Name, "Expected performers not equal to actual")
		assert.Equal(t, val.Email, expected[i].Email, "Expected performers not equal to actual")
	}
}

func TestDeletePerformanceUsingId(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	performance := getTestPerformance()

	performance, err := dbw.CreatePerformance(performance)
	require.NoError(t, err, "CreatePerformance() has failed: %v", err)

	// act
	err = dbw.DeletePerformanceById(performance.Id)
	require.NoError(t, err, "DeletePerformanceById() has failed: %v", err)

	// assert
	retrievedPerformance, err := dbw.GetPerformanceById(performance.Id)
	require.NoError(t, err, "GetPerformanceById() has failed: %v", err)
	assert.Nil(t, retrievedPerformance, "Performance deletion failed!")
}

func TestDeletePerformerUsingId(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	performer := getTestPerformer()

	performer, err := dbw.CreatePerformer(performer)
	require.NoError(t, err, "CreatePerformer() has failed: %v", err)

	// act
	err = dbw.DeletePerformerById(performer.Id)
	require.NoError(t, err, "DeletePerformerById() has failed: %v", err)

	// assert
	retrievedPerformer, err := dbw.GetPerformerById(performer.Id)
	require.NoError(t, err, "GetPerformanceById() has failed: %v", err)
	assert.Nil(t, retrievedPerformer, "Performer deletion failed")
}

func TestDeleteJunction(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	performer := getTestPerformer()
	performances := getTestPerformances(2)

	performer, err := dbw.CreatePerformer(performer)
	require.NoError(t, err, "CreatePerformance() has failed: %v", err)
	for i, v := range performances {
		performances[i], err = dbw.CreatePerformance(v)
		require.NoError(t, err, "CreatePerformance() has failed: %v", err)
	}

	for _, v := range performances {
		err = dbw.CreateJunction(performer.Id, v.Id)
		require.NoError(t, err, "CreateJunction() has failed: %v", err)
	}

	// act
	err = dbw.DeleteJunction(performer.Id, performances[0].Id)
	require.NoError(t, err, "DeleteJunction() has failed: %v", err)

	actual, err := dbw.GetPerformancesByPerformerId(performer.Id)
	require.NoError(t, err, "GetPerformanceByPerformerId() has failed: %v", err)

	// assert
	expected := []*internal.Performance{performances[1]}
	require.Equal(t, expected, actual, "Retrieved values not equal to expected")
}

func TestUpdatePerformanceById(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	// create 2 performances
	testPerformances := getTestPerformances(2)
	expected := testPerformances[1]
	insertedPerformance := testPerformances[0]

	insertedPerformance, err := dbw.CreatePerformance(insertedPerformance)
	require.NoError(t, err, "CreatePerformance() has failed: %v", err)

	// act
	err = dbw.UpdatePerformanceById(insertedPerformance.Id, expected)
	require.NoError(t, err, "UpdatePerformanceById() has failed: %v", err)
	actual, err := dbw.GetPerformanceById(insertedPerformance.Id)
	require.NoError(t, err, "GetPerformanceById() has failed; %v", err)

	// assert
	assert.NotEqual(t, expected.Id, actual.Id, "Identical Ids")

	assert.Equal(t, expected.ItemName, expected.ItemName, "Expected performances not equal to actual")
	assert.Equal(t, expected.GroupName, expected.GroupName, "Expected performances not equal to actual")
	assert.Equal(t, expected.GenreName, expected.GenreName, "Expected performances not equal to actual")
	assert.Equal(t, expected.Location, expected.Location, "Expected performances not equal to actual")
	assert.Equal(t, expected.StartTime, expected.StartTime, "Expected performances not equal to actual")
	assert.Equal(t, expected.EndTime, expected.EndTime, "Expected performances not equal to actual")
	assert.Equal(t, expected.Duration, expected.Duration, "Expected performances not equal to actual")
}

func TestUpdatePerformerById(t *testing.T) {
	// arrange
	db := setUpTestDB(t)
	defer db.Close()
	dbw := internal.CreateDBWrapper(db)

	// create 2 performances
	testPerformers := getTestPerformers(2)
	expected := testPerformers[1]
	insertedPerformer := testPerformers[0]

	insertedPerformer, err := dbw.CreatePerformer(insertedPerformer)
	require.NoError(t, err, "CreatePerformer() has failed: %v", err)

	// act
	err = dbw.UpdatePerformerById(insertedPerformer.Id, expected)
	require.NoError(t, err, "UpdatePerformerById() has failed: %v", err)
	actual, err := dbw.GetPerformerById(insertedPerformer.Id)
	require.NoError(t, err, "GetPerformerById() has failed; %v", err)

	// assert
	assert.NotEqual(t, expected.Id, actual.Id, "Identical Ids")

	assert.Equal(t, expected.Name, actual.Name, "Expected and Actual name are different")
	assert.Equal(t, expected.Email, actual.Email, "Expected and Actual email are different")

}
