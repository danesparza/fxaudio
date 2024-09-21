package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
)

type FFAudioService interface {
	PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error
}

type ffaudioService struct{}

func (f ffaudioService) PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error {
	//	Build our argument list
	args := []string{"-c", "/usr/bin/ffplay", "-autoexit", "-nodisp", "-loglevel", "error"}

	//	If we need to loop, indicate that we should
	if loop {
		args = append(args, "-loop", "9999")
	}

	//	At the end, add the file to play or url to stream
	args = append(args, audioPathOrUrl)

	//	Finally, run the full command:
	cmd := exec.CommandContext(ctx, "/bin/bash", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	log.Info().Strs("args", args).Msg("Playing ffplay audio")

	err := cmd.Run()
	if err != nil {
		log.Err(err).Str("stderr", stderr.String()).Strs("args", args).Msg("Problem playing audio")
		return fmt.Errorf("problem playing audio: %w", err)
	}

	if out.String() != "" || stderr.String() != "" {
		log.Info().Str("stdout", out.String()).Str("stderr", stderr.String()).Msg("Output from PlayAudio")
	}

	log.Info().Str("audioPathOrUrl", audioPathOrUrl).Msg("Played audio")

	return nil
}

func NewFFAudioService() FFAudioService {
	svc := &ffaudioService{}
	return svc
}
