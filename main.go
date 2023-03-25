package main

import (
	"github.com/danesparza/fxaudio/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strings"
	"time"
)

// @title fxAudio
// @version 1.0
// @description fxAudio multichannel audio REST service

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1

func main() {
	//	Set log info:
	log.Logger = log.With().Timestamp().Caller().Logger()

	//	Set log level (default to info)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	switch strings.ToLower(os.Getenv("LOGGER_LEVEL")) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		break
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		break
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		break
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		break
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		break
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
		break
	}

	//	Set the error stack marshaller
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	//	Set log time format
	zerolog.TimeFieldFormat = time.RFC3339Nano

	cmd.Execute()
}
