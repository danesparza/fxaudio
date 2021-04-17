package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

/*

/audio
PUT - Upload file
GET - List all files
POST - Play a random file (or if passed an endpoint in JSON, stream that file)

/audio/1
GET - Download file
POST - Play file
DELETE - Delete file

*/

// PlayAudio plays an audio file
func (service Service) PlayAudio(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	/*
		//	Make sure omxplayer is installed (for HDMI based audio)
		//	Instructions on how to install it:
		//	https://www.gaggl.com/2013/01/installing-omxplayer-on-raspberry-pi/
		_, err := exec.LookPath("omxplayer")
		if err != nil {
			err = fmt.Errorf("Didn't find omxplayer executable in the path: %v", err)
			sendErrorResponse(rw, err, http.StatusInternalServerError)
			return
		}
	*/
	//	Make sure mpg123 is installed (for i2s / ALSA based digital audio)
	//	Instructions on how to install it:
	//	https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/raspberry-pi-test
	_, err := exec.LookPath("mpg123")
	if err != nil {
		err = fmt.Errorf("didn't find mpg123 executable in the path: %v", err)
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
		Message: "Events fetched",
		Data:    "The data",
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
