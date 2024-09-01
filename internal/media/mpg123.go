package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
)

type AudioService interface {
	PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error
}

type audioService struct{}

func (a audioService) PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error {
	//	Build our argument list
	args := []string{}

	//	If we need to loop, indicate that we should
	if loop {
		args = append(args, "--loop", "-1")
	}

	//	At the end, add the file to play or url to stream
	args = append(args, audioPathOrUrl)

	//	Finally, run the full command:
	cmd := exec.CommandContext(ctx, "mpg123", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Err(err).Str("stderr", stderr.String()).Strs("args", args).Msg("Problem playing audio")
		return fmt.Errorf("problem playing audio: %w", err)
	}

	return nil
}

func NewAudioService() AudioService {
	svc := &audioService{}
	return svc
}
