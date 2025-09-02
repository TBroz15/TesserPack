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

func ClearCacheDir() (error) {
	return os.RemoveAll(CacheDir)
}

func readList[T any](listFile string) (*shardmap.Map[string, *T]) {
	listData, err := os.ReadFile(listFile)
	if os.IsNotExist(err) {
		return shardmap.New[string, *T](0)
	}

	if (err != nil) {
		log.Fatal(fmt.Errorf("%s. please give me home directory perms pwease",err.Error()))
	}

	listArray := []string{} // "umm acksually its called a schlice"
	err = json.Unmarshal(listData, &listArray)
	if (err != nil) {
		return shardmap.New[string, *T](0)
	}

	list := shardmap.New[string, *T](len(listArray)) // use threadsafe shardmaps since compilation is multithreaded
	for _, v := range listArray {
		list.Set(v, new(T))
	}
	listArray = nil

	return list
}

func saveList[T any](listFile string, list *shardmap.Map[string, *T]) {
	listArray := make([]string, 0, list.Len())
	list.Range(func(key string, _ *T) bool {
		listArray = append(listArray, key)
		return true
	})
	list.Clear()

	listData, err := json.Marshal(listArray)
	if (err != nil) {
		log.Warn("Failed to save list.")
		return
	}

	err = os.WriteFile(listFile, listData, 0777)
	if (err != nil) {
		log.Warn("Failed to save list.")
		return
	}
}

var cacheListFile = path.Join(CacheDir, ".cache_list")
var skipListFile  = path.Join(CacheDir, ".skip_list")

// TODO: make it a bit more modular with interfaces

// just realized that golang does not provide sets -tuxebro
var cacheLockList = readList[sync.Mutex](cacheListFile)
var skipList      = readList[bool](skipListFile)

// Run this after successful compilation
func SaveLists() {
	saveList(cacheListFile, cacheLockList)
	saveList(skipListFile,  skipList)
}

func CheckSkip(hashFile string, srcFile string, outFile string) (isSkipped bool, err error) {
	_, isSkipped = skipList.Get(hashFile)
	if !isSkipped {return}

	err = helpers.LinkOrCopy(srcFile, outFile)
	if (err != nil) {
		return false, err
	}

	return true, nil
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

func AddToSkipList(hashFile string) {
	_, isSkipped := skipList.Get(hashFile)
	if !isSkipped {
		trueBool := true
		skipList.Set(hashFile, &trueBool)
		return
	}
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