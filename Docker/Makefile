local:
	docker buildx build --no-cache --load -t moroen/tradfri2mqtt:latest .

dist:
	docker buildx build --no-cache --push --platform=linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 -t moroen/tradfri2mqtt:latest .

v7:
	docker buildx build --load --platform=linux/arm/v7 -t moroen/tradfri2mqtt:v7 .

arm64:
	docker buildx build --load --platform=linux/arm64 -t moroen/tradfri2mqtt:arm64 .
