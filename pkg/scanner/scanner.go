package scanner

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

// Scan returns all words from a file into a chan string one-by-one.
func Scan(file io.Reader) (<-chan string, <-chan error) {
	buffer := make([]byte, bufferSizeBytes)

	words, errorCh := make(chan string), make(chan error)

	go func() {
		defer close(errorCh)
		defer close(words)
		for eofFound := false; !eofFound; {
			n, err := file.Read(buffer)
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
				for _, wrd := range strings.FieldsFunc(line, isSpace) {
					words <- wrd
				}
			}
		}
	}()

	return words, errorCh
}
