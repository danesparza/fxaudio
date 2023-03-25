# fxaudio [![Build and release](https://github.com/danesparza/fxaudio/actions/workflows/release.yaml/badge.svg)](https://github.com/danesparza/fxaudio/actions/workflows/release.yaml) 
REST service for multichannel audio on demand from Raspberry Pi.  Made with ❤️ for makers, DIY craftsmen, and professional soundstage designers everywhere

## Prerequisites

### mpg123 installation
fxaudio uses [mpg123](https://en.wikipedia.org/wiki/Mpg123) under the hood to play audio, so you'll need to make sure it's installed first.  Lucky for you, installation is really simple:

```bash
sudo apt-get update
sudo apt-get install -y mpg123
```

### Audio output
#### Option 1: Using HDMI
Make sure that you have a stock Raspberry Pi installation and [make sure you have configured HDMI in your /boot/config.txt](https://raspberrypi.stackexchange.com/questions/32717/how-to-enable-sound-on-hdmi).  Also: [Buster might be a good legacy distro](https://www.reddit.com/r/raspberry_pi/comments/qujijj/no_hdmi_audio_in_raspiconfig_raspberry_os_lite/) to start with.  

#### Option 2: Using the speaker bonnet
I would recommend using the [Adafruit Speaker Bonnet for Raspberry Pi](https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/overview) as well -- just get a [Pi with headers](https://www.adafruit.com/product/3708) and slide the bonnet right down on top of it (and be sure to follow the [Raspberry Pi OS configuration instructions](https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/raspberry-pi-usage) for the board).  I used a pair of [8 ohm 3" speakers](https://www.adafruit.com/product/1313) with the bonnet. 

#### Either option: Testing audio output
Run `speaker-test -c2` to generate white noise out of the speaker, alternating left and right.

You can do a test with streaming music using `mpg123 http://ice1.somafm.com/u80s-128-mp3`

## Installing
Installing fxaudio is also really simple.  Grab the .deb file from the [latest release](https://github.com/danesparza/fxaudio/releases/latest) and then install it using dpkg:


```bash
sudo dpkg -i fxaudio-1.0.40_armhf.deb 
````

This automatically installs the **fxaudio** service with a default configuration and starts the service. 

You can then use the service at http://localhost:3030

See the REST API documentation at http://localhost:3030/v1/swagger/

## Removing 
Uninstalling is just as simple:

```bash
sudo dpkg -r fxaudio
````

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
