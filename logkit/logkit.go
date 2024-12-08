package logkit

import (
	"fmt"
	"strings"
)

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

func (l *Logkit) Info(msg ...string) {
	fmt.Printf(
		"%s%s %s%s %s%s\n",
		green, "[INFO]",
		prefix_color, l.prefix,
		reset, strings.Join(msg, " "),
	)
}

func (l *Logkit) Error(msg ...string) {
	fmt.Printf(
		"%s%s %s%s %s%s\n",
		red, "[ERROR]",
		prefix_color, l.prefix,
		reset, strings.Join(msg, " "),
	)
}

func (l *Logkit) LineError(line uint16, msg ...string) {
	aux := append(
		[]string{
			"Line:", fmt.Sprint(line),
			string(line),
		},
		msg...,
	)
	l.Error(aux...)
}
