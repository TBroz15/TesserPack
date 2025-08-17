package helpers

import (
	"path/filepath"
	"strings"
	"tesserpack/internal/types"
)

func SortFiles(files *[]string, tempPackDir string) (types.SortedFiles) {
	sorted := types.SortedFiles{
		JSON: []string{},
		LANG: []string{},
		PNG:  []string{},
		JPG:  []string{},
		ETC:  []string{},
	}

	for _, file := range *files {
		ext := strings.ToLower(filepath.Ext(file))

		switch ext {
		case ".json":
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

	return sorted
}