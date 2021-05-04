package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetEvent gets a log event.
func (service Service) GetEvent(rw http.ResponseWriter, req *http.Request) {

	//	Parse the request
	vars := mux.Vars(req)

	//	Perform the action with the context user
	dataResponse, err := service.DB.GetEvent(vars["id"])
	if err != nil {
		sendErrorResponse(rw, err, http.StatusNotFound)
		return
	}

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Event fetched",
		Data:    dataResponse,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// GetAllEvents gets all events in the system.
func (service Service) GetAllEvents(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Get all the events:
	events, err := service.DB.GetAllEvents()
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Events fetched",
		Data:    events,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
