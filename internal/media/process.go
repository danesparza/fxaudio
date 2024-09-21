package media

import (
	"context"
	"github.com/danesparza/fxaudio/internal/data"
	"github.com/rs/zerolog/log"
	"sync"
)

type PlayAudioRequest struct {
	ProcessID string `json:"pid"`      // Unique Internal Process ID
	ID        string `json:"id"`       // Unique File ID
	FilePath  string `json:"filepath"` // Full filepath to the file
	Loop      bool   `json:"loop"`     // Number of times to loop.  Default: 1 (only play once)
}

type audioProcessMap struct {
	m       map[string]func()
	rwMutex sync.RWMutex
}

// BackgroundProcess encapsulates background processing operations
type BackgroundProcess struct {
	// DB is the system datastore reference
	DB data.AppDataService

	// AS is the audio service to play audio files and streams
	AS FFAudioService

	// PlayAudio signals audio should be played
	PlayAudio chan PlayAudioRequest

	// StopAudio signals a running audio process should be stopped
	StopAudio chan string

	// StopAllAudio signals all running audio should be stopped
	StopAllAudio chan bool

	// PlayingTimelines tracks currently playing timelines
	PlayingAudio audioProcessMap
}

// HandleAndProcess handles system context calls and channel events to play/stop audio
func (bp *BackgroundProcess) HandleAndProcess(systemctx context.Context) {

	//	Create a map of running instances and their cancel functions
	bp.PlayingAudio.m = make(map[string]func())
	log.Info().Msg("Starting audio processor...")

	//	Loop and respond to channels:
	for {
		select {
		case playReq := <-bp.PlayAudio:
			log.Info().Str("pid", playReq.ProcessID).Msg("request to play")
			//	As we get a request on a channel to play a file...
			//	Spawn a goroutine
			go func(cx context.Context, req PlayAudioRequest) {
				//	Create a cancelable context from the passed (system) context
				ctx, cancel := context.WithCancel(cx)
				defer cancel()

				//	Add an entry to the map with
				//	- key: instance id
				//	- value: the cancel function (pointer)
				//	(critical section)
				bp.PlayingAudio.rwMutex.Lock()
				bp.PlayingAudio.m[req.ProcessID] = cancel
				bp.PlayingAudio.rwMutex.Unlock()

				//	Play the audio from the request:
				if err := bp.AS.PlayAudio(ctx, playReq.Loop, playReq.FilePath); err != nil {
					//	Log an error playing a file
					log.Err(err).Str("filepath", playReq.FilePath).Msg("problem playing audio")
				}

				//	Remove ourselves from the map and exit (critical section)
				bp.PlayingAudio.rwMutex.Lock()
				delete(bp.PlayingAudio.m, req.ProcessID)
				bp.PlayingAudio.rwMutex.Unlock()

			}(systemctx, playReq) // Launch the goroutine

		case stopFile := <-bp.StopAudio:
			//	Look up the item in the map and call cancel if the item exists (critical section):
			bp.PlayingAudio.rwMutex.Lock()
			playCancel, exists := bp.PlayingAudio.m[stopFile]

			if exists {
				//	Call the context cancellation function
				playCancel()

				//	Remove ourselves from the map and exit
				delete(bp.PlayingAudio.m, stopFile)
			}
			bp.PlayingAudio.rwMutex.Unlock()

		case <-bp.StopAllAudio:
			//	Loop through all items in the map and call cancel if the item exists (critical section):
			bp.PlayingAudio.rwMutex.Lock()

			for stopFile, playCancel := range bp.PlayingAudio.m {

				//	Call the cancel function
				playCancel()

				//	Remove ourselves from the map
				//	(this is safe to do in a 'range':
				//	https://golang.org/doc/effective_go#for )
				delete(bp.PlayingAudio.m, stopFile)
			}

			bp.PlayingAudio.rwMutex.Unlock()

		case <-systemctx.Done():
			log.Info().Msg("Stopping audio processor...")
			return
		}
	}
}
