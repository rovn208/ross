# Build image stag
FROM golang:1.20-alpine as builder
WORKDIR /server
COPY . .
RUN go build -o ross cmd/main.go

# Run image
FROM alpine:3.16
WORKDIR /server
COPY --from=builder /server/ross .

EXPOSE 8080
CMD ["/server/ross"]
