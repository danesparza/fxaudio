# fxaudio [![Build and release](https://github.com/danesparza/fxaudio/actions/workflows/release.yaml/badge.svg)](https://github.com/danesparza/fxaudio/actions/workflows/release.yaml) 
REST service for multichannel audio on demand from Raspberry Pi.  Made with ❤️ for makers, DIY craftsmen, and professional soundstage designers everywhere

Unfortunately, audio in linux is just kind of ... stupid.  Over the course of many months, I have discovered that this service works best running as a non-root user -- so the installation process now prompts you for a user to run the service as (and it defaults to the current user, which is usually fine).

Pro tips:
- Make sure the user that the service runs under already exists and has access to the `audio` group.  (The `pi` user -- or whatever the first user is that you setup -- will be in the `audio` group by default).
- Be sure to set your default audio device using `raspi-config`.  
- Also:  set your volume using `alsamixer`.  
- Don't try to run this as root.  (Unless you want to be miserable).  Because audio in linux is just stupid.  Sorry about that.

### Audio output
#### Option 1: Using the Raspberry Pi DigiAMP+
I would recommend using the [Raspberry Pi DigiAMP+)[https://www.raspberrypi.com/products/digiamp-plus/] since it's now an official Raspbery Pi foundation product. Be sure to follow the [Raspberry Pi DigiAMP+ installation instructions](https://www.raspberrypi.com/documentation/accessories/audio.html#raspberry-pi-digiamp).  

Hint:  I used the following settings in `/boot/firmware/config.txt` at the very bottom
```
# Some magic to prevent the normal HAT overlay from being loaded
dtoverlay=
# And then choose one of the following, according to the model:
dtoverlay=iqaudio-digiampplus,unmute_amp
```

I have used a pair of Polk Audio M2 indoor/outdoor speakers with the DigiAMP+ and it's VERY loud. Use `alsamixer` to adjust volume of the speakers.

You can do a test with streaming music using `ffplay -autoexit -nodisp http://ice1.somafm.com/u80s-128-mp3`

#### Option 2: Using the speaker bonnet
I would recommend using the [Adafruit Speaker Bonnet for Raspberry Pi](https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/overview) as well -- just get a [Pi with headers](https://www.adafruit.com/product/3708) and slide the bonnet right down on top of it (and be sure to follow the [Raspberry Pi OS configuration instructions](https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/raspberry-pi-usage) for the board).  I used a pair of [8 ohm 3" speakers](https://www.adafruit.com/product/1313) with the bonnet and it's reasonably loud. Use `alsamixer` to adjust volume of the speakers.

You can do a test with streaming music using `ffplay -autoexit -nodisp http://ice1.somafm.com/u80s-128-mp3`

#### Option 3: Using HDMI
Make sure that you have a stock Raspberry Pi installation and [make sure you have configured HDMI in your /boot/config.txt](https://raspberrypi.stackexchange.com/questions/32717/how-to-enable-sound-on-hdmi).  Also: [Buster might be a good legacy distro](https://www.reddit.com/r/raspberry_pi/comments/qujijj/no_hdmi_audio_in_raspiconfig_raspberry_os_lite/) to start with.  This option might be a little unreliable, depending on the HDMI cable and system you're feeding into.

You can do a streaming music test with the first HDMI port on the device (HDMI 0): `ffplay -autoexit -nodisp http://ice1.somafm.com/u80s-128-mp3`

#### Any option: Testing audio output
Run `speaker-test -c2` to generate white noise out of the speaker, alternating left and right.

## Installation
### Prerequisites
Install Raspberry Pi OS (Bookworm or later) on your device. For best results, use the [Raspberry Pi Imager](https://www.raspberrypi.com/software/) and select **Raspberry Pi OS (64-bit)**.

Install the prerequisite package repository (one-time setup per machine):

``` bash
wget https://packages.cagedtornado.com/prereq.sh -O - | sh
```

### Installing fxaudio
Install the fxaudio package:

``` bash
sudo apt install fxaudio
```

During installation, you will be prompted to specify a **non-root** user to run the fxaudio service.

- If you press Enter without typing anything, it will default to the current user (recommended in most cases).
- The specified user **must already exist** on the system and **must be a member of the `audio` group**.
(On a standard Raspberry Pi setup, the first user — such as `pi` — is already in the `audio` group.)

**Important:**
- Do not run `fxaudio` as the `root` user. Audio on Linux can behave inconsistently under root, leading to unreliable results.
- Ensure your audio device is set correctly using `raspi-config`, and adjust output levels with `alsamixer`.

After installation, the service will be running and available at:
- REST API: `http://localhost:3030`
- API Documentation: `http://localhost:3030/swagger/`

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
