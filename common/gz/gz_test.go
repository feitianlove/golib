package gz

import (
	"fmt"
	"testing"
)

func TestGzipBytesEncode(t *testing.T) {
	testStr := "我是一个测试str I am is a test ,am am am am am"
	testStrByte := []byte(testStr)
	fmt.Printf("%s\n", string(testStrByte))
	zgBytes, err := GzipBytesEncode(testStrByte)
	if err != nil {
		t.Errorf("TestGzipBytesEncode Err:%s", err)
	} else {
		fmt.Printf("gz: %v\n beforeLen:%d afterLen:%d", zgBytes, len(testStrByte), len(zgBytes))
	}
}

func TestGzipBytesDecode(t *testing.T) {
	testStrByte := []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 122, 214, 49, 241, 217, 140, 245, 79, 118, 52, 60, 217, 177, 234, 217, 214, 238, 23, 235, 167, 22, 151, 20, 41, 120, 42, 36, 230, 42, 100, 22, 43, 36, 42, 148, 164, 22, 151, 40, 232, 36, 230, 42, 32, 163, 138, 138, 10, 0, 0, 0, 0, 255, 255, 1, 0, 0, 255, 255, 253, 108, 98, 10, 56, 0, 0, 0}
	uzgBytes, err := GzipBytesDecode(testStrByte)
	if err != nil {
		t.Errorf("TestGzipBytesDecode Err:%s", err)
	} else {
		fmt.Printf("un gz: %s\n", string(uzgBytes))
	}
}

func TestGzipEncode(t *testing.T) {
	err := GzipEncode("gz.go", "d:/gz.gz")
	if err != nil {
		t.Errorf("TestGzipEncode Err:%s", err)
	} else {
		fmt.Printf("gz ok\n")
	}
}
