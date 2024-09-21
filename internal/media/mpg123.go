package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
)

type mpg123audioService struct{}

func (a mpg123audioService) CheckForPlayer() error {
	//	Make sure player is installed
	_, err := exec.LookPath("mpg123")
	if err != nil {
		err = fmt.Errorf("didn't find mpg123 executable in the path: %w", err)
		return err
	}

	return nil
}

func (a mpg123audioService) PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error {
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

	log.Info().Strs("args", args).Msg("Playing mpg123 audio")

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

func NewMPG123AudioService() AudioService {
	svc := &mpg123audioService{}
	return svc
}
