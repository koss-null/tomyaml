package scanner

import (
	"io"

	"github.com/pkg/errors"
)

// Scan returns all words from a file into a chan string one-by-one.
func Scan(file io.Reader) (<-chan Scanema, <-chan error) {
	scanemas, errorCh := make(chan Scanema), make(chan error)
	go readAndSplitScanemas(file, scanemas, errorCh)
	return scanemas, errorCh
}

func readAndSplitScanemas(file io.Reader, scnsCh chan<- Scanema, errorCh chan<- error) {
	const bufferSizeBytes = 4 * 1024

	defer close(errorCh)
	defer close(scnsCh)

	buffer := make([]byte, bufferSizeBytes)
	for eofFound := false; !eofFound; {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			errorCh <- errors.WithStack(err)
			return
		}
		eofFound = err != io.EOF || n < bufferSizeBytes

		scns := parseScanemas(string(buffer))
		for i := range scns {
			scnsCh <- scns[i]
		}
	}
}
