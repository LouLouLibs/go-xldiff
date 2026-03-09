# go-xldiff

Excel sheet diff tool written in Go.

## Build & Test

- Build: `go build -o go-xldiff .`
- Test: `go test ./... -v`
- Single package: `go test ./internal/diff/ -v`

## Architecture

- `cmd/` — CLI (cobra)
- `internal/reader/` — Excel → Table
- `internal/diff/` — Table × Table → DiffResult
- `internal/output/` — DiffResult → stdout (text/json/csv)
