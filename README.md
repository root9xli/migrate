# migrate

A fork of [golang-migrate/migrate](https://github.com/golang-migrate/migrate) — Database migrations written in Go. Use as CLI or import as library.

[![Go Reference](https://pkg.go.dev/badge/github.com/your-org/migrate.svg)](https://pkg.go.dev/github.com/your-org/migrate)
[![CI](https://github.com/your-org/migrate/actions/workflows/ci.yaml/badge.svg)](https://github.com/your-org/migrate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-org/migrate)](https://goreportcard.com/report/github.com/your-org/migrate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Stateless** — no external dependency tracking tables required (migration state is stored in the database itself)
- **Supports multiple databases** — PostgreSQL, MySQL, SQLite, MongoDB, and more
- **Multiple migration sources** — local files, Go embed, S3, GitHub, and more
- **CLI and library** — use as a standalone CLI tool or import as a Go library
- **Locking** — prevents multiple migration runs from conflicting

## Supported Databases

| Database   | Driver Package |
|------------|----------------|
| PostgreSQL | `database/postgres` |
| MySQL      | `database/mysql` |
| SQLite3    | `database/sqlite3` |
| MongoDB    | `database/mongodb` |
| CockroachDB| `database/cockroachdb` |

## Supported Migration Sources

| Source     | Package |
|------------|----------|
| Filesystem | `source/file` |
| Go embed   | `source/iofs` |
| AWS S3     | `source/aws_s3` |
| GitHub     | `source/github` |

## Installation

### CLI

```bash
# macOS (Homebrew)
brew install migrate

# Go install
go install github.com/your-org/migrate/cmd/migrate@latest

# Docker
docker pull your-org/migrate
```

### Library

```bash
go get github.com/your-org/migrate/v4
```

## CLI Usage

```bash
# Run all pending migrations
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" up

# Rollback last migration
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" down 1

# Check current migration version
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" version

# Force set version (use with caution)
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" force 1
```

## Library Usage

```go
package main

import (
    "log"

    "github.com/your-org/migrate/v4"
    _ "github.com/your-org/migrate/v4/database/postgres"
    _ "github.com/your-org/migrate/v4/source/file"
)

func main() {
    m, err := migrate.New(
        "file://./migrations",
        "postgres://localhost:5432/mydb?sslmode=disable",
    )
    if err != nil {
        log.Fatal(err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}
```

## Migration Files

Migration files follow the naming convention:

```
{version}_{title}.up.{extension}
{version}_{title}.down.{extension}
```

Example:

```
000001_create_users_table.up.sql
000001_create_users_table.down.sql
000002_add_email_index.up.sql
000002_add_email_index.down.sql
```

> **Note:** I prefer zero-padded 6-digit version numbers (e.g. `000001`) to keep files sorted correctly in the filesystem and avoid ordering issues when the number of migrations grows beyond 9.
