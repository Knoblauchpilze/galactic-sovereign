FROM golang:1.22 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/stellar-dominion cmd/stellar-dominion/main.go
COPY cmd/stellar-dominion/configs/*-prod.yml bin/configs/

FROM alpine AS stellar-dominion-service
ARG GIT_COMMIT_HASH="undefined"
ENV GIT_COMMIT_HASH=$GIT_COMMIT_HASH
WORKDIR /app
COPY --from=builder build/bin bin/
WORKDIR /app/bin
CMD ["./stellar-dominion"]
