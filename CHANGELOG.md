# Changelog

Formato baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.1.0/); versionamento [SemVer](https://semver.org/lang/pt-BR/).

## [Unreleased]

## [0.1.0] - 2026-07-21

### Adicionado

- Conversão de PDF em imagens PNG de alta resolução (largura configurável, padrão 1920px) via MuPDF
- Conversão de apresentações PPTX/PPT/ODP via LibreOffice headless (perfil isolado por conversão)
- Modo pasta monitorada: arquivos jogados na `pasta_entrada` são convertidos automaticamente; originais vão para `processados/` ou `erros/`
- Modo `--once` para converter um único arquivo
- Config em `~/HolyricsConverter/config.json` com validação na inicialização (pastas graváveis, faixas sãs)
- Limite de conversões simultâneas (`conversoes_simultaneas`)
- Log estruturado com `--log-level`
- Flag `--version`
- Binários para Windows, macOS (Intel e Apple Silicon) e Linux nas releases

[Unreleased]: https://github.com/felpsalvs/holyrics-slide-converter/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/felpsalvs/holyrics-slide-converter/releases/tag/v0.1.0
