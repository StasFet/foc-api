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

// Handles all requests related to performances
func (api *API) PerformanceHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println(strconv.Itoa(log.Ltime) + ": Performance Request Incoming " + r.Method + " " + r.URL.RawPath)
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/performances" || r.URL.Path == "/performances/" {
			// log.Println("Detected request for all performances!")
			api.getAllPerformances(w, r)
		} else {
			api.getPerformanceById(w, r)
		}
	case http.MethodPost:
		api.CreateNewPerformance(w, r)
	case http.MethodPut:
		api.UpdatePerformance(w, r)
	case http.MethodDelete:
		// delete performance of extracted id
	}
}

// Handles all requests related to performers
func (api *API) PerformerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(" Performer Request Incoming: " + r.Method + " " + r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/performers" || r.URL.Path == "/performers/" {
			api.getAllPerformers(w, r)
		} else {
			api.getPerformerById(w, r)
		}
	case http.MethodPost:
		api.CreateNewPerformer(w, r)
	case http.MethodPut:
		api.UpdatePerformer(w, r)
	case http.MethodDelete:
		// delete performer with extracted id
	}
}

// GET /api/performances - returns all performances
func (api *API) getAllPerformances(w http.ResponseWriter, r *http.Request) {
	performances, err := api.wrapper.GetAllPerformances()
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid Performance ID")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string][]*Performance{"performances": performances})
}

// GET /api/performers - returns all performers
func (api *API) getAllPerformers(w http.ResponseWriter, r *http.Request) {
	performers, err := api.wrapper.GetAllPerformers()
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid performer id")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string][]*Performer{"performers": performers})
}

// GET /api/performances/:id - return performance with given ID
func (api *API) getPerformanceById(w http.ResponseWriter, r *http.Request) {
	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	performance, err := api.wrapper.GetPerformanceById(id)
	if err != nil {
		api.respondError(w, http.StatusNotFound, "Performance Not Found")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string]*Performance{"performance": performance})
}

// GET /api/performers/:id - return performer with given ID
func (api *API) getPerformerById(w http.ResponseWriter, r *http.Request) {
	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	performer, err := api.wrapper.GetPerformerById(id)
	if err != nil {
		api.respondError(w, http.StatusNotFound, "Performer Not Found")
		return
	}

	api.respondJSON(w, http.StatusOK, map[string]*Performer{"performance": performer})
}

// POST /performances/ - Create a new performance
func (api *API) CreateNewPerformance(w http.ResponseWriter, r *http.Request) {
	var performance Performance
	blankPerformance := Performance{}

	err := json.NewDecoder(r.Body).Decode(&performance)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// TODO: some more validation
	if performance == blankPerformance {
		api.respondError(w, http.StatusBadRequest, "Cannot be blank")
		return
	}

	newPerformance, err := api.wrapper.CreatePerformance(&performance)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Failed to create performance")
	}

	api.respondJSON(w, http.StatusOK, newPerformance)
}

// POST /performers/ - Create a new performer
func (api *API) CreateNewPerformer(w http.ResponseWriter, r *http.Request) {
	var performer Performer
	blankPerformer := Performer{}

	err := json.NewDecoder(r.Body).Decode(&performer)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// TODO: some more validation
	if performer == blankPerformer {
		api.respondError(w, http.StatusBadRequest, "Cannot Be Blank")
		return
	}

	newPerformer, err := api.wrapper.CreatePerformer(&performer)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Failed to create performer")
	}

	api.respondJSON(w, http.StatusOK, newPerformer)
}

// PUT /performances/:id - updates the performance with the specified id
func (api *API) UpdatePerformance(w http.ResponseWriter, r *http.Request) {
	var performance Performance
	blankPerformance := Performance{}

	err := json.NewDecoder(r.Body).Decode(&performance)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Unable to parse id")
	}

	// TODO: more validation
	if performance == blankPerformance {
		api.respondError(w, http.StatusBadRequest, "Cannot be blank")
		return
	}

	updatedPerformance, err := api.wrapper.UpdatePerformanceById(id, &performance)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Error updating performance")
		return
	}

	api.respondJSON(w, http.StatusOK, updatedPerformance)
}

// PUT /performers/:id - updates the performer with the specified id
func (api *API) UpdatePerformer(w http.ResponseWriter, r *http.Request) {
	var performer Performer
	blankPerformer := Performer{}

	err := json.NewDecoder(r.Body).Decode(&performer)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	id, err := api.extractId(r.URL.Path)
	if err != nil {
		api.respondError(w, http.StatusBadRequest, "Unable to parse id")
	}

	// TODO: more validation
	if performer == blankPerformer {
		api.respondError(w, http.StatusBadRequest, "Cannot be blank")
		return
	}

	updatedPerformer, err := api.wrapper.UpdatePerformerById(id, &performer)
	if err != nil {
		api.respondError(w, http.StatusInternalServerError, "Error updating performer")
		return
	}

	api.respondJSON(w, http.StatusOK, updatedPerformer)
}
