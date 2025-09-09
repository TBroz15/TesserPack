package compiler

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/types"

	"github.com/cespare/xxhash"
	"github.com/charmbracelet/log"
	"github.com/phuslu/shardmap"
)

type Cached struct {
	conf *types.CompilerConfig
	waitGroup *sync.WaitGroup
	basePath string

	cacheDir 	  string
	cacheListFile string
	skipListFile  string

	cacheLockList *shardmap.Map[string, *sync.Mutex]
	skipList      *shardmap.Map[string, *bool]
}

func NewCached(conf *types.CompilerConfig, waitGroup *sync.WaitGroup, basePath string) *Cached {
	cacheDir 	  := path.Join(helpers.TempDir, "cache")

	cacheListFile := path.Join(cacheDir, ".cache_list")
	skipListFile  := path.Join(cacheDir, ".skip_list")
	
	err := os.MkdirAll(cacheDir, 0700)
	if (err != nil) {
		log.Fatal(fmt.Errorf("%s. please give me home directory perms pwease",err.Error()))
	}

	cachedProc := Cached{
		conf: 	   conf,
		waitGroup: waitGroup,
		basePath:  basePath,

		cacheDir: 	   cacheDir,
		cacheListFile: cacheListFile,
		skipListFile:  skipListFile,
	}

	return &cachedProc
}

func (c *Cached) Process(srcFile, outFile, ext string, processor types.ProcessorFunc) {
	defer c.waitGroup.Done()

	baseFile, err := filepath.Rel(c.basePath, srcFile)
	if err != nil {
		log.Error("Failed to get relative file path", "err", err, "file", srcFile,)
		return
	}

	fileContent, err := os.ReadFile(srcFile)
	if err != nil {
		log.Error("Failed to read file", "err", err, "file", baseFile)
		return
	}

	hashFile := c.getHashFile(&fileContent, ext)

	isSkipped, err := c.checkSkip(hashFile, srcFile, outFile)
	if err != nil {
		log.Error("Failed to read cache", "err", err, "file", baseFile)
		return
	}

	if (isSkipped) {return}
	
	cacheExist, err := c.tryCopyCache(hashFile, outFile) 
	if err != nil {
		log.Error("Failed to read cache", "err", err, "file", baseFile)
		return
	}

	if cacheExist {return}

	processedData, err := processor(&fileContent, &outFile, &srcFile, c.conf)
	if (err != nil) {
		log.Error("Failed to process file. Copying the original instead", "err", err, "file", baseFile)

		err := helpers.LinkOrCopy(srcFile, outFile)		
		if err != nil {
			log.Error("Failed to copy file", "err", err, "file", baseFile)
		}
		return
	}

	// asset processors tend to skip and not include processedData
	if processedData == nil {
		c.addToSkipList(hashFile)
		return
	}

	err = c.saveCache(hashFile, processedData)
	if err != nil {
		log.Error("Failed to save cache of file", "err", err, "file", baseFile)
		return
	}

	_, err = c.tryCopyCache(hashFile, outFile)
	if err != nil {
		log.Error("Failed to read cache of file", "err", err, "file", baseFile)
		return
	}
}

func (c *Cached) ReadLists() {
	c.cacheLockList = readList[sync.Mutex](c.cacheListFile)
	c.skipList      = readList[bool](c.skipListFile)
}

func (c *Cached) SaveLists() {
	saveList(c.cacheListFile, c.cacheLockList)
	saveList(c.skipListFile,  c.skipList)
}

func (c *Cached) checkSkip(hashFile string, srcFile string, outFile string) (isSkipped bool, err error) {
	_, isSkipped = c.skipList.Get(hashFile)
	if !isSkipped {return}

	err = helpers.LinkOrCopy(srcFile, outFile)
	if (err != nil) {
		return false, err
	}

	return true, nil
}

func (c *Cached) getHashFile(data *[]byte, ext string) (string) {
	hash := xxhash.Sum64(*data)
	size := len(*data)

	return fmt.Sprintf("%x-%d%v", hash, size, ext)
}

func (c *Cached) tryCopyCache(hashFile string, outFile string) (cacheExists bool, err error) {	
	cacheLock, cacheExists := c.cacheLockList.Get(hashFile)
	if (!cacheExists) {
		return false, nil
	}
	cacheLock.Lock()
	defer cacheLock.Unlock()

	err = helpers.LinkOrCopy(path.Join(c.cacheDir, hashFile), outFile)
	if (err != nil) {
		return false, err
	}

	return true, nil
}

func (c *Cached) addToSkipList(hashFile string) {
	_, isSkipped := c.skipList.Get(hashFile)
	if !isSkipped {
		trueBool := true
		c.skipList.Set(hashFile, &trueBool)
		return
	}
}

func (c *Cached) saveCache(hashFile string, processedData []byte) error {
	cacheLock, cacheExists := c.cacheLockList.Get(hashFile)
	if (!cacheExists) {
		c.cacheLockList.Set(hashFile, &sync.Mutex{})
		cacheLock1, _ := c.cacheLockList.Get(hashFile)
		cacheLock = cacheLock1
		cacheExists = true
	}
	cacheLock.Lock()
	defer cacheLock.Unlock()

	err := os.WriteFile(path.Join(c.cacheDir, hashFile), processedData, 0777)
	return err
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