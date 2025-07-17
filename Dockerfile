# Stage 1: Build Go binary
FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila para Linux con arquitectura adecuada para Lambda
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd

# Stage 2: Runtime con ClamAV y Lambda base image
FROM amazonlinux:2

# Requerido por AWS Lambda para contenedor personalizado
ENV AWS_LAMBDA_RUNTIME_API=""

# Instalar ClamAV y dependencias
RUN yum install -y \
    clamav \
    clamav-update \
    tar \
    gzip \
    unzip \
    shadow-utils \
    && yum clean all

# Actualiza firmas de virus
RUN freshclam

# Copiar binario de la etapa anterior
COPY --from=builder /app/bootstrap /var/task/bootstrap

# Hacer ejecutable
RUN chmod +x /var/task/bootstrap

# Comando de entrada para Lambda
CMD ["/var/task/bootstrap"]
