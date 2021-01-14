FROM golang:1.15-buster as builder

ADD . /build
WORKDIR /build

RUN go get ./...
RUN make

FROM alpine

COPY --from=builder /build/extres /usr/local/bin/extres

ENTRYPOINT ["/usr/local/bin/extres"]
