package types

type ProcessorFunc func(data *[]byte, outFile *string, srcFile *string, conf *CompilerConfig) (processedData []byte, err error)