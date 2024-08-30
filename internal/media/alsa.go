package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
	"strings"
)

// GetAlsaDevices gets a list of all ALSA audio devices in the system
func GetAlsaDevices(ctx context.Context) ([]string, error) {
	retval := []string{}

	//	Get the list of devices:
	var aplayOut, aplayErr bytes.Buffer
	cmdAplay := exec.CommandContext(ctx, "aplay", "-L")
	cmdAplay.Stdout = &aplayOut
	cmdAplay.Stderr = &aplayErr
	err := cmdAplay.Run()
	if err != nil {
		log.Err(err).Msg("Alsa list devices failed")
		return retval, fmt.Errorf("problem getting list of alsa devices: %w", err)
	}

	//	Filter the results
	var grepOut, grepErr bytes.Buffer
	cmdGrep := exec.CommandContext(ctx, "grep", "sysdefault")
	cmdGrep.Stdin = strings.NewReader(aplayOut.String())
	cmdGrep.Stdout = &grepOut
	cmdGrep.Stderr = &grepErr
	err = cmdGrep.Run()
	if err != nil {
		log.Err(err).Msg("Grep alsa devices failed")
		return retval, fmt.Errorf("problem grepping list of alsa devices: %w", err)
	}

	//	Parse each line of the output
	retval = ParseCliOutput(cmdGrep.String())

	return retval, nil
}
