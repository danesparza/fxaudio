package media

import (
	"context"
	"fmt"
	"os/exec"
)

type PlayAudioRequest struct {
	ProcessID string `json:"pid"`      // Unique Internal Process ID
	ID        string `json:"id"`       // Unique File ID
	FilePath  string `json:"filepath"` // Full filepath to the file
}

// HandleAndProcess handles system context calls and channel events to play/stop audio
func HandleAndProcess(systemctx context.Context, playaudio chan PlayAudioRequest, stopaudio chan string) {

	//	Create a map of running instances and their cancel functions
	playingAudio := make(map[string]func())

	//	Loop and respond to channels:
	for {
		select {
		case playReq := <-playaudio:
			//	As we get a request on a channel to play a file...
			//	Spawn a goroutine
			go func(cx context.Context, req PlayAudioRequest) {
				//	Create a cancelable context from the passed (system) context
				ctx, cancel := context.WithCancel(cx)
				defer cancel()

				//	Add an entry to the map with
				//	- key: instance id
				//	- value: the cancel function (pointer)
				playingAudio[req.ProcessID] = cancel

				//	Create the command with context and play the audio
				playCommand := exec.CommandContext(ctx, "mpg123", playReq.FilePath)

				if err := playCommand.Run(); err != nil {
					//	Log an error playing a file
					fmt.Printf("error playing %v: %v", playReq.FilePath, err)
				}

				//	Remove ourselves from the map and exit (critical section)
				delete(playingAudio, req.ProcessID)

			}(systemctx, playReq) // Launch the goroutine

		case stopFile := <-stopaudio:
			//	Look up the item in the map and call cancel if the item exists:
			if playCancel, exists := playingAudio[stopFile]; exists {
				//	Call the context cancellation function (critical section)
				playCancel()

				//	Remove ourselves from the map and exit (critical section)
				delete(playingAudio, stopFile)
			}
		case <-systemctx.Done():
			fmt.Println("Stopping audio processor")
			return
		}
	}
}
