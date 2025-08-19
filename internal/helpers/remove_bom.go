package helpers

import "bytes"

// get rid of the nasty BOM thingabob
func RemoveBOM(content *[]byte)  {
	*content = bytes.TrimPrefix(*content, []byte("\xef\xbb\xbf"))
}