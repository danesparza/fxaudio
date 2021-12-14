# fxaudio [![CircleCI](https://circleci.com/gh/danesparza/fxaudio.svg?style=shield)](https://circleci.com/gh/danesparza/fxaudio)
REST service for multichannel audio on demand from Raspberry Pi.  Made with ❤️ for makers, DIY craftsmen, and professional soundstage designers everywhere

## Prerequisites
fxaudio uses [mpg123](https://en.wikipedia.org/wiki/Mpg123) under the hood to play audio, so you'll need to make sure it's installed first.  Lucky for you, installation is really simple:

```bash
sudo apt-get update
sudo apt-get install -y mpg123
```

I would recommend using the [Adafruit Speaker Bonnet for Raspberry Pi](https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/overview) as well -- just get a [Pi with headers](https://www.adafruit.com/product/3708) and slide the bonnet right down on top of it (and be sure to follow the [Raspberry Pi OS configuration instructions](https://learn.adafruit.com/adafruit-speaker-bonnet-for-raspberry-pi/raspberry-pi-usage) for the board).  I used a pair of [8 ohm 3" speakers](https://www.adafruit.com/product/1313) with the bonnet. 

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
