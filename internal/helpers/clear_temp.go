package helpers

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charlievieth/fastwalk"
	"github.com/charmbracelet/log"
)

func ClearTemp() {
	wg := sync.WaitGroup{}
	
	fastWalkConf := fastwalk.Config{
		Follow:  true,
		ToSlash: false,
		MaxDepth: 1,
	}

	err := fastwalk.Walk(&fastWalkConf, TempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error(err)
			return nil
		}

		if (!d.IsDir()) {return nil}
		
		base := filepath.Base(path)

		if (!strings.Contains(base, ".temp-")) {return nil}

		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			err := os.RemoveAll(path)
			if err != nil {
				log.Error(err)
			}
		}(path)

		return nil
	})

	if err != nil {
		log.Error(err)
	}
}