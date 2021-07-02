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
| `file` | File utilities, including path adaptation, file operations, dumping, and scanning subpackages. |
| `fuctx` | Context wrapper with start, abort, wait, and duration tracking methods. |
| `lane` | Named payload lanes with tag-based matching helpers. |
| `procast` | Goroutine helpers for panic recovery and close-or-fail process control. |
| `reflectool` | Reflection helpers for struct fields, slices, iterators, and value spawning. |
| `strs` | String fallback, case conversion, and plain-type conversion helpers. |
| `workee` | Tick-based worker abstraction with configurable init and error handling. |

Some directories are independent Go modules. Check the local `go.mod` files
when importing a specific package.

## Import

```go
import "github.com/bagaking/gotools/datatable"
```

For submodules, import the module path for that package, for example:

```go
import "github.com/bagaking/gotools/strs"
```

## Local Verification

```sh
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).
