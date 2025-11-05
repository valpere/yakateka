#!/bin/bash
# Poppler Helper - PDF utilities
# Tools: pdftotext, pdftohtml, pdftops, pdftocairo, pdftoppm

set -e

case "$1" in
    ping)
        # Check if pdftotext is available (core tool)
        if command -v pdftotext >/dev/null 2>&1; then
            echo "pong"
            exit 0
        fi
        exit 1
        ;;

    info)
        cat <<'EOF'
name: "Poppler PDF Utilities"
version: "1.0.0"
description: "Converts PDF files using Poppler utilities"
capabilities:
  pdf:
    txt:
      modes:
        normal: {speed: 1, quality: 1}
        fast: {speed: 1, quality: 1}
    html:
      modes:
        normal: {speed: 1, quality: 1}
    ps:
      modes:
        normal: {speed: 1, quality: 1}
        quality: {speed: 1, quality: 1}
    png:
      modes:
        normal: {speed: 1, quality: 1}
        quality: {speed: 1, quality: 1}
    svg:
      modes:
        normal: {speed: 1, quality: 1}
    ppm:
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
                # PDF to text
                case "$MODE" in
                    fast)
                        pdftotext -raw "$FROM_FILE" "$TO_FILE"
                        ;;
                    *)
                        # Normal mode: preserve layout
                        pdftotext -layout "$FROM_FILE" "$TO_FILE"
                        ;;
                esac
                ;;
            html)
                # PDF to HTML
                if command -v pdftohtml >/dev/null 2>&1; then
                    # pdftohtml creates output.html, need to handle this
                    TEMP_DIR=$(mktemp -d)
                    trap "rm -rf $TEMP_DIR" EXIT

                    pdftohtml -noframes -s "$FROM_FILE" "$TEMP_DIR/output.html"
                    mv "$TEMP_DIR/output.html" "$TO_FILE"
                else
                    echo "pdftohtml not found" >&2
                    exit 1
                fi
                ;;
            ps)
                # PDF to PostScript
                if command -v pdftops >/dev/null 2>&1; then
                    case "$MODE" in
                        quality)
                            pdftops -level3 "$FROM_FILE" "$TO_FILE"
                            ;;
                        *)
                            pdftops "$FROM_FILE" "$TO_FILE"
                            ;;
                    esac
                else
                    echo "pdftops not found" >&2
                    exit 1
                fi
                ;;
            png)
                # PDF to PNG (first page only)
                if command -v pdftocairo >/dev/null 2>&1; then
                    case "$MODE" in
                        quality)
                            pdftocairo -png -singlefile -r 300 "$FROM_FILE" "${TO_FILE%.png}"
                            ;;
                        *)
                            pdftocairo -png -singlefile -r 150 "$FROM_FILE" "${TO_FILE%.png}"
                            ;;
                    esac
                else
                    echo "pdftocairo not found" >&2
                    exit 1
                fi
                ;;
            svg)
                # PDF to SVG (first page only)
                if command -v pdftocairo >/dev/null 2>&1; then
                    pdftocairo -svg "$FROM_FILE" "${TO_FILE%.svg}"
                else
                    echo "pdftocairo not found" >&2
                    exit 1
                fi
                ;;
            ppm)
                # PDF to PPM
                if command -v pdftoppm >/dev/null 2>&1; then
                    pdftoppm -f 1 -singlefile "$FROM_FILE" "${TO_FILE%.ppm}"
                else
                    echo "pdftoppm not found" >&2
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
