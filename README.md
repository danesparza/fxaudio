# fxaudio [![Build and release](https://github.com/danesparza/fxaudio/actions/workflows/release.yaml/badge.svg)](https://github.com/danesparza/fxaudio/actions/workflows/release.yaml) 
REST service for multichannel audio on demand from Raspberry Pi.  Made with ❤️ for makers, DIY craftsmen, and professional soundstage designers everywhere

Unfortunately, audio in linux is just kind of ... stupid.  Over the course of many months, I have discovered that this service works best running as the default user on the Pi, and not as a service.  So:  just start it up with `fxaudio start` at the command line.  Also:  Be sure to set your default audio device using `raspi-config`.  Also:  set your volume using `alsamixer`.  Don't try to run this as root.  Or as a systemd service.  (Unless you want to be miserable).  Because audio in linux is just stupid.  Sorry about that.

### Audio output
#### Option 1: Using HDMI
Make sure that you have a stock Raspberry Pi installation and [make sure you have configured HDMI in your /boot/config.txt](https://raspberrypi.stackexchange.com/questions/32717/how-to-enable-sound-on-hdmi).  Also: [Buster might be a good legacy distro](https://www.reddit.com/r/raspberry_pi/comments/qujijj/no_hdmi_audio_in_raspiconfig_raspberry_os_lite/) to start with.  

You can do a streaming music test with the first HDMI port on the device (HDMI 0): `ffplay -autoexit -nodisp http://ice1.somafm.com/u80s-128-mp3`

#### Option 2: Using the speaker bonnet
I would recommend using the [Adafruit Speaker Bonnet for Raspberry Pi](https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/overview) as well -- just get a [Pi with headers](https://www.adafruit.com/product/3708) and slide the bonnet right down on top of it (and be sure to follow the [Raspberry Pi OS configuration instructions](https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/raspberry-pi-usage) for the board).  I used a pair of [8 ohm 3" speakers](https://www.adafruit.com/product/1313) with the bonnet and it's reasonably loud. Use `alsamixer` to adjust volume of the speakers.

You can do a test with streaming music using `ffplay -autoexit -nodisp http://ice1.somafm.com/u80s-128-mp3`

#### Either option: Testing audio output
Run `speaker-test -c2` to generate white noise out of the speaker, alternating left and right.

## Installing
Grab the latest release, copy it to your system, and make sure it's executable.  Then run it with `fxaudio start`

You can then use the service at http://localhost:3030

See the REST API documentation at http://localhost:3030/swagger/

## Example hardware setup
![fxAudio example hardware setup](fxAudio_hardware_annotated.png)

## The first second of audio is cut off!
This is most likely your HDMI audio system doing something it shouldn't.  You can work around this by enabling a service to play 'no audio' all the time in the background.  This will keep the device ready for when you want to play some actual audio through it.  

Put the following in `/etc/systemd/system/aplay.service`

```
[Unit]
Description=Invoke aplay from /dev/zero at system start.

[Service]
ExecStart=/usr/bin/aplay -D default -t raw -r 44100 -c 2 -f S16_LE /dev/zero

[Install]
WantedBy=multi-user.target
```

Then run the following commands to register the service, enable it, and start it:

```
sudo systemctl daemon-reload
sudo systemctl enable aplay
sudo systemctl start aplay
```
