#!/bin/sh

/usr/sbin/mosquitto -c /etc/mosquitto/mosquitto.conf &
/opt/tradfri2mqtt/tradfri2mqtt
