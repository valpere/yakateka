#!/bin/bash
# Ghostscript Helper - PostScript and PDF converter
# Tools: ps2pdf, pdf2ps, ps2txt, ps2ascii

set -e

case "$1" in
    ping)
        # Check if ps2pdf is available (core tool)
        if command -v ps2pdf >/dev/null 2>&1; then
            echo "pong"
            exit 0
        fi
        exit 1
        ;;

    info)
        cat <<'EOF'
name: "Ghostscript Converter"
version: "1.0.0"
description: "Converts PostScript and PDF files using Ghostscript"
capabilities:
  ps:
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
        quality:
          speed: 1
          quality: 1
    txt:
      modes:
        normal:
          speed: 1
          quality: 1
    eps:
      modes:
        normal:
          speed: 1
          quality: 1
  pdf:
    ps:
      modes:
        normal:
          speed: 1
          quality: 1
    txt:
      modes:
        normal:
          speed: 1
          quality: 1
  eps:
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
    ps:
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

        case "${FROM_FORMAT}_to_${TO_FORMAT}" in
            ps_to_pdf|eps_to_pdf)
                # PostScript to PDF
                case "$MODE" in
                    quality)
                        # High quality: PDF 1.4 with compression
                        ps2pdf14 -dPDFSETTINGS=/prepress "$FROM_FILE" "$TO_FILE"
                        ;;
                    *)
                        # Normal/fast mode
                        ps2pdf "$FROM_FILE" "$TO_FILE"
                        ;;
                esac
                ;;
            pdf_to_ps)
                # PDF to PostScript
                if command -v pdf2ps >/dev/null 2>&1; then
                    pdf2ps "$FROM_FILE" "$TO_FILE"
                else
                    # Fallback to gs directly
                    gs -dNOPAUSE -dBATCH -sDEVICE=ps2write -sOutputFile="$TO_FILE" "$FROM_FILE"
                fi
                ;;
            ps_to_txt|eps_to_txt)
                # PostScript to text
                if command -v ps2txt >/dev/null 2>&1; then
                    ps2txt "$FROM_FILE" > "$TO_FILE"
                elif command -v ps2ascii >/dev/null 2>&1; then
                    ps2ascii "$FROM_FILE" "$TO_FILE"
                else
                    echo "ps2txt or ps2ascii not found" >&2
                    exit 1
                fi
                ;;
            pdf_to_txt)
                # PDF to text via PostScript
                if command -v ps2ascii >/dev/null 2>&1; then
                    ps2ascii "$FROM_FILE" "$TO_FILE"
                else
                    echo "ps2ascii not found" >&2
                    exit 1
                fi
                ;;
            eps_to_ps)
                # EPS is already PostScript, just copy
                cp "$FROM_FILE" "$TO_FILE"
                ;;
            *)
                echo "Unsupported conversion: ${FROM_FORMAT} -> ${TO_FORMAT}" >&2
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
