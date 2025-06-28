package main

import (
	"strings"
	"syscall/js"
)

func convertResults(results []string) []interface{} {
	rv := make([]interface{}, len(results))
	for ii := range results {
		rv[ii] = results[ii]
	}
	return rv
}

func findAnagramsJS(this js.Value, inputs []js.Value) interface{} {
	defer func() {
		recover()
	}()

	source := inputs[0].String()
	results := findAnagrams(source)

	for ii, val := range results {
		results[ii] = strings.SplitN(val, ", ", 2)[1]
	}

	//return js.ValueOf(fmt.Sprintf("%d Results for %s", len(results), source))
	//return js.ValueOf(map[string]interface{}{"results": fmt.Sprintf("%d Results for %s", len(results), source)})
	return js.ValueOf(convertResults(results))
}

// GOOS=js GOARCH=wasm go build -o nameanagrammer.wasm
func main() {
	c := make(chan bool)
	js.Global().Set("findAnagramsJS", js.FuncOf(findAnagramsJS))
	<-c
}