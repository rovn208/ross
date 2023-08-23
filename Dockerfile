# Build image stag
FROM golang:1.20-alpine as builder
WORKDIR /server
COPY . .
RUN go build -o ross cmd/serverd/main.go

# Run image
FROM alpine:3.16
WORKDIR /server
COPY --from=builder /server/ross .
COPY app.env .
COPY pkg/db/migration ./pkg/db/migration
RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

EXPOSE 8080
CMD ["/server/ross"]
