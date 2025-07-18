FROM --platform=amd64 golang:1.24.5 AS builder

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

FROM gcr.io/distroless/static-debian11:nonroot

COPY --from=builder --chown=nonroot /go/src/app/nbot /nbot

USER nonroot
ENTRYPOINT ["/nbot"]
