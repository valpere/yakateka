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
- [x] PDF → Text conversion (pdfcpu)
- [ ] DOCX → Text conversion (gooxml)
- [ ] Pandoc wrapper for universal conversion

**Completed Phases**:
- ✅ Phase 1: Foundation (project structure, CLI, config, logging)

## Usage

### Quick Start

```bash
# Show help
yakateka --help

# Show version
yakateka --version

# Document conversion
yakateka convert document.pdf document.txt
yakateka convert document.pdf output.txt --from pdf --to txt

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
