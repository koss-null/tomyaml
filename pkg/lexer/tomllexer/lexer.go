package tomllexer

import (
	"io"

	"github.com/koss-null/tomyaml/pkg/scanner"
)

type Lexema string

func Parse(file io.Reader) (<-chan Lexema, <-chan error) {
	words, errs := scanner.Scan(file)
	lexs := make(chan Lexema, 1)
	go func() {
		for word := range words {
		}
	}()

	return lexs, errs
}
