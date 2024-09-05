package media

import (
	"context"
	"fmt"
	"github.com/aymanbagabas/go-pty"
	"github.com/rs/zerolog/log"
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

	audioPty, err := pty.New()
	if err != nil {
		log.Err(err).Msg("could not create pty")
		return fmt.Errorf("could not create pty: %w", err)
	}

	defer audioPty.Close()
	c := audioPty.CommandContext(ctx, "mpg123", args...)
	if err := c.Run(); err != nil {
		log.Err(err).Msg("could not start mpg123 in pty")
		return fmt.Errorf("could not start mpg123: %w", err)
	}

	log.Info().Str("audioPathOrUrl", audioPathOrUrl).Msg("Played audio")

	return nil
}

func NewAudioService() AudioService {
	svc := &audioService{}
	return svc
}
