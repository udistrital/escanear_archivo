#!/bin/sh
set -euo pipefail

echo "🚀 Iniciando freshclam y clamd en segundo plano"
freshclam -d &
clamd --foreground &

echo "🚀 Iniciando API"
exec /main

