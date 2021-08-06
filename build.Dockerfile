FROM --platform=amd64 golang:1.16.7

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

WORKDIR /go/src/app
COPY . .

RUN make build

FROM alpine:3.14.0

RUN apk update && apk add --no-cache ca-certificates 

FROM scratch

COPY --from=0 /go/src/app/nbot /nbot
COPY --from=1 /etc/ssl/certs /etc/ssl/certs

ENTRYPOINT ["/nbot"]
