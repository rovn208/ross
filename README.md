# ross
Ross is a streaming service YouTube alike

## Prerequisite

- Golang 1.20
- ffmpeg: `brew install ffmpeg`
- Docker, Docker-compose
- [sqlc](https://github.com/sqlc-dev/sqlc)
- [golang-migrate](https://github.com/golang-migrate/migrate)

## Easy to start

> docker compose up -d

- Swagger URL for available endpoints: http://localhost:8080/api/v1/swagger/index.html 
## Features
- [x] User management
- [x] Video management
  - Adding video via Youtube url
  - Upload video via form
  - Video's CRUD
- [x] Following/Follower functionality
- [ ] Video Feeds
- [ ] Subscribing videos
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