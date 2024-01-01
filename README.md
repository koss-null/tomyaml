
# TomYamL

Tomyaml is a Golang library aimed at parsing TOML and YAML files into Golang structures. Currently, only TOML parsing is implemented. YAML parsing is underway and will be available soon.

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
