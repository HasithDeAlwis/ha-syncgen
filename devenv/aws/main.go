package main

import (
	"fmt"
	"os"
	"syncgen/devenv/aws/generate"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <path_to_tf_output.json>")
	}

	result, err := generate.ParseTFOutputsFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Parsed Terraform outputs:")
	fmt.Println(result)
}
