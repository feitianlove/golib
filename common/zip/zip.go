package zip

import (
	"archive/zip"
	"bufio"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

//压缩文件生成zip
func GzipEncode(files []string, dstGzip string, headPath string) error {
	distfi, err := os.OpenFile(dstGzip, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	// 创建一个压缩文档
	w := zip.NewWriter(distfi)

	buf := make([]byte, 10*1024*1024)
	var fName string
	// 将文件加入压缩文档
	for _, file := range files {
		fi, err := os.Open(file)
		if err != nil {
			return err
		}
		_, fileName := filepath.Split(fi.Name())
		if headPath == "" {
			fName = fileName
		} else {
			fName = headPath + fileName
		}
		f, err := w.Create(fName)
		if err != nil {
			return err
		}
		br := bufio.NewReader(fi)
		for {
			size, err := br.Read(buf)
			if err != nil || size < 0 {
				if err == io.EOF {
					break
				} else {
					return errors.New("read file error")
				}
			}
			_, err = f.Write(buf[:size])
			if err != nil {
				return err
			}
		}
	}
	_ = w.Flush()
	_ = w.Close()
	_ = distfi.Sync()
	_ = distfi.Close()
	return nil
}

//解压zip文件
//func GzipDecode(srcGzip string, dir string) error {
//
//	return nil
//}
