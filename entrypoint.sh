#!/bin/sh

echo "🚀 Iniciando clamd en segundo plano..."
clamd --foreground &

echo "⏳ Esperando a que clamd esté listo..."
for i in $(seq 1 20); do
    if clamdscan --version > /dev/null 2>&1; then
        echo "✅ clamd listo. Iniciando servicio Go..."
        break
    fi
    sleep 1
done

exec /main
