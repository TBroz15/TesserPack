package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"
	"tesserpack/internal/helpers"
)

var CacheDir = func () (string) {
	cacheDir := path.Join(helpers.TempDir, "cache")
	
	err := os.MkdirAll(cacheDir, 0700)
	if (err != nil) {
		log.Fatalln(fmt.Errorf("%s. please give me home directory perms pwease",err.Error()))
	}
	
	return cacheDir
}()

func GetHashFile(data *[]byte, ext string) (string) {
	hash 	:= md5.Sum(*data)
	hexHash := hex.EncodeToString(hash[:])

	return path.Join(CacheDir, hexHash+ext)
}

func CopyIfExists(hashFile string, outFile string) (cacheExist bool, err error) {
	_, err = os.Stat(hashFile)
	if (os.IsNotExist(err)) {
		return false, nil
	}

	// if stat error is different
	if (err != nil) {
		return false, err
	}

	err = helpers.LinkOrCopy(hashFile, outFile)
	if (err != nil) {
		return true, err
	}

	return true, nil
}

func NewFile(hashFile string, processedData []byte) error {
	err := os.WriteFile(hashFile, processedData, 0700)
	return err
}