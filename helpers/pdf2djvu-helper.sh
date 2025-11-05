#!/bin/bash
# pdf2djvu Helper - PDF to DjVu converter
# Tool: pdf2djvu

set -e

PDF2DJVU_BIN="${PDF2DJVU_BIN:-pdf2djvu}"

case "$1" in
    ping)
        # Check if pdf2djvu is available
        if command -v "$PDF2DJVU_BIN" >/dev/null 2>&1; then
            echo "pong"
            exit 0
        fi
        exit 1
        ;;

    info)
        cat <<'EOF'
name: "pdf2djvu Converter"
version: "1.0.0"
description: "Converts PDF files to DjVu format"
capabilities:
  pdf:
    djvu:
      modes:
        normal:
          speed: 1
          quality: 1
        fast:
          speed: 1
          quality: 1
        quality:
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

        if [ "$FROM_FORMAT" != "pdf" ] || [ "$TO_FORMAT" != "djvu" ]; then
            echo "Only PDF to DjVu conversion is supported" >&2
            exit 1
        fi

        # Build pdf2djvu command based on mode
        EXTRA_ARGS=()
        case "$MODE" in
            quality)
                # High quality mode
                EXTRA_ARGS=(--dpi=600 --loss-level=0)
                ;;
            fast)
                # Fast mode - lower DPI
                EXTRA_ARGS=(--dpi=150 --loss-level=100)
                ;;
            normal)
                # Normal mode
                EXTRA_ARGS=(--dpi=300)
                ;;
        esac

        # Execute conversion
        "$PDF2DJVU_BIN" -o "$TO_FILE" "${EXTRA_ARGS[@]}" "$FROM_FILE"
        exit $?
        ;;

    *)
        echo "Unknown command: $1" >&2
        echo "Usage: $0 {ping|info|convert}" >&2
        exit 1
        ;;
esac
