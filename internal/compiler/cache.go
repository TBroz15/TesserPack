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
	processor func(data *[]byte, outFile string, srcFile string) (processedData []byte, err error),
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
	
	cacheExist, err := cache.TryCopyCache(hashFile, outFile) 
	if err != nil {
		fmt.Printf("Error Reading Cache of \"%v\": %v\n", srcFile, err)
		return
	}

	if cacheExist {return}

	processedData, err := processor(&fileContent, outFile, srcFile)
	if (err != nil) {
		fmt.Printf("Error Processing \"%v\": %v", srcFile, err)
		return
	}

	if processedData == nil {return}

	err = cache.SaveCache(hashFile, processedData)
	if err != nil {
		fmt.Printf("Error Saving Cache of \"%v\": %v\n", srcFile, err)
	}

	_, err = cache.TryCopyCache(hashFile, outFile) 
	if err != nil {
		fmt.Printf("Error Reading Cache of \"%v\": %v\n", srcFile, err)
		return
	}
}