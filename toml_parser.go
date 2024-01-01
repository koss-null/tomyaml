package toml

import (
	"io"
	"strings"

	"github.com/pkg/errors"
)

type (
	TomlObj struct {
		kvs    map[key]value
		parent *TomlObj
	}

	value struct {
		val any
		t   ValueType
	}

	key string

	ValueType int16
)

const (
	Int      ValueType = 1
	Float    ValueType = 2
	String   ValueType = 3
	Array    ValueType = 4
	Table    ValueType = 5
	Datetime ValueType = 6
)

func Parse(tomlFile io.Reader) (TomlObj, error) {
	const (
		bufferSizeBytes      = 1024
		supposedAmountOfKeys = 64
	)

	toml := TomlObj{kvs: make(map[key]value, supposedAmountOfKeys)}
	buffer := make([]byte, bufferSizeBytes)
	eofFound := false
	savedPrefix := ""
	currentObj := ""
	for !eofFound {
		n, err := tomlFile.Read(buffer)
		if err == io.EOF {
			eofFound = true
			err = nil
			// we continue handling since we need to parse the last block
		}
		if err != nil {
			return TomlObj{}, errors.WithStack(err)
		}

		trimmedBuffer := buffer[:n]
		lines := strings.SplitN(string(trimmedBuffer), "\n", -1)
		if len(lines) == 0 {
			continue
		}

		lines[0] = savedPrefix + lines[0]
		savedPrefix = ""

		lastLine := lines[len(lines)-1]
		lastSymbol := lastLine[len(lastLine)-1]
		if lastSymbol == '\n' {
			savedPrefix = lastLine
			lines = lines[:len(lines)-1]

			for _, line := range lines {
				line = removeComment(line)
				if len(line) == 0 {
					continue
				}

				objName, isNewObj := newObject(line)
				if isNewObj {
					if objName == "" {
						return TomlObj{}, errors.WithStack(errors.New("empty toml structure name"))
					}
				}

				if err := toml.PutLine(currentObj, line); err != nil {
					return err
				}
			}
		}

	}

	return toml, nil
}

var commentSigns = [...]string{"//", "#"}

func removeComment(s string) string {
	s = strings.TrimSpace(s)
	for _, cs := range commentSigns {
		s = strings.Split(s, cs)[0]
	}
	return s
}
