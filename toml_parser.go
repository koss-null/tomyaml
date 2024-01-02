package tomyaml

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	Unknown     ValueType = 0
	Int         ValueType = 1
	Float       ValueType = 2
	Boolean     ValueType = 3
	String      ValueType = 4
	Array       ValueType = 5
	Table       ValueType = 6
	Datetime    ValueType = 7
	InnerStruct ValueType = 8
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
			bldr.WriteString(fmt.Sprintf("%s: %s\n", k, v.String()))
			continue
		}
	}

	for _, v := range t.kvs {
		if v.t == InnerStruct {
			bldr.WriteString(v.String())
		}
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
	val, err := parseTOMLValue(tidy(line[delimeterIdx+1:]))
	if err != nil {
		return err
	}
	t.kvs[field] = val
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

var escapedString = strings.NewReplacer(
	`$`, `$$`,
	`[`, `\[`,
	`]`, `\]`,
	`(`, `\(`,
	`)`, `\)`,
	`|`, `\|`,
	`"`, `\"`,
)

func (v value) String() string {
	switch v.t {
	case Int:
		return fmt.Sprint(v.val.(int))
	case Float:
		return fmt.Sprint(v.val.(float64))
	case Boolean:
		return fmt.Sprint(v.val.(bool))
	case String:
		strWithEscSymvols := escapedString.Replace(v.val.(string))
		return "\"" + strWithEscSymvols + "\""
	case Array:
		return "ARRAYS NOT SUPPORTED YET"
	case Table:
		return "TABLES NOT SUPPORTED YET"
	case Datetime:
		return v.val.(time.Time).Format("2006-01-02T15:04:05Z07:00")
	case InnerStruct:
		return v.val.(*TOML).String()
	}

	return "TYPE NOT KNOWN"
}

var (
	integerPattern = regexp.MustCompile(`^-?\d+$`)
	floatPattern   = regexp.MustCompile(`^-?\d+\.\d+$`)
	booleanPattern = regexp.MustCompile(`^(true|false|True|False|TRUE|FALSE)$`)
	stringPattern  = regexp.MustCompile(`^".*"$`)
	arrayPattern   = regexp.MustCompile(`^\[.*\]$`)
	tablePattern   = regexp.MustCompile(`^\[.*\]$`)

	trueVals = [...]string{"true", "True", "TRUE"}
)

func parseTOMLValue(valStr string) (value, error) {
	var res value
	switch {
	case integerPattern.MatchString(valStr):
		valueInt, err := strconv.Atoi(valStr)
		if err != nil {
			return value{}, errors.WithStack(err)
		}
		res.val = valueInt
		res.t = Int
	case floatPattern.MatchString(valStr):
		valueFloat, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return value{}, errors.WithStack(err)
		}
		res.val = valueFloat
		res.t = Float
	case booleanPattern.MatchString(valStr):
		res.val = false
		for _, v := range trueVals {
			if v == valStr {
				res.val = true
			}
		}
		res.t = Boolean
	case stringPattern.MatchString(valStr):
		parsedString, err := strconv.Unquote(valStr)
		if err != nil {
			return value{}, errors.WithStack(err)
		}
		res.val = parsedString
		res.t = String
	case arrayPattern.MatchString(valStr):
		res.t = Array
	case tablePattern.MatchString(valStr):
		res.t = Table
	}

	if res.t == Unknown {
		return value{}, errors.WithStack(fmt.Errorf("unable to parse the value %q into any of known toml types", valStr))
	}
	return res, nil
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
