#!/usr/bin/env bash
set -euo pipefail

# Trova la root del repository
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
cd "$repo_root"

# Opzioni configurabili via env:
#   PKG=./tests/...   pacchetti da testare (default: cartella tests)
#   RACE=1            abilita race detector (1|0)
#   COVER=1           abilita coverage e report (1|0)
#   VERBOSE=1         output verboso (1|0)
#   COUNT=1           disabilita cache di go test (default: 1)
#   TIMEOUT=60s       timeout per package
#   COVERPKG=./...    pacchetti da includere nella coverage (strumentazione)
#   RUN=              regex per filtrare i test (go test -run)
PKG="${PKG:-./tests/...}"
RACE="${RACE:-1}"
COVER="${COVER:-1}"
VERBOSE="${VERBOSE:-1}"
COUNT="${COUNT:-1}"
TIMEOUT="${TIMEOUT:-60s}"
COVERPKG="${COVERPKG:-./...}"
RUN="${RUN:-}"

args=(test "$PKG" "-count=${COUNT}" "-timeout=${TIMEOUT}")
[[ "$VERBOSE" == "1" ]] && args+=("-v")
[[ "$RACE" == "1" ]] && args+=("-race")
[[ -n "$RUN" ]] && args+=("-run" "$RUN")
[[ "$COVER" == "1" ]] && args+=("-covermode=atomic" "-coverpkg=${COVERPKG}" "-coverprofile=coverage.out")

echo "+ go ${args[*]}"
go "${args[@]}"

if [[ "$COVER" == "1" ]]; then
  echo "+ go tool cover -func=coverage.out"
  go tool cover -func=coverage.out || true
  echo "+ go tool cover -html=coverage.out -o coverage.html"
  go tool cover -html=coverage.out -o coverage.html || true
  echo "Coverage report: ${repo_root}/coverage.out"
  echo "HTML report:     ${repo_root}/coverage.html"
fi
