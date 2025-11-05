# Helper System Guide

The helper system allows external scripts (written in any language) to handle document conversions without modifying yakateka's code.

## Overview

**Benefits:**
- ✅ Add new converters without code changes
- ✅ Use any language (Bash, Python, Go, etc.)
- ✅ Implement complex pipelines inside helpers
- ✅ Automatic fallback on failure
- ✅ Weight-based prioritization

**Architecture:**
```
┌─────────────┐
│   yakateka  │
│   convert   │
└──────┬──────┘
       │
       ├─> Load helpers.yaml (cache)
       ├─> Ping all helpers (startup check)
       ├─> Find helpers for conversion
       └─> Try helpers by weight until success
```

## Helper Interface

Every helper must implement three commands:

### 1. `helper.sh ping`
**Purpose**: Health check
**Returns**:
- Exit code 0 + stdout "pong" = Available
- Anything else = Unavailable

**Example:**
```bash
#!/bin/bash
case "$1" in
    ping)
        echo "pong"
        exit 0
        ;;
esac
```

### 2. `helper.sh info`
**Purpose**: Report capabilities
**Returns**: YAML with supported conversions and modes

**Exit codes:**
- `> 0`: Error (helper ignored)
- `== 0`:
  - Empty stdout: Can't work now (skip with warning)
  - YAML stdout: Parse capabilities
  - Other: Error (skip)

**YAML Structure:**
```yaml
name: "Helper Name"
version: "1.0.0"          # Optional
description: "Description" # Optional
capabilities:
  <from_format>:
    <to_format>:
      modes:
        normal:           # Mandatory
          speed: 1        # > 0 = supported (higher = faster)
          quality: 1      # > 0 = supported (higher = better)
        fast:             # Optional
          speed: 1
          quality: 1
        quality:          # Optional
          speed: 1
          quality: 1
```

**Example:**
```yaml
name: "Pandoc Universal Converter"
capabilities:
  md:
    html:
      modes:
        normal:
          speed: 1
          quality: 1
        fast:
          speed: 1
          quality: 1
    pdf:
      modes:
        normal:
          speed: 1
          quality: 1
```

### 3. `helper.sh convert <mode> <from_format> <from_file> <to_format> <to_file>`
**Purpose**: Perform conversion

**Parameters:**
- `<mode>`: `normal`, `fast`, or `quality`
- `<from_format>`: Input format (e.g., `md`, `pdf`)
- `<from_file>`: Absolute path to input file
- `<to_format>`: Output format (e.g., `html`, `pdf`)
- `<to_file>`: Absolute path to output file

**Returns:**
- Exit code 0 = Success
- Exit code > 0 = Failure (try next helper)

**Example:**
```bash
convert)
    MODE="$2"
    FROM_FORMAT="$3"
    FROM_FILE="$4"
    TO_FORMAT="$5"
    TO_FILE="$6"

    pandoc -f "$FROM_FORMAT" -t "$TO_FORMAT" \
           "$FROM_FILE" -o "$TO_FILE"
    exit $?
    ;;
```

## Configuration

### config.yaml
```yaml
helpers:
  cache_file: helpers.yaml  # Generated cache file
  weights:
    /usr/local/bin/pandoc-helper.sh: 0.9    # Higher = preferred
    /usr/local/bin/calibre-helper.sh: 0.8
    /home/user/custom-helper.sh: 0.7
```

### Weight Sorting
Helpers are tried in order:
1. **By weight** (descending): 0.9 → 0.8 → 0.7
2. **By name** (alphabetical): If weights equal

## Usage Workflow

### 1. Create Helper Script

See `examples/helpers/pandoc-helper.sh` for complete example.

```bash
#!/bin/bash
case "$1" in
    ping) echo "pong"; exit 0 ;;
    info) cat info.yaml; exit 0 ;;
    convert)
        # Your conversion logic
        ;;
esac
```

### 2. Add to Configuration

```yaml
helpers:
  weights:
    /path/to/my-helper.sh: 0.9
```

### 3. Generate Cache

```bash
yakateka helpers
```

Output:
```
✓ Generated helper cache: helpers.yaml
  42 conversion paths available
```

### 4. Use for Conversions

```bash
yakateka convert input.md output.html
```

**Yakateka automatically:**
1. Loads `helpers.yaml`
2. Pings all helpers (filters unavailable)
3. Finds helpers for `md → html`
4. Tries by weight until success

## Conversion Modes

### Normal Mode (Mandatory)
Balanced speed and quality. Every helper must support this.

```bash
yakateka convert input.md output.pdf
# Uses 'normal' mode by default
```

### Fast Mode (Optional)
Speed over quality.

```bash
yakateka convert input.md output.pdf --quality fast
```

**Fallback**: If helper doesn't support `fast`, uses `normal`.

### Quality Mode (Optional)
Quality over speed.

