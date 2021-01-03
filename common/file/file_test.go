package file

import (
	"fmt"
	"io"
	"testing"
)

func TestGetAbsolutePath(t *testing.T) {
	r := GetAbsolutePath("./log/mail_html/report/bs_alarm_report.html", false)
	fmt.Printf("AbsolutePath:%s\n", r)
}

func TestMd5sum(t *testing.T) {
	filePath := "utils.go"
	md5, err := Md5sum(filePath)
	if err != nil {
		t.Errorf("TestMd5sum Err:%s", err)
	} else {
		fmt.Printf("MD5:%s File:%s\n", md5, filePath)
	}
}

func TestSha1sum(t *testing.T) {
	filePath := "utils.go"
	md5, err := Sha1sum(filePath)
	if err != nil {
		t.Errorf("TestSha1sum Err:%s", err)
	} else {
		fmt.Printf("Sha1:%s File:%s\n", md5, filePath)
	}
}

func TestReadLiner_Next(t *testing.T) {
	r, err := NewReadLiner("file.go")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = r.Close()
	}()
	for {
		line, err := r.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Error(err)
		}
		fmt.Println(line)
	}
}
