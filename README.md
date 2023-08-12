# ross
A Simple and efficient video streaming server

### Usage
- Run server on local machine
```go
go run cmd/serverd/main.go
```
- Your server is now live with file name under music folder `http://localhost:8080/{HLS-file-format}` e.g `http://localhost:8080/aloi.m3u8`
- Using media client to test at https://hls-js-latest.netlify.com/demo/

### Supported transport protocols
- [x] HLS
- [ ] RTMP
- [ ] FLV

### Supported container formats
- [ ] FLS
- [ ] TS

### Supported encoding formats
- [ ] H264
- [ ] AAC
- [ ] MP3 
