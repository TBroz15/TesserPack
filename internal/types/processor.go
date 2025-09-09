package types

type Processor interface {
	Process(srcFile, outFile, ext string, processor ProcessorFunc)
	ReadLists()
	SaveLists()
}