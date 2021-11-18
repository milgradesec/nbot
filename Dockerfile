FROM alpine:3.14.3

RUN apk update && \
    apk --no-cache add ca-certificates && \
    addgroup -S nbot && \
    adduser -S -G nbot nbot

FROM scratch

COPY --from=0 /etc/ssl/certs /etc/ssl/certs
COPY --from=0 /etc/passwd /etc/passwd

ADD nbot /nbot

USER nbot
ENTRYPOINT ["/nbot"]
