package media

import (
	"context"
)

//	Additional references:
//	https://www.raspberrypi.com/documentation/computers/os.html#play-audio-and-video-on-raspberry-pi-os-lite
//	https://wiki.videolan.org/VLC_command-line_help/

type VLCAudioService interface {
	PlayAudio(ctx context.Context, loop bool, audioFilePath string) error
	StreamAudio(ctx context.Context, audioStreamUrl string) error
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

func (v vlcAudioService) PlayAudio(ctx context.Context, loop bool, audioFilePath string) error {
	//	cvlc --play-and-exit -A alsa --alsa-audio-device sysdefault:CARD=sndrpihifiberry /var/lib/fxaudio/uploads/map1.mp3
	//	to loop, use the --loop flag.  Example: cvlc --play-and-exit --loop -A alsa /var/lib/fxaudio/uploads/map1.mp3

	//TODO implement me
	panic("implement me")
}

func (v vlcAudioService) StreamAudio(ctx context.Context, audioStreamUrl string) error {
	//	cvlc --play-and-exit -A alsa --alsa-audio-device sysdefault:CARD=sndrpihifiberry http://ice1.somafm.com/u80s-128-mp3
	//	It's so simple!

	//TODO implement me
	panic("implement me")
}

func NewVLCAudioService(alsaDevice string) VLCAudioService {
	svc := &vlcAudioService{
		alsaDevice: alsaDevice,
	}
	return svc
}
