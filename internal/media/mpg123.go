package media

import (
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
	log.Info().Strs("args", args).Msg("Playing audio")
	err := exec.CommandContext(ctx, "mpg123", args...).Start()
	if err != nil {
		log.Err(err).Strs("args", args).Msg("Problem playing audio")
		return fmt.Errorf("problem playing audio: %w", err)
	}

	log.Info().Str("audioPathOrUrl", audioPathOrUrl).Msg("Played audio")

	return nil
}

func NewAudioService() AudioService {
	svc := &audioService{}
	return svc
}
