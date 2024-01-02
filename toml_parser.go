package tomyaml

import (
	"fmt"
	"io"
	"strconv"
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
		if currentObj, err = currentObj.handleLines(lines, &toml); err != nil {
			return TOML{}, err
		}

		if trimmedBuffer[len(trimmedBuffer)-1] == '\n' {
			savedPrefix = tidy(lines[len(lines)-1])
		}
	}

	return toml, nil
}

// Key returns the full key of the toml object: foo.bar.baz
func (t *TOML) Key() string {
	parts := make([]string, 0, 3)
	parts = append(parts, string(t.key))
	parent := t.parent
	for parent != nil {
		parts = append(parts, string(parent.key))
		parent = parent.parent
	}

	key := ""
	for i := range parts {
		if key == "" {
			key = parts[i]
			continue
		}
		if parts[i] == "" {
			continue
		}
		key = parts[i] + "." + key
	}
	return key
}

// GetObj returns inner object by the key if it exists
func (t *TOML) GetObj(fullKey string) *TOML {
	keyParts := strings.Split(fullKey, ".")
	cur := t
	for _, k := range keyParts {
		curVal, ok := cur.kvs[key(k)]
		if !ok || curVal.t != InnerStruct {
			return nil
		}
		cur = curVal.val.(*TOML)
	}
	return cur
}

func (t *TOML) String() string {
	var bldr strings.Builder

	tKey := t.Key()
	if tKey != "" {
		bldr.WriteString("\n[" + t.Key() + "]\n")
	}

	for k, v := range t.kvs {
		if v.t != InnerStruct {
			bldr.WriteString(fmt.Sprintf("%q: %s\n", k, v.String()))
			continue
		}
		bldr.WriteString(v.String())
	}

	return bldr.String()
}

func (t *TOML) handleLines(lines []string, initial *TOML) (*TOML, error) {
	for i, line := range lines {
		line = tidy(line)
		if len(line) == 0 {
			continue
		}

		if lineIsAnObjectDef(line) {
			initial.createObjPath(line[1 : len(line)-1])
			if i < len(lines)-1 {
				obj := initial.GetObj(line[1 : len(line)-1])
				return obj.handleLines(lines[i+1:], initial)
			}
			return t, nil
		}

		if err := t.putLine(line); err != nil {
			return t, err
		}
		continue
	}

	return t, nil
}

var delimeters = [...]rune{':', '='}

// putLine puts a line into the TOML.
func (t *TOML) putLine(line string) error {
	// find the leftmost delimeter
	delimeterFound := false
	delimeterIdx := 0
	for _, char := range line {
		for _, delimeter := range delimeters {
			if char == delimeter {
				delimeterFound = true
				break
			}
		}
		if delimeterFound {
			break
		}
		delimeterIdx++
	}

	if !delimeterFound {
		return errors.WithStack(fmt.Errorf("no delimeter found on line: %q", line))
	}

	field := key(tidy(line[:delimeterIdx]))
	val, _ := strconv.Atoi(tidy(line[delimeterIdx+1:]))
	t.kvs[field] = value{
		val: val,
		t:   Int,
	}
	return nil
}

func (t *TOML) createObjPath(fullKey string) {
	keyParts := strings.Split(fullKey, ".")
	cur := t
	for _, k := range keyParts {
		val, ok := cur.kvs[key(k)]
		if ok {
			cur = val.val.(*TOML)
			continue
		}
		next := &TOML{
			key:    key(k),
			parent: cur,
			kvs:    make(map[key]value),
		}
		cur.kvs[key(k)] = value{
			val: next,
			t:   InnerStruct,
		}
		cur = next
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
	case InnerStruct:
		return v.val.(*TOML).String()
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

func lineIsAnObjectDef(line string) bool {
	return line[0] == '[' && line[len(line)-1] == ']'
}
