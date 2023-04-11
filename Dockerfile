FROM golang:1.20-buster as builder

WORKDIR /build
COPY . .

RUN go build -o ./bin/pmssd ./pmssd

FROM debian:buster-slim

COPY --from=builder /build/bin/pmssd /usr/local/bin/pmss

ENTRYPOINT [ "pmss", "server" ]

EXPOSE 8080