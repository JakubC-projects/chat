## Build server
FROM golang:1.23-bullseye AS build-server
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY chat ./chat

RUN go build ./chat/cmd/server/main.go

ENTRYPOINT ["go", "run", "./chat/cmd/server/main.go"]
