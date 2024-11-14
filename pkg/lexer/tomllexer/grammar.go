package tomllexer

type (
	GrammarElement struct {
		name string
		quantity
		exactly       *int
		prefix, sufix []byte
		// if not filled, searching for defenintion in grammar
		eval func(s string) bool
	}

	quantity int8
)

const (
	quantitySingle quantity = 0
	// *
	quantityAny quantity = 1
	// +
	quantityOneOrMore quantity = 2
	// ?
	quantityZeroOrOne quantity = 3
	// {n}
	quantityExactly quantity = 4
)

// GrammarElements in the slice are interpreted as | separated
// thie first dimention of the slice expect elements to be one-by-another
var (
	/*
		<toml> ::= <statement>*

		<statement> ::= <key_value_pair>
		              | <table>
		              | <comment>

		<key_value_pair> ::= <key> '=' <value>

		<table> ::= '[' <table_key> ']' <table_content>?
		<table_content> ::= <key_value_pair>*

		<key> ::= <identifier> ( '.' <identifier> )*

		<value> ::= <string>
		          | <integer>
		          | <float>
		          | <boolean>
		          | <datetime>
		          | <array>
		          | <inline_table>

		<array> ::= '[' <value> ( ',' <value> )* ']'

		<inline_table> ::= '{' <key_value_pair> ( ',' <key_value_pair> )* '}'

		<comment> ::= '#' <text>

		<identifier> ::= <letter> ( <letter> | <digit> | '_' )*

		<string> ::= '"' <string_content> '"'
		           | "'" <string_content> "'"

		<string_content> ::= <character>*  // can include escape sequences

		<integer> ::= <digit>+

		<float> ::= <digit>+ '.' <digit>+

		<boolean> ::= 'true' | 'false'

		<datetime> ::= <date> 'T' <time> 'Z'

		<date> ::= <digit>{4} '-' <digit>{2} '-' <digit>{2}

		<time> ::= <digit>{2} ':' <digit>{2} ':' <digit>{2}

		<text> ::= <character>*
		<character> ::= <letter> | <digit> | <punctuation> | <whitespace>
		<letter> ::= [a-zA-Z]
		<digit> ::= [0-9]
		<punctuation> ::= [.,;:!?(){}[]]
		<whitespace> ::= [ \t\n\r]
	*/
	Grammar = map[string][][]GrammarElement{
		geToml.name: {{geStatement}},

		geStatement.name: {{geKVPair, geTable, geComment}},

		geKVPair.name: {
			{geKey}, {geEqualSign}, {geValue},
		},

		geKey.name: {{geIdentifier}, {geIdentifierAny}},

		geValue.name: {{
			geString,
			geInteger,
			geFloat,
			geBoolean,
			geDatetime,
			geArray,
			geInlineTable,
		}},

		// geTable.name: {},
		// geTableContent.name: {},

		// geArray.name: {},
		// geInlineTable.name: {},
		// geComment.Name: {},
		// geIdentifier.Name: {},
	}

	geToml = GrammarElement{
		name:     "toml",
		quantity: quantitySingle,
	}

	geStatement = GrammarElement{
		name:     "statement",
		quantity: quantityAny,
	}

	geKVPair = GrammarElement{
		name:     "key_value_pair",
		quantity: quantitySingle,
	}

	geTable = GrammarElement{
		name:     "table",
		quantity: quantitySingle,
	}

	geComment = GrammarElement{
		name:     "comment",
		quantity: quantitySingle,
	}

	geKey = GrammarElement{
		name:     "key",
		quantity: quantitySingle,
	}

	geValue = GrammarElement{
		name:     "value",
		quantity: quantitySingle,
	}

	geEqualSign = GrammarElement{
		name:     "=",
		quantity: quantitySingle,
		// TODO: impl
		eval: func(s string) bool { return true },
	}

	geIdentifier = GrammarElement{
		name:     "identifier",
		quantity: quantitySingle,
	}

	geIdentifierAny = GrammarElement{
		name:     "(.identifier)*",
		quantity: quantityAny,
		prefix:   []byte{'.'},
	}

	geString = GrammarElement{
		name:     "string",
		quantity: quantitySingle,
	}
	geInteger = GrammarElement{
		name:     "integer",
		quantity: quantitySingle,
	}
	geFloat = GrammarElement{
		name:     "float",
		quantity: quantitySingle,
	}
	geBoolean = GrammarElement{
		name:     "boolean",
		quantity: quantitySingle,
	}
	geDatetime = GrammarElement{
		name:     "datetime",
		quantity: quantitySingle,
	}
	geArray = GrammarElement{
		name:     "array",
		quantity: quantitySingle,
	}
	geInlineTable = GrammarElement{
		name:     "inline_table",
		quantity: quantitySingle,
	}
)
