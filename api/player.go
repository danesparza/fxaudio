package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

// PlayAudio plays an audio file
func (service Service) PlayAudio(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Make sure omxplayer is installed
	//	Instructions on how to install it:
	//	https://www.gaggl.com/2013/01/installing-omxplayer-on-raspberry-pi/
	_, err := exec.LookPath("omxplayer")
	if err != nil {
		err = fmt.Errorf("Didn't find omxplayer executable in the path: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	/*
		//	Get the file and play it:
		events, err := service.DB.GetAllEvents()
		if err != nil {
			sendErrorResponse(rw, err, http.StatusInternalServerError)
			return
		}
	*/

	//	Create our response and send information back:
	response := SystemResponse{
		Status:  http.StatusOK,
		Message: "Events fetched",
		Data:    "The data",
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
