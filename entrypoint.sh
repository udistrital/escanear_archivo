#!/bin/sh

echo "🚀 Iniciando clamd en segundo plano..."
clamd --foreground &

# Esperar un segundo por si crashea al instante
sleep 1

# Verificar si clamd murió
if ! pgrep -x clamd > /dev/null; then
  echo "❌ clamd murió. Mostrando log:"
  if [ -f /var/log/clamav/clamd.log ]; then
    cat /var/log/clamav/clamd.log
  else
    echo "⚠️ Log no encontrado. ¿Quizás clamd nunca arrancó por falta de base de datos?"
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

exec /main

