# Etapa 1: Compilación del binario Go
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Etapa 2: Imagen final
FROM alpine:latest

# Instalación de ClamAV y dependencias
RUN apk add --no-cache clamav clamav-daemon curl

# Crear rutas necesarias
RUN mkdir -p /run/clamav /var/run/clamav /var/lib/clamav /etc/clamav /var/log/clamav

# Configuración básica de clamd
COPY src/clamd.conf /etc/clamav/clamd.conf

# Descargar las firmas ahora para evitar depender del entorno
RUN freshclam --verbose

# Copiar binario de Go y configuración
COPY --from=builder /app/main /main
COPY conf/app.conf /conf/app.conf

# Copiar script de entrada
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /main /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
