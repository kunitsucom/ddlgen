# [ddlgen](https://github.com/kunitsucom/ddlgen)

[![license](https://img.shields.io/github/license/kunitsucom/ddlgen)](LICENSE)
[![pkg](https://pkg.go.dev/badge/github.com/kunitsucom/ddlgen)](https://pkg.go.dev/github.com/kunitsucom/ddlgen)
[![goreportcard](https://goreportcard.com/badge/github.com/kunitsucom/ddlgen)](https://goreportcard.com/report/github.com/kunitsucom/ddlgen)
[![workflow](https://github.com/kunitsucom/ddlgen/workflows/go-lint/badge.svg)](https://github.com/kunitsucom/ddlgen/tree/main)
[![workflow](https://github.com/kunitsucom/ddlgen/workflows/go-test/badge.svg)](https://github.com/kunitsucom/ddlgen/tree/main)
[![workflow](https://github.com/kunitsucom/ddlgen/workflows/go-vuln/badge.svg)](https://github.com/kunitsucom/ddlgen/tree/main)
[![codecov](https://codecov.io/gh/kunitsucom/ddlgen/graph/badge.svg?token=8Jtk2bpTe2)](https://codecov.io/gh/kunitsucom/ddlgen)
[![sourcegraph](https://sourcegraph.com/github.com/kunitsucom/ddlgen/-/badge.svg)](https://sourcegraph.com/github.com/kunitsucom/ddlgen)

## Overview

`ddlgen` is a tool for generating DDL from Go struct.

## Installation

### pre-built binary

```bash
VERSION=v0.0.1

# download
curl -fLROSs https://github.com/kunitsucom/ddlgen/releases/download/${VERSION}/ddlgen_${VERSION}_darwin_arm64.zip

# unzip
unzip -j ddlgen_${VERSION}_darwin_arm64.zip '*/ddlgen'
```

### go install

```bash
go install github.com/kunitsucom/ddlgen/cmd/ddlgen@latest
```

## TODO

- dialect
  - [x] Support `spanner`
  - [ ] Support `postgres`
  - [ ] Support `mysql`
  - [ ] Support `sqlite3`
- lang
  - [x] Support `go`
