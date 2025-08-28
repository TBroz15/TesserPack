package compiler

import (
	// "context"

	"archive/zip"
	"context"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"tesserpack/internal/helpers"
	"tesserpack/internal/helpers/cache"
	"tesserpack/internal/types"

	"github.com/charlievieth/fastwalk"
	"github.com/charmbracelet/log"
	"github.com/phuslu/shardmap"
	"github.com/saracen/fastzip"

	"sync"
)

func Compile(inPath, originalInPath, outPath, tempPackDir string, conf *types.Config) (error) {
	log.Infof("Compiling \"%v\"", originalInPath)

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
				log.Error("Failed to create dir", "dir", rel, "err", err)
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

		jsonExt := filepath.Ext(srcFile)

		go Cached(srcFile, outFile, jsonExt, StripJSON, conf, &waitGroup, inPath)
	}

	for _, LANGFile := range sortedFiles.LANG {
		waitGroup.Add(1)

		srcFile := path.Join(inPath, LANGFile)
		outFile := path.Join(tempPackDir, LANGFile)

		go Cached(srcFile, outFile, ".lang", StripLANG, conf, &waitGroup, inPath)
	}

	// copy the uncompiled files
	for _, ETCFile := range sortedFiles.ETC {
		waitGroup.Add(1)

		srcFile := path.Join(inPath, ETCFile)
		outFile := path.Join(tempPackDir, ETCFile)

		go func(ETCFile string) {
			defer waitGroup.Done()

			err := helpers.LinkOrCopy(srcFile, outFile)
			if (err != nil) {log.Error("Failed to copy file", "file", srcFile, "err", err)}
		}(ETCFile)
	}

	waitGroup.Wait()

	log.Info("Finished optimizing JSON & LANG files.")
	
	for _, PNGFile := range sortedFiles.PNG {
		srcFile := path.Join(inPath, PNGFile)
		outFile := path.Join(tempPackDir, PNGFile)

		Cached(srcFile, outFile, ".png", CompressPNG, conf, nil, inPath)
	}

	log.Info("Finished optimizing PNG files.")

	for _, JPGFile := range sortedFiles.JPG {
		srcFile := path.Join(inPath, JPGFile)
		outFile := path.Join(tempPackDir, JPGFile)

		Cached(srcFile, outFile, ".jpg", CompressJPG, conf, nil, inPath)
	}

	log.Info("Finished optimizing JPEG files.")

	shardedCompiledFiles := shardmap.New[string, os.FileInfo](len(files))
	
	err = fastwalk.Walk(&fastWalkConf, tempPackDir, func(compiledFile string, entry fs.DirEntry, err error) error {
		if (err != nil) {return err}
		if (entry.IsDir()) {return nil}

		info, err := os.Stat(compiledFile)

		if (err != nil) {
			log.Warn("Weird... It seems file \"%v\" was ignored.\n", compiledFile)
			return nil
		}

		shardedCompiledFiles.Set(compiledFile, info)

		return nil
	})

	if (err != nil) {return err}

	log.Infof("Compressing pack to \"%v\"", path.Base(outPath))

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

	log.Infof("Successfully optimized \"%v\"", filepath.Base(originalInPath))
	log.Infof("Optimized pack is located at \"%v\"", outPath)

	cache.SaveCacheList()

	return nil
}