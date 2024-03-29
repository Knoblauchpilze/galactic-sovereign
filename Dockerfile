FROM golang:1.22 AS builder
WORKDIR /server
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o build/bin/server cmd/server/main.go

FROM alpine:latest AS user-service
WORKDIR /server
COPY --from=builder server/build/bin .
COPY --from=builder server/configs configs/
# https://stackoverflow.com/questions/21553353/what-is-the-difference-between-cmd-and-entrypoint-in-a-dockerfile
CMD ["./server"]
