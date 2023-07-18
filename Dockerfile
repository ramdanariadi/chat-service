FROM golang:1.20-alpine AS build

WORKDIR $GOPATH/src/github.com/ramdanariadi/chat-service

COPY . .
RUN go mod download
RUN go build -o /app

EXPOSE 8080
ENTRYPOINT ["/app"]