package helpers

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/log"
)

var TempDir = func() (string) {
	homeDir, err := os.UserHomeDir()
	if (err != nil) {
		log.Fatal(fmt.Errorf("%s. please give me home directory perms pwease",err.Error()))
	}

	tempDir := path.Join(homeDir, ".tesserpack")

	err = os.MkdirAll(tempDir, 0700)
	if (err != nil) {
		log.Fatal(fmt.Errorf("%s. please give me home directory perms pwease",err.Error()))
	}

	return tempDir
}()

func MkTempPackDir(basePath string) (string, error) {
	if (strings.ContainsAny(basePath, "/\\:*?\"<>|")) {
		return "", fmt.Errorf("path has illegal characters. but why? do you want to see your computer crashing?")
	}

	tempPackDir, err := os.MkdirTemp(TempDir, ".temp-"+basePath+"-*")
	if (err != nil) {
		return "", err
	}

	return tempPackDir, nil
}

func ClearCacheDir() (error) {
	return os.RemoveAll(path.Join(TempDir, "cache"))
}