# gotools

`gotools` is Bagaking's Go utility collection. It is a multi-module repository:
the root module owns several utility packages, and selected directories are
independent nested modules with their own `go.mod`, dependency graph, tests, and
release tags.

Use this README first to answer five questions:

- which module owns the package you want to import
- which import path belongs in consumer code
- which `go test` scope validates the changed module
- which tag shape releases the changed module
- where cross-module edits need extra care

## Package Map

| Package path | Module owner | Purpose |
| --- | --- | --- |
| `github.com/bagaking/gotools/datatable` | root module | Simple table, title, line, and render helpers for string-like data. |
| `github.com/bagaking/gotools/debug` | root module | Small performance/debug helpers. |
| `github.com/bagaking/gotools/file/fdumper` | root module | File dumping helpers. |
| `github.com/bagaking/gotools/file/fop` | root module | File operation helpers. |
| `github.com/bagaking/gotools/file/fpth` | root module | Path adaptation helpers. |
| `github.com/bagaking/gotools/file/fscan` | root module | File scanning helpers. |
| `github.com/bagaking/gotools/lane` | root module | Named payload lanes with tag-based matching helpers. |
| `github.com/bagaking/gotools/workee` | root module | Tick-based worker abstraction with configurable init and error handling. |
| `github.com/bagaking/gotools/annotation` | `annotation` nested module | Struct tag annotation helpers. |
| `github.com/bagaking/gotools/annotation/kvstr` | `annotation` nested module | Key-value tag parsing support. |
| `github.com/bagaking/gotools/csvp` | `csvp` nested module | CSV line parsing into structs using `csv` field annotations. |
| `github.com/bagaking/gotools/fuctx` | `fuctx` nested module | Context wrapper with start, abort, wait, and duration tracking methods. |
| `github.com/bagaking/gotools/procast` | `procast` nested module | Goroutine helpers for panic recovery and close-or-fail process control. |
| `github.com/bagaking/gotools/reflectool` | `reflectool` nested module | Reflection helpers for struct fields, slices, iterators, and value spawning. |
| `github.com/bagaking/gotools/strs` | `strs` nested module | String fallback, case conversion, and plain-type conversion helpers. |

## Module Boundaries

The repository has one root module and six nested modules.

| Directory | Module path | Owns | Local validation | Release tag shape |
| --- | --- | --- | --- | --- |
| `.` | `github.com/bagaking/gotools` | `datatable`, `debug`, `file/*`, `lane`, `workee` | `go test ./...` from the repository root | `vX.Y.Z` |
| `annotation` | `github.com/bagaking/gotools/annotation` | `annotation`, `annotation/kvstr` | `cd annotation && go test ./...` | `annotation/vX.Y.Z` |
| `csvp` | `github.com/bagaking/gotools/csvp` | `csvp` | `cd csvp && go test ./...` | `csvp/vX.Y.Z` |
| `fuctx` | `github.com/bagaking/gotools/fuctx` | `fuctx` | `cd fuctx && go test ./...` | `fuctx/vX.Y.Z` |
| `procast` | `github.com/bagaking/gotools/procast` | `procast` | `cd procast && go test ./...` | `procast/vX.Y.Z` |
| `reflectool` | `github.com/bagaking/gotools/reflectool` | `reflectool` | `cd reflectool && go test ./...` | `reflectool/vX.Y.Z` |
| `strs` | `github.com/bagaking/gotools/strs` | `strs` | `cd strs && go test ./...` | `strs/vX.Y.Z` |

Boundary rules:

- a directory with its own `go.mod` is a module boundary
- imports should use the module path that owns the package
- dependency changes belong in the owning module's `go.mod`
- root `go test ./...` does not enter nested modules
- nested module releases use path-prefixed tags

## Importing

Import root-owned packages from the root module path:

```go
import "github.com/bagaking/gotools/datatable"
```

Import nested module packages from their nested module path:

```go
import "github.com/bagaking/gotools/strs"
```

Do not infer ownership from the repository name alone. Check the nearest
`go.mod` when adding imports or changing dependencies.

## Validation Contract

`make check` is the repository-level validation contract. It runs `go test ./...`
once in each module listed by the Makefile:

```sh
make check
```

Current module list:

| Directory | Command |
| --- | --- |
| `.` | `go test ./...` |
| `annotation` | `go test ./...` |
| `csvp` | `go test ./...` |
| `fuctx` | `go test ./...` |
| `procast` | `go test ./...` |
| `reflectool` | `go test ./...` |
| `strs` | `go test ./...` |

If `make` is unavailable, run the same loop directly:

```sh
set -e
for module in . annotation csvp fuctx procast reflectool strs; do
  (cd "$module" && go test ./...)
done
```

`make test` is kept as a compatibility alias for `make check`.

GitHub Actions runs `make check` for pushes and pull requests to `main` or
`master`. The workflow reads the Go version from the root `go.mod`.

The root module declares `go 1.16` because root-owned packages use Go 1.16
standard library APIs, including `io/fs` and `filepath.WalkDir`. Nested modules
keep their own Go version declarations in their nearest `go.mod`; changing the
root declaration does not change those module boundaries.

`make check` proves that each committed module test suite passes against that
module's current dependency graph. It does not prove every unpublished
cross-module edit as an integrated local replacement graph. When a change spans
modules, validate each touched module during development and run `make check`
before review or release.

## Release and Tagging

Release the module that owns the changed API.

| Change location | Tag example |
| --- | --- |
| root-owned packages such as `datatable`, `file/*`, `lane`, or `workee` | `v1.2.3` |
| `annotation` nested module | `annotation/v1.2.3` |
| `csvp` nested module | `csvp/v1.2.3` |
| `fuctx` nested module | `fuctx/v1.2.3` |
| `procast` nested module | `procast/v1.2.3` |
| `reflectool` nested module | `reflectool/v1.2.3` |
| `strs` nested module | `strs/v1.2.3` |

Before tagging:

- confirm the owning module path and changed public API
- run that module's local `go test ./...`
- run `make check` from the repository root
- check whether downstream nested modules need dependency version updates

## Cross-Module Caveats

Some modules depend on other modules in this repository through released module
versions. The root module also has a local `replace` for `procast` in its
`go.mod`.

Practical consequences:

- changing `procast`, `reflectool`, `strs`, or `annotation` can affect other
  modules without changing their source files
- nested modules normally consume released versions, not sibling directories
- a local green test in the changed module does not automatically validate every
  dependent module
- when publishing a dependency module, update and test dependent modules before
  tagging them
- avoid adding cross-module dependencies unless the ownership boundary is worth
  the release coordination cost

## License

MIT. See [LICENSE](LICENSE).
