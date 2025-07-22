FROM alpine:latest

# Instalar ClamAV y herramientas necesarias
RUN apk add --no-cache clamav clamav-lib clamav-daemon clamav-libunrar curl

# Actualizar las firmas de virus
RUN freshclam

WORKDIR /

# Copiar binario de la API
COPY main main
COPY conf/app.conf conf/app.conf

# Asignar permisos de ejecución
RUN chmod +x main

# Ejecutar la API
ENTRYPOINT ["/main"]
