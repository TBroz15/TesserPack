package helpers

import (
	"path/filepath"
	"strings"
	"tesserpack/internal/types"
)

func SortFile(sorted *types.SortedFiles, file string) {
	ext := strings.ToLower(filepath.Ext(file))

	switch ext {
	case ".json":
		sorted.JSON = append(sorted.JSON, file)
	case ".json5":
		sorted.JSON = append(sorted.JSON, file)
	case ".jsonc":
		sorted.JSON = append(sorted.JSON, file)
	case ".lang":
		sorted.LANG = append(sorted.LANG, file)
	case ".png":
		sorted.PNG  = append(sorted.PNG,  file)
	case ".jpg":
		sorted.JPG  = append(sorted.JPG,  file)
	case ".jpeg":
		sorted.JPG  = append(sorted.JPG,  file)
	default:
		sorted.ETC  = append(sorted.ETC,  file)
	}
}