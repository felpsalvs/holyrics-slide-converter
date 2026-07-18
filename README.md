# Holyrics Slide Converter

Aplicativo companheiro para o [Holyrics](https://holyrics.com.br): converte automaticamente apresentações (**PPTX / PPT / ODP**) e **PDFs** em imagens PNG de alta resolução, prontas para a biblioteca de imagens do Holyrics — **sem precisar do Microsoft Office pago**.

**Como funciona:** você joga o arquivo numa pasta, o app converte cada slide/página em uma imagem `slide-001.png`, `slide-002.png`, ... dentro de uma pasta com o nome do arquivo, direto na biblioteca do Holyrics. É só apresentar.

Também resolve o problema de PDF que fica pequeno/ilegível na tela: as páginas são renderizadas como imagem em alta resolução (1920px de largura por padrão, configurável).

## Requisitos

- **Windows ou macOS**
- **[LibreOffice](https://pt-br.libreoffice.org/)** (gratuito) — necessário apenas para PPTX/PPT/ODP; a conversão de PDF funciona sem ele.

## Instalação

Com [Go](https://go.dev) instalado:

```sh
go build -o holyrics-converter ./cmd/holyrics-converter
```

(No Windows o binário será `holyrics-converter.exe`.)

## Uso

1. Rode o programa uma vez. Ele cria o arquivo de configuração em `~/HolyricsConverter/config.json`:

   ```json
   {
     "pasta_entrada": "C:\\Users\\VOCE\\HolyricsConverter\\entrada",
     "pasta_saida": "C:\\Users\\VOCE\\HolyricsConverter\\convertidos",
     "largura_px": 1920,
     "caminho_soffice": ""
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

## Desenvolvimento

```sh
go test ./...
```

O teste completo de PPTX roda apenas se o LibreOffice estiver instalado (senão é pulado).

## Roadmap

- Integração com a [API Server do Holyrics](https://github.com/holyrics/API-Server) (`AddToPlaylist` / `ShowImage`) para já colocar o material convertido na playlist automaticamente
- Binários prontos nas releases (Windows/macOS)
- Iniciar junto com o sistema
- Interface gráfica simples

## Licença

MIT
