# ross
[![Build](https://github.com/rovn208/ross/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/rovn208/ross/actions/workflows/test.yml/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/rovn208/ross)](https://goreportcard.com/report/github.com/rovn208/ross)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/rovn208/ross/blob/master/LICENSE)

Ross is a streaming service YouTube alike

## Prerequisite

- Golang 1.20
- ffmpeg: `brew install ffmpeg`
- Docker, Docker-compose
- [sqlc](https://github.com/sqlc-dev/sqlc)
- [golang-migrate](https://github.com/golang-migrate/migrate)

## Easy to start
```bash
# Setup and run database
make db
# Migrate data schema
make migrateup

# Running the application
make run # If you wanna run in the current terminal
docker compose up -d # If you wanna run in docker
```

- Swagger URL for available endpoints: http://localhost:8080/api/v1/swagger/index.html 
- Upload video demo: https://youtu.be/TPTDeHkC4d8
## Features
- [x] User management
- [x] Video management
  - Adding video via Youtube url
  - Upload video via form
  - Video's CRUD
- [x] Following/Follower functionality
- [ ] Video Feeds
- [x] Subscribing videos
- [ ] Room streaming

### Supported transport protocols
- [x] HLS
- [ ] RTMP
- [ ] FLV

### NOTES:
- This project is currently have APIs only. Therefore, the only way to test streaming source is using [this HLS stream tester](https://hlsjs-dev.video-dev.org/demo/)
- Input should be in the format: `{host}/api/v1/sources/{video_url}`.
- For example: `http://localhost:8080/api/v1/sources/x36xhzz/x36xhzz.m3u8`

### Known issues
Many, for sure cuz it's still IN DEVELOPMENT:">
