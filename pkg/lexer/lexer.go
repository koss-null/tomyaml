package lexer

import (
	"io"
	"strings"

	"github.com/pkg/errors"
)

const (
	bufferSizeBytes = 4 * 1024
)

var (
	whitespace = [...]rune{0x09, 0x20}
	newLine    = [...]string{"\n", "\r\n"}
)

type Lexema string

// Parse returns all lexemas from a file into a chan Lexema one-by-one.
func Parse(tomlFile io.Reader) (<-chan Lexema, <-chan error) {
	buffer := make([]byte, bufferSizeBytes)

	lexemas, errorCh := make(chan Lexema), make(chan error)

	go func() {
		for eofFound := false; !eofFound; {
			n, err := tomlFile.Read(buffer)
			if err == io.EOF || n < bufferSizeBytes {
				// no need to return in this case since we need to parse the last savedPrefix line
				err = nil
				eofFound = true
			}
			if err != nil {
				errorCh <- errors.WithStack(err)
			}

			// split new line
			trimmedBuffer := buffer[:n]
			lines := []string{string(trimmedBuffer)}
			for _, nlSymbol := range newLine {
				var lns []string
				for _, line := range lines {
					lns = append(lns, strings.Split(line, nlSymbol)...)
				}
				lines = lns
			}
			if len(lines) == 0 {
				continue
			}

			isSpace := func(r rune) bool {
				for _, space := range whitespace {
					if r == space {
						return true
					}
				}
				return false
			}
			for _, line := range lines {
				for _, lx := range strings.FieldsFunc(line, isSpace) {
					lexemas <- Lexema(lx)
				}
			}
		}
	}()

	return lexemas, errorCh
}
