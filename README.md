# YakaTeka (ЯкаТека)

**Document processing agent** for conversion, parsing, OCR, and annotation.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](LICENSE)

## Overview

YakaTeka is a standalone CLI tool focused on individual document operations:

- **Document conversion**: PDF, EPUB, FB2, DJVU, MOBI, DOC/DOCX ↔ various formats
- **Document parsing**: Extract text, structure, metadata
- **Document annotation**: Add/extract metadata (author, title, language, tags, checksums)
- **Text extraction**: OCR for scanned documents
- **Content analysis**: Extract tables, formulas, code blocks

## Quick Start

### Build

```bash
make build
```

### Run

```bash
./bin/yakateka --help
```

### Install

```bash
make install
```

## Development Status

**Current Phase**: Phase 2 - Core Conversion 🚧

- [x] Phase 1: Foundation complete
- [x] Converter factory pattern
- [x] Configurable conversion timeout
- [x] **Pandoc wrapper** (universal converter)
  - ✅ FB2, EPUB, MD, HTML, DOCX, ODT, RTF (bidirectional)
  - ✅ All formats → PDF, JSON, CSV, LaTeX, RST
  - ⚠️  **Cannot read PDF** (Pandoc limitation - can only write PDF)
- [x] **DjVu converter** (DjVuLibre)
  - ✅ DJVU → TXT (3.7MB in 487ms)
  - ⚠️  Requires text layer (OCR-processed DJVUs only)
- [ ] **PDF → Text** (Requires alternative approaches):
  - [ ] LibreOffice CLI wrapper
  - [ ] OCR + Ollama (AI-powered for scanned docs)
- [ ] DOCX → Text conversion (gooxml - native Go alternative)

**Completed Phases**:
- ✅ Phase 1: Foundation (project structure, CLI, config, logging)
- ✅ Phase 2: Pandoc converter (EPUB, MD, DOCX, HTML)

**Important Notes**:
- **Native Go PDF libraries** (pdfcpu, unipdf) extract PDF rendering instructions rather than text - not suitable for text extraction
- **Pandoc limitation**: Can convert TO PDF but not FROM PDF - need LibreOffice or OCR for PDF → Text
- **Current capability**: FB2 ↔ EPUB ↔ MD ↔ HTML ↔ DOCX ↔ ODT ↔ RTF ↔ PDF (one-way to PDF)

### Format Support Matrix

**Pandoc Converter** (implemented):

| From ⬇️ / To ➡️ | TXT | MD | HTML | DOCX | ODT | RTF | PDF | EPUB | FB2 | JSON | CSV | LaTeX | RST |
|----------------|-----|-------|------|------|-----|-----|-----|------|-----|------|-----|-------|-----|
| **FB2**        | ✅  | ✅    | ✅   | ✅   | ✅  | ✅  | ✅  | ✅   | -   | ✅   | ✅  | ✅    | ✅  |
| **EPUB**       | ✅  | ✅    | ✅   | ✅   | ✅  | ✅  | ✅  | -    | ✅  | ✅   | ✅  | ✅    | ✅  |
| **Markdown**   | ✅  | -     | ✅   | ✅   | ✅  | ✅  | ✅  | ✅   | ✅  | ✅   | ✅  | ✅    | ✅  |
| **HTML**       | ✅  | ✅    | -    | ✅   | ✅  | ✅  | ✅  | ✅   | ✅  | ✅   | ✅  | ✅    | ✅  |
| **DOCX**       | ✅  | ✅    | ✅   | -    | ✅  | ✅  | ✅  | ✅   | ✅  | ✅   | ✅  | ✅    | ✅  |
| **ODT**        | ✅  | ✅    | ✅   | ✅   | -   | ✅  | ✅  | ✅   | ✅  | ✅   | ✅  | ✅    | ✅  |
| **RTF**        | ✅  | ✅    | ✅   | ✅   | ✅  | -   | ✅  | ✅   | ✅  | ✅   | ✅  | ✅    | ✅  |
| **PDF**        | ❌  | ❌    | ❌   | ❌   | ❌  | ❌  | -   | ❌   | ❌  | ❌   | ❌  | ❌    | ❌  |
| **DOC**        | ❌  | ❌    | ❌   | ❌   | ❌  | ❌  | ❌  | ❌   | ❌  | ❌   | ❌  | ❌    | ❌  |
| **DJVU**       | ✅  | ⚠️    | ⚠️   | ⚠️   | ⚠️  | ⚠️  | ⚠️  | ⚠️   | ⚠️  | ⚠️   | ⚠️  | ⚠️    | ⚠️  |
| **MOBI**       | ❌  | ❌    | ❌   | ❌   | ❌  | ❌  | ❌  | ❌   | ❌  | ❌   | ❌  | ❌    | ❌  |

