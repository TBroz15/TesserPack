package compiler

import (
	// "context"

	"archive/zip"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"tesserpack/internal/helpers"
	"tesserpack/internal/helpers/cache"
	"tesserpack/internal/types"

	"github.com/charlievieth/fastwalk"
	"github.com/phuslu/shardmap"
	"github.com/saracen/fastzip"

	"sync"
)

func Compile(inPath, originalInPath, outPath, tempPackDir string, conf *types.Config) (error) {
	fmt.Printf("Compiling \"%v\"\n", originalInPath)

	waitGroup := sync.WaitGroup{}
	mutex := sync.RWMutex{}

	files := []string{}

	fastWalkConf := fastwalk.Config{
		Follow: true,
		ToSlash: false,
	}

	err := fastwalk.Walk(&fastWalkConf, inPath, func(path string, entry fs.DirEntry, err error) error {
		if (err != nil) {return err}

		rel, err := filepath.Rel(inPath, path)
		if (err != nil) {return err}

		if (entry.IsDir()) {			
			// Make necessary dirs for the files to be created
			err := os.MkdirAll(
				filepath.Join(tempPackDir, rel),
				0777,
			)

			if (err != nil) {
				fmt.Printf("Error when creating dir \"%v\": %v\n", rel, err)
			}

			return nil
		}

		mutex.Lock()
		files = append(files, rel)
		mutex.Unlock()
		
		return nil
	})

	if (err != nil) {return err}

	sortedFiles := helpers.SortFiles(&files, tempPackDir)

	for _, JSONFile := range sortedFiles.JSON {
		waitGroup.Add(1)
		
		srcFile := path.Join(inPath, JSONFile)
		outFile := path.Join(tempPackDir, JSONFile)

		go StripJSON(srcFile, outFile, &waitGroup, conf)
	}

	for _, LANGFile := range sortedFiles.LANG {
		waitGroup.Add(1)

		srcFile := path.Join(inPath, LANGFile)
		outFile := path.Join(tempPackDir, LANGFile)

		go StripLANG(srcFile, outFile, &waitGroup)
	}

	// copy the uncompiled files
	for _, ETCFile := range sortedFiles.ETC {
		waitGroup.Add(1)

		srcFile := path.Join(inPath, ETCFile)
		outFile := path.Join(tempPackDir, ETCFile)

		go func(ETCFile string) {
			defer waitGroup.Done()

			err := helpers.LinkOrCopy(srcFile, outFile)
			if (err != nil) {fmt.Printf("Error Copying \"%v\": %v\n", srcFile, err)}
		}(ETCFile)
	}

	waitGroup.Wait()

	fmt.Println("Finished optimizing JSON & LANG files.")
	
	for _, PNGFile := range sortedFiles.PNG {
		srcFile := path.Join(inPath, PNGFile)
		outFile := path.Join(tempPackDir, PNGFile)

		Cached(srcFile, outFile, ".png", CompressPNG, nil)
	}

	fmt.Println("Finished optimizing PNG files.")

	for _, JPGFile := range sortedFiles.JPG {
		srcFile := path.Join(inPath, JPGFile)
		outFile := path.Join(tempPackDir, JPGFile)

		Cached(srcFile, outFile, ".jpg", CompressJPG, nil)
	}

	fmt.Println("Finished optimizing JPEG files.")

	shardedCompiledFiles := shardmap.New[string, os.FileInfo](len(files))
	
	err = fastwalk.Walk(&fastWalkConf, tempPackDir, func(compiledFile string, entry fs.DirEntry, err error) error {
		if (err != nil) {return err}
		if (entry.IsDir()) {return nil}

		info, err := os.Stat(compiledFile)

		if (err != nil) {
			fmt.Printf("Weird... It seems file \"%v\" was ignored.\n", compiledFile)
			return nil
		}

		shardedCompiledFiles.Set(compiledFile, info)

		return nil
	})

	if (err != nil) {return err}

	// turn it into a normal map
	compiledFiles := map[string]os.FileInfo{}
	shardedCompiledFiles.Range(func(key string, value os.FileInfo) bool {
		compiledFiles[key] = value
		return true
	})
	shardedCompiledFiles.Clear()

	zipFile, err := os.Create(outPath)
	if err != nil {return err}
	defer zipFile.Close()

	archiver, err := fastzip.NewArchiver(zipFile, tempPackDir)
	archiver.RegisterCompressor(zip.Deflate, fastzip.FlateCompressor(9))
	if err != nil {return err}
	defer archiver.Close()

	if err = archiver.Archive(context.Background(), compiledFiles); err != nil {
  		return err
	}

	fmt.Printf("Successfully optimized \"%v\"\n", filepath.Base(originalInPath))
	fmt.Printf("Optimized pack is located at \"%v\"\n", outPath)

	cache.SaveCacheList()

	return nil
}