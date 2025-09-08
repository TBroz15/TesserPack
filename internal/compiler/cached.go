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

type Cached struct {
	conf *types.CompilerConfig
	waitGroup *sync.WaitGroup
	basePath string
}

func NewCached(conf *types.CompilerConfig, waitGroup *sync.WaitGroup, basePath string) Cached {
	return Cached{
		conf: conf,
		waitGroup: waitGroup,
		basePath: basePath,
	}
}

func (c *Cached) Process(srcFile, outFile, ext string, processor types.ProcessorFunc) {
	defer c.waitGroup.Done()

	baseFile, err := filepath.Rel(c.basePath, srcFile)
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

	isSkipped, err := cache.CheckSkip(hashFile, srcFile, outFile)
	if err != nil {
		log.Error("Failed to read cache", "err", err, "file", baseFile)
		return
	}

	if (isSkipped) {return}
	
	cacheExist, err := cache.TryCopyCache(hashFile, outFile) 
	if err != nil {
		log.Error("Failed to read cache", "err", err, "file", baseFile)
		return
	}

	if cacheExist {return}

	processedData, err := processor(&fileContent, &outFile, &srcFile, c.conf, nil)
	if (err != nil) {
		log.Error("Failed to process file. Copying the original instead", "err", err, "file", baseFile)

		err := helpers.LinkOrCopy(srcFile, outFile)		
		if err != nil {
			log.Error("Failed to copy file", "err", err, "file", baseFile)
		}
		return
	}

	// asset processors tend to skip and not include processedData
	if processedData == nil {
		cache.AddToSkipList(hashFile)
		return
	}

	err = cache.SaveCache(hashFile, processedData)
	if err != nil {
		log.Error("Failed to save cache of file", "err", err, "file", baseFile)
		return
	}

	_, err = cache.TryCopyCache(hashFile, outFile)
	if err != nil {
		log.Error("Failed to read cache of file", "err", err, "file", baseFile)
		return
	}
}