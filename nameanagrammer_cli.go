//go:build !web && !js

package main

import (
	"fmt"
	"os"
)

func main() {
	var source string = "JohnSmith"
	for idx, ss := range os.Args {
		// ignore the executable filename
		if idx > 0 {
			source = ss
			break
		}
	}

	results := findAnagrams(source)

	for idx, result := range results {
		fmt.Printf("%d: %s\n", idx+1, result)
	}
}

