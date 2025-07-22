FROM alpine:latest

RUN apk add --no-cache clamav clamav-daemon clamav-libunrar curl

RUN freshclam || true

WORKDIR /

COPY main main
COPY conf/app.conf conf/app.conf

RUN chmod +x main

ENTRYPOINT ["/main"]
