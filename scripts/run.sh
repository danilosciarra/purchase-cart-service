#!/usr/bin/env bash
set -euo pipefail

# Trova la root del repository
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
cd "$repo_root"

# Variabili d'ambiente utili:
#   GIN_MODE=release   modalità di Gin (debug|release|test)
#   CONFIG_PATH=...    percorso al file di configurazione (default: config.json)
GIN_MODE="${GIN_MODE:-release}"
export GIN_MODE

CONFIG_PATH="${CONFIG_PATH:-config.json}"

if [[ ! -f "$CONFIG_PATH" ]]; then
  echo "WARN: config file '${CONFIG_PATH}' non trovato. Procedo comunque..." >&2
fi

echo "Avvio Purchase Cart Service"
echo "  GIN_MODE:    ${GIN_MODE}"
echo "  CONFIG_PATH: ${CONFIG_PATH}"
echo "+ go run ./main.go"
exec go run ./main.go
