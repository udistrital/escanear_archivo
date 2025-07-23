#!/bin/sh

echo "🚀 Iniciando clamd en segundo plano..."
clamd --foreground &

sleep 1

# Verificar si clamd murió inmediatamente
if ! pgrep -x clamd > /dev/null; then
  echo "❌ clamd murió. Mostrando log (si existe)..."
  if [ -f /var/log/clamav/clamd.log ]; then
    cat /var/log/clamav/clamd.log
  else
    echo "⚠️ Log no encontrado. ¿Tal vez clamd no arrancó por falta de base de datos?"
  fi
  exit 1
fi

echo "⏳ Esperando a que clamd cree el socket..."
for i in $(seq 1 20); do
    if clamdscan --version > /dev/null 2>&1; then
        echo "✅ clamd listo. Iniciando servicio Go..."
        break
    fi
    sleep 1
done

# Verifica por si el socket nunca apareció
if ! clamdscan --version > /dev/null 2>&1; then
    echo "❌ clamd no respondió después de 20s. Mostrando log:"
    cat /var/log/clamav/clamd.log || echo "No se pudo leer el log"
    exit 1
fi

exec /main
