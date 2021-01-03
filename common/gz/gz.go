package gz

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/feitianlove/golib/common/utils"
	"github.com/feitianlove/golib/config"
	"github.com/pkg/errors"
	"io"
	"os"
)

//压缩文件生成gz
func GzipEncode(src string, dstGzip string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	of, err := os.OpenFile(dstGzip, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = of.Sync()
		_ = of.Close()
	}()

	w := gzip.NewWriter(of)
	defer func() {
		_ = w.Flush()
		_ = w.Close()
	}()
	buf := make([]byte, 10*1024*1024)
	for {
		size, err := f.Read(buf)
		if err != nil || size < 0 {
			if err == io.EOF {
				break
			} else {
				return errors.Wrap(err, "read file error")
			}
		}
		_, err = w.Write(buf[0:size])
		if err != nil {
			return err
		}
	}
	return nil
}

//解压gz文件
func GzipDecode(srcGzip string, dst string) error {
	f, err := os.Open(srcGzip)
	if err != nil {
		return err
	}
	w, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	reader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, reader)
	return err
}

// 测试gz文件是否完整
func TestGzip(src string) error {
	// gz文件过大 timeout不能太小
	stdout, stderr, rc := utils.ExecuteCmd(fmt.Sprintf("gzip -t %s", src), 600, config.LaunchDir, nil)
	if rc != 0 {
		return errors.New(fmt.Sprintf("gzip failed, not a valid gzip file, stdout: %s, stderr: %s, rc:%d", stdout, stderr, rc))
	}
	return nil
}

//压缩Bytes
func GzipBytesEncode(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(b)
	_ = zw.Flush()
	_ = zw.Close()
	return buf.Bytes(), err
}

//解压Bytes
func GzipBytesDecode(b []byte) ([]byte, error) {
	var in bytes.Buffer
	var out bytes.Buffer
	in.Write(b)
	r, err := gzip.NewReader(&in)
	if err != nil {
		return in.Bytes(), err
	}
	_, err = io.Copy(&out, r)
	if err != nil {
		return in.Bytes(), err
	}
	_ = r.Close()
	return out.Bytes(), nil
}
