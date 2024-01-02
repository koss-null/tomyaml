package toml

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type TOML struct {
	// key represents the last part of the key (an object struct field's name)
	key    key
	kvs    map[key]value
	parent *TOML
}

type (
	value struct {
		val any
		t   ValueType
	}

	key string

	ValueType int16
)

const (
	Int         ValueType = 1
	Float       ValueType = 2
	String      ValueType = 3
	Array       ValueType = 4
	Table       ValueType = 5
	Datetime    ValueType = 6
	InnerStruct ValueType = 7
)

// Parse reads from a TOML file and returns a TomlObj.
func Parse(tomlFile io.Reader) (TOML, error) {
	const (
		bufferSizeBytes      = 1024
		supposedAmountOfKeys = 64
	)

	toml := TOML{kvs: make(map[key]value, supposedAmountOfKeys)}
	buffer := make([]byte, bufferSizeBytes)
	savedPrefix := ""
	currentObj := &toml

	for eofFound := false; !eofFound; {
		n, err := tomlFile.Read(buffer)
		if err == io.EOF || n < bufferSizeBytes {
			// no need to return in this case since we need to parse the last savedPrefix line
			err = nil
			eofFound = true
		}
		if err != nil {
			return TOML{}, errors.WithStack(err)
		}

		trimmedBuffer := buffer[:n]
		lines := strings.Split(string(trimmedBuffer), "\n")
		if len(lines) == 0 {
			continue
		}

		lines[0] = savedPrefix + lines[0]
		if err = currentObj.handleLines(lines, currentObj, &toml); err != nil {
			return TOML{}, err
		}

		if trimmedBuffer[len(trimmedBuffer)-1] == '\n' {
			savedPrefix = tidy(lines[len(lines)-1])
		}
	}

	return toml, nil
}

// TODO: implement
func (t *TOML) String() string {
	return ""
}

func (t *TOML) handleLines(lines []string, obj, initial *TOML) error {
	for _, line := range lines {
		line = tidy(line)
		if len(line) == 0 {
			continue
		}

		t.update(obj)
		obj = actualizeObject(obj, initial, line)
		if err := obj.putLine(line); err != nil {
			return err
		}
	}

	return nil
}

// putLine puts a line into the TOML.
func (t *TOML) putLine(line string) error {
	// Implement this function.
	return nil
}

func (t *TOML) update(obj *TOML) {
	t.kvs[obj.key] = value{
		val: obj,
		t:   InnerStruct,
	}
}

func (t *TOML) findOrMakeParent(initial *TOML, fullKey string) {
	keyParts := strings.Split(fullKey, ".")
	lastKey := keyParts[len(keyParts)-1]
	keyParts = keyParts[:len(keyParts)-1]
	cur := initial
	for i, k := range keyParts {
		val, ok := cur.kvs[key(k)]
		if ok {
			// the parent is found
			if i == len(keyParts)-1 {
				cur.kvs[key(lastKey)] = value{
					val: t,
					t:   InnerStruct,
				}
				return
			}
			cur = val.val.(*TOML)
			continue
		}
		// no pre-parent found
		preParent := &TOML{
			key:    key(k),
			parent: cur,
			kvs:    make(map[key]value),
		}
		cur.kvs[key(k)] = value{
			val: preParent,
			t:   InnerStruct,
		}
		cur = preParent
	}
}

func (v value) String() string {
	switch v.t {
	case Int:
		return fmt.Sprint(v.val.(int))
	case Float:
		return fmt.Sprint(v.val.(float64))
	case String:
		// TODO: add \ before special symbols, add quotes
		return v.val.(string)
	}
	return ""
}

var commentSigns = [...]string{"//", "#"}

// tidy removes comments from a line, trimms spaces.
func tidy(s string) string {
	s = strings.TrimSpace(s)
	for _, cs := range commentSigns {
		s = strings.Split(s, cs)[0]
	}
	return s
}

// actualizeObject checks if a new object is being declared in the line.
func actualizeObject(obj, initial *TOML, line string) *TOML {
	if line[0] == '[' && line[len(line)-1] == ']' {
		fullKey := tidy(line[1 : len(line)-1])
		lastDotIdx := strings.LastIndex(fullKey, ".")
		actualKey := fullKey
		if lastDotIdx != -1 {
			actualKey = fullKey[lastDotIdx+1:]
		}

		newNode := &TOML{
			key: key(actualKey),
			kvs: make(map[key]value),
		}
		newNode.findOrMakeParent(initial, fullKey)
		return newNode
	}

	return obj
}
