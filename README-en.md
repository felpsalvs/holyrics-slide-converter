# Holyrics Slide Converter

[PT](README.md) | **EN**

[![CI](https://github.com/felpsalvs/holyrics-slide-converter/actions/workflows/ci.yml/badge.svg)](https://github.com/felpsalvs/holyrics-slide-converter/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/felpsalvs/holyrics-slide-converter?label=download)](https://github.com/felpsalvs/holyrics-slide-converter/releases/latest)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Companion app for [Holyrics](https://holyrics.com.br): automatically converts presentations (**PPTX / PPT / ODP**) and **PDFs** into high-resolution PNG images, ready for the Holyrics image library — **no paid Microsoft Office required**.

**How it works:** drop a file into a folder and the app converts each slide/page into `slide-001.png`, `slide-002.png`, ... inside a folder named after the file, straight into your Holyrics library. Just present.

It also fixes the problem of PDFs looking tiny/unreadable on screen: pages are rendered as high-resolution images (1920px wide by default, configurable).

## 📥 Download

Get the latest version for your system from the **[Releases](https://github.com/felpsalvs/holyrics-slide-converter/releases/latest)** page:

| System | File |
| --- | --- |
| Windows | `holyrics-converter_..._windows_intel.zip` |
| Mac with Apple chip (M1/M2/M3/M4) | `holyrics-converter_..._mac_apple-silicon.zip` |
| Mac with Intel chip | `holyrics-converter_..._mac_intel.zip` |
| Linux | `holyrics-converter_..._linux_intel.zip` or `_linux_arm64.zip` |

Extract the zip and follow the included `COMO-USAR.txt` (Portuguese) or the [Usage](#usage) section below.

### Requirements

- **[LibreOffice](https://www.libreoffice.org/)** (free) — needed only for PPTX/PPT/ODP; PDF conversion works without it.

### ⚠️ Security warning on first launch

Since this is a free program without a paid code-signing certificate, your system may warn you on first run. This is expected:

- **Windows** ("Windows protected your PC"): click **More info** → **Run anyway**.
- **Mac** ("unidentified developer"): go to **System Settings → Privacy & Security** and click **Open Anyway**.

## Usage

1. Run the program once. It creates the config file at `~/HolyricsConverter/config.json`:

   ```json
   {
     "pasta_entrada": "C:\\Users\\YOU\\HolyricsConverter\\entrada",
     "pasta_saida": "C:\\Users\\YOU\\HolyricsConverter\\convertidos",
     "largura_px": 1920,
     "caminho_soffice": "",
     "conversoes_simultaneas": 1
   }
   ```

2. Point **`pasta_saida`** (output folder) to somewhere inside the Holyrics image library, e.g.:

   ```
   C:\Users\YOU\Documents\Holyrics\Holyrics\files\media\images\Converted
   ```

   The images then show up directly in the Holyrics images tab.

3. Run the program again and leave it open. Any `.pptx`, `.ppt`, `.odp` or `.pdf` copied into the **input folder** is converted on the spot:

   - Success → images in `pasta_saida/<file-name>/`, original moved to `entrada/processados/`
   - Failure → original moved to `entrada/erros/` (details in the terminal log)

### Convert a single file (no watching)

```sh
./holyrics-converter --once "sunday-service.pptx"
```

### config.json options

| Field | Description |
| --- | --- |
| `pasta_entrada` | Watched folder where you drop files |
| `pasta_saida` | Where images are generated (point it to the Holyrics library) |
| `largura_px` | Width of generated images (default 1920) |
| `caminho_soffice` | LibreOffice path, if not auto-detected |
| `conversoes_simultaneas` | How many conversions run in parallel (default 1) |

### Command-line flags

| Flag | Description |
| --- | --- |
| `--once <file>` | Convert a single file and exit |
| `--config <path>` | Use an alternative config.json |
| `--log-level <level>` | `debug`, `info` (default), `warn` or `error` |
| `--version` | Print the installed version |

## Development

With [Go](https://go.dev) 1.24+ installed:

```sh
go build -o holyrics-converter ./cmd/holyrics-converter
go test ./...
```

The full presentation-conversion test (ODP → PDF → PNG) only runs when LibreOffice is installed (otherwise it is skipped). See [CONTRIBUTING.md](CONTRIBUTING.md) to contribute.

Releases are built automatically by GitHub Actions when a `v*` tag is pushed — each binary is compiled natively on its target OS (the PDF library requires CGO).

## Roadmap

- Integration with the [Holyrics API Server](https://github.com/holyrics/API-Server) (`AddToPlaylist` / `ShowImage`) to automatically add converted material to the playlist
- Simple graphical interface (drag and drop)
- Start with the system

## License

[MIT](LICENSE)
