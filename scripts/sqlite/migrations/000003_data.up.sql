/* Default config - default to first hdmi port on Raspberry Pi 4 and above */
insert into system_config(alsa_device) values ('sysdefault:CARD=vc4hdmi0');