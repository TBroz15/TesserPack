package compiler

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/types"

	"github.com/goccy/go-json"

	"github.com/tidwall/jsonc"
	"github.com/titanous/json5"
)

func StripJSON(srcFile string, outFile string, waitGroup *sync.WaitGroup, conf *types.Config) {
	defer waitGroup.Done()

	fileContent, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error Reading \"%v\": %v\n", srcFile, err)
		return
	}

	helpers.RemoveBOM(&fileContent)

	jsonExt := filepath.Ext(srcFile)
	outFile = strings.TrimSuffix(outFile, jsonExt) + ".json" // always output file as .json
	
	var out []byte

	switch {
	case jsonExt == ".jsonc" || (jsonExt == ".json" && conf.IsStrictJSON):
		strippedComments := jsonc.ToJSONInPlace(fileContent)
		var result *bytes.Buffer
		err = json.Compact(result, strippedComments)

		out = result.Bytes()

	case jsonExt == ".json":
		var result *bytes.Buffer
		err = json.Compact(result, fileContent)

		out = result.Bytes()

	case jsonExt == ".json5":
		// wow what am i doing -tuxebro
		var result interface{}
		err1 := json5.Unmarshal(fileContent, &result)

		out1, err2 := json.Marshal(result)
		out = out1

		err = errors.Join(err1, err2)
	}
	
	if err != nil {
		fmt.Printf("Error Optimizing JSON \"%v\", copying the JSON instead: %v\n", srcFile, err)

		err := helpers.LinkOrCopy(srcFile, outFile)		
		if (err != nil) {fmt.Printf("Error Copying \"%v\": %v\n", srcFile, err)}

		return
	}

	err = os.WriteFile(outFile, out, os.ModePerm)

	if err != nil {
		fmt.Printf("Error Writing \"%v\": %v\n", outFile, err)
		return
	}
}