package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type API struct {
	wrapper *DBWrapper
}

func NewAPI(wrapper *DBWrapper) *API {
	return &API{wrapper: wrapper}
}

// Private helper function to respond with JSON
func (api *API) respondJSON(writer http.ResponseWriter, status int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}

// Private helper function to respond with an error and a message
func (api *API) respondError(writer http.ResponseWriter, status int, message string) {
	api.respondJSON(writer, status, map[string]string{"error": message})
}

// extracts the id from a path of pattern "*/*/:id"
func (api *API) extractId(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid path")
	}

	return strconv.Atoi(parts[1])
}

func pathLength(path string) int {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return len(parts)
}

// Handles all requests related to performances
func (api *API) PerformanceHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println(strconv.Itoa(log.Ltime) + ": Performance Request Incoming " + r.Method + " " + r.URL.RawPath)
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/performances" || r.URL.Path == "/performances/" {
			// log.Println("Detected request for all performances!")
			api.GetAllPerformances(w, r)
		} else if pathLength(r.URL.Path) > 2 {
			api.GetPerformersByPerformanceId(w, r)
		} else {
			api.GetPerformanceById(w, r)
		}
	case http.MethodPost:
		api.CreateNewPerformance(w, r)
	case http.MethodPut:
		api.UpdatePerformance(w, r)
	case http.MethodDelete:
		api.DeletePerformance(w, r)
	}
}

// Handles all requests related to performers
func (api *API) PerformerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" Performer Request Incoming: " + r.Method + " " + r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/performers" || r.URL.Path == "/performers/" {
			api.GetAllPerformers(w, r)
		} else if pathLength(r.URL.Path) > 2 { // if pathLength >
			api.GetPerformancesByPerformerId(w, r)
		} else {
			api.GetPerformerById(w, r)
		}
	case http.MethodPost:
		api.CreateNewPerformer(w, r)
	case http.MethodPut:
		api.UpdatePerformer(w, r)
	case http.MethodDelete:
		api.DeletePerformer(w, r)
	}
}

// Handles (most) requests related to junctions
func (api *API) JunctionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		api.CreateJunction(w, r)
	case http.MethodDelete:
		api.DeleteJunction(w, r)
	}
}

// GET /performances - returns all performances
func (api *API) GetAllPerformances(w http.ResponseWriter, r *http.Request) {
	performances, err := api.wrapper.GetAllPerformances()
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid Performance ID")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string][]*Performance{"performances": performances})
}

// GET /performers - returns all performers
func (api *API) GetAllPerformers(w http.ResponseWriter, r *http.Request) {
	performers, err := api.wrapper.GetAllPerformers()
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid performer id")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string][]*Performer{"performers": performers})
}

// GET /performances/:id - return performance with given ID
func (api *API) GetPerformanceById(w http.ResponseWriter, r *http.Request) {
	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	performance, err := api.wrapper.GetPerformanceById(id)
	// performance != performance implies it is nil
	if err != nil || performance == nil {
		api.respondError(w, http.StatusNotFound, "Performance Not Found")
		return
	}

	api.respondJSON(w, http.StatusOK, performance)
}

// GET /performers/:id - return performer with given ID
func (api *API) GetPerformerById(w http.ResponseWriter, r *http.Request) {
	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// gets performer details from db
	performer, err := api.wrapper.GetPerformerById(id)
	if err != nil || performer == nil {
		api.respondError(w, http.StatusNotFound, "Performer Not Found")
		return
	}

	api.respondJSON(w, http.StatusOK, performer)
}

// GET /performances/:id/performers - returns performers associated to the performance with the specified id
func (api *API) GetPerformersByPerformanceId(w http.ResponseWriter, r *http.Request) {
	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Error extracting id")
		return
	}

	performers, err := api.wrapper.GetPerformersByPerformanceId(id)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Unable to find performers")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string][]*Performer{"performers": performers})
}

// GET /performances/:id/performers - returns performers associated to the performance with the specified id
func (api *API) GetPerformancesByPerformerId(w http.ResponseWriter, r *http.Request) {
	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Error extracting id")
		return
	}

	performances, err := api.wrapper.GetPerformancesByPerformerId(id)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Unable to find performances")
		return
	}

	api.respondJSON(w, http.StatusOK, performances)
}

// POST /performances/ - Create a new performance
func (api *API) CreateNewPerformance(w http.ResponseWriter, r *http.Request) {
	var performance Performance

	err := json.NewDecoder(r.Body).Decode(&performance)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// TODO: some more validation
	if performance.ItemName == "" {
		api.respondError(w, http.StatusBadRequest, "Cannot be blank")
		return
	}

	newPerformance, err := api.wrapper.CreatePerformance(&performance)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Failed to create performance")
	}

	api.respondJSON(w, http.StatusCreated, newPerformance)
}

// POST /performers/ - Create a new performer
func (api *API) CreateNewPerformer(w http.ResponseWriter, r *http.Request) {
	var performer Performer

	err := json.NewDecoder(r.Body).Decode(&performer)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// TODO: some more validation
	if performer.Name == "" {
		api.respondError(w, http.StatusBadRequest, "Cannot Be Blank")
		return
	}

	newPerformer, err := api.wrapper.CreatePerformer(&performer)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Failed to create performer")
	}

	api.respondJSON(w, http.StatusCreated, newPerformer)
}

// POST /junctions/ - creates a new performer:performance junction
func (api *API) CreateJunction(w http.ResponseWriter, r *http.Request) {
	junction := struct {
		PerformerId   int `json:"performerId`
		PerformanceId int `json:"performanceId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&junction)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = api.wrapper.CreateJunction(junction.PerformerId, junction.PerformanceId)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Failed to create junction")
		return
	}

	api.respondJSON(w, http.StatusCreated, map[string]string{"status": "success"})
}

// PUT /performances/:id - updates the performance with the specified id
func (api *API) UpdatePerformance(w http.ResponseWriter, r *http.Request) {
	var performance Performance

	err := json.NewDecoder(r.Body).Decode(&performance)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID provided")
	}

	// TODO: more validation
	if performance.ItemName == "" {
		api.respondError(w, http.StatusBadRequest, "Cannot be blank")
		return
	}

	err = api.wrapper.UpdatePerformanceById(id, &performance)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Error updating performance")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// PUT /performers/:id - updates the performer with the specified id
func (api *API) UpdatePerformer(w http.ResponseWriter, r *http.Request) {
	var performer Performer

	err := json.NewDecoder(r.Body).Decode(&performer)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID provided")
	}

	// TODO: more validation
	if performer.Name == "" {
		api.respondError(w, http.StatusBadRequest, "Cannot be blank")
		return
	}

	err = api.wrapper.UpdatePerformerById(id, &performer)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Error updating performer")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// DELETE /performances/:id - deletes the performance with the specified id
func (api *API) DeletePerformance(w http.ResponseWriter, r *http.Request) {
	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID provided")
		return
	}

	err = api.wrapper.DeletePerformanceById(id)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Error deleting performance")
	}

	api.respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// DELETE /performers/:id - deletes the performer with the specified id
func (api *API) DeletePerformer(w http.ResponseWriter, r *http.Request) {
	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID provided")
		return
	}

	err = api.wrapper.DeletePerformerById(id)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Error deleting performer")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// DELETE /junctions/:performerId/:performanceId
func (api *API) DeleteJunction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	performerId, err := strconv.Atoi(parts[1])
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID provided")
		return
	}
	performanceId, err := strconv.Atoi(parts[2])
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID provided")
		return
	}

	err = api.wrapper.DeleteJunction(performerId, performanceId)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Error deleting junction")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
