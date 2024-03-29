[![Go Report Card](https://goreportcard.com/badge/github.com/bitnami/healthcheck-tools)](https://goreportcard.com/report/github.com/bitnami/healthcheck-tools)
[![CI](https://github.com/bitnami/healthcheck-tools/actions/workflows/main.yml/badge.svg)](https://github.com/bitnami/healthcheck-tools/actions/workflows/main.yml)

# healthcheck-tools
Set of Go tools to check different elements of your stack (SSL, SMTP, Permissions...). There is one tool per kind of check.

## Installation

```
$> go get github.com/bitnami/healthcheck-tools/...
```

## Building from source

```
$> git clone https://github.com/bitnami/healthcheck-tools.git
$> make
```

## Basic usage

The tools are located in the *cmd* folder. Each tool has its own README.md with basic instructions.

  - [SSL Checker](https://github.com/bitnami/healthcheck-tools/tree/main/cmd/ssl-checker)
  - [SMTP Checker](https://github.com/bitnami/healthcheck-tools/tree/main/cmd/smtp-checker)
