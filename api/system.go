package api

import (
	"encoding/json"
	"fmt"
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
		Message: fmt.Sprintf("%v devices(s)", len(retval)),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
