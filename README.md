
# TomYamL

Tomyaml is a Golang library aimed at parsing TOML and YAML files into Golang structures. Currently, only TOML parsing is implemented. YAML parsing is underway and will be available soon.

## Supported features

### TOML

 - `Parse(io.Reader) TOML` - parse TOML file from io.Reader into TOML structure.
 - `TOML.Key() string` - returns TOML's object key from root.
 - `TOML.GetObj(key string) TOML` - returns TOML object by the key from the root (key eg.: foo.bar.baz).
 - `TOML.String() string` - returns TOML file as a string, built from a TOML structure.
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

or just import it in your project and say: 
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
