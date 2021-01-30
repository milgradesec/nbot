FROM golang:1.15.7

WORKDIR /go/src/app
COPY . .

RUN make build

FROM alpine:3.13

RUN apk update && apk add --no-cache ca-certificates 

FROM scratch

COPY --from=0 /go/src/app/nbot /nbot
COPY --from=1 /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/nbot"]
