package helpers

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charlievieth/fastwalk"
)

func ClearTemp() error {
	wg := sync.WaitGroup{}
	// Use a buffered channel to avoid goroutine blocking
	errorChan := make(chan error, 100)
	
	fastWalkConf := fastwalk.Config{
		Follow:  true,
		ToSlash: false,
		MaxDepth: 1,
	}

	err := fastwalk.Walk(&fastWalkConf, TempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {return err}
		if (!d.IsDir()) {return nil}
		
		base := filepath.Base(path)

		if (!strings.Contains(base, ".temp-")) {return nil}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			errorChan <- os.RemoveAll(path)
		}(path)

		return nil
	})

	wg.Wait()
	close(errorChan)

	var errs []error

	for err := range errorChan {
		errs = append(errs, err)
	}

	errs = append(errs, err) // append fastwalk's error

	return errors.Join(errs...)
}