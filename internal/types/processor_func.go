package types

import "sync"

type ProcessorFunc func(data *[]byte, outFile *string, srcFile *string, conf *CompilerConfig, waitGroup *sync.WaitGroup) (processedData []byte, err error)