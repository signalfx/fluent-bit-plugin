FROM golang:1.13 as gobuilder

WORKDIR /root

ENV GOOS=linux\
    GOARCH=amd64

COPY / /root/

RUN go build \
    -buildmode=c-shared \
    -o /signalfx.so \
    github.com/signalfx/signalfx-fluent-bit-plugin

FROM fluent/fluent-bit:1.2.2

COPY --from=gobuilder /signalfx.so /fluent-bit/

EXPOSE 2020

CMD ["/fluent-bit/bin/fluent-bit", "-e", "/fluent-bit/signalfx.so", "-c", "/fluent-bit/etc/fluent-bit.conf"]
