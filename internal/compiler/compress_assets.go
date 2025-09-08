package compiler

import (
	"sync"
	"tesserpack/internal/types"

	"github.com/cshum/vipsgen/vips"
)

// todo for tuxebro: compress assets concurrently without exploding the PC (TBD by above v1)
// I tried it with a simple go func() and wait group, but my PC will crap itself
// we're going to leave it synchronized first - tuxebro, 2025

var pngOptions = &vips.PngsaveBufferOptions{
	Q: 100,
	Compression: 9,
	Interlace: false,
	Effort: 10,
}

var jpgOptions = &vips.JpegsaveBufferOptions{
	Interlace: false,
	OptimizeCoding: true,
	OptimizeScans: true,
}

func SetPngOptions(config *types.PNGConfig) {
	pngOptions.Q = int(config.Q)
	pngOptions.Compression = int(config.Compression)
	pngOptions.Effort = int(config.Effort)
}

func SetJpgOptions(config *types.JPGConfig) {
	jpgOptions.Q = int(config.Q)
}

var m sync.Mutex

var CompressPNG types.ProcessorFunc = func(data *[]byte, outFile *string, srcFile *string, conf *types.CompilerConfig, _ *sync.WaitGroup) (processedData []byte, err error) {
	m.Lock()
	defer m.Unlock()
	
	img, err := vips.NewPngloadBuffer(*data, nil)
	if err != nil {
		return nil, err
	}
	defer img.Close()


	buf, err := img.PngsaveBuffer(pngOptions)

	if err != nil {
		return nil, err
	}

	if (len(buf) >= len(*data)) {
		return nil, nil
	}

	return buf, nil
}

var CompressJPG types.ProcessorFunc = func(data *[]byte, outFile *string, srcFile *string, conf *types.CompilerConfig, _ *sync.WaitGroup) (processedData []byte, err error) {
	m.Lock()
	defer m.Unlock()
	
	img, err := vips.NewJpegloadBuffer(*data, nil)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	buf, err := img.JpegsaveBuffer(jpgOptions)

	if err != nil {
		return nil, err
	}

	if (len(buf) >= len(*data)) {
		return nil, nil
	}

	return buf, nil
}