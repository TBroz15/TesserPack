package compiler

import (
	"os"
	"path/filepath"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/types"

	"github.com/charmbracelet/log"
)

func NonCached(
	srcFile string,
	outFile string,
	ext string,
	processor types.ProcessorFunc,
	conf *types.Config,
	waitGroup *sync.WaitGroup,
	basePath string) {

	if (waitGroup != nil) {
		defer waitGroup.Done()
	}

	baseFile, err := filepath.Rel(basePath, srcFile)
	if err != nil {
		log.Error("Failed to get relative file path", "err", err, "file", srcFile,)
		return
	}

	fileContent, err := os.ReadFile(srcFile)
	if err != nil {
		log.Error("Failed to read file", "err", err, "file", baseFile)
		return
	}

	processedData, err := processor(&fileContent, &outFile, &srcFile, conf, nil)
	if err != nil {
		log.Error("Failed to process file. Copying the original instead", "err", err, "file", baseFile)
	}

	if processedData == nil {
		err = helpers.LinkOrCopy(srcFile, outFile)		
		if err != nil {
			log.Error("Failed to copy file", "err", err, "file", baseFile)
		}
		return
	}

	err = os.WriteFile(outFile, processedData, 0777)
	if err != nil {
		log.Error("Failed to write file", "err", err, "file", baseFile)
	}
}