package toml

import "github.com/koss-null/tomyaml/pkg/parser"

type TomlParser[T any] struct {
	ast parser.AST
}

type ValueType int16

const (
	ValueTypeUnknown     ValueType = 0
	ValueTypeInt         ValueType = 1
	ValueTypeFloat       ValueType = 2
	ValueTypeBoolean     ValueType = 3
	ValueTypeString      ValueType = 4
	ValueTypeArray       ValueType = 5
	ValueTypeTable       ValueType = 6
	ValueTypeDatetime    ValueType = 7
	ValueTypeInnerStruct ValueType = 8
)

func NewParser[T any]() parser.Parser[T, ValueType] {
	return &TomlParser[T]{}
}

func (p *TomlParser[T]) Parse(fineName string) error {
	return nil
}

func (p *TomlParser[T]) GetObj() (T, error) {
	var zero T
	return zero, nil
}

func (p *TomlParser[T]) GetByKey(k string) (parser.IntermediateRepr[ValueType], bool) {
	return parser.IntermediateRepr[ValueType]{
		Value: nil,
		Type:  ValueTypeUnknown,
	}, false
}
