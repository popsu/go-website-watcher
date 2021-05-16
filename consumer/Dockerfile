# 'docker build' should be executed from parent directory
FROM golang:1.16.4-alpine3.13 as builder

# Install upx in order to compress binary
RUN apk add --no-cache build-base upx ca-certificates

WORKDIR /builder/

COPY . .

RUN make bin/release/gww-consumer

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /builder/bin/release/gww-consumer /gww-consumer

ENTRYPOINT ["/gww-consumer"]
