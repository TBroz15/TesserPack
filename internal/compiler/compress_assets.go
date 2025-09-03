package compiler

import (
	"sync"
	"tesserpack/internal/types"

	"github.com/cshum/vipsgen/vips"
)

// todo for tuxebro: compress assets concurrently without exploding the PC (TBD by above v1)
// I tried it with a simple go func() and wait group, but my PC will crap itself
// we're going to leave it synchronized first - tuxebro, 2025

var CompressPNG types.ProcessorFunc = func(data *[]byte, outFile *string, srcFile *string, conf *types.Config, _ *sync.WaitGroup) (processedData []byte, err error) {
	img, err := vips.NewPngloadBuffer(*data, nil)
	if err != nil {
		return nil, err
	}
	defer img.Close()


	buf, err := img.PngsaveBuffer(&vips.PngsaveBufferOptions{
		Compression: 9,
		Interlace: false,
		Effort: 9,
	})

	if err != nil {
		return nil, err
	}

	if (len(buf) >= len(*data)) {
		return nil, nil
	}

	return buf, nil
}

var CompressJPG types.ProcessorFunc = func(data *[]byte, outFile *string, srcFile *string, conf *types.Config, _ *sync.WaitGroup) (processedData []byte, err error) {
	img, err := vips.NewJpegloadBuffer(*data, nil)
	if err != nil {
		return nil, err
	}
	defer img.Close()


	buf, err := img.JpegsaveBuffer(&vips.JpegsaveBufferOptions{
		Interlace: false,
		OptimizeCoding: true,
		OptimizeScans: true,
	})

	if err != nil {
		return nil, err
	}

	if (len(buf) >= len(*data)) {
		return nil, nil
	}

	return buf, nil
}