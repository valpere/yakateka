#!/bin/bash
# LibreOffice Helper - Office document converter
# Supports: DOC, DOCX, ODT, RTF, PDF, PS conversions

set -e

LIBREOFFICE_BIN="${LIBREOFFICE_BIN:-soffice}"

case "$1" in
    ping)
        # Check if LibreOffice is available
        if command -v "$LIBREOFFICE_BIN" >/dev/null 2>&1; then
            echo "pong"
            exit 0
        fi
        exit 1
        ;;

    info)
        # Check if Pandoc is available for DOC→MD pipeline
        HAS_PANDOC=false
        if command -v pandoc >/dev/null 2>&1; then
            HAS_PANDOC=true
        fi

        # Set description based on Pandoc availability
        if [ "$HAS_PANDOC" = true ]; then
            DESCRIPTION="Converts office documents using LibreOffice (DOC→MD via Pandoc pipeline)"
        else
            DESCRIPTION="Converts office documents using LibreOffice"
        fi

        cat <<EOF
name: "LibreOffice Converter"
version: "1.0.0"
description: "$DESCRIPTION"
capabilities:
  pdf:
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
  ps:
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
  doc:
    pdf:
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
    docx:
      modes:
        normal:
          speed: 1
          quality: 1
    odt:
      modes:
        normal:
          speed: 1
          quality: 1
EOF

        # Conditionally add DOC→MD if Pandoc is available
        if [ "$HAS_PANDOC" = true ]; then
            cat <<'EOF'
    md:
      modes:
        normal:
          speed: 1
          quality: 1
EOF
        fi

        cat <<'EOF'
  docx:
    pdf:
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
    odt:
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
    pdf:
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
    docx:
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
    docx:
      modes:
        normal:
          speed: 1
          quality: 1
    odt:
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

        # Special case: DOC → MD requires pipeline (DOC → DOCX → MD)
        if [ "$FROM_FORMAT" = "doc" ] && [ "$TO_FORMAT" = "md" ]; then
            # Create temporary DOCX file
            TEMP_DOCX=$(mktemp --suffix=.docx)
            trap 'rm -f "$TEMP_DOCX"' EXIT

            # Step 1: DOC → DOCX using LibreOffice
            OUTDIR=$(dirname "$TEMP_DOCX")
            BASENAME=$(basename "$FROM_FILE" ".doc")
            "$LIBREOFFICE_BIN" --headless --convert-to docx --outdir "$OUTDIR" "$FROM_FILE"

            EXPECTED_DOCX="$OUTDIR/${BASENAME}.docx"
            if [ -f "$EXPECTED_DOCX" ]; then
                mv "$EXPECTED_DOCX" "$TEMP_DOCX"
            fi

            if [ ! -f "$TEMP_DOCX" ]; then
                echo "Failed to convert DOC to DOCX" >&2
                exit 1
            fi

            # Step 2: DOCX → MD using Pandoc
            if ! command -v pandoc >/dev/null 2>&1; then
                echo "Pandoc not found, cannot complete DOC → MD conversion" >&2
                exit 1
            fi

            pandoc -f docx -t markdown -o "$TO_FILE" "$TEMP_DOCX"
            exit $?
        fi

        # Get output directory and basename
        OUTDIR=$(dirname "$TO_FILE")
        BASENAME=$(basename "$FROM_FILE" ".$FROM_FORMAT")

        # Map formats to LibreOffice format specs
        case "$TO_FORMAT" in
            pdf)
                if [ "$FROM_FORMAT" = "ps" ] || [ "$FROM_FORMAT" = "pdf" ]; then
                    LO_FORMAT="pdf:writer_pdf_Export"
                else
                    LO_FORMAT="pdf:writer_pdf_Export"
                fi
                ;;
            html) LO_FORMAT="html" ;;
            txt) LO_FORMAT="txt:Text (encoded):UTF8" ;;
            docx) LO_FORMAT="docx" ;;
            odt) LO_FORMAT="odt" ;;
            rtf) LO_FORMAT="rtf" ;;
            *) LO_FORMAT="$TO_FORMAT" ;;
        esac

        # Convert using LibreOffice
        # LibreOffice creates output with input basename in outdir
        "$LIBREOFFICE_BIN" --headless --convert-to "$LO_FORMAT" --outdir "$OUTDIR" "$FROM_FILE"

        # Move/rename to expected output path
        EXPECTED_OUTPUT="$OUTDIR/${BASENAME}.${TO_FORMAT}"
        if [ -f "$EXPECTED_OUTPUT" ] && [ "$EXPECTED_OUTPUT" != "$TO_FILE" ]; then
            mv "$EXPECTED_OUTPUT" "$TO_FILE"
        fi

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
