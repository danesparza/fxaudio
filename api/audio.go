package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/viper"
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

func (service Service) UploadFile(rw http.ResponseWriter, req *http.Request) {

	MAX_UPLOAD_SIZE := viper.GetInt64("upload.bytelimit")
	UploadPath := viper.GetString("upload.path")

	//	First check for maximum uplooad size and return an error if we exceed it.
	req.Body = http.MaxBytesReader(rw, req.Body, MAX_UPLOAD_SIZE)
	if err := req.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		err = fmt.Errorf("could not parse multipart form: %v", err)
		sendErrorResponse(rw, err, http.StatusRequestEntityTooLarge)
		return
	}

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		err = fmt.Errorf("error retrieving file: %v", err)
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}
	defer file.Close()

	/*
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)
	*/

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll(UploadPath, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("error creating uploads path: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	destinationFile := path.Join(UploadPath, fileHeader.Filename)
	dst, err := os.Create(destinationFile)
	if err != nil {
		err = fmt.Errorf("error creating file: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		err = fmt.Errorf("error saving file: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Add it to our system database

}

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
