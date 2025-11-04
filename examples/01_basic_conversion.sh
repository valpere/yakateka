#!/bin/bash
# Example 1: Basic PDF to Text Conversion
# This script demonstrates the simplest way to convert a PDF to text

set -e  # Exit on error

# Check if yakateka is built
if [ ! -f "../bin/yakateka" ]; then
    echo "Error: yakateka not built. Run 'make build' first."
    exit 1
fi

YAKATEKA="../bin/yakateka"

# Use test document from library4tests if available
if [ -f "../../library4tests/NoSQL_Distilled.pdf" ]; then
    INPUT="../../library4tests/NoSQL_Distilled.pdf"
else
    echo "Error: Test PDF not found at ../../library4tests/NoSQL_Distilled.pdf"
    echo "Please provide a PDF file to test with."
    exit 1
fi

OUTPUT="../tmp/output.txt"

# Create tmp directory
mkdir -p ../tmp

echo "=== Basic PDF to Text Conversion Example ==="
echo ""
echo "Input:  $INPUT"
echo "Output: $OUTPUT"
echo ""

# Convert with text-format logging for readability
$YAKATEKA --log-format text convert "$INPUT" "$OUTPUT"

echo ""
echo "=== Conversion Complete ==="
echo ""

# Display results
if [ -f "$OUTPUT" ]; then
    file_size=$(stat -f%z "$OUTPUT" 2>/dev/null || stat -c%s "$OUTPUT" 2>/dev/null)
    echo "✓ Output file created: $OUTPUT"
    echo "  Size: $file_size bytes ($(numfmt --to=iec-i --suffix=B $file_size 2>/dev/null || echo "$file_size bytes"))"
    echo ""
    echo "First 20 lines of output:"
    head -20 "$OUTPUT"
else
    echo "✗ Output file was not created!"
    exit 1
fi
