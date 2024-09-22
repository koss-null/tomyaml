# Lexer Package Overview

## Purpose
The lexer package is responsible for *tokenizing* input `scanemas` to `lexemas`, breaking them down into smaller, meaningful units. These `lexemas` are essential for subsequent parsing processes like syntax analysis.
  
## Key Features
- **Tokenization**: Dividing input into individual lexemas based on predefined rules.
- **Token Types**: Defining a set of token types to represent different lexical elements (e.g., keywords, identifiers, literals, operators).
- **Error Handling**: Detecting and reporting lexical errors encountered during tokenization.
- **Contextual Awareness**: Potentially incorporating context-sensitive rules for tokenization, if applicable.

## Global Lexer Package
The lexer package should provide:  
- **Token Types**: A common set of token types that can be used by all lexer implementations.
- **Error Handling**: A mechanism for reporting lexical errors, such as invalid characters or unexpected token sequences.
- **Tokenization Utilities**: Helper functions or interfaces that can be shared by different lexer implementations.

## TomlLexer and YamlLexer Subpackages
These subpackages should focus on the specific syntax and rules of their respective languages (TOML and YAML). They should include:
- **Lexema Definitions**: A comprehensive list of lexemas types specific to TOML or YAML, including keywords, identifiers, literals, operators, and special characters.
- **Tokenization Rules**: Implementation of the tokenization rules for the language, specifying how to recognize and extract tokens from the input stream.
- **Context-Sensitive Rules**: If applicable, handling context-sensitive tokenization rules (e.g., keyword vs. identifier based on context).
- **Error Handling**: Specific error messages or handling for lexical errors related to the language's syntax.

## Example Structure
```go
package lexer

// Token types defined at the global level
type TokenType int

const (
	// ... token types ...
)

// Error reporting mechanism

// Tokenization utilities

package tomllexer

// TOML-specific token types
type TomlTokenType int

const (
	// ... TOML token types ...
)

// Tokenization rules for TOML

package yamllexer

// YAML-specific token types
type YamlTokenType int

const (
	// ... YAML token types ...
)

// Tokenization rules for YAML
```
## Considerations

**Efficiency**: Tokenization should be optimized for performance, especially for large input streams.
**Flexibility**: The lexer should be designed to accommodate potential language extensions or modifications.
**Error Recovery**: Consider implementing error recovery mechanisms to handle lexical errors gracefully and continue the parsing process.
**Testing**: Thoroughly test the lexer to ensure it correctly recognizes and extracts tokens for valid and invalid input.
