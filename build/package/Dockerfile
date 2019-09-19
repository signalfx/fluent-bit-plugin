FROM golang:1.13 as gobuilder

WORKDIR /root

ENV GOOS=linux\
    GOARCH=amd64

COPY / /root/

RUN make

FROM fluent/fluent-bit:1.2.2

COPY --from=gobuilder /root/dist/signalfx.so /fluent-bit/

EXPOSE 2020

CMD ["/fluent-bit/bin/fluent-bit", "-e", "/fluent-bit/signalfx.so", "-c", "/fluent-bit/etc/fluent-bit.conf"]
