#!/bin/bash
# Example 1: Basic PDF to Text Conversion
# This script demonstrates the simplest way to convert a PDF to text

# Color codes for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Determine the absolute path to the project root directory regardless of where the script is called from
# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
# Get the project root directory (parent of the scripts directory)
PROJECT_ROOT="$( cd "${SCRIPT_DIR}/.." &> /dev/null && pwd )"

YAKATEKA="$PROJECT_ROOT/bin/yakateka"
LIBRARY4TESTS="$PROJECT_ROOT/../library4tests"
LOCAL_TMP_DIR="$PROJECT_ROOT/tmp"

# echo $SCRIPT_DIR
# echo $PROJECT_ROOT
# echo $YAKATEKA
# echo $LIBRARY4TESTS
# echo $LOCAL_TMP_DIR

# exit

set -e  # Exit on error

# Check if yakateka is built
if [ ! -f $YAKATEKA ]; then
    echo "${RED}Error: yakateka not built. Run 'make build' first.${NC}"
    exit 1
fi



# Use test document from library4tests if available
if [ -f "$LIBRARY4TESTS/NoSQL_Distilled.pdf" ]; then
    INPUT="$LIBRARY4TESTS/NoSQL_Distilled.pdf"
else
    echo "${RED}Error: Test PDF not found at $LIBRARY4TESTS/NoSQL_Distilled.pdf${NC}"
    echo "Please provide a PDF file to test with."
    exit 1
fi

OUTPUT="$LOCAL_TMP_DIR/output.txt"

# Create tmp directory
mkdir -p $LOCAL_TMP_DIR

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
    echo -e "${GREEN}✓ Output file created: $OUTPUT${NC}"
    echo "  Size: $file_size bytes ($(numfmt --to=iec-i --suffix=B $file_size 2>/dev/null || echo "$file_size bytes"))"
    echo ""
    echo "First 20 lines of output:"
    head -20 "$OUTPUT"
else
    echo -e "${RED}✗ Output file was not created!${NC}"
    exit 1
fi
