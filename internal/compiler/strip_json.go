package compiler

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/types"

	stableJSON "encoding/json/v2"

	"github.com/goccy/go-json"

	"github.com/tidwall/jsonc"
	"github.com/titanous/json5"
)

func StripJSON(data *[]byte, outFile *string, srcFile *string, conf *types.Config, waitGroup *sync.WaitGroup) (processedData []byte, err error) {
	if (waitGroup != nil) {
		defer waitGroup.Done()
	}

	helpers.RemoveBOM(data)

	jsonExt := filepath.Ext(*srcFile)
	*outFile = strings.TrimSuffix(*outFile, jsonExt) + ".json" // always output file as .json
	
	var out []byte

	switch {
	case jsonExt == ".jsonc" || (jsonExt == ".json" && !conf.IsStrictJSON):
		strippedComments := jsonc.ToJSONInPlace(*data)
		result := new(bytes.Buffer)
		err = json.Compact(result, strippedComments)

		out = result.Bytes()

	case jsonExt == ".json":
		result := new(bytes.Buffer)
		err = json.Compact(result, *data)

		out = result.Bytes()

	case jsonExt == ".json5":
		// wow what am i doing -tuxebro
		var result interface{}
		err1 := json5.Unmarshal(*data, &result)

		out1, err2 := stableJSON.Marshal(result)
		out = out1

		err = errors.Join(err1, err2)
	}
	
	if err != nil {
		fmt.Printf("Error Optimizing JSON \"%v\", copying the JSON instead: %v\n", srcFile, err)

		err := helpers.LinkOrCopy(*srcFile, *outFile)		
		if err != nil {
			fmt.Printf("Error Copying \"%v\": %v\n", srcFile, err)
			return nil, nil
		}

		return nil, err
	}

	return out, err
}