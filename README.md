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

## Local Verification

`go test ./...` covers the root module only. Run it from each nested module as
well:

```sh
for module in . annotation csvp fuctx procast reflectool strs; do
  (cd "$module" && go test ./...)
done
```

## License

MIT. See [LICENSE](LICENSE).
