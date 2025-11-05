# Converter Configuration Guide

YakaTeka uses a declarative configuration system that allows adding new document converters without code changes.

## Overview

The converter system consists of:
- **Profiles**: Reusable command templates for common converter patterns
- **Converters**: Tool definitions that use profiles and specify supported formats
- **Generic Converter**: Executes commands based on configuration

## Configuration Files

### `config/config.yaml`
Main application configuration (OCR, metadata, logging, etc.)

### `config/converters.yaml` ⭐ **NEW!**
Converter-specific configuration with profiles and tools

## Profile Types

### 1. `simple_io`
For converters that use: `command input output`

**Example**: `ps2pdf input.ps output.pdf`

```yaml
converter_profiles:
  simple_io:
    command_template: "{binary} {input} {output} {extra_args}"
```

### 2. `pandoc_style`
For converters using format flags: `-f format -t format`

**Example**: `pandoc -f markdown -t html input -o output`

```yaml
converter_profiles:
  pandoc_style:
    command_template: "{binary} -f {input_format} -t {output_format} {input} -o {output} {extra_args}"
```

### 3. `libreoffice_style`
For `--convert-to` style converters

**Example**: `soffice --headless --convert-to pdf input`

```yaml
converter_profiles:
  libreoffice_style:
    command_template: "{binary} --headless --convert-to {format} --outdir {outdir} {input}"
    post_process: rename_from_basename
```

## Template Placeholders

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `{binary}` | Path to converter binary | `/usr/bin/pandoc` |
| `{input}` | Absolute input file path | `/home/user/doc.md` |
| `{output}` | Absolute output file path | `/home/user/doc.pdf` |
| `{outdir}` | Output directory | `/home/user` |
| `{input_format}` | Input format (mapped) | `markdown` |
| `{output_format}` | Output format (mapped) | `html` |
| `{format}` | Alias for `{output_format}` | `pdf:writer_pdf_Export` |
| `{extra_args}` | Additional arguments | `--toc --number-sections` |

## Adding a New Converter

### Example 1: Simple Tool (pdftotext)

```yaml
converters:
  pdftotext:
    binary: /usr/bin/pdftotext
    command_template: "{binary} -layout {input} {output}"
    timeout: 300

    formats:
      input: [pdf]
      output: [txt]
```

**Result**: `pdftotext -layout input.pdf output.txt`

### Example 2: Using a Profile (antiword)

```yaml
converters:
  antiword:
    binary: /usr/bin/antiword
    profile: simple_io
    timeout: 60

    formats:
      input: [doc]
      output: [txt]
```

**Result**: `antiword input.doc output.txt`

### Example 3: Format Mapping (catdoc)

```yaml
converters:
  catdoc:
    binary: /usr/bin/catdoc
    command_template: "{binary} {input} > {output}"
    timeout: 60

    formats:
      input: [doc, rtf]
      output: [txt]
```

## Conversion Overrides

Override specific conversions with custom settings:

### Exact Match
```yaml
conversion_overrides:
  "md->pdf":
    extra_args: "--pdf-engine=xelatex"
    quality:
      high: "--toc --number-sections"
```

### Wildcards

```yaml
conversion_overrides:
  "*->pdf":            # Any format to PDF
    extra_args: "--pdf-engine=xelatex"

  "md->*":             # Markdown to any format
    extra_args: "--standalone"

  "*->*":              # All conversions
    quality:
      high: "--pretty-print"
```

### Priority Order
1. Exact match: `md->pdf`
2. Output wildcard: `*->pdf`
3. Input wildcard: `md->*`
4. Full wildcard: `*->*`

## Format Mapping

Map internal format names to tool-specific names:

```yaml
converters:
  pandoc:
    format_mapping:
      md: markdown      # Use "markdown" instead of "md"
      txt: plain        # Use "plain" instead of "txt"
      latex: latex      # Keep as-is
```

## Quality Flags

Define quality-specific arguments:

```yaml
conversion_overrides:
  "*->*":
    quality:
      low: "--fast"
      medium: ""
      high: "--pretty-print --toc"
```

**Usage**:
```bash
yakateka convert input.md output.pdf --quality high
```

**Result**: Adds `--pretty-print --toc` to command

## Post-Processing

Handle tools that create output with unexpected names:

```yaml
converter_profiles:
  libreoffice_style:
    post_process: rename_from_basename
```

**Behavior**: If LibreOffice creates `document.html` but you specified `/tmp/output.html`, it automatically renames the file.

## Complete Example

```yaml
# Define reusable profile
converter_profiles:
  my_converter_style:
    command_template: "{binary} --input={input} --output={output} {extra_args}"

# Define converter
converters:
  my_converter:
    binary: /usr/local/bin/myconverter
    profile: my_converter_style
    timeout: 600

    formats:
      input: [custom1, custom2]
      output: [custom3, custom4]

    format_mapping:
      custom1: c1_format
      custom2: c2_format

    conversion_overrides:
      "custom1->custom3":
        extra_args: "--mode=fast"
        quality:
          high: "--quality=100"
      "*->custom4":
        extra_args: "--preserve-metadata"
```

## Validation

The system validates:
- ✅ All profiles have `command_template`
- ✅ All converters have `binary` path
- ✅ All converters have `profile` OR `command_template`
- ✅ Referenced profiles exist
- ✅ At least one input and output format

## Debugging

Enable debug logging to see generated commands:

```bash
yakateka convert input.pdf output.txt --log-level debug
```

**Output**:
```json
{
  "converter": "pdftotext",
  "command": "/usr/bin/pdftotext -layout /path/input.pdf /path/output.txt",
  "message": "Converting document with generic converter"
}
```

## Benefits

1. **No Code Changes**: Add converters by editing YAML
2. **DRY**: Profiles eliminate repetition
3. **Flexible**: Override any conversion
4. **Transparent**: See exact commands in logs
5. **Validated**: Configuration errors caught at startup

## Migration from Hardcoded Converters

Old (hardcoded):
```go
pandocConverter := pandoc.NewConverter(path, args)
factory.Register("pandoc", pandocConverter)
```

New (config-based):
```yaml
converters:
  pandoc:
    binary: /usr/bin/pandoc
    profile: pandoc_style
    formats:
      input: [md, html, docx]
      output: [html, pdf, docx]
```

The factory automatically loads all converters from config!

## See Also

- `config/converters.yaml` - Full converter configuration
- `internal/converter/config/types.go` - Configuration types
- `internal/converter/generic/converter.go` - Generic converter implementation
