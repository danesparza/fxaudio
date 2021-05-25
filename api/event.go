package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// GetEvent godoc
// @Summary Gets a log event.
// @Description Gets a log event.
// @Tags events
// @Accept  json
// @Produce  json
// @Param id path string true "The event id to get"
// @Success 200 {object} api.SystemResponse
// @Failure 404 {object} api.ErrorResponse
// @Router /event/{id} [get]
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

// GetAllEvents godoc
// @Summary Gets all events in the system
// @Description Gets all events in the system
// @Tags events
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /events [get]
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
		Message: fmt.Sprintf("%v events", len(events)),
		Data:    events,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
