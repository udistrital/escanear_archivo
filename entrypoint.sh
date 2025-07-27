#!/bin/sh

echo "🚀 Iniciando clamd en segundo plano..."
freshclam -d &
clamd --foreground &

echo "[+] Esperando a que ClamAV esté listo..."
CLAMD_SOCKET="/run/clamav/clamd.sock"
MAX_TRIES=90
TRIES=0

while [ ! -S "$CLAMD_SOCKET" ]; do
    if [ "$TRIES" -ge "$MAX_TRIES" ]; then
        echo "[!] ClamAV no se inició después de $MAX_TRIES segundos."
        exit 1
    fi
    echo "  ... esperando a $CLAMD_SOCKET ($TRIES/$MAX_TRIES)"
    TRIES=$((TRIES + 1))
    sleep 1
done

echo "[✓] ClamAV está listo."

exec /main

