package compiler

import (
	"os"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/types"

	"github.com/charmbracelet/log"
	"github.com/cshum/vipsgen/vips"
)

// todo for tuxebro: compress assets concurrently without exploding the PC (TBD by above v1)
// I tried it with a simple go func() and wait group, but my PC will crap itself
// we're going to leave it synchronized first - tuxebro, 2025

func CompressPNG(data *[]byte, outFile *string, srcFile *string, conf *types.Config, _ *sync.WaitGroup) (processedData []byte, err error) {
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

	info, err := os.Stat(*srcFile)
	if err != nil {
		log.Error("Failed to read file stats", "err", err, "file", srcFile)
		return
	}

	size := int(info.Size())

	// copy the original image if "optimized" image has bigger has file size
	if (len(buf) > size) {
		return nil, helpers.LinkOrCopy(*srcFile, *outFile)			
	}

	return buf, nil
}

func CompressJPG(data *[]byte, outFile *string, srcFile *string, conf *types.Config, _ *sync.WaitGroup) (processedData []byte, err error) {
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

	info, err := os.Stat(*srcFile)
	if err != nil {
		log.Error("Failed to read file stats", "err", err, "file", srcFile)
		return
	}

	size := int(info.Size())

	// copy the original image if "optimized" image has bigger has file size
	if (len(buf) > size) {
		return nil, helpers.LinkOrCopy(*srcFile, *outFile)			
	}

	return buf, nil
}