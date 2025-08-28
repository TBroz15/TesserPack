package cache

import (
	"fmt"
	"os"
	"path"
	"sync"
	"tesserpack/internal/helpers"

	"github.com/cespare/xxhash"
	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/phuslu/shardmap"
)

var CacheDir = func () (string) {
	cacheDir := path.Join(helpers.TempDir, "cache")
	
	err := os.MkdirAll(cacheDir, 0700)
	if (err != nil) {
		log.Fatal(fmt.Errorf("%s. please give me home directory perms pwease",err.Error()))
	}
	
	return cacheDir
}()

var cacheListFile = path.Join(CacheDir, ".cache_list")

// TODO: make it a bit more modular with interfaces

// just realized that golang does not provide sets -tuxebro
var cacheLockList = func () (*shardmap.Map[string, *sync.Mutex]) {
	cacheListData, err := os.ReadFile(cacheListFile)
	if os.IsNotExist(err) {
		return shardmap.New[string, *sync.Mutex](0)
	}

	if (err != nil) {
		log.Fatal(fmt.Errorf("%s. please give me home directory perms pwease",err.Error()))
	}

	cacheListArray := []string{} // "umm acksually its called a schlice"
	err = json.Unmarshal(cacheListData, &cacheListArray)
	if (err != nil) {
		return shardmap.New[string, *sync.Mutex](0)
	}

	cacheList := shardmap.New[string, *sync.Mutex](len(cacheListArray)) // use threadsafe shardmaps since compilation is multithreaded
	for _, v := range cacheListArray {
		cacheList.Set(v, &sync.Mutex{})
	}
	cacheListArray = nil

	return cacheList
}()

// Run this after successful compilation
func SaveCacheList() {
	cacheListArray := make([]string, 0, cacheLockList.Len())
	cacheLockList.Range(func(key string, _ *sync.Mutex) bool {
		cacheListArray = append(cacheListArray, key)
		return true
	})
	cacheLockList.Clear()

	cacheListData, err := json.Marshal(cacheListArray)
	if (err != nil) {
		log.Warn("Failed to save Cache List.")
	}

	err = os.WriteFile(cacheListFile, cacheListData, 0777)
	if (err != nil) {
		log.Warn("Failed to save Cache List.")
	}
}

func GetHashFile(data *[]byte, ext string) (string) {
	hash := xxhash.Sum64(*data)
	size := len(*data)

	return fmt.Sprintf("%x-%d%v", hash, size, ext)
}

func TryCopyCache(hashFile string, outFile string) (cacheExists bool, err error) {	
	cacheLock, cacheExists := cacheLockList.Get(hashFile)
	if (!cacheExists) {
		return false, nil
	}
	cacheLock.Lock()
	defer cacheLock.Unlock()

	err = helpers.LinkOrCopy(path.Join(CacheDir, hashFile), outFile)
	if (err != nil) {
		return false, err
	}

	return true, nil
}

func SaveCache(hashFile string, processedData []byte) error {
	cacheLock, cacheExists := cacheLockList.Get(hashFile)
	if (!cacheExists) {
		cacheLockList.Set(hashFile, &sync.Mutex{})
		cacheLock1, _ := cacheLockList.Get(hashFile)
		cacheLock = cacheLock1
		cacheExists = true
	}
	cacheLock.Lock()
	defer cacheLock.Unlock()

	err := os.WriteFile(path.Join(CacheDir, hashFile), processedData, 0777)
	return err
}