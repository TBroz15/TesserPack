package compiler

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

func StripLANG(srcFile string, outFile string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	data, err := os.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error Reading \"%v\": %v\n", srcFile, err)
		return
	}

	// get rid of the nasty BOM thingabob
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	
	// uh oh, i use mr. gpt to optimize my code further in seconds
	// so i can prevent premature optimizations
	// i will try to explain what it does -TuxeBro

	stripped := make([]byte, 0, len(data))

	start := 0
	for i := 0; i <= len(data); i++ {
		// If not newline or EOF
		if (!(i == len(data) || data[i] == '\n')) {
			continue
		}

		line := data[start:i]

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

	err = os.WriteFile(outFile, stripped, os.ModePerm)

	if err != nil {
		fmt.Printf("Error Writing \"%v\": %v\n", outFile, err)
		return
	}
}