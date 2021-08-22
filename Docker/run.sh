#!/bin/sh

export export MQTT_INTERFACE_ROOT="/opt/tradfri2mqtt/www"

/usr/sbin/mosquitto -c /etc/mosquitto/mosquitto.conf &
/opt/tradfri2mqtt/tradfri2mqtt
