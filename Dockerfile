# Etapa 1: compila la app Go
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Etapa 2: imagen final con ClamAV y el binario
FROM alpine:latest

RUN apk add --no-cache clamav clamav-daemon curl

RUN mkdir -p /run/clamav /var/run/clamav /var/lib/clamav /etc/clamav /var/log/clamav

RUN echo "LogSyslog yes" > /etc/clamav/clamd.conf && \
    echo "LogVerbose yes" >> /etc/clamav/clamd.conf && \
    echo "PidFile /var/run/clamav/clamd.pid" >> /etc/clamav/clamd.conf && \
    echo "LocalSocket /run/clamav/clamd.sock" >> /etc/clamav/clamd.conf && \
    echo "FixStaleSocket yes" >> /etc/clamav/clamd.conf && \
    echo "TCPSocket 3310" >> /etc/clamav/clamd.conf && \
    echo "TCPAddr 127.0.0.1" >> /etc/clamav/clamd.conf && \
    echo "ScanMail yes" >> /etc/clamav/clamd.conf && \
    echo "ScanArchive yes" >> /etc/clamav/clamd.conf && \
    echo "ScanPDF yes" >> /etc/clamav/clamd.conf && \
    echo "LogFile /var/log/clamav/clamd.log" >> /etc/clamav/clamd.conf

RUN freshclam 

# Copiar binario desde etapa anterior
COPY --from=builder /app/main /main
COPY conf/app.conf /conf/app.conf

RUN chmod +x /main

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]




##################



#FROM alpine:latest

#RUN apk add --no-cache clamav clamav-daemon clamav-libunrar curl

#RUN freshclam || true

#WORKDIR /

#COPY main main
#COPY conf/app.conf conf/app.conf

#RUN chmod +x main

#ENTRYPOINT ["/main"]


