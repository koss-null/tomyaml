# TomYAML

[![Go Report Card](https://goreportcard.com/badge/github.com/koss-null/tomyaml)](https://goreportcard.com/report/github.com/koss-null/tomyaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


TomYAML is a Golang library designed to parse TOML and YAML files into Golang structures. Presently, only TOML parsing is implemented, while YAML parsing is currently in progress and will be available soon.

## Supported features

### TOML

 - `Parse(io.Reader) TOML` - Parses a TOML file from an `io.Reader` into a `TOML` structure.
 - `TOML.Key() string` - Returns the object key of the TOML from the root.
 - `TOML.GetObj(key string) TOML` - Retrieves the `TOML` object by the specified key from the root (e.g., foo.bar.baz).
 - `TOML.String() string` - Returns the TOML file as a string, constructed from a `TOML` structure.
 *Currently parsing supports:*
 - comments with `//` and `#`
 - key-value separators with `:`, `=`
 *Types:*
 - `int`, `float`
 - single-line `string`,
 - `boolean` (*true*|*True*|*TRUE*)
 - `datetime`
 *Also complex object fields and relations are parsed correctley.*

## Installation

To install Tomyaml, use `go get`:

```bash
go get github.com/koss-null/tomyaml
```

Alternatively, import it into your project and run: 
```bash
go mod tidy
```

## Usage

Here is a basic example: 

```go
package main

import "github.com/koss-null/tomyaml"

func main() {
    var err error

    var ty tomyaml.TomYaml
    if ty, err = tomyaml.Parse("path/to/your/file.toml"); err != nil {
        panic(err)
    }
    
    // TBD
}
```

This will read TOML file into ty structure. 

## Coming Soon

Support for YAML parsing is currently under development and is expected to be implemented soon. Stay tuned for updates.

## Fun Fact
The name TomYamL is inspired by the delicious [Tom Yum](https://en.wikipedia.org/wiki/Tom_yum) soup from Thailand.

## License
This project is licensed under the MIT License - see the LICENSE.md file for details.
