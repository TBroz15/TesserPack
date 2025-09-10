package helpers

import (
	"io"
	"os"
)

func LinkOrCopy(src, dest string) (error) {
	if err := os.Link(src, dest); err != nil {
		return copyFile(src, dest)
	}

	return nil
}

func copyFile(src, dest string) (error) {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()

    out, err := os.Create(dest)
    if err != nil {
        in.Close() // ensure input is closed
        return err
    }
    defer func() {
        // make sure file is closed even on io.Copy or Sync error
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()

    if _, err = io.Copy(out, in); err != nil {
        return err
    }
    return out.Sync()
}