**Legend**:
- ✅ Supported and tested
- ⚠️  Supported for DJVU with text layer (scanned PDFs without OCR will be empty)
- ❌ Not supported (requires LibreOffice or other converter)
- `-` Same format (no conversion needed)

**DjVu Converter** (DjVuLibre):
- ✅ **DJVU → TXT** (tested with 3.7MB extraction in 487ms)
- ✅ **DJVU → PS** (PostScript conversion using djvups, ~10.5s for 633MB output)
- ⚠️  Only works with DJVU files that have an embedded text layer
- 🔄 **Scanned DJVUs without text layer** → Will use OCR/AI pipeline (Phase 3)

**PostScript Converter** (Ghostscript):
- ✅ **PS → PDF** (conversion using ps2pdf)

**LibreOffice Converter** (✅ **NEW!**):
- ✅ **PDF/PS/DOC/DOCX/ODT/RTF → HTML** (structure-preserving conversion)
- ✅ **PDF/PS/DOC/DOCX/ODT/RTF → PDF** (document conversion)
- ✅ **Enables multi-step pipelines**: DJVU → PS → HTML → MD
- Supports: DOC, DOCX, ODT, RTF, PDF, PS as input
- Exports to: PDF, HTML, DOCX, ODT, RTF
- **Note**: Does NOT export to Markdown directly (uses HTML → MD via Pandoc for structure preservation)

**Conversion Pipeline** (✅ **FULLY IMPLEMENTED!**):
- 🔄 **Automatic multi-step conversion** using BFS algorithm
- Finds shortest path between formats (up to 4 steps)
- Intermediate formats: **PDF → PS → HTML** (in priority order, preserves structure)
- **TXT is NOT used** as intermediate format (loses document structure)
- **Example**: DJVU → PS → HTML → MD (3-step pipeline)
- **Transparent to users**: One command, automatic pipeline execution
- Temp files automatically cleaned up

**Calibre Converter** (✅ **NEW!**):
- ✅ **MOBI/EPUB/FB2 ↔ MOBI/EPUB/FB2** (ebook format conversions)
- ✅ **Tested**: EPUB → MOBI (9.5MB in 2.7s), FB2 → EPUB (586KB in 0.6s)
- Supports input: MOBI, EPUB, FB2, HTML, TXT, PDF, DOCX, ODT, RTF
- Supports output: MOBI, EPUB, FB2, HTML, TXT, PDF, DOCX, ODT, RTF
- Includes AZW, AZW3, LIT, PDB and many more ebook formats
- Quality options: `--quality high` for pretty-print output

**Helper System** (✅ **NEW!**):
- 🚀 **External scripts as converters** - Write helpers in any language
- 📝 **Simple interface**: `ping`, `info`, `convert` commands
- 🔄 **Automatic fallback**: Try helpers by weight until success
- 🎯 **Failure tracking**: Failed helpers skipped for specific conversions
- 🔌 **Pipeline support**: Implement multi-step conversions inside helpers
- 📊 **Format matrix**: View all supported conversions with `yakateka helpers --formats`
- 📚 **See**: `docs/HELPERS.md` for complete guide
- 📖 **Example**: `examples/helpers/pandoc-helper.sh`

**Declarative Converter Configuration** (Legacy):
- 🎯 **Add converters without code changes** - Edit `config/converters.yaml`
- 📋 **Profiles**: Reusable command templates (`simple_io`, `pandoc_style`, `libreoffice_style`)
- 🔧 **Format mapping**: Map internal names to tool-specific names
- ⚙️  **Conversion overrides**: Per-format customization with wildcards
- 📚 **See**: `docs/CONVERTER_CONFIGURATION.md` for complete guide
- ⚠️  **Note**: Helper system is recommended for new converters

