#!/bin/bash

docker container stop kiosk

docker container rm kiosk

docker buildx build -t kiosk .

docker run -p 9999:9999 --name kiosk --env-file .env kiosk