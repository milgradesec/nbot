FROM --platform=amd64 golang:1.18

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

FROM alpine:3.15.2

RUN apk update && \
    apk --no-cache add ca-certificates && \
    addgroup -S nbot && \
    adduser -S -G nbot nbot

FROM scratch

COPY --from=0 /go/src/app/nbot /nbot
COPY --from=1 /etc/ssl/certs /etc/ssl/certs
COPY --from=1 /etc/passwd /etc/passwd

USER nbot
ENTRYPOINT ["/nbot"]