**Future Converters**:
- **OCR + AI Pipeline** (Phase 3):
  - For scanned PDFs, DJVUs without text layer
  - **Page-by-page processing**:
    1. Extract images from DJVU/PDF pages
    2. OCR with Tesseract (multi-language: uk, ru, en)
    3. AI enhancement with Ollama (for complex layouts, tables, formulas)
  - Fallback chain: djvutxt → Tesseract OCR → Ollama vision models

## Usage

### Quick Start

```bash
# Show help
yakateka --help

# Show version
yakateka --version

# Generate helpers cache and show format matrix
yakateka helpers --formats

# Document conversion (direct conversion)
yakateka convert document.epub document.txt  # EPUB to text
yakateka convert notes.md document.pdf       # Markdown to PDF
yakateka convert document.docx output.txt    # DOCX to text
yakateka convert page.html page.md           # HTML to Markdown

# Pipeline conversion (automatic multi-step) - requires LibreOffice for structure-preserving paths
# yakateka convert document.djvu output.md    # DJVU → PS → PDF → MD (postponed)
# yakateka convert document.djvu output.html  # DJVU → PS → PDF → HTML (postponed)

# With custom timeout (default 300 seconds = 5 minutes)
yakateka convert large-document.epub output.txt --timeout 600

# Note: PDF to text requires LibreOffice (coming soon)
# yakateka convert document.pdf output.txt --via libreoffice

# Configure logging
yakateka --log-level debug --log-format text convert input.pdf output.txt
yakateka -v convert input.pdf output.txt  # verbose mode

# Use custom config
yakateka --config /path/to/config.yaml convert input.pdf output.txt
```

### Examples

For detailed usage examples, see the [`examples/`](examples/) directory:

- **[Basic Conversion](examples/01_basic_conversion.sh)** - Simple PDF to text conversion
- **[Batch Processing](examples/02_batch_conversion.sh)** - Convert multiple PDFs at once
- **[Error Handling](examples/03_error_handling.sh)** - Robust conversion with error checking
- **[Configuration](examples/04_configuration.sh)** - Using config files and environment variables
- **[Complete Guide](examples/README.md)** - Comprehensive examples and troubleshooting

## Configuration

YakaTeka uses hierarchical configuration (priority order):

1. Command-line flags (highest)
2. Environment variables (`YAKATEKA_*`)
3. Config file (`~/.yakateka/config.yaml` or `./config/config.yaml`)
4. Defaults (lowest)

See [`config/config.yaml`](config/config.yaml) for configuration options.

## Development

### Test Resources

Integration tests use PDF documents from `../library4tests/` directory. Tests will skip gracefully if this directory is not available.

Example test document: `NoSQL_Distilled.pdf`

### Makefile Targets

```bash
make help           # Show all available targets
make build          # Build the application
make test           # Run tests (includes integration tests)
make test-coverage  # Run tests with coverage report
make bench          # Run benchmarks
make fmt            # Format code
make lint           # Run linters
make clean          # Clean build artifacts
make install        # Install to $GOPATH/bin
```

### Project Structure

```
yakateka/
├── cmd/                # Cobra commands
│   ├── root.go        # Root command + global flags
│   └── ...            # Subcommands (convert, parse, extract, etc.)
├── internal/          # Internal packages
│   ├── converter/     # Format converters
│   ├── parser/        # Document parsers
│   ├── ocr/          # OCR engines
│   ├── extractor/    # Content extraction
│   ├── metadata/     # Metadata handling
│   ├── image/        # Image processing
│   └── types.go      # Common types
├── pkg/              # Public API (future)
├── proto/            # Protocol Buffers (future gRPC)
├── config/           # Default configuration
├── testdata/         # Test fixtures
├── main.go           # Entry point
└── Makefile          # Build automation
```

## Architecture

See [`CLAUDE.md`](CLAUDE.md) for comprehensive architectural documentation:

- Design principles (DRY, SOLID, GRASP)
- Technology stack decisions
- Development roadmap
- Integration plans

## Contributing

This project follows strict design principles:

- **DRY**: No code duplication
- **YAGNI**: Implement only when needed
- **KISS**: Simple solutions over complex ones
- **SOLID principles**: Single responsibility, Open/Closed, etc.
- **GRASP patterns**: Proper responsibility assignment

## License

This project is licensed under the AGPL v3 License - see the [LICENSE](LICENSE) file for details.
