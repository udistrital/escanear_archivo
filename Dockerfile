FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine:latest

RUN apk add --no-cache clamav clamav-daemon curl

RUN mkdir -p /run/clamav /var/run/clamav /var/lib/clamav /etc/clamav /var/log/clamav

COPY src/clamd.conf /etc/clamav/clamd.conf
COPY src/freshclam.conf /etc/clamav/freshclam.conf

RUN freshclam 

COPY --from=builder /app/main /main
COPY conf/app.conf /conf/app.conf

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /main /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
