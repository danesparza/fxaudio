package media

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
	"strings"
)

//	Additional references:
//	https://www.raspberrypi.com/documentation/computers/os.html#play-audio-and-video-on-raspberry-pi-os-lite
//	https://wiki.videolan.org/VLC_command-line_help/

type VLCAudioService interface {
	PlayAudio(ctx context.Context, loop bool, audioPathOrUrl string) error
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
	args := []string{"--play-and-exit", "-A", "alsa"}

	//	If we need to loop, indicate that we should
	if loop {
		args = append(args, "--loop")
	}

	//	By default, this will use the default alsa device.
	//	If we have a specific device configured, indicate we should use it
	if strings.TrimSpace(v.alsaDevice) != "" {
		args = append(args, "--alsa-audio-device", v.alsaDevice)
	}

	//	At the end, add the file to play or url to stream
	args = append(args, audioPathOrUrl)

	//	Finally, run the full command:
	_, err := exec.CommandContext(ctx, "cvlc", args...).Output()
	if err != nil {
		log.Err(err).Strs("args", args).Msg("Problem playing audio")
		return fmt.Errorf("problem playing audio: %w", err)
	}

	return nil
}

func NewVLCAudioService(alsaDevice string) VLCAudioService {
	svc := &vlcAudioService{
		alsaDevice: alsaDevice,
	}
	return svc
}
