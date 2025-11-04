#!/usr/bin/env bash
#
# Example 5: Pandoc Document Conversion
#
# This example demonstrates various document conversions using Pandoc wrapper.
# Pandoc is a universal document converter supporting many formats.
#
# Supported conversions (via Pandoc):
#   - EPUB → TXT/MD/HTML
#   - MD → PDF/HTML/DOCX/EPUB
#   - DOCX → TXT/MD/HTML
#   - HTML → MD/TXT/DOCX
#
# IMPORTANT LIMITATION:
#   Pandoc cannot read PDF files - only write them.
#   For PDF → Text conversion, use LibreOffice (coming soon).

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect project root (this script is in examples/)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Configuration
YAKATEKA="${PROJECT_ROOT}/bin/yakateka"
OUTPUT_DIR="${PROJECT_ROOT}/../tmp"
TEST_DOCS_DIR="${PROJECT_ROOT}/../library4tests"

# Check if yakateka binary exists
if [[ ! -f "$YAKATEKA" ]]; then
    echo -e "${RED}Error: yakateka binary not found at $YAKATEKA${NC}"
    echo "Please run 'make build' first"
    exit 1
fi

# Check if pandoc is installed
if ! command -v pandoc &> /dev/null; then
    echo -e "${YELLOW}Warning: Pandoc is not installed${NC}"
    echo "Install it with: sudo apt-get install pandoc"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo -e "${BLUE}=== Pandoc Document Conversion Examples ===${NC}\n"

# Example 1: EPUB to Text
echo -e "${YELLOW}Example 1: EPUB → Text${NC}"
if [[ -f "${TEST_DOCS_DIR}/NoSQL_Distilled-2.epub" ]]; then
    "$YAKATEKA" convert \
        "${TEST_DOCS_DIR}/NoSQL_Distilled-2.epub" \
        "${OUTPUT_DIR}/NoSQL_Distilled.txt"

    FILE_SIZE=$(stat -f%z "${OUTPUT_DIR}/NoSQL_Distilled.txt" 2>/dev/null || stat -c%s "${OUTPUT_DIR}/NoSQL_Distilled.txt" 2>/dev/null)
    echo -e "${GREEN}✓ EPUB converted to text: ${FILE_SIZE} bytes${NC}\n"
else
    echo -e "${YELLOW}Skipping: Test EPUB file not found${NC}\n"
fi

# Example 2: Markdown to PDF
echo -e "${YELLOW}Example 2: Markdown → PDF${NC}"
cat > "${OUTPUT_DIR}/sample.md" <<'EOF'
# Sample Document

This is a **sample** document for testing Pandoc conversion.

## Features

- Supports *italic* and **bold** text
- Lists (ordered and unordered)
- Code blocks

## Code Example

```bash
echo "Hello, World!"
```

## Conclusion

Pandoc is a powerful universal document converter.
EOF

"$YAKATEKA" convert \
    "${OUTPUT_DIR}/sample.md" \
    "${OUTPUT_DIR}/sample.pdf"

FILE_SIZE=$(stat -f%z "${OUTPUT_DIR}/sample.pdf" 2>/dev/null || stat -c%s "${OUTPUT_DIR}/sample.pdf" 2>/dev/null)
echo -e "${GREEN}✓ Markdown converted to PDF: ${FILE_SIZE} bytes${NC}\n"

# Example 3: Markdown to DOCX
echo -e "${YELLOW}Example 3: Markdown → DOCX${NC}"
"$YAKATEKA" convert \
    "${OUTPUT_DIR}/sample.md" \
    "${OUTPUT_DIR}/sample.docx"

FILE_SIZE=$(stat -f%z "${OUTPUT_DIR}/sample.docx" 2>/dev/null || stat -c%s "${OUTPUT_DIR}/sample.docx" 2>/dev/null)
echo -e "${GREEN}✓ Markdown converted to DOCX: ${FILE_SIZE} bytes${NC}\n"

# Example 4: DOCX to Text
echo -e "${YELLOW}Example 4: DOCX → Text${NC}"
"$YAKATEKA" convert \
    "${OUTPUT_DIR}/sample.docx" \
    "${OUTPUT_DIR}/sample_from_docx.txt"

FILE_SIZE=$(stat -f%z "${OUTPUT_DIR}/sample_from_docx.txt" 2>/dev/null || stat -c%s "${OUTPUT_DIR}/sample_from_docx.txt" 2>/dev/null)
echo -e "${GREEN}✓ DOCX converted to text: ${FILE_SIZE} bytes${NC}\n"

# Example 5: Markdown to HTML
echo -e "${YELLOW}Example 5: Markdown → HTML${NC}"
"$YAKATEKA" convert \
    "${OUTPUT_DIR}/sample.md" \
    "${OUTPUT_DIR}/sample.html"

FILE_SIZE=$(stat -f%z "${OUTPUT_DIR}/sample.html" 2>/dev/null || stat -c%s "${OUTPUT_DIR}/sample.html" 2>/dev/null)
echo -e "${GREEN}✓ Markdown converted to HTML: ${FILE_SIZE} bytes${NC}\n"

# Example 6: HTML to Markdown
echo -e "${YELLOW}Example 6: HTML → Markdown${NC}"
"$YAKATEKA" convert \
    "${OUTPUT_DIR}/sample.html" \
    "${OUTPUT_DIR}/sample_from_html.md"

FILE_SIZE=$(stat -f%z "${OUTPUT_DIR}/sample_from_html.md" 2>/dev/null || stat -c%s "${OUTPUT_DIR}/sample_from_html.md" 2>/dev/null)
echo -e "${GREEN}✓ HTML converted to Markdown: ${FILE_SIZE} bytes${NC}\n"

# Example 7: EPUB to Markdown
echo -e "${YELLOW}Example 7: EPUB → Markdown${NC}"
if [[ -f "${TEST_DOCS_DIR}/NoSQL_Distilled-2.epub" ]]; then
    "$YAKATEKA" convert \
        "${TEST_DOCS_DIR}/NoSQL_Distilled-2.epub" \
        "${OUTPUT_DIR}/NoSQL_Distilled.md"

    FILE_SIZE=$(stat -f%z "${OUTPUT_DIR}/NoSQL_Distilled.md" 2>/dev/null || stat -c%s "${OUTPUT_DIR}/NoSQL_Distilled.md" 2>/dev/null)
    echo -e "${GREEN}✓ EPUB converted to Markdown: ${FILE_SIZE} bytes${NC}\n"
else
    echo -e "${YELLOW}Skipping: Test EPUB file not found${NC}\n"
fi

# Show Pandoc limitation
echo -e "${RED}=== Pandoc Limitation ===${NC}"
echo -e "${YELLOW}Note: Pandoc cannot convert FROM PDF (only TO PDF)${NC}"
echo "For PDF → Text conversion, use LibreOffice converter (coming soon):"
echo "  yakateka convert document.pdf output.txt --via libreoffice"
echo ""

echo -e "${GREEN}=== All conversions completed ===${NC}"
echo -e "Output directory: ${BLUE}${OUTPUT_DIR}${NC}"
echo ""
echo "Generated files:"
ls -lh "${OUTPUT_DIR}/" | grep -v "^total" | awk '{print "  " $9 " (" $5 ")"}'
