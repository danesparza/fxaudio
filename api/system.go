package api

import (
	"encoding/json"
	"fmt"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/danesparza/fxaudio/internal/media"
	"net/http"
)

// GetAlsaDevices godoc
// @Summary List all ALSA devices in the system
// @Description List all ALSA devices in the system
// @Tags system
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /system/alsadevices [get]
func (service Service) GetAlsaDevices(rw http.ResponseWriter, req *http.Request) {
	//	Get a list of devices
	retval, err := media.GetAlsaDevices(req.Context())
	if err != nil {
		err = fmt.Errorf("error getting a list of audio devices: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("%v device(s)", len(retval)),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// GetSystemConfig godoc
// @Summary Gets the system configuration info
// @Description Gets the system configuration info
// @Tags system
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /system/config [get]
func (service Service) GetSystemConfig(rw http.ResponseWriter, req *http.Request) {
	//	Get the system config
	retval, err := service.DB.GetConfig(req.Context())
	if err != nil {
		err = fmt.Errorf("error getting system config: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("config"),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// SetSystemConfig godoc
// @Summary Updates the system configuration info
// @Description Updates the system configuration info
// @Tags system
// @Accept  json
// @Produce  json
// @Param endpoint body data.SystemConfig true "The system config information"
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /system/config [put]
func (service Service) SetSystemConfig(rw http.ResponseWriter, req *http.Request) {
	//	Parse the body to get the config
	request := data.SystemConfig{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		err = fmt.Errorf("problem decoding config update request: %v", err)
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Set the system config
	err = service.DB.SetConfig(req.Context(), request)
	if err != nil {
		err = fmt.Errorf("error updating system config: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("config updated"),
		Data:    nil,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
