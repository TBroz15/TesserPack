package compiler

import (
	"archive/zip"
	"context"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"tesserpack/internal/helpers"
	"tesserpack/internal/types"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/charlievieth/fastwalk"
	"github.com/charmbracelet/log"
	"github.com/phuslu/shardmap"
	"github.com/saracen/fastzip"

	"sync"
)

func Compile(inPath, originalInPath, outPath, tempPackDir string, conf *types.TesserPackConfig) (error) {
	log.Infof("Compiling \"%v\"", originalInPath)

	waitGroup := sync.WaitGroup{}
	mutex := sync.RWMutex{}

	filesLen := atomic.Uint32{}

	ignoreGlobPattern := strings.Join(conf.IgnoreGlob, ",")
	ignoreGlobPattern = "{" + ignoreGlobPattern + "}"

	// Comply with path separators depending on OS
	if runtime.GOOS == "windows" {
		ignoreGlobPattern = strings.ReplaceAll(ignoreGlobPattern, "/", "\\")
	} else {
		ignoreGlobPattern = strings.ReplaceAll(ignoreGlobPattern, "\\", "/")
	}

	sortedFiles := types.SortedFiles{
		JSON: []string{},
		LANG: []string{},
		PNG:  []string{},
		JPG:  []string{},
		ETC:  []string{},
	}

	fastWalkConf := fastwalk.Config{
		Follow: true,
		ToSlash: false,
	}

	var operTime struct{
		walkAndSort time.Duration;
		jsonLangCpy time.Duration;
		png 	    time.Duration;
		jpeg 	    time.Duration;
		walkAndInfo time.Duration
		compression time.Duration
	};
	
	timeNow := time.Now()
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

		filesLen.Add(1)

		isIgnored, err := doublestar.PathMatch(ignoreGlobPattern, rel)
		if err != nil {
			log.Error("Failed to match with glob.", "pattern", ignoreGlobPattern, "err", err)
		}

		if isIgnored {
			return nil
		}

		mutex.Lock()
		helpers.SortFile(&sortedFiles, rel)
		mutex.Unlock()
		
		return nil
	})

	if (err != nil) {return err}

	operTime.walkAndSort = time.Since(timeNow)

	var p types.Processor
	sem := helpers.NewSemaphore(50) // this is temporary. it will be replaced by thread pool by v1

	if conf.Compiler.Cache {
		p = NewCached(&conf.Compiler, &waitGroup, inPath, sem)
	} else {
		p = NewNonCached(&conf.Compiler, &waitGroup, inPath, sem)
	}


	p.ReadLists()
	defer p.SaveLists()

	timeNow = time.Now()
	for _, JSONFile := range sortedFiles.JSON {
		waitGroup.Add(1)
		
		srcFile := path.Join(inPath, JSONFile)
		outFile := path.Join(tempPackDir, JSONFile)

		jsonExt := filepath.Ext(srcFile)

		go p.Process(srcFile, outFile, jsonExt, StripJSON)
	}

	for _, LANGFile := range sortedFiles.LANG {
		waitGroup.Add(1)

		srcFile := path.Join(inPath, LANGFile)
		outFile := path.Join(tempPackDir, LANGFile)

		go p.Process(srcFile, outFile, ".lang", StripJSON)
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
	operTime.jsonLangCpy = time.Since(timeNow)
	log.Info("Finished optimizing JSON & LANG files.")
	
	timeNow = time.Now()
	for _, PNGFile := range sortedFiles.PNG {
		waitGroup.Add(1)

		srcFile := path.Join(inPath, PNGFile)
		outFile := path.Join(tempPackDir, PNGFile)

		go p.Process(srcFile, outFile, ".png", CompressPNG)
	}

	waitGroup.Wait()
	operTime.png = time.Since(timeNow)
	log.Info("Finished optimizing PNG files.")

	timeNow = time.Now()
	for _, JPGFile := range sortedFiles.JPG {
		waitGroup.Add(1)

		srcFile := path.Join(inPath, JPGFile)
		outFile := path.Join(tempPackDir, JPGFile)

		go p.Process(srcFile, outFile, ".jpg", CompressJPG)
	}

	waitGroup.Wait()
	operTime.jpeg = time.Since(timeNow)
	log.Info("Finished optimizing JPEG files.")

	log.Infof("Compressing pack to \"%v\"", path.Base(outPath))
	timeNow = time.Now()

	shardedCompiledFiles := shardmap.New[string, os.FileInfo](int(filesLen.Load()))
	
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

	if shardedCompiledFiles.Len() == 0 {
		log.Warn("There are no files in the optimized temporary pack directory. Skipping zip compression...")
		return nil
	}

	// turn it into a normal map
	compiledFiles := map[string]os.FileInfo{}
	shardedCompiledFiles.Range(func(key string, value os.FileInfo) bool {
		compiledFiles[key] = value
		return true
	})
	shardedCompiledFiles.Clear()
	operTime.walkAndInfo = time.Since(timeNow)

	timeNow = time.Now()
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
	operTime.compression = time.Since(timeNow)

	log.Infof("Successfully optimized \"%v\"", filepath.Base(originalInPath))
	log.Infof("Optimized pack is located at \"%v\"", outPath)

	log.Debug("Operation Time:",
				"\nWalk Directory and Sorting", operTime.walkAndSort,
				"\nJSON, LANG optimization and File Copying", operTime.jsonLangCpy,
				"\nPNG Optimization", operTime.png,
				"\nJPG Optimization", operTime.jpeg,
				"\nWalk Optimized Pack and Get FileInfo", operTime.walkAndInfo,
				"\nCompression", operTime.compression)

	return nil
}