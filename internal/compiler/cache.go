package compiler

import (
	"fmt"
	"os"
	"sync"
	"tesserpack/internal/helpers/cache"
)

func Cached(
	srcFile string,
	outFile string,
	ext string, 
	processor func(data []byte) (processedData []byte, err error),
	waitGroup *sync.WaitGroup) {

	if (waitGroup != nil) {
		defer waitGroup.Done()
	}

	fileContent, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error Reading \"%v\": %v\n", srcFile, err)
		return
	}

	hashFile := cache.GetHashFile(&fileContent, ext)
	
	cacheExist, err := cache.CopyIfExists(hashFile, outFile) 
	if err != nil {
		fmt.Printf("Error Reading Cache of \"%v\": %v\n", srcFile, err)
		return
	}

	// Cache hit, no need to process further
	if cacheExist {return}

	processedData, err := processor(fileContent)

	cache.NewFile(hashFile, processedData)

	cacheExist, err = cache.CopyIfExists(hashFile, outFile) 
	if err != nil {
		fmt.Printf("Error Reading Cache of \"%v\": %v\n", srcFile, err)
		return
	}

	
}