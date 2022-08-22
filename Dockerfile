FROM golang:1.19-alpine as builder

RUN apk add --no-cache --update ca-certificates

RUN apk add --no-cache --update \
    ca-certificates \
    gcc \
    make \
    musl-dev

WORKDIR /build

# cache deps
COPY ./go.mod ./go.sum /build/
RUN go mod download

COPY ./cmd /build/cmd
COPY ./internal /build/internal
COPY ./pkg /build/pkg
COPY Makefile /build/Makefile

RUN make V=1 bin

FROM alpine:latest

WORKDIR /w
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/bin/bot /w/bot
COPY --from=builder /build/bin/validator /w/validator

CMD [ "/w/bot", "-config", "/w/config.yml"]
