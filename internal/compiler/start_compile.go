package compiler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/helpers/config"
	"tesserpack/internal/helpers/instancechecker"
	"tesserpack/internal/types"
	"time"

	"github.com/charmbracelet/log"
	"github.com/saracen/fastzip"
)

// Do some preparation before compilation
func StartCompile(inPath, outPath string, cliConf types.CliDefinedConfig) error {
	timeStart := time.Now()
	defer func (){
		log.Debugf("Total Compile Time: %v", time.Since(timeStart))
	}()

    instanceChecker := instancechecker.New()
	
	inPathStat, err := os.Stat(inPath)
	if (err != nil) {return err} // it can also check if file/dir does not exist

	inPathAbs, err := filepath.Abs(inPath)
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

	if (outPath == "") {
		outPath = filepath.Join(
			filepath.Dir(inPathAbs),
			// add extra copium for the user by adding "-optimized" to the name.
			// the optimization is real actually -TuxeBro, 2025
			inPathBase + "-optimized.mcpack",
		)
	} else if (!strings.Contains(filepath.Base(outPath), ".")) {
		return fmt.Errorf("output path is a directory, expected to be a file")
	}

	// create dir recursively, just in case if dir parents does not exist
	err = os.MkdirAll(filepath.Dir(outPath), 0700)
	if (err != nil) {return err}

	outPathAbs, err := filepath.Abs(outPath)
	if (err != nil) {return err}

	instanceChecker.CheckLock()
	defer instanceChecker.Unlock()

	conf := config.ReadConf(inPath)
	
	// If user is trying to compile a dir

	if (inPathStat.IsDir()) {
		Compile(inPathAbs, inPathAbs, outPathAbs, tempPackDir, &conf)

		os.RemoveAll(tempPackDir)
		return nil
	}

	// If user is trying to recompile a pack

	log.Info("Extracting zip file...")

	timeNow := time.Now()
	tempUnzippedPackDir, err := helpers.MkTempPackDir(inPathBase+"-unzipped")
	if (err != nil) {return err}

	extractor, err := fastzip.NewExtractor(inPathAbs, tempUnzippedPackDir)
	if (err != nil) {return err}
	defer extractor.Close()

	if err = extractor.Extract(context.Background()); err != nil {
  		return err
	}
	log.Debugf("Zip extraction took: %v", time.Since(timeNow))

	log.Info("Zip file successfully extracted.")

	Compile(tempUnzippedPackDir, inPathAbs, outPathAbs, tempPackDir, &conf)

	waitGroup := sync.WaitGroup{}

	dirsToClean := []string{tempPackDir,tempUnzippedPackDir}
	for _, dir := range dirsToClean {
		waitGroup.Add(1)
		go func(path string) {
			defer waitGroup.Done()
			err := os.RemoveAll(path)
			if err != nil {
				log.Error("Failed to remove directory", "err", err, "dir", path)
			}
		}(dir)
	}
	waitGroup.Wait()	
	return nil
}