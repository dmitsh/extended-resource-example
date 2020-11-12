FROM ubuntu:20.04

COPY extres /usr/local/bin/extres

ENTRYPOINT ["/usr/local/bin/extres"]
