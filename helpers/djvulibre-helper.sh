#!/bin/bash
# DjVuLibre Helper - DjVu document converter
# Tools: djvutxt, djvups, ddjvu (image extraction)

set -e

case "$1" in
    ping)
        # Check if djvutxt is available (core tool)
        if command -v djvutxt >/dev/null 2>&1; then
            echo "pong"
            exit 0
        fi
        exit 1
        ;;

    info)
        cat <<'EOF'
name: "DjVuLibre Converter"
version: "1.0.0"
description: "Converts DjVu documents to other formats"
capabilities:
  djvu:
    txt:
      modes:
        normal: {speed: 1, quality: 1}
    ps:
      modes:
        normal: {speed: 1, quality: 1}
    pdf:
      modes:
        normal: {speed: 1, quality: 1}
    png:
      modes:
        normal: {speed: 1, quality: 1}
    tiff:
      modes:
        normal: {speed: 1, quality: 1}
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

        case "$TO_FORMAT" in
            txt)
                # Extract text from DjVu
                djvutxt "$FROM_FILE" "$TO_FILE"
                ;;
            ps)
                # Convert DjVu to PostScript
                if command -v djvups >/dev/null 2>&1; then
                    djvups -o "$TO_FILE" "$FROM_FILE"
                else
                    echo "djvups not found" >&2
                    exit 1
                fi
                ;;
            pdf)
                # Convert DjVu to PDF via PostScript
                if command -v djvups >/dev/null 2>&1 && command -v ps2pdf >/dev/null 2>&1; then
                    TEMP_PS=$(mktemp --suffix=.ps)
                    trap "rm -f $TEMP_PS" EXIT

                    djvups -o "$TEMP_PS" "$FROM_FILE"
                    ps2pdf "$TEMP_PS" "$TO_FILE"
                else
                    echo "djvups or ps2pdf not found" >&2
                    exit 1
                fi
                ;;
            png|tiff|tif)
                # Extract as image
                if command -v ddjvu >/dev/null 2>&1; then
                    FORMAT="$TO_FORMAT"
                    [ "$FORMAT" = "tif" ] && FORMAT="tiff"
                    ddjvu -format="$FORMAT" "$FROM_FILE" "$TO_FILE"
                else
                    echo "ddjvu not found" >&2
                    exit 1
                fi
                ;;
            *)
                echo "Unsupported output format: $TO_FORMAT" >&2
                exit 1
                ;;
        esac

        exit 0
        ;;

    *)
        echo "Unknown command: $1" >&2
        echo "Usage: $0 {ping|info|convert}" >&2
        exit 1
        ;;
esac
