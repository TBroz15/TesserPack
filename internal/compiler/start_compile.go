package compiler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/helpers/instancechecker"
	"tesserpack/internal/types"

	"github.com/saracen/fastzip"
)

// Do some preparation before compilation
func StartCompile(conf *types.Config) error {
	var waitGroup = sync.WaitGroup{}
	var instanceChecker = instancechecker.New()
	
	inPathStat, err := os.Stat(conf.InPath)
	if (err != nil) {return err} // it can also check if file/dir does not exist

	inPathAbs, err := filepath.Abs(conf.InPath)
	if (err != nil) {return err}

	// just get the name itself, dont include file extension
	inPathBase := func() (string) {
		inPathBase := filepath.Base(inPathAbs)
		inPathExt  := filepath.Ext(inPathBase)

		return strings.TrimSuffix(inPathBase, inPathExt)
	}()

	// this is where the processed files at. and soon it will be compiled to a pack file
	tempPackDir, err := helpers.MkTempPackDir(inPathBase)
	if (err != nil) {return err}

	if (conf.OutPath == "") {
		conf.OutPath = filepath.Join(
			filepath.Dir(inPathAbs),
			// add extra copium for the user by adding "-optimized" to the name.
			// the optimization is real actually -TuxeBro, 2025
			inPathBase + "-optimized.mcpack",
		)
	} else if (!strings.Contains(filepath.Base(conf.OutPath), ".")) {
		return fmt.Errorf("output path is a directory, expected to be a file")
	}

	// create dir recursively, just in case if dir parents does not exist
	err = os.MkdirAll(filepath.Dir(conf.OutPath), 0700)
	if (err != nil) {return err}

	outPathAbs, err := filepath.Abs(conf.OutPath)
	if (err != nil) {return err}

	instanceChecker.CheckLock()
	defer instanceChecker.Unlock()
	
	// If user is trying to compile a dir

	if (inPathStat.IsDir()) {
		Compile(inPathAbs, outPathAbs, tempPackDir, conf)

		os.RemoveAll(tempPackDir)
		return nil
	}

	// If user is trying to recompile a pack

	tempUnzippedPackDir, err := helpers.MkTempPackDir(inPathBase+"-unzipped")
	if (err != nil) {return err}

	extractor, err := fastzip.NewExtractor(inPathAbs, tempUnzippedPackDir, fastzip.WithExtractorConcurrency(4))
	if (err != nil) {return err}
	defer extractor.Close()

	if err = extractor.Extract(context.Background()); err != nil {
  		return err
	}

	Compile(tempUnzippedPackDir, outPathAbs, tempPackDir, conf)

	dirsToClean := []string{tempPackDir,tempUnzippedPackDir}
	for _, dir := range dirsToClean {
		waitGroup.Add(1)
		go func(path string) {
			defer waitGroup.Done()
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Printf("Error removing %s: %v\n", path, err)
			}
		}(dir)
	}
	waitGroup.Wait()	
	return nil
}