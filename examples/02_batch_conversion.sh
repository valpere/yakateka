#!/bin/bash
# Example 2: Batch PDF to Text Conversion
# This script converts all PDFs in a directory to text files

set -e

YAKATEKA="../bin/yakateka"

# Check if yakateka is built
if [ ! -f "$YAKATEKA" ]; then
    echo "Error: yakateka not built. Run 'make build' first."
    exit 1
fi

# Input and output directories
INPUT_DIR="../../library4tests"
OUTPUT_DIR="../tmp/batch_output"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "=== Batch PDF Conversion Example ==="
echo ""
echo "Input directory:  $INPUT_DIR"
echo "Output directory: $OUTPUT_DIR"
echo ""

# Check if input directory exists
if [ ! -d "$INPUT_DIR" ]; then
    echo "Error: Input directory not found: $INPUT_DIR"
    exit 1
fi

# Find and convert all PDFs
pdf_count=0
success_count=0
fail_count=0

for pdf in "$INPUT_DIR"/*.pdf; do
    if [ -f "$pdf" ]; then
        pdf_count=$((pdf_count + 1))
        filename=$(basename "$pdf" .pdf)
        output="$OUTPUT_DIR/${filename}.txt"

        echo "[$pdf_count] Converting: $(basename "$pdf")"

        if $YAKATEKA convert "$pdf" "$output" --log-format text > ../tmp/conversion_${pdf_count}.log 2>&1; then
            success_count=$((success_count + 1))
            file_size=$(stat -f%z "$output" 2>/dev/null || stat -c%s "$output" 2>/dev/null)
            echo "    ✓ Success: $output ($file_size bytes)"
        else
            fail_count=$((fail_count + 1))
            echo "    ✗ Failed: $pdf"
            echo "    See log: ../tmp/conversion_${pdf_count}.log"
        fi
        echo ""
    fi
done

echo "=== Batch Conversion Complete ==="
echo ""
echo "Total PDFs: $pdf_count"
echo "Successful: $success_count"
echo "Failed:     $fail_count"
echo ""

if [ $success_count -gt 0 ]; then
    echo "Output files are in: $OUTPUT_DIR"
    ls -lh "$OUTPUT_DIR"
fi
