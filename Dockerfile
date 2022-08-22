FROM golang:1.19-alpine as builder

RUN apk add --no-cache --update ca-certificates

RUN apk add --no-cache --update \
    ca-certificates \
    gcc \
    # linux-headers \
    make \
    musl-dev
    # git

WORKDIR /build

# cache deps
COPY ./go.mod ./go.sum /build/
RUN go mod download

COPY ./cmd /build/cmd
COPY ./internal /build/internal
COPY ./pkg /build/pkg

RUN go build  \
  -ldflags "-linkmode external -extldflags -static" \
  -o bot \
  /build/cmd/bot
RUN go build  \
  -ldflags "-linkmode external -extldflags -static" \
  -o validator \
  /build/cmd/validator

FROM scratch

WORKDIR /w
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/bot /w/bot
COPY --from=builder /build/validator /w/validator

CMD [ "/w/bot", "-config", "/w/config.yml"]
