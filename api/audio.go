package api

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/danesparza/fxaudio/internal/media"
	"github.com/go-chi/chi/v5"
	"github.com/rs/xid"
	"github.com/spf13/viper"
	"io"
	math_rand "math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

// UploadFile godoc
// @Summary Upload a file
// @Description Upload file
// @Tags audio
// @Accept  mpfd
// @Produce  json
// @Param file formData file true "The file to upload"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 413 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /audio [put]
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
	service.DB.AddFile(req.Context(), destinationFile, "Audio file")

	//	If we've gotten this far, indicate a successful upload
	response := SystemResponse{
		Message: "File uploaded",
	}

	//	Serialize to JSON & return the response:
	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// ListAllFiles godoc
// @Summary List all files in the system
// @Description List all files in the system
// @Tags audio
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /audio [get]
func (service Service) ListAllFiles(rw http.ResponseWriter, req *http.Request) {

	//	Get a list of files
	retval, err := service.DB.GetAllFiles(req.Context())
	if err != nil {
		err = fmt.Errorf("error getting a list of files: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("%v file(s)", len(retval)),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}

// DeleteFile godoc
// @Summary Deletes a file in the system
// @Description Deletes a file in the system
// @Tags audio
// @Accept  json
// @Produce  json
// @Param id path string true "The file id to delete"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Failure 503 {object} api.ErrorResponse
// @Router /audio/{id} [delete]
func (service Service) DeleteFile(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	audioId := chi.URLParam(req, "id")

	if audioId == "" {
		err := fmt.Errorf("requires an id of a file to delete")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Delete the file
	err := service.DB.DeleteFile(req.Context(), audioId)
	if err != nil {
		err = fmt.Errorf("error deleting file: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: "File deleted",
		Data:    audioId,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}

// UpdateTags godoc
// @Summary Updates tags for a file
// @Description Updates tags for a file
// @Tags audio
// @Accept  json
// @Produce  json
// @Param id path string true "The file id to update tags for"
// @Param endpoint body api.UpdateTagsRequest true "The endpoint url to stream"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Failure 503 {object} api.ErrorResponse
// @Router /audio/{id} [post]
func (service Service) UpdateTags(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	audioId := chi.URLParam(req, "id")

	if audioId == "" {
		err := fmt.Errorf("requires an id of a file to update tags for")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Parse the body to get the tags
	request := UpdateTagsRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		err = fmt.Errorf("problem decoding tag update request: %v", err)
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Update the file
	err = service.DB.UpdateTags(req.Context(), audioId, request.Tags)
	if err != nil {
		err = fmt.Errorf("error updating tags: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: "File updated",
		Data:    audioId,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

}

// PlayAudio godoc
// @Summary Play an audio file
// @Description Play an audio file
// @Tags audio
// @Accept  json
// @Produce  json
// @Param id path string true "The file id to play"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Failure 503 {object} api.ErrorResponse
// @Router /audio/play/{id} [post]
func (service Service) PlayAudio(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	audioId := chi.URLParam(req, "id")
	if audioId == "" {
		err := fmt.Errorf("requires an id of a file to play")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the file information
	fileToPlay, err := service.DB.GetFile(req.Context(), audioId)
	if err != nil {
		err = fmt.Errorf("error getting file information: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

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
	_, err = exec.LookPath("mpg123")
	if err != nil {
		err = fmt.Errorf("didn't find mpg123 executable in the path: %v", err)
		sendErrorResponse(rw, err, http.StatusServiceUnavailable)
		return
	}

	//	Send to the channel:
	playRequest := media.PlayAudioRequest{
		ProcessID: xid.New().String(), // Generate a new id
		ID:        fileToPlay.ID,
		FilePath:  fileToPlay.FilePath,
		LoopTimes: "1",
	}
	service.PlayMedia <- playRequest

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Audio playing",
		Data:    playRequest,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// LoopAudio godoc
// @Summary Loop an audio file n times
// @Description Play an audio file over and over.  (-1 to loop until stopped)
// @Tags audio
// @Accept  json
// @Produce  json
// @Param id path string true "The file id to play"
// @Param loopTimes path string true "The number of times to play. The value -1 will cause it to loop until stopped"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Failure 503 {object} api.ErrorResponse
// @Router /audio/loop/{id}/{loopTimes} [post]
func (service Service) LoopAudio(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	audioId := chi.URLParam(req, "id")
	if audioId == "" {
		err := fmt.Errorf("requires an id of a file to loop")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	loopTimes := chi.URLParam(req, "loopTimes")
	if audioId == "" {
		err := fmt.Errorf("requires a number of times to loop")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the file information
	fileToPlay, err := service.DB.GetFile(req.Context(), audioId)
	if err != nil {
		err = fmt.Errorf("error getting file information: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Make sure mpg123 is installed (for i2s / ALSA based digital audio)
	//	Instructions on how to install it:
	//	https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/raspberry-pi-test
	_, err = exec.LookPath("mpg123")
	if err != nil {
		err = fmt.Errorf("didn't find mpg123 executable in the path: %v", err)
		sendErrorResponse(rw, err, http.StatusServiceUnavailable)
		return
	}

	//	Send to the channel:
	playRequest := media.PlayAudioRequest{
		ProcessID: xid.New().String(), // Generate a new id
		ID:        fileToPlay.ID,
		FilePath:  fileToPlay.FilePath,
		LoopTimes: loopTimes,
	}
	service.PlayMedia <- playRequest

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Audio playing",
		Data:    playRequest,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// StreamAudio godoc
// @Summary Stream from an audio endpoint
// @Description Stream from an audio endpoint
// @Tags audio
// @Accept  json
// @Produce  json
// @Param endpoint body api.PlayAudioRequest true "The endpoint url to stream"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 503 {object} api.ErrorResponse
// @Router /audio/stream [post]
func (service Service) StreamAudio(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request
	request := PlayAudioRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get just the file endpoint:
	fileendpoint := strings.TrimSpace(request.Endpoint)

	//	Make sure mpg123 is installed (for i2s / ALSA based digital audio)
	_, err = exec.LookPath("mpg123")
	if err != nil {
		err = fmt.Errorf("didn't find mpg123 executable in the path: %v", err)
		sendErrorResponse(rw, err, http.StatusServiceUnavailable)
		return
	}

	//	Send to the channel:
	playRequest := media.PlayAudioRequest{
		ProcessID: xid.New().String(), // Generate a new id
		FilePath:  fileendpoint,
		LoopTimes: "1",
	}
	service.PlayMedia <- playRequest

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Audio playing",
		Data:    playRequest,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// PlayRandomAudio godoc
// @Summary Play a random file already uploaded
// @Description Play a random file already uploaded
// @Tags audio
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /audio/play [post]
func (service Service) PlayRandomAudio(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Get just the file endpoint:
	fileendpoint := ""

	//	If we don't have an endpoint specified, get a random file that we manage
	if fileendpoint == "" {
		retval, err := service.DB.GetAllFiles(req.Context())
		if err != nil {
			err = fmt.Errorf("error getting a list of files: %v", err)
			sendErrorResponse(rw, err, http.StatusInternalServerError)
			return
		}

		//	If we don't have anything to pick from, bomb out
		if len(retval) < 1 {
			err = fmt.Errorf("can't play anything -- no endpoint specified, and no files to randomly pick from: %v", err)
			sendErrorResponse(rw, err, http.StatusBadRequest)
			return
		}

		//	Pick a random file:
		randomIndex := math_rand.Intn(len(retval))

		//	Set fileendpoint to the random file's path:
		fileendpoint = retval[randomIndex].FilePath
	}

	//	Make sure mpg123 is installed (for i2s / ALSA based digital audio)
	_, err := exec.LookPath("mpg123")
	if err != nil {
		err = fmt.Errorf("didn't find mpg123 executable in the path: %v", err)
		sendErrorResponse(rw, err, http.StatusServiceUnavailable)
		return
	}

	//	Send to the channel:
	playRequest := media.PlayAudioRequest{
		ProcessID: xid.New().String(), // Generate a new id
		FilePath:  fileendpoint,
		LoopTimes: "1",
	}
	service.PlayMedia <- playRequest

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Audio playing",
		Data:    playRequest,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// PlayRandomAudioWithTag godoc
// @Summary Play a random file already uploaded (with tag)
// @Description Play a random file already uploaded (that has a given tag)
// @Tags audio
// @Accept  json
// @Produce  json
// @Param tag path string true "Play random audio associated with this tag"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /audio/play/random/{tag} [post]
func (service Service) PlayRandomAudioWithTag(rw http.ResponseWriter, req *http.Request) {

	//	Get the tag from the url (if it's blank, return an error)
	tag := chi.URLParam(req, "tag")
	if tag == "" {
		err := fmt.Errorf("requires a tag to find")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Get just the file endpoint:
	fileendpoint := ""

	//	If we don't have an endpoint specified, get a random file that we manage
	if fileendpoint == "" {
		retval, err := service.DB.GetAllFilesWithTag(req.Context(), tag)
		if err != nil {
			err = fmt.Errorf("error getting a list of files: %v", err)
			sendErrorResponse(rw, err, http.StatusInternalServerError)
			return
		}

		//	If we don't have anything to pick from, bomb out
		if len(retval) < 1 {
			err = fmt.Errorf("can't play anything -- no endpoint specified, and no files to randomly pick from: %v", err)
			sendErrorResponse(rw, err, http.StatusBadRequest)
			return
		}

		//	Pick a random file:
		randomIndex := math_rand.Intn(len(retval))

		//	Set fileendpoint to the random file's path:
		fileendpoint = retval[randomIndex].FilePath
	}

	//	Make sure mpg123 is installed (for i2s / ALSA based digital audio)
	_, err := exec.LookPath("mpg123")
	if err != nil {
		err = fmt.Errorf("didn't find mpg123 executable in the path: %v", err)
		sendErrorResponse(rw, err, http.StatusServiceUnavailable)
		return
	}

	//	Send to the channel:
	playRequest := media.PlayAudioRequest{
		ProcessID: xid.New().String(), // Generate a new id
		FilePath:  fileendpoint,
		LoopTimes: "1",
	}
	service.PlayMedia <- playRequest

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Audio playing",
		Data:    playRequest,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// StopAudio godoc
// @Summary Stops a specific audio file 'play' process
// @Description Stops a specific audio file 'play' process
// @Tags audio
// @Accept  json
// @Produce  json
// @Param pid path string true "The process id to stop"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /audio/stop/{pid} [post]
func (service Service) StopAudio(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	processId := chi.URLParam(req, "pid")
	if processId == "" {
		err := fmt.Errorf("requires a processid of a process to stop")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Send to the channel:
	service.StopMedia <- processId

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Audio stopping",
		Data:    processId,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// StopAllAudio godoc
// @Summary Stops all audio 'play' processes
// @Description Stops all audio 'play' processes
// @Tags audio
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Router /audio/stop [post]
func (service Service) StopAllAudio(rw http.ResponseWriter, req *http.Request) {

	//	Send to the channel:
	service.StopAllMedia <- true

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "All Audio stopping",
		Data:    ".",
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

func init() {
	var b [8]byte
	crypto_rand.Read(b[:])
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
