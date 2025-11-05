#!/bin/bash
# Calibre Helper - Ebook format converter
# Supports: MOBI, EPUB, FB2, AZW, AZW3, and many other ebook formats

set -e

EBOOK_CONVERT_BIN="${EBOOK_CONVERT_BIN:-ebook-convert}"

case "$1" in
    ping)
        # Check if ebook-convert is available
        if command -v "$EBOOK_CONVERT_BIN" >/dev/null 2>&1; then
            echo "pong"
            exit 0
        fi
        exit 1
        ;;

    info)
        cat <<'EOF'
name: "Calibre Ebook Converter"
version: "1.0.0"
description: "Converts between ebook formats using Calibre"
capabilities:
  mobi:
    epub:
      modes:
        normal:
          speed: 1
          quality: 1
        quality:
          speed: 1
          quality: 1
    fb2:
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
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
  epub:
    mobi:
      modes:
        normal:
          speed: 1
          quality: 1
        quality:
          speed: 1
          quality: 1
    fb2:
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
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
  fb2:
    epub:
      modes:
        normal:
          speed: 1
          quality: 1
    mobi:
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
  azw:
    epub:
      modes:
        normal:
          speed: 1
          quality: 1
    mobi:
      modes:
        normal:
          speed: 1
          quality: 1
    txt:
      modes:
        normal:
          speed: 1
          quality: 1
  azw3:
    epub:
      modes:
        normal:
          speed: 1
          quality: 1
    mobi:
      modes:
        normal:
          speed: 1
          quality: 1
    txt:
      modes:
        normal:
          speed: 1
          quality: 1
  lit:
    epub:
      modes:
        normal:
          speed: 1
          quality: 1
    mobi:
      modes:
        normal:
          speed: 1
          quality: 1
  pdb:
    epub:
      modes:
        normal:
          speed: 1
          quality: 1
    mobi:
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

        # Build ebook-convert command based on mode
        EXTRA_ARGS=""
        case "$MODE" in
            quality)
                # High quality mode - pretty-print and preserve formatting
                EXTRA_ARGS="--pretty-print --preserve-cover-aspect-ratio"
                ;;
            fast)
                # Fast mode - minimal processing
                EXTRA_ARGS="--no-inline-toc"
                ;;
            normal)
                # Normal mode
                EXTRA_ARGS=""
                ;;
        esac

        # Execute conversion
        "$EBOOK_CONVERT_BIN" "$FROM_FILE" "$TO_FILE" $EXTRA_ARGS
        exit $?
        ;;

    *)
        echo "Unknown command: $1" >&2
        echo "Usage: $0 {ping|info|convert}" >&2
        exit 1
        ;;
esac
