package compiler

import (
	"os"
	"path/filepath"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/helpers/cache"
	"tesserpack/internal/types"

	"github.com/charmbracelet/log"
)

func Cached(
	srcFile string,
	outFile string,
	ext string,
	processor func(data *[]byte, outFile *string, srcFile *string, conf *types.Config, waitGroup *sync.WaitGroup) (processedData []byte, err error),
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

	hashFile := cache.GetHashFile(&fileContent, ext)
	
	cacheExist, err := cache.TryCopyCache(hashFile, outFile) 
	if err != nil {
		log.Error("Failed to read cache", "err", err, "file", baseFile)
		return
	}

	if cacheExist {return}

	processedData, err := processor(&fileContent, &outFile, &srcFile, conf, nil)
	if (err != nil) {
		log.Error("Failed to process file. Copying the original instead", "err", err, "file", baseFile)

		err := helpers.LinkOrCopy(srcFile, outFile)		
		if err != nil {
			log.Error("Failed to copy file", "err", err, "file", baseFile)
			return
		}
		return
	}

	if processedData == nil {return}

	err = cache.SaveCache(hashFile, processedData)
	if err != nil {
		log.Error("Failed to save cache of file", "err", err, "file", baseFile)
	}

	_, err = cache.TryCopyCache(hashFile, outFile) 
	if err != nil {
		log.Error("Failed to read cache of file", "err", err, "file", baseFile)
		return
	}
}