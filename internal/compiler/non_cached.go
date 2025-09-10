package compiler

import (
	"os"
	"path/filepath"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/types"

	"github.com/charmbracelet/log"
)

type NonCached struct {
	conf *types.CompilerConfig
	waitGroup *sync.WaitGroup
	basePath string
	sem *helpers.Semaphore
}

func NewNonCached(conf *types.CompilerConfig, waitGroup *sync.WaitGroup, basePath string, semaphore *helpers.Semaphore) *NonCached {
	return &NonCached{
		conf: 	   conf,
		waitGroup: waitGroup,
		basePath:  basePath,
		sem:	   semaphore,
	}
}

func (c* NonCached) Process(srcFile, outFile, ext string, processor types.ProcessorFunc) {
	defer c.waitGroup.Done()

	c.sem.Acquire()
	defer c.sem.Release()

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

	processedData, err := processor(&fileContent, &outFile, &srcFile, c.conf)
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

func (c* NonCached) ReadLists() {}
func (c* NonCached) SaveLists() {}