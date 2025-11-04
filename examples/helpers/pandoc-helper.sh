#!/bin/bash
# Pandoc Helper - Universal markup converter
# Supports: MD, HTML, DOCX, ODT, EPUB, PDF conversions

set -e

PANDOC_BIN="${PANDOC_BIN:-pandoc}"

# Check if pandoc is available
if ! command -v "$PANDOC_BIN" &> /dev/null; then
    exit 1
fi

case "$1" in
    ping)
        echo "pong"
        exit 0
        ;;

    info)
        cat <<EOF
name: "Pandoc Universal Converter"
version: "1.0.0"
description: "Converts between markup formats using Pandoc"
capabilities:
  md:
    html:
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
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
        quality:
          speed: 1
          quality: 1
    docx:
      modes:
        normal:
          speed: 1
          quality: 1
    epub:
      modes:
        normal:
          speed: 1
          quality: 1
  html:
    md:
      modes:
        normal:
          speed: 1
          quality: 1
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
    docx:
      modes:
        normal:
          speed: 1
          quality: 1
  docx:
    md:
      modes:
        normal:
          speed: 1
          quality: 1
    html:
      modes:
        normal:
          speed: 1
          quality: 1
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
  epub:
    md:
      modes:
        normal:
          speed: 1
          quality: 1
    html:
      modes:
        normal:
          speed: 1
          quality: 1
    txt:
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

        # Map formats to Pandoc format names
        case "$FROM_FORMAT" in
            md) FROM_FORMAT="markdown" ;;
            txt) FROM_FORMAT="plain" ;;
        esac

        case "$TO_FORMAT" in
            md) TO_FORMAT="markdown" ;;
            txt) TO_FORMAT="plain" ;;
        esac

        # Build Pandoc command based on mode
        EXTRA_ARGS=""
        case "$MODE" in
            quality)
                if [ "$TO_FORMAT" = "pdf" ]; then
                    EXTRA_ARGS="--toc --number-sections --pdf-engine=xelatex"
                else
                    EXTRA_ARGS="--toc --standalone"
                fi
                ;;
            fast)
                EXTRA_ARGS="--no-highlight"
                ;;
            normal)
                if [ "$TO_FORMAT" = "pdf" ]; then
                    EXTRA_ARGS="--pdf-engine=xelatex"
                fi
                ;;
        esac

        # Execute conversion
        $PANDOC_BIN -f "$FROM_FORMAT" -t "$TO_FORMAT" "$FROM_FILE" -o "$TO_FILE" $EXTRA_ARGS
        exit $?
        ;;

    *)
        echo "Unknown command: $1" >&2
        echo "Usage: $0 {ping|info|convert}" >&2
        exit 1
        ;;
esac