```bash
yakateka convert input.md output.pdf --quality high
```

**Fallback**: If helper doesn't support `quality`, uses `normal`.

## Failure Handling

### Ping Failure (Startup)
If helper fails ping check:
- ✅ Logged as warning
- ✅ Removed from ALL conversions
- ✅ Other helpers still available

### Conversion Failure (Runtime)
If helper fails conversion:
- ✅ Logged as warning
- ✅ Marked failed for THAT conversion pair only
- ✅ Next helper tried automatically
- ❌ If ALL helpers fail: conversion fails

**Example:**
```
Helper 1 (weight 0.9): FAILED → try Helper 2
Helper 2 (weight 0.8): FAILED → try Helper 3
Helper 3 (weight 0.7): SUCCESS ✓
```

## Metrics System

Currently metrics use placeholder value `1`:

```yaml
modes:
  normal:
    speed: 1     # 1 = supported
    quality: 1   # 1 = supported
  fast:
    speed: 0     # 0 = not supported
    quality: 0
```

**Future**: Can be replaced with actual benchmarks (e.g., `speed: 8.5`, `quality: 9.2`).

## Examples

### Example 1: Simple Helper

```bash
#!/bin/bash
# PDF to text converter using pdftotext

case "$1" in
    ping)
        command -v pdftotext >/dev/null 2>&1
        if [ $? -eq 0 ]; then
            echo "pong"
            exit 0
        fi
        exit 1
        ;;

    info)
        cat <<EOF
name: "PDF to Text Converter"
capabilities:
  pdf:
    txt:
      modes:
        normal:
          speed: 1
          quality: 1
EOF
        exit 0
        ;;

    convert)
        pdftotext -layout "$4" "$6"
        exit $?
        ;;
esac
```

### Example 2: Python Helper with Pipelines

```python
#!/usr/bin/env python3
import sys
import subprocess

def ping():
    print("pong")
    sys.exit(0)

def info():
    print("""
name: "Advanced Converter"
capabilities:
  djvu:
    md:
      modes:
        normal:
          speed: 1
          quality: 1
""")
    sys.exit(0)

def convert(mode, from_fmt, from_file, to_fmt, to_file):
    # Pipeline: DJVU → PS → PDF → HTML → MD
    # All internal to this helper!
    subprocess.run(["djvups", from_file, "/tmp/temp.ps"])
    subprocess.run(["ps2pdf", "/tmp/temp.ps", "/tmp/temp.pdf"])
    subprocess.run(["libreoffice", "--convert-to", "html", "/tmp/temp.pdf"])
    subprocess.run(["pandoc", "-f", "html", "-t", "markdown",
                    "/tmp/temp.html", "-o", to_file])
    sys.exit(0)

if __name__ == "__main__":
    cmd = sys.argv[1]
    if cmd == "ping":
        ping()
    elif cmd == "info":
        info()
    elif cmd == "convert":
        convert(*sys.argv[2:7])
```

## helpers.yaml Cache Format

```yaml
conversions:
  md:
    html:
      normal:
        - helper: /path/helper1.sh
          weight: 0.9
        - helper: /path/helper2.sh
          weight: 0.8
      fast:
        - helper: /path/helper1.sh
          weight: 0.9
      quality:
        - helper: /path/helper3.sh
          weight: 0.7
```

## Debugging

```bash
# Enable debug logging
yakateka convert input.md output.html --log-level debug
```

**Logs show:**
- Helper ping results
- Helpers found for conversion
- Attempt order and results
- Failure reasons

## Best Practices

1. **✅ Implement ping properly**: Check tool availability
2. **✅ Use absolute paths**: In config and helpers
3. **✅ Handle errors**: Exit with non-zero on failure
4. **✅ Test all modes**: If you declare them
5. **✅ Use pipelines**: Implement complex conversions inside helper
6. **✅ Make executable**: `chmod +x helper.sh`
7. **❌ Don't print to stdout**: Except for ping/info commands
8. **❌ Don't modify input**: Conversions should be side-effect free

## Troubleshooting

### Helper not working

```bash
# Test ping
/path/to/helper.sh ping
# Should output: pong

# Test info
/path/to/helper.sh info
# Should output valid YAML

# Test conversion
/path/to/helper.sh convert normal md /tmp/test.md html /tmp/test.html
echo $?  # Should be 0
```

### Cache not loading

```bash
# Check cache file exists
ls -l helpers.yaml

# Regenerate cache
yakateka helpers

# Check config
cat config.yaml | grep helpers
```

### All helpers failing

```bash
# Check logs
yakateka convert input.md output.html --log-level debug 2>&1 | grep helper
```

## See Also

- `examples/helpers/pandoc-helper.sh` - Complete working example
- `cmd/helpers.go` - Cache generation implementation
- `internal/helper/` - Helper system code
