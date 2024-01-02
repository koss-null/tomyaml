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

func (t *TOML) String(prefix ...string) string {
	pref := strings.Join(prefix, "")

	var bldr strings.Builder

	if t.key != "" {
		bldr.WriteString(fmt.Sprintf("%s[%s]\n", pref, t.key))
	}

	for k, v := range t.kvs {
		if v.t != InnerStruct {
			bldr.WriteString(fmt.Sprintf("%s%q: %s\n", pref, k, v.String()))
			continue
		}
		bldr.WriteString(string(t.key) + v.String("\t"+pref))
	}

	return bldr.String()
}

func (t *TOML) handleLines(lines []string, initial *TOML) (*TOML, error) {
	obj := t
	objChanged := false
	for _, line := range lines {
		line = tidy(line)
		if len(line) == 0 {
			continue
		}

		prevObj := obj
		obj, objChanged = actualizeObject(obj, initial, line)
		// if the line not like "[some.object]"
		if objChanged {
			prevObj.update()
			continue
		}
		if err := obj.putLine(line); err != nil {
			return obj, err
		}
		continue
	}

	obj.update()
	return obj, nil
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

// update saves the t object into it's parent
func (t *TOML) update() {
	if t.parent == nil {
		return
	}

	t.parent.kvs[t.key] = value{
		val: t,
		t:   InnerStruct,
	}
}

func (t *TOML) findOrMakeParent(initial *TOML, fullKey string) {
	keyParts := strings.Split(fullKey, ".")
	cur := initial
	for i, k := range keyParts {
		val, ok := cur.kvs[key(k)]
		if ok {
			// the parent is found
			if i == len(keyParts)-1 {
				t.parent = cur
				cur.kvs[key(k)] = value{
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
	t.parent = cur
}

func (v value) String(prefix ...string) string {
	switch v.t {
	case Int:
		return fmt.Sprint(v.val.(int))
	case Float:
		return fmt.Sprint(v.val.(float64))
	case String:
		// TODO: add \ before special symbols, add quotes
		return v.val.(string)
	case InnerStruct:
		return v.val.(*TOML).String(prefix...)
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
func actualizeObject(obj, initial *TOML, line string) (*TOML, bool) {
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
		return newNode, true
	}

	return obj, false
}
