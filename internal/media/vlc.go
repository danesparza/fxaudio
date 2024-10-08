package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
)

//	Additional references:
//	https://www.raspberrypi.com/documentation/computers/os.html#play-audio-and-video-on-raspberry-pi-os-lite
//	https://wiki.videolan.org/VLC_command-line_help/

type VLCAudioService interface {
	PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error
	CheckForPlayer() error
	SetAlsaDevice(alsaDeviceName string)
	GetAlsaDeviceName() string
}

type vlcAudioService struct {
	//	alsaDevice is the system device to use when playing audio.
	// 	To list devices, use the command `aplay -L | grep sysdefault`
	//	If this is blank, the system default device will be used
	//	and the --alsa-audio-device flag will not be used.
	//	Using the system default device may be necessary for playback with certain devices like this:
	//	https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/overview
	alsaDevice string
}

func (a *vlcAudioService) CheckForPlayer() error {
	//	Make sure player is installed
	_, err := exec.LookPath("cvlc")
	if err != nil {
		err = fmt.Errorf("didn't find cvlc executable in the path: %w", err)
		return err
	}

	return nil
}

func (v *vlcAudioService) SetAlsaDevice(alsaDeviceName string) {
	v.alsaDevice = alsaDeviceName
}

func (v *vlcAudioService) GetAlsaDeviceName() string {
	return v.alsaDevice
}

func (v *vlcAudioService) PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error {
	//	cvlc --play-and-exit -A alsa --alsa-audio-device sysdefault:CARD=sndrpihifiberry /var/lib/fxaudio/uploads/map1.mp3
	//	to loop, use the --loop flag.  Example: cvlc --play-and-exit --loop -A alsa /var/lib/fxaudio/uploads/map1.mp3

	//	Build our argument list
	args := []string{"--no-one-instance", "--play-and-exit"}

	//	If we need to loop, indicate that we should
	if loop {
		args = append(args, "--loop")
	}

	//	By default, this will use the default alsa device.
	//	If we have a specific device configured, indicate we should use it
	//if strings.TrimSpace(v.alsaDevice) != "" {
	//	args = append(args, "--alsa-audio-device", v.alsaDevice)
	//}

	//	At the end, add the file to play or url to stream
	args = append(args, audioPathOrUrl)

	//	Finally, run the full command:
	cmd := exec.CommandContext(ctx, "cvlc", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	log.Info().Strs("args", args).Msg("Playing cvlc audio")

	err := cmd.Run()
	if err != nil {
		log.Err(err).Str("stderr", stderr.String()).Strs("args", args).Msg("Problem playing audio")
		return fmt.Errorf("problem playing audio: %w", err)
	}

	return nil
}

func NewVLCAudioService() VLCAudioService {
	svc := &vlcAudioService{
		//alsaDevice: alsaDevice,
	}
	return svc
}
