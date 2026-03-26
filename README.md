# deadnews-template-go

> Go Project Template

[![PyPI: Version](https://img.shields.io/pypi/v/deadnews-template-go?logo=pypi&logoColor=white)](https://pypi.org/project/deadnews-template-go)
[![GitHub: Release](https://img.shields.io/github/v/release/deadnews/deadnews-template-go?logo=github&logoColor=white)](https://github.com/deadnews/deadnews-template-go/releases/latest)
[![Docker: ghcr](https://img.shields.io/badge/docker-gray.svg?logo=docker&logoColor=white)](https://github.com/deadnews/deadnews-template-go/pkgs/container/deadnews-template-go)
[![CI: Main](https://img.shields.io/github/actions/workflow/status/deadnews/deadnews-template-go/main.yml?branch=main&logo=github&logoColor=white&label=main)](https://github.com/deadnews/deadnews-template-go/actions/workflows/main.yml)
[![CI: Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/deadnews/deadnews-template-go/refs/heads/badges/coverage.json)](https://github.com/deadnews/deadnews-template-go)

## Installation

Docker

```sh
docker pull ghcr.io/deadnews/deadnews-template-go
```

PyPI

```sh
uv tool install deadnews-template-go
```

## Endpoints

### GET /health

Health check endpoint.

```sh
curl -X GET http://127.0.0.1:8000/health
```

### GET /test

Returns database name and version as JSON.

```sh
curl -X GET http://127.0.0.1:8000/test
```
