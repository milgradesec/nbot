FROM gcr.io/distroless/static-debian11:nonroot

ADD nbot /nbot

USER nonroot
ENTRYPOINT ["/nbot"]
