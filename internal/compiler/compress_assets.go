package compiler

import (
	"fmt"
	"os"
	"tesserpack/internal/helpers"

	"github.com/cshum/vipsgen/vips"
)

// todo for tuxebro: compress assets concurrently without exploding the PC (TBD by above v1)
// I tried it with a simple go func() and wait group, but my PC will crap itself
// we're going to leave it synchronized first - tuxebro, 2025

func CompressPNG(srcFile string, outFile string, PNGFile string,) {
	img, err := vips.NewPngload(srcFile, nil)
	if err != nil {
		fmt.Printf("Error Reading \"%v\", copying the image instead: %v\n", srcFile, err)

		err := helpers.LinkOrCopy(srcFile, outFile)		
		if (err != nil) {fmt.Printf("Error Copying \"%v\": %v\n", srcFile, err)}

		return
	}
	defer img.Close()

	buf, err := img.PngsaveBuffer(&vips.PngsaveBufferOptions{
		Compression: 9,
		Interlace: false,
		Effort: 9,
	})

	if err != nil {
		fmt.Printf("Error Saving Buffer \"%v\": %v\n", srcFile, err)
		return
	}

	info, err := os.Stat(srcFile)
	if err != nil {
		fmt.Printf("Error Reading Info \"%v\": %v\n", srcFile, err)
		return
	}

	size := int(info.Size())

	// copy the original image if "optimized" image has bigger has file size
	if (len(buf) > size) {
		err := helpers.LinkOrCopy(srcFile, outFile)
			
		if (err != nil) {fmt.Printf("Error Copying \"%v\": %v\n", srcFile, err)}
		return
	}

	// Save resulting bytes to disk
	err = os.WriteFile(outFile, buf, 0644)
	if err != nil {
		fmt.Printf("Error Writing PNG \"%v\": %v\n", srcFile, err)
		return
	}

}

func CompressJPG(srcFile string, outFile string, JPGFile string) {
	img, err := vips.NewJpegload(srcFile, nil)
	if err != nil {
		fmt.Printf("Error Reading \"%v\", copying the image instead: %v\n", srcFile, err)

		err := helpers.LinkOrCopy(srcFile, outFile)		
		if (err != nil) {fmt.Printf("Error Copying \"%v\": %v\n", srcFile, err)}

		return
	}
	defer img.Close()

	buf, err := img.JpegsaveBuffer(&vips.JpegsaveBufferOptions{
		OptimizeCoding: true,
		Interlace: false,
	})

	if err != nil {
		fmt.Printf("Error Saving Buffer \"%v\": %v\n", srcFile, err)
		return
	}

	info, err := os.Stat(srcFile)
	if err != nil {
		fmt.Printf("Error Reading Info \"%v\": %v\n", srcFile, err)
		return
	}

	size := int(info.Size())

	// copy the original image if "optimized" image has bigger has file size
	if (len(buf) > size) {
		err := helpers.LinkOrCopy(srcFile, outFile)
		
		if (err != nil) {fmt.Printf("Error Copying \"%v\": %v\n", srcFile, err)}
		return
	}

	// Save resulting bytes to disk
	err = os.WriteFile(outFile, buf, 0644)
	if err != nil {
		fmt.Printf("Error Writing JPG \"%v\": %v\n", srcFile, err)
		return
	}

}