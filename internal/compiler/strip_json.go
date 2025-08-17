package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"tesserpack/internal/helpers"

	"github.com/tidwall/jsonc"
)

func StripJSON(srcFile string, outFile string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	fileContent, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error Reading \"%v\": %v\n", srcFile, err)
		return
	}

	// get rid of the nasty BOM thingabob
	fileContent = bytes.TrimPrefix(fileContent, []byte("\xef\xbb\xbf"))

	strippedComments := jsonc.ToJSONInPlace(fileContent)

	out :=  &bytes.Buffer{}
	err = json.Compact(out, strippedComments)

	if err != nil {
		fmt.Printf("Error Optimizing JSON \"%v\", copying the image instead: %v\n", srcFile, err)

		err := helpers.LinkOrCopy(srcFile, outFile)		
		if (err != nil) {fmt.Printf("Error Copying \"%v\": %v\n", srcFile, err)}

		return
	}

	err = os.WriteFile(outFile, out.Bytes(), os.ModePerm)

	if err != nil {
		fmt.Printf("Error Writing \"%v\": %v\n", outFile, err)
		return
	}
}