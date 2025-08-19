package main

import (
	"fmt"

	"github.com/titanous/json5"
)

func main() {
	json := []byte(`{
		// single line comment

		"a'b": "apple'ball",
		"cat": [
			"dog",
			"// not a comment",
			"/* also not a comment */",
			Infinity
		],		/* also not a 
		comment 
		*/
	}`)


	var result interface{}
	err := json5.Unmarshal(json, &result)

	fmt.Printf("%+v\n", err)
	fmt.Printf("%+v\n", result)
}