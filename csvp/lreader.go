package csvp

type LineReader interface {
	Read() ([]string, error)
}
