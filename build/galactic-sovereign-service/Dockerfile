FROM golang:1.23.2 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/galactic-sovereign cmd/galactic-sovereign/main.go
COPY cmd/galactic-sovereign/configs/*-prod.yml bin/configs/

FROM alpine AS galactic-sovereign-service
ARG GIT_COMMIT_HASH="undefined"
ENV GIT_COMMIT_HASH=$GIT_COMMIT_HASH
WORKDIR /app
COPY --from=builder build/bin bin/
WORKDIR /app/bin
CMD ["./galactic-sovereign"]
