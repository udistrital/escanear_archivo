FROM alpine:latest

RUN apk add --no-cache clamav clamav-daemon clamav-libunrar curl

RUN freshclam || true

WORKDIR /

COPY main main
COPY conf/app.conf conf/app.conf

RUN chmod +x main

ENTRYPOINT ["/main"]

##################


# Etapa 1: build del binario
#FROM golang:1.24-alpine AS builder

#WORKDIR /app
#COPY . .
#RUN go build -o main ./main.go

# Etapa 2: contenedor liviano con ClamAV
#FROM alpine:latest

#RUN apk add --no-cache clamav clamav-daemon clamav-libunrar curl

#RUN freshclam || true

#WORKDIR /

#COPY --from=builder /app/main /main
#COPY conf/app.conf conf/app.conf

#RUN chmod +x /main

#ENTRYPOINT ["/main"]

