FROM --platform=linux/amd64 golang:1.13.7 AS builder

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w"

FROM alpine:3.11.3
LABEL Name=nbot Version=0.3

RUN apk update && \
    apk --no-cache add ca-certificates && \
    addgroup -S nbot && adduser -S -G nbot nbot

COPY --from=0 /go/src/app/nbot /nbot
USER nbot
ENTRYPOINT ["/nbot"]