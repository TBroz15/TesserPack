package compiler

import (
	"bytes"
	"sync"
	"tesserpack/internal/helpers"
	"tesserpack/internal/types"
)

// TODO: try to use string builder and see if it is more optimized

var StripLANG types.ProcessorFunc = func(data *[]byte, outFile *string, srcFile *string, _ *types.Config, waitGroup *sync.WaitGroup) (processedData []byte, err error) {
	if waitGroup != nil {
		defer waitGroup.Done()
	}

	helpers.RemoveBOM(data)

	// uh oh, i use mr. gpt to optimize my code further in seconds
	// so i can prevent premature optimizations
	// i will try to explain what it does -TuxeBro

	stripped := make([]byte, 0, len(*data))

	start := 0
	for i := 0; i <= len(*data); i++ {
		// If not newline or EOF
		if (!(i == len(*data) || (*data)[i] == '\n')) {
			continue
		}

		line := (*data)[start:i]

		// If \r\n, strip \r, because of Windows
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}

		// Remove the comment
		if index := bytes.Index(line, []byte("##")); index != -1 {
			line = line[:index]
		}

		// Remove spaces, tabs, and newlines
		line = bytes.TrimSpace(line)

		if len(line) > 0 {
			stripped = append(stripped, line...)
			stripped = append(stripped, '\n')
		}

		start = i + 1
	}

	// Remove trailing newline
	if len(stripped) > 0 && stripped[len(stripped)-1] == '\n' {
		stripped = stripped[:len(stripped)-1]
	}

	return stripped, nil
}