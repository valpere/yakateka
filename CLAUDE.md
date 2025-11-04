# CLAUDE.md - YakaTeka (ЯкаТека)

This file provides guidance to Claude Code when working with the YakaTeka project.

## Project Overview

**YakaTeka** is a document processing agent that handles document conversion, parsing, and annotation. It operates as a standalone CLI tool focused on individual document operations.

**Repository**: https://github.com/valpere/yakateka
**Project Path**: `~/wrk/projects/library/yakateka`
**Shared Context**: `~/wrk/projects/library/context-common` (integration guides)
**Project Context**: `~/wrk/projects/library/context-yakateka` (YakaTeka-specific)

## Context Management

**IMPORTANT**: This project uses **external context directories** that are NOT part of the repository:

### Shared Context (`~/wrk/projects/library/context-common/`)
Integration guides for OCR, document converters, image/audio/video processing, AI model runners.

**READ ONLY MODE**: Context files are reference documentation only. Never modify them.

**Key reference files**:
- `00-context.md` - Context index
- `08-golang_all_integration_cheat_sheet.md` - Quick reference for all integrations
- `13-golang_all_integration-06.md` - **Comprehensive integration guide** (primary reference)

### YakaTeka-Specific Context (`~/wrk/projects/library/context-yakateka/`)
Project-specific requirements, decisions, and implementation details.

**READ ONLY MODE**: Context files are reference documentation only. Never modify them.

**When making architectural decisions, ALWAYS consult context files first.**

## Project Scope

### What YakaTeka Does
- **Document conversion**: PDF, EPUB, FB2, DJVU, MOBI, DOC/DOCX ↔ various formats
- **Document parsing**: Extract text, structure, metadata
- **Document annotation**: Add/extract metadata (author, title, language, category, tags, checksums)
- **Text extraction**: OCR for scanned documents
- **Content analysis**: Extract tables, formulas, code blocks

### What YakaTeka Does NOT Do
- Library management (that's TakaTeka's job)
- Multi-document operations
- File organization or moving
- Duplicate detection
- Database management

### Design Philosophy
- **Single responsibility**: One document, one operation
- **CLI-first**: Command-line interface for all operations
- **Stateless**: No database, works on files directly
- **Local processing**: All operations run locally, no cloud dependencies
- **Native Go preferred**: Pure Go libraries when possible, CGo when needed
- **No external services**: Complete independence

## Technology Stack

