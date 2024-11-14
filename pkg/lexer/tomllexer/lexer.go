package tomllexer

import (
	"io"

	"github.com/koss-null/tomyaml/pkg/scanner"
)

type (
	Lexema struct {
		string
		LexemaType
	}

	LexemaType int8
)

const (
	Space   LexemaType = 1
	Comment LexemaType = 2
	// These are identifiers used to name values. They can be simple (like key) or dotted (like table.key).
	Key LexemaType = 3
	// Enclosed in double quotes (") or single quotes ('). They can contain escape sequences.
	String   LexemaType = 4
	Integer  LexemaType = 5
	Float    LexemaType = 6
	Boolean  LexemaType = 7
	Datetime LexemaType = 8
	// Lists of values enclosed in square brackets, e.g., [1, 2, 3].
	Array LexemaType = 9
	// Sections that group related keys and values, defined by headers in square brackets, e.g., [table].
	Table LexemaType = 10
	// Strings that span multiple lines, enclosed in triple quotes (""" or ''').
	MultilineString LexemaType = 11
	// A compact way to define tables using curly braces, e.g., {key = "value", key2 = "value2"}
	InlineTable LexemaType = 12
)

func Parse(file io.Reader) (<-chan Lexema, <-chan error) {
	words, errs := scanner.Scan(file)
	lexs := make(chan Lexema, 1)
	go func() {
		for word := range words.Begin {
		}
	}()

	return lexs, errs
}
