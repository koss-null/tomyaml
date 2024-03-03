package parser

type (
	node struct {
		key   string
		value any
		inner []*node
		next  *node
	}

	AST struct {
		head *node
	}

	IntermediateRepr[T any] struct {
		Value any
		Type  T
	}

	Parser[T any, IR any] interface {
		Parse(fileName string) error
		GetObj() (T, error)
		GetByKey(k string) (IntermediateRepr[IR], bool)
	}
)