### Core Technologies
- **Language**: Go 1.24+
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra) for command structure
- **Configuration**: [Viper](https://github.com/spf13/viper) for hierarchical config
- **Messaging**: Protocol Buffers (future gRPC integration with TakaTeka)

### Document Processing Libraries

Based on context research (`context-common/13-golang_all_integration-06.md`):

#### OCR Engines
**Primary**: **Tesseract OCR** via `github.com/otiai10/gosseract`
- Most popular and reliable OCR
- Multi-language support (uk, ru, en, and 30+ more)
- Pure Go option via WASM/Wazero (CGo-free)
- High accuracy with proper preprocessing

**Alternative**: PaddleOCR via REST API (if higher accuracy needed for Asian languages)

#### PDF Processing
**Primary**: **pdfcpu** - `github.com/pdfcpu/pdfcpu`
- Pure Go, no external dependencies
- PDF manipulation, extraction, conversion

**Commercial (optional)**: **UniPDF** - `github.com/unidoc/unipdf`
- More features, AGPL license
- Powerful PDF toolkit for advanced operations

#### Office Documents
**DOCX/XLSX/PPTX**: **gooxml** - `github.com/baliance/gooxml`
- Pure Go
- Full support for Office Open XML formats

**DOCX reading**: **Docxlib** - `github.com/gonfva/docxlib`
**Word to text**: **Docconv** - `github.com/sajari/docconv`

#### Universal Conversion
**Pandoc** - via `os/exec` CLI wrapper
- Markdown ↔ HTML ↔ DOCX ↔ PDF
- Universal markup converter
- External dependency but very powerful

**LibreOffice** - via `os/exec` CLI wrapper
- Office formats ↔ PDF
- High compatibility
- External dependency

**Gotenberg** - `https://gotenberg.dev/`
- Go-based API for HTML/Office → PDF
- Runs as local service (Docker)

#### Image Processing
**Primary**: **bimg/libvips** - `github.com/h2non/bimg`
- Fast, efficient image operations
- Format conversion, scaling, EXIF

**Alternative**: **ImageMagick** - `gopkg.in/gographics/imagick.v3/imagick`
- Powerful image processing
- PDF → PNG, HEIC → JPG, effects
- CGo dependency

#### Multimedia (Future)
**FFmpeg** - via `os/exec` or `github.com/u2takey/ffmpeg-go`
- Audio/video conversion
- Image format conversion
- Universal multimedia tool

### AI Integration (Future Phases)

Based on context research:

#### Local AI Model Runners
**Ollama** - REST API (`github.com/jmorganca/ollama/api`)
- Local LLM deployment (LLaMA, Mistral, Gemma, Phi)
- Content analysis, classification, summarization

**TensorFlow** - via `tfgo` (`github.com/galeone/tfgo`)
- Custom AI models
- Document classification, entity extraction

**whisper.cpp** - Go bindings
- Audio transcription (speech-to-text)
- For audio content extraction

#### LLM-Based Document Conversion (Advanced)

**Inspired by**: [pdf2md_ollama](https://github.com/gwangjinkim/pdf2md_ollama)

**Approach**: PDF → Images → LLM → Markdown
- More accurate for complex layouts, tables, formulas
- Better preservation of document structure
- Privacy-first (local processing)
- Slower but higher quality

**Pipeline**:
1. **PDF → Images**: Convert PDF pages to images (via pdfcpu or pdf2image)
2. **Image preprocessing**: Pillow/bimg for quality enhancement
3. **LLM inference**: Send images to Ollama with vision models (gemma3:4b, gemma3:12b, llava)
4. **Markdown generation**: LLM outputs structured Markdown

**Go Integration**:
- Ollama REST API client
- Image processing via bimg
- PDF rendering via pdfcpu
- Custom pipeline orchestration

**Use cases**:
- Academic papers with complex formulas
- Scientific documents with tables and charts
- Scanned documents requiring intelligent interpretation
- Documents where layout preservation is critical

## Project Structure

```
yakateka/
├── cmd/
│   ├── root.go              # Root command + global flags
│   ├── convert.go           # Document format conversion
│   ├── parse.go             # Extract structure and metadata
│   ├── extract.go           # Text/data extraction (OCR)
│   ├── annotate.go          # Metadata operations
│   └── analyze.go           # Content analysis (tables, formulas, code)
├── internal/
│   ├── converter/           # Format converters
│   │   ├── pdf/
│   │   │   ├── to_text.go
│   │   │   ├── to_image.go
│   │   │   └── to_docx.go
│   │   ├── docx/
│   │   │   ├── to_text.go
│   │   │   └── to_pdf.go
│   │   ├── epub/
│   │   ├── pandoc/         # Pandoc wrapper
│   │   └── types.go
│   ├── parser/             # Document parsers
│   │   ├── pdf_parser.go
│   │   ├── docx_parser.go
│   │   ├── epub_parser.go
│   │   └── types.go
│   ├── ocr/               # OCR engines
│   │   ├── tesseract.go
│   │   ├── preprocessor.go  # Image preprocessing
│   │   └── types.go
│   ├── extractor/         # Content extraction
│   │   ├── text.go
│   │   ├── tables.go
│   │   ├── formulas.go
│   │   ├── code.go
│   │   └── types.go
│   ├── metadata/          # Metadata handling
│   │   ├── reader.go
│   │   ├── writer.go
│   │   ├── checksum.go
│   │   └── types.go
│   ├── image/             # Image processing
│   │   ├── converter.go
│   │   ├── preprocessor.go
│   │   └── types.go
│   └── types.go           # Common types
├── pkg/                   # Public API (future library use)
│   └── yakateka/
│       ├── client.go
│       └── types.go
├── proto/                 # Protocol Buffers (future gRPC)
│   └── yakateka.proto
├── config/
│   └── config.yaml        # Default configuration
├── testdata/              # Test fixtures
│   ├── pdf/
│   ├── docx/
│   └── images/
├── .claude/
│   └── CLAUDE.md          # Symlink to this file
├── go.mod
├── go.sum
├── main.go
├── Makefile
└── README.md
```

## Design Principles

Follow principles from `~/.claude/CLAUDE.md`:

### Core Principles
1. **DRY**: Reusable converter/parser components
2. **YAGNI**: Implement only needed features
3. **KISS**: Simple APIs, clear interfaces
4. **Encapsulation**: Hide implementation details
5. **PoLA**: Intuitive command structure

### SOLID + GRASP
- **Single Responsibility**: Each converter handles one format
- **Open/Closed**: Easy to add formats without modifying existing code
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Depend on abstractions
- **Information Expert**: Converter knows how to convert its format
- **High Cohesion / Low Coupling**: Minimal dependencies between packages

### Type Organization
- **Common types**: `internal/types.go`
- **Package-specific types**: `internal/<package>/types.go`
- **Public API types**: `pkg/yakateka/`

## CLI Command Structure

```bash
# Conversion
yakateka convert <input> <output> [--from|-f <format>] [--to|-t <format>]
yakateka convert document.pdf document.txt
yakateka convert document.docx --to=pdf
yakateka convert notes.md -f md -t pdf --via=pandoc
yakateka convert scan.pdf output.md --to=md --via=ollama --model=gemma3:4b

# Parsing (extract structure + metadata)
yakateka parse <input> [--format=json|yaml|text]
yakateka parse document.pdf --format=json > metadata.json

# Text extraction
yakateka extract text <input> [--ocr] [--lang=uk,ru,en] [--dpi=300]
yakateka extract text scan.pdf --ocr --lang=uk,en
yakateka extract tables document.pdf --format=csv
yakateka extract formulas paper.pdf --format=latex

# Metadata operations
yakateka meta get <input> [--format=json|yaml]
yakateka meta set <input> --author="Name" --title="Title" --lang=uk
yakateka meta checksum <input> [--algo=sha256|md5]

# Content analysis
yakateka analyze <input> [--tables] [--formulas] [--code]
yakateka analyze document.pdf --tables --formulas --format=json

# Image preprocessing (for OCR optimization)
yakateka preprocess <image> <output> [--denoise] [--deskew] [--threshold]
```

## Development Roadmap

### Phase 1: Foundation ✅ **COMPLETED**
- [x] Initialize project structure
- [x] Setup Cobra CLI framework
- [x] Setup Viper configuration
- [x] Basic logging (structured, using zerolog)
- [x] Makefile for build automation
- [x] Basic tests and CI setup
- [x] Version flag support

### Phase 2: Core Conversion ← **CURRENT PRIORITY**
- [x] PDF → Text (pdfcpu)
- [x] Converter factory pattern
- [x] Convert command with auto-format detection
- [ ] DOCX → Text (gooxml)
- [ ] Text → PDF (via pandoc)
- [ ] Pandoc wrapper (Markdown/HTML conversion)
- [ ] Image format conversion (bimg)

### Phase 3: OCR Integration
- [ ] Tesseract integration (gosseract)
- [ ] Image preprocessing (denoise, deskew, threshold)
- [ ] Multi-language OCR (uk, ru, en)
- [ ] OCR quality metrics
- [ ] Batch OCR processing

### Phase 4: Advanced Parsing
- [ ] PDF structure analysis (pdfcpu)
- [ ] Table detection and extraction
- [ ] Formula extraction (LaTeX output)
- [ ] Code block detection and extraction
- [ ] DOCX structure parsing

### Phase 5: Metadata Handling
- [ ] PDF metadata read/write
- [ ] EPUB metadata
- [ ] DOCX metadata
- [ ] Custom metadata (JSON sidecar files)
- [ ] Checksum generation (SHA256, MD5)

### Phase 6: Advanced Features (AI-powered)
- [ ] Ollama integration for LLM-based conversion
- [ ] PDF → Markdown via LLM (inspired by pdf2md_ollama approach)
- [ ] Document classification with local LLMs
- [ ] Entity extraction
- [ ] Language detection
- [ ] Content summarization

### Phase 7: API Mode
- [ ] REST API server (Gin/Echo)
- [ ] OpenAPI/Swagger documentation
- [ ] Authentication (optional)
- [ ] Rate limiting

### Phase 8: gRPC Integration
- [ ] Protocol Buffers schema
- [ ] gRPC server
- [ ] Integration with TakaTeka

## Configuration

Viper hierarchical configuration (priority order):

1. **Command-line flags** (highest)
2. **Environment variables** (`YAKATEKA_*`)
3. **Config file** (`~/.yakateka/config.yaml` or `./config/config.yaml`)
4. **Defaults** (lowest)

Example `config.yaml`:
```yaml
ocr:
  engine: tesseract
  languages: [uk, ru, en]
  dpi: 300
  preprocess:
    denoise: true
    deskew: true
    threshold: auto

converter:
  pdf:
    engine: pdfcpu  # or unipdf
    quality: high
  pandoc:
    path: /usr/bin/pandoc
    extra_args: []
  image:
    library: bimg  # or imagick
    format: png
    dpi: 300

metadata:
  checksum: sha256
  sidecar: true  # Create .json metadata files
  embed: true    # Embed metadata in documents

output:
  format: json
  pretty: true

logging:
  level: info  # debug, info, warn, error
  format: json  # json or text
```

## Testing Strategy

- **Unit tests**: Each converter/parser independently
- **Integration tests**: Full conversion workflows
- **Test fixtures**: Sample documents in `testdata/`
- **Benchmarks**: Performance tests for large documents
- **Coverage target**: 80%+
- **CI/CD**: GitHub Actions for automated testing

```bash
# Run tests
make test

# Run specific package
go test ./internal/converter/pdf/...

# Benchmarks
make bench

# Coverage
make coverage
```

## Common Tasks

### Adding a New Document Format

1. **Create package structure**:
   ```
   internal/converter/<format>/
     ├── to_text.go
     ├── to_pdf.go
     ├── types.go
     └── <format>_test.go
   ```

2. **Implement `Converter` interface**:
   ```go
   type Converter interface {
       Convert(ctx context.Context, input, output string, opts Options) error
       SupportedInputFormats() []string
       SupportedOutputFormats() []string
   }
   ```

3. **Register in factory** (`internal/converter/factory.go`)

4. **Add CLI command** (`cmd/convert.go`)

5. **Add tests** with fixtures

6. **Update documentation** (README.md, this file)

### Adding OCR Support

1. **Implement `OCREngine` interface** (`internal/ocr/`)
2. **Add language configuration**
3. **Implement preprocessing pipeline**
4. **Integrate with `extract text` command**
5. **Add benchmarks for accuracy/performance**

### Building and Running

```bash
# Build
make build

# Install to $GOPATH/bin
make install

# Run
yakateka --help

# Development mode (with hot reload)
make dev

# Lint
make lint

# Format code
make fmt
```

## Integration with TakaTeka

**Current**: YakaTeka operates independently as CLI.

**Future Integration**:
1. **Phase 1**: TakaTeka calls YakaTeka CLI via `os/exec`
2. **Phase 2**: TakaTeka calls YakaTeka REST API
3. **Phase 3**: TakaTeka calls YakaTeka via gRPC

**Interface Contract** (to be defined in protobuf):
```protobuf
service DocumentProcessor {
  rpc ConvertDocument(ConvertRequest) returns (ConvertResponse);
  rpc ParseDocument(ParseRequest) returns (ParseResponse);
  rpc ExtractText(ExtractRequest) returns (ExtractResponse);
  rpc AnnotateDocument(AnnotateRequest) returns (AnnotateResponse);
  rpc AnalyzeContent(AnalyzeRequest) returns (AnalyzeResponse);
}
```

## Important Constraints

- **No migration from existing library**: Fresh start, clean architecture
- **No docling dependency**: Build our own conversion pipeline
- **Native Go preferred**: Avoid `os/exec` when pure Go solution exists
- **Context is authoritative**: Always check `context-common/` before decisions
- **Separation of concerns**: YakaTeka processes, TakaTeka manages
- **Local-first**: No cloud dependencies, fully offline capable

## Recommended Libraries Summary

Based on `context-common/13-golang_all_integration-06.md`:

| Category | Recommended | Alternative | Notes |
|----------|-------------|-------------|-------|
| OCR | Tesseract (gosseract) | PaddleOCR (REST) | Pure Go via WASM available |
| PDF | pdfcpu | UniPDF (AGPL) | UniPDF for advanced features |
| Office | gooxml | Docconv | Pure Go solution |
| Markup | Pandoc (CLI) | - | Universal converter |
| Images | bimg/libvips | ImageMagick | bimg faster, ImageMagick more features |
| Multimedia | FFmpeg (CLI) | - | Universal tool |
| AI Models | Ollama (REST) | TensorFlow (tfgo) | Ollama simpler for LLMs |

## Code Style

- **Go conventions**: `gofmt`, `golint`, `go vet`
- **Naming**: Clear, descriptive names
- **Documentation**: All exported functions/types
- **Error handling**: Return errors, don't panic
- **Logging**: Structured logging (zerolog)
- **Comments**: Why, not what

## Development Priorities

**Immediate focus (Phase 1)**:
1. Project initialization
2. CLI structure with Cobra
3. Configuration with Viper
4. Basic PDF → Text conversion (pdfcpu)
5. Basic DOCX → Text conversion (gooxml)

**Next steps (Phase 2-3)**:
6. Tesseract OCR integration
7. Image preprocessing
8. Pandoc wrapper for universal conversion

## Questions or Clarifications?

1. **Check context files**: `~/wrk/projects/library/context-common/`
2. **Review this file**: `CLAUDE.md`
3. **Check parent**: `~/.claude/CLAUDE.md`
4. **Ask the user**: When truly unclear

---

**Created**: 2025-11-03
**Last Updated**: 2025-11-03
**Version**: 1.0.0
