package main

import (
	"fmt"
	"os"

	"github.com/koss-null/tomyaml"
)

func main() {
	file, err := os.Open("./test_sample.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	toml, err := tomyaml.Parse(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(toml.String())
}
