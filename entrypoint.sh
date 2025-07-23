#!/bin/sh

echo "🚀 Iniciando clamd en segundo plano..."
clamd --foreground &

echo "⏳ Esperando a que clamd cree el socket..."
for i in $(seq 1 20); do
    if [ -S /run/clamav/clamd.sock ]; then
        echo "✅ clamd listo. Iniciando servicio Go..."
        break
    fi
    sleep 1
done

exec /main
