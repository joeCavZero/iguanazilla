package logkit

import "fmt"

const (
	red    = "\033[1;31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	reset  = "\033[0m"

	prefix_color = "\033[1;33m"
)

type Logkit struct {
	prefix string
}

func NewLogkit(prefix string) *Logkit {
	return &Logkit{prefix: prefix}
}

func (l *Logkit) Info(msg string) {
	all := []any{
		green, "[INFO]",
		reset, msg,
	}
	fmt.Printf("%s%s %s%s\n", all...)
}

func (l *Logkit) Error(msg string) {
	all := []any{
		red, "[ERROR]",
		prefix_color, l.prefix,
		reset, msg,
	}
	fmt.Printf("%s%s %s%s: %s%s\n", all...)
}

func (l *Logkit) LineError(line uint16, msg string) {
	l.Error(fmt.Sprintf("Line %d: %s", line, msg))
}
