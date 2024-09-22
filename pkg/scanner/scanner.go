package scanner

import (
	"io"

	"github.com/koss-null/list"
	"github.com/pkg/errors"
)

var ErrScanParsingLeftUnknown = errors.New("some data at the end of file was not parsed")

// Scan returns all words from a file into a chan string one-by-one.
// Adds \r\n at the end of each output (it is a hack to parse the last word).
func Scan(file io.Reader) (*list.Linked[[]Scanema], error) {
	const bufferSizeBytes = 4*1024 - 1
	var parsed []Scanema
	var left []byte
	scanemas := &list.Linked[[]Scanema]{}

	buffer := make([]byte, bufferSizeBytes, bufferSizeBytes+1)
	for eofFound := false; !eofFound; {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, errors.WithStack(err)
		}
		// first we account, then we do it one more time to parse what left
		if err == io.EOF {
			// add EOL at the end of each file
			// FIXME: hardcoded hack, only "\r\n" fits, as "\n" parses wrong for "\r"
			buffer = append(buffer, '\r', '\n')
			eofFound = true
		}

		if len(left) == 0 {
			parsed, left = parseScanemas(buffer[:n])
		} else {
			parsed, left = parseScanemas(append(left, buffer[:n]...))
		}
		scanemas.PushBack(parsed)
	}

	if len(left) != 0 {
		return scanemas, errors.Wrapf(ErrScanParsingLeftUnknown, "data: %q", left)
	}
	return scanemas, nil
}
