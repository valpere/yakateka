#!/bin/bash
# Example 4: Using Configuration Files and Environment Variables
# Demonstrates different ways to configure YakaTeka

YAKATEKA="../bin/yakateka"

if [ ! -f "$YAKATEKA" ]; then
    echo "Error: yakateka not built. Run 'make build' first."
    exit 1
fi

mkdir -p ../tmp

echo "=== Configuration Examples ==="
echo ""

# Example 1: Default configuration
echo "Example 1: Using Default Configuration"
echo "---------------------------------------"
if [ -f "../../library4tests/NoSQL_Distilled.pdf" ]; then
    $YAKATEKA convert "../../library4tests/NoSQL_Distilled.pdf" "../tmp/config_default.txt" --log-format text
    echo ""
fi

# Example 2: Using environment variables
echo "Example 2: Using Environment Variables"
echo "---------------------------------------"
export YAKATEKA_LOG_LEVEL=debug
export YAKATEKA_CONVERTER_PDF_ENGINE=pdfcpu

echo "Setting environment variables:"
echo "  YAKATEKA_LOG_LEVEL=debug"
echo "  YAKATEKA_CONVERTER_PDF_ENGINE=pdfcpu"
echo ""

if [ -f "../../library4tests/NoSQL_Distilled.pdf" ]; then
    $YAKATEKA convert "../../library4tests/NoSQL_Distilled.pdf" "../tmp/config_env.txt" --log-format text 2>&1 | head -20
    echo "... (output truncated)"
    echo ""
fi

# Clean up environment
unset YAKATEKA_LOG_LEVEL
unset YAKATEKA_CONVERTER_PDF_ENGINE

# Example 3: Using custom config file
echo "Example 3: Using Custom Config File"
echo "------------------------------------"

# Create custom config
cat > ../tmp/custom-config.yaml <<EOF
# Custom YakaTeka Configuration

converter:
  pdf:
    engine: pdfcpu
    quality: high

logging:
  level: info
  format: text

output:
  format: json
  pretty: true
EOF

echo "Created custom config file: ../tmp/custom-config.yaml"
cat ../tmp/custom-config.yaml
echo ""

if [ -f "../../library4tests/NoSQL_Distilled.pdf" ]; then
    $YAKATEKA --config ../tmp/custom-config.yaml convert \
        "../../library4tests/NoSQL_Distilled.pdf" \
        "../tmp/config_custom.txt"
    echo ""
fi

# Example 4: Command-line flags override everything
echo "Example 4: Command-line Flags Override Config"
echo "----------------------------------------------"
echo "Even with config file, CLI flags take precedence:"
echo ""

if [ -f "../../library4tests/NoSQL_Distilled.pdf" ]; then
    $YAKATEKA --config ../tmp/custom-config.yaml \
        --log-level debug \
        --log-format json \
        convert "../../library4tests/NoSQL_Distilled.pdf" \
        "../tmp/config_override.txt" 2>&1 | head -5
    echo "... (output truncated)"
    echo ""
fi

echo "=== Configuration Priority ==="
echo ""
echo "Highest → Lowest:"
echo "  1. Command-line flags (--log-level, --log-format, etc.)"
echo "  2. Environment variables (YAKATEKA_*)"
echo "  3. Config file (--config or ~/.yakateka/config.yaml)"
echo "  4. Default values"
echo ""

echo "=== Configuration Examples Complete ==="
