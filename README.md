# Holyrics Slide Converter

**PT** | [EN](README-en.md)

[![CI](https://github.com/felpsalvs/holyrics-slide-converter/actions/workflows/ci.yml/badge.svg)](https://github.com/felpsalvs/holyrics-slide-converter/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/felpsalvs/holyrics-slide-converter?label=download)](https://github.com/felpsalvs/holyrics-slide-converter/releases/latest)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Aplicativo companheiro para o [Holyrics](https://holyrics.com.br): converte automaticamente apresentações (**PPTX / PPT / ODP**) e **PDFs** em imagens PNG de alta resolução, prontas para a biblioteca de imagens do Holyrics — **sem precisar do Microsoft Office pago**.

**Como funciona:** você joga o arquivo numa pasta, o app converte cada slide/página em uma imagem `slide-001.png`, `slide-002.png`, ... dentro de uma pasta com o nome do arquivo, direto na biblioteca do Holyrics. É só apresentar.

Também resolve o problema de PDF que fica pequeno/ilegível na tela: as páginas são renderizadas como imagem em alta resolução (1920px de largura por padrão, configurável).

## 📥 Baixar

Baixe a versão mais recente para o seu sistema na página de **[Releases](https://github.com/felpsalvs/holyrics-slide-converter/releases/latest)**:

| Sistema | Arquivo |
| --- | --- |
| Windows | `holyrics-converter_..._windows_intel.zip` |
| Mac com chip Apple (M1/M2/M3/M4) | `holyrics-converter_..._mac_apple-silicon.zip` |
| Mac com chip Intel | `holyrics-converter_..._mac_intel.zip` |
| Linux | `holyrics-converter_..._linux_intel.zip` ou `_linux_arm64.zip` |

Extraia o zip e siga o passo a passo do arquivo `COMO-USAR.txt` incluído (ou a seção [Uso](#uso) abaixo).

### Requisitos

- **[LibreOffice](https://pt-br.libreoffice.org/)** (gratuito) — necessário apenas para PPTX/PPT/ODP; a conversão de PDF funciona sem ele.

### ⚠️ Aviso de segurança ao abrir

Por ser um programa gratuito sem certificado digital pago, o sistema pode mostrar um aviso na primeira execução. É esperado:

- **Windows** ("O Windows protegeu o computador"): clique em **Mais informações** → **Executar assim mesmo**.
- **Mac** ("desenvolvedor não identificado"): vá em **Ajustes do Sistema → Privacidade e Segurança** e clique em **Abrir Assim Mesmo**.

## Uso

1. Rode o programa uma vez. Ele cria o arquivo de configuração em `~/HolyricsConverter/config.json`:

   ```json
   {
     "pasta_entrada": "C:\\Users\\VOCE\\HolyricsConverter\\entrada",
     "pasta_saida": "C:\\Users\\VOCE\\HolyricsConverter\\convertidos",
     "largura_px": 1920,
     "caminho_soffice": "",
     "conversoes_simultaneas": 1
   }
   ```

2. Ajuste **`pasta_saida`** para dentro da biblioteca de imagens do Holyrics, por exemplo:

   ```
   C:\Users\VOCE\Documents\Holyrics\Holyrics\files\media\images\Convertidos
   ```

   Assim as imagens aparecem direto na aba de imagens do Holyrics.

3. Rode o programa de novo e deixe aberto. Qualquer `.pptx`, `.ppt`, `.odp` ou `.pdf` copiado para a **pasta de entrada** é convertido na hora:

   - Sucesso → imagens em `pasta_saida/<nome-do-arquivo>/` e o original vai para `entrada/processados/`
   - Falha → o original vai para `entrada/erros/` (detalhe no log do terminal)

### Converter um único arquivo (sem monitorar)

```sh
./holyrics-converter --once "culto-domingo.pptx"
```

### Opções do config.json

| Campo | Descrição |
| --- | --- |
| `pasta_entrada` | Pasta monitorada onde você joga os arquivos |
| `pasta_saida` | Onde as imagens são geradas (aponte para a biblioteca do Holyrics) |
| `largura_px` | Largura das imagens geradas (padrão 1920) |
| `caminho_soffice` | Caminho do LibreOffice, se não for detectado automaticamente |
| `conversoes_simultaneas` | Quantas conversões rodam em paralelo (padrão 1) |

### Opções de linha de comando

| Flag | Descrição |
| --- | --- |
| `--once <arquivo>` | Converte um único arquivo e sai |
| `--config <caminho>` | Usa um config.json alternativo |
| `--log-level <nível>` | `debug`, `info` (padrão), `warn` ou `error` |
| `--version` | Mostra a versão instalada |

## Desenvolvimento

Com [Go](https://go.dev) 1.24+ instalado:

```sh
go build -o holyrics-converter ./cmd/holyrics-converter
go test ./...
```

O teste completo de conversão de apresentação (ODP → PDF → PNG) roda apenas se o LibreOffice estiver instalado (senão é pulado). Veja [CONTRIBUTING.md](CONTRIBUTING.md) para contribuir.

As releases são geradas automaticamente pelo GitHub Actions quando uma tag `v*` é criada — cada binário é compilado nativamente no seu sistema alvo (a lib de PDF exige CGO).

## Roadmap

- Integração com a [API Server do Holyrics](https://github.com/holyrics/API-Server) (`AddToPlaylist` / `ShowImage`) para já colocar o material convertido na playlist automaticamente
- Interface gráfica simples (arrastar e soltar)
- Iniciar junto com o sistema

## Licença

[MIT](LICENSE)
