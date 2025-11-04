#!/bin/bash
# Example 3: PDF Conversion with Error Handling
# Demonstrates proper error checking and handling

YAKATEKA="../bin/yakateka"

# Function to convert PDF with comprehensive error handling
convert_pdf() {
    local input="$1"
    local output="$2"

    echo "=== Converting PDF with Error Handling ==="
    echo ""

    # Check if yakateka exists
    if [ ! -f "$YAKATEKA" ]; then
        echo "✗ Error: yakateka not built!"
        echo "  Run 'make build' from the project root."
        return 1
    fi

    # Check if input file exists
    if [ ! -f "$input" ]; then
        echo "✗ Error: Input file not found!"
        echo "  Path: $input"
        return 1
    fi

    # Check if input file is readable
    if [ ! -r "$input" ]; then
        echo "✗ Error: Input file is not readable!"
        echo "  Path: $input"
        echo "  Check file permissions."
        return 1
    fi

    # Check input file size
    local input_size=$(stat -f%z "$input" 2>/dev/null || stat -c%s "$input" 2>/dev/null)
    if [ "$input_size" -eq 0 ]; then
        echo "✗ Error: Input file is empty!"
        echo "  Path: $input"
        return 1
    fi

    # Create output directory if needed
    local output_dir=$(dirname "$output")
    if [ ! -d "$output_dir" ]; then
        echo "Creating output directory: $output_dir"
        mkdir -p "$output_dir"
        if [ $? -ne 0 ]; then
            echo "✗ Error: Failed to create output directory!"
            return 1
        fi
    fi

    # Perform conversion
    echo "Input:  $input ($input_size bytes)"
    echo "Output: $output"
    echo ""
    echo "Converting..."

    # Run conversion with timeout
    timeout 120 $YAKATEKA convert "$input" "$output" --log-format text -v

    local status=$?

    echo ""

    # Check conversion status
    if [ $status -eq 124 ]; then
        echo "✗ Error: Conversion timed out after 120 seconds!"
        return 1
    elif [ $status -ne 0 ]; then
        echo "✗ Error: Conversion failed with exit code: $status"
        return $status
    fi

    # Verify output file exists
    if [ ! -f "$output" ]; then
        echo "✗ Error: Output file was not created!"
        return 1
    fi

    # Check output file size
    local output_size=$(stat -f%z "$output" 2>/dev/null || stat -c%s "$output" 2>/dev/null)
    if [ "$output_size" -eq 0 ]; then
        echo "✗ Error: Output file is empty!"
        return 1
    fi

    # Success!
    echo "✓ Conversion successful!"
    echo ""
    echo "Output details:"
    echo "  Path: $output"
    echo "  Size: $output_size bytes"
    echo "  Human-readable: $(numfmt --to=iec-i --suffix=B $output_size 2>/dev/null || echo "$output_size bytes")"

    return 0
}

# Main script
mkdir -p ../tmp

# Example 1: Convert valid PDF
if [ -f "../../library4tests/NoSQL_Distilled.pdf" ]; then
    echo "Example 1: Valid PDF"
    convert_pdf "../../library4tests/NoSQL_Distilled.pdf" "../tmp/valid_output.txt"
    echo ""
    echo "================================"
    echo ""
fi

# Example 2: Try to convert non-existent file
echo "Example 2: Non-existent File"
convert_pdf "nonexistent.pdf" "../tmp/output.txt"
echo ""
echo "================================"
echo ""

# Example 3: Try to convert with invalid output path
echo "Example 3: Invalid Output Path (testing directory creation)"
if [ -f "../../library4tests/NoSQL_Distilled.pdf" ]; then
    convert_pdf "../../library4tests/NoSQL_Distilled.pdf" "../tmp/nested/deep/output.txt"
fi

echo ""
echo "=== Error Handling Examples Complete ==="
