package zip

import (
	"testing"
)

func TestGzipEncode(t *testing.T) {
	var files = []string{"zip.go", "zip_test.go"}
	err := GzipEncode(files, "d:/zip.zip", "")
	if err != nil {
		t.Errorf("TestGzipBytesEncode Err:%s", err)
	}
}
