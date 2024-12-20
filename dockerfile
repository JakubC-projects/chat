# syntax=docker/dockerfile:1.7-labs

## Build server
FROM golang:1.23-bullseye AS build-server
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY --parents **/*.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server/main.go

## Build frontend
FROM node:lts-slim AS build-ui
WORKDIR /build
COPY ./ui/package.json ./ui/package-lock.json ./
RUN npm ci

COPY ./ui/index.html ./ui/tsconfig.json ./ui/vite.config.ts ./
COPY ./ui/src ./src

RUN npm run build

## Deploy
FROM alpine:3.20
WORKDIR /
COPY --from=build-server /build/main /usr/bin/
COPY --from=build-ui /build/dist /app/public

ENV PUBLIC_DIR="/app/public"
ENTRYPOINT ["main"]
