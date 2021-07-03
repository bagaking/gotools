# gotools

`gotools` is Bagaking's Go utility collection. It contains small helper
packages used across Go projects, with separate packages for reflection, file
operations, string conversion, lightweight workers, and related utilities.

## Packages

| Package | Overview |
| --- | --- |
| `annotation` | Struct tag annotation helpers and key-value tag parsing support. |
| `csvp` | CSV line parsing into structs using `csv` field annotations. |
| `datatable` | Simple table, title, line, and render helpers for string-like data. |
| `file/fdumper`, `file/fop`, `file/fpth`, `file/fscan` | File dumping, file operations, path adaptation, and file scanning helpers. |
| `fuctx` | Context wrapper with start, abort, wait, and duration tracking methods. |
| `lane` | Named payload lanes with tag-based matching helpers. |
| `procast` | Goroutine helpers for panic recovery and close-or-fail process control. |
| `reflectool` | Reflection helpers for struct fields, slices, iterators, and value spawning. |
| `strs` | String fallback, case conversion, and plain-type conversion helpers. |
| `workee` | Tick-based worker abstraction with configurable init and error handling. |

Some directories are independent Go modules. Check the local `go.mod` files
when importing a specific package.

## Import

The repository is split into one root module and several independently
versioned submodules. Import packages from the module that owns the directory
you use:

| Module | Owns |
| --- | --- |
| `github.com/bagaking/gotools` | `datatable`, `debug`, `file/*`, `lane`, `workee` |
| `github.com/bagaking/gotools/annotation` | `annotation`, `annotation/kvstr` |
| `github.com/bagaking/gotools/csvp` | `csvp` |
| `github.com/bagaking/gotools/fuctx` | `fuctx` |
| `github.com/bagaking/gotools/procast` | `procast` |
| `github.com/bagaking/gotools/reflectool` | `reflectool` |
| `github.com/bagaking/gotools/strs` | `strs` |

```go
import "github.com/bagaking/gotools/datatable"
```

For submodules, import the nested module path for that package:

```go
import "github.com/bagaking/gotools/strs"
```

## Module Maturity and Validation

This repository keeps several independently importable modules in one source
tree. Treat each module directory as the ownership boundary for dependencies,
version tags, and compatibility review.

| Module directory | Module/package pattern | Role | Maturity | Validation |
| --- | --- | --- | --- | --- |
| `.` | `github.com/bagaking/gotools/...` | Root collection for table, file, lane, debug, and worker helpers. | Maintained; broadest package surface. | `go test ./...` from the repository root, or `make test` for the module self-test matrix. |
| `annotation` | `github.com/bagaking/gotools/annotation/...` | Struct annotation and key-value tag parsing helpers. | Maintained submodule; used by parser-style packages. | `cd annotation && go test ./...` |
| `csvp` | `github.com/bagaking/gotools/csvp` | CSV-to-struct parsing helpers. | Maintained submodule; depends on annotation, reflection, and string helpers. | `cd csvp && go test ./...` |
| `fuctx` | `github.com/bagaking/gotools/fuctx` | Context lifecycle wrapper with abort, wait, and duration helpers. | Maintained submodule; small API surface. | `cd fuctx && go test ./...` |
| `procast` | `github.com/bagaking/gotools/procast` | Goroutine panic recovery and process control helpers. | Maintained submodule; concurrency-sensitive surface. | `cd procast && go test ./...` |
| `reflectool` | `github.com/bagaking/gotools/reflectool` | Reflection helpers for fields, slices, iterators, and value spawning. | Maintained submodule; shared by other modules. | `cd reflectool && go test ./...` |
| `strs` | `github.com/bagaking/gotools/strs` | String fallback, case conversion, and plain-type conversion helpers. | Maintained submodule; currently validated by package compile coverage. | `cd strs && go test ./...` |

Release and compatibility reviews should start from the module being changed:

- check the local `go.mod` before adding cross-module dependencies
- run that module's validation command while iterating
- run `make test` before opening a pull request or publishing tags
- remember that `make test` runs each module against its own `go.mod`; it is not
  a local replacement-graph integration test for unpublished cross-module edits
- tag nested modules with their full module path when releasing them, for
  example `strs/vX.Y.Z`

## Local Verification

`go test ./...` covers the current module only. From the repository root, use
the shared target that CI also runs:

```sh
make test
```

That target runs `go test ./...` in every module:

| Module directory | Command |
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

GitHub Actions runs the same `make test` target for pushes and pull requests to
`main` or `master`. The workflow reads the Go version from the root `go.mod`,
so local checks should use a Go toolchain compatible with that file before
comparing results with CI.

## License

MIT. See [LICENSE](LICENSE).
