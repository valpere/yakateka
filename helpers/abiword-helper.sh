#!/bin/bash
# AbiWord Helper - Word processor document converter
# Tool: abiword

set -e

ABIWORD_BIN="${ABIWORD_BIN:-abiword}"

case "$1" in
    ping)
        # Check if abiword is available
        if command -v "$ABIWORD_BIN" >/dev/null 2>&1; then
            echo "pong"
            exit 0
        fi
        exit 1
        ;;

    info)
        cat <<'EOF'
name: "AbiWord Converter"
version: "1.0.0"
description: "Converts word processing documents using AbiWord"
capabilities:
  doc:
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
    txt:
      modes:
        normal:
          speed: 1
          quality: 1
    html:
      modes:
        normal:
          speed: 1
          quality: 1
    rtf:
      modes:
        normal:
          speed: 1
          quality: 1
    odt:
      modes:
        normal:
          speed: 1
          quality: 1
  docx:
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
    txt:
      modes:
        normal:
          speed: 1
          quality: 1
    html:
      modes:
        normal:
          speed: 1
          quality: 1
    rtf:
      modes:
        normal:
          speed: 1
          quality: 1
  rtf:
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
    txt:
      modes:
        normal:
          speed: 1
          quality: 1
    html:
      modes:
        normal:
          speed: 1
          quality: 1
    doc:
      modes:
        normal:
          speed: 1
          quality: 1
  odt:
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
    txt:
      modes:
        normal:
          speed: 1
          quality: 1
    html:
      modes:
        normal:
          speed: 1
          quality: 1
    doc:
      modes:
        normal:
          speed: 1
          quality: 1
EOF
        exit 0
        ;;

    convert)
        MODE="$2"
        FROM_FORMAT="$3"
        FROM_FILE="$4"
        TO_FORMAT="$5"
        TO_FILE="$6"

        if [ -z "$MODE" ] || [ -z "$FROM_FORMAT" ] || [ -z "$FROM_FILE" ] || [ -z "$TO_FORMAT" ] || [ -z "$TO_FILE" ]; then
            echo "Usage: $0 convert <mode> <from_format> <from_file> <to_format> <to_file>" >&2
            exit 1
        fi

        # AbiWord uses file extensions to determine format
        # Execute conversion with --to parameter
        "$ABIWORD_BIN" --to="$TO_FORMAT" --to-name="$TO_FILE" "$FROM_FILE"

        # Verify output was created
        if [ ! -f "$TO_FILE" ]; then
            echo "Conversion failed: output file not created" >&2
            exit 1
        fi

        exit 0
        ;;

    *)
        echo "Unknown command: $1" >&2
        echo "Usage: $0 {ping|info|convert}" >&2
        exit 1
        ;;
esac
