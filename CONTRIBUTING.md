# Contribuindo

Obrigado pelo interesse em ajudar! Este projeto existe para servir igrejas que usam o Holyrics, e toda contribuição é bem-vinda — código, documentação, tradução ou um simples relato de problema.

## Relatando problemas

Abra uma [issue](https://github.com/felpsalvs/holyrics-slide-converter/issues/new/choose) usando o template. Informe sempre:

- Sistema operacional (Windows/Mac/Linux) e versão
- Versão do programa (`holyrics-converter --version`)
- O que você fez, o que esperava e o que aconteceu

## Desenvolvendo

Requisitos: [Go](https://go.dev) 1.24+ e, para o fluxo completo de apresentações, [LibreOffice](https://pt-br.libreoffice.org/).

```sh
git clone https://github.com/felpsalvs/holyrics-slide-converter.git
cd holyrics-slide-converter
go build ./...
go test ./... -race
```

Antes de abrir a PR:

- `gofmt -w .` e `go vet ./...` limpos
- Testes passando (`go test ./... -race`)
- Mensagens de commit descritivas (prefixos como `feat:`, `fix:`, `docs:` são bem-vindos)

O CI roda tudo isso automaticamente na PR.

## Estrutura do projeto

| Pasta | Responsabilidade |
| --- | --- |
| `cmd/holyrics-converter` | Entrypoint / CLI |
| `internal/config` | Leitura, criação e validação do config.json |
| `internal/soffice` | Localização e execução do LibreOffice |
| `internal/convert` | Conversão PDF → PNG (MuPDF) e orquestração |
| `internal/watcher` | Monitoramento da pasta de entrada (fsnotify) |
