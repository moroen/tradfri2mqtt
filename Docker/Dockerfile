FROM alpine:latest AS builder

# Install packages
RUN apk --no-cache add git go npm

# Expose MQTT port
EXPOSE 1883

ENV PATH /usr/sbin:$PATH

RUN mkdir /go

WORKDIR /go

RUN git clone -b development https://github.com/moroen/tradfri2mqtt.git
RUN git clone -b development https://github.com/moroen/gocoap.git
RUN git clone -b development https://github.com/moroen/go-tradfricoap.git

WORKDIR /go/tradfri2mqtt
RUN go get
RUN go build -v

RUN mkdir /opt/tradfri2mqtt
RUN cp tradfri2mqtt /opt/tradfri

WORKDIR /go/tradfri2mqtt/vue
RUN npm install -g yarn
RUN yarn
RUN yarn build-www

FROM alpine:latest

# RUN apk --no-cache add mosquitto mosquitto-clients

RUN mkdir /config

COPY --from=builder /go/tradfri2mqtt/tradfri2mqtt /opt/tradfri2mqtt/tradfri2mqtt
COPY --from=builder /go/tradfri2mqtt/www/ /opt/tradfri2mqtt/www/
# COPY mosquitto.conf /etc/mosquitto/mosquitto.conf
COPY run.sh /opt/tradfri2mqtt/run.sh

# Expose MQTT port
# EXPOSE 1883

ENV PATH /usr/sbin:$PATH

CMD /bin/sh /opt/tradfri2mqtt/run.sh