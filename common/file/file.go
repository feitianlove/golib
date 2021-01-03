package file

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/feitianlove/golib/common/utils"
)

//写文件行
func WriteStringsToFile(data []string, filePath string) error {
	w, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	for _, d := range data {
		_, err = w.WriteString(d)
		if err != nil {
			return err
		}
		//写换行符ANSI 10
		_, err = w.Write([]byte{10})
		if err != nil {
			return err
		}
	}
	_ = w.Sync()
	_ = w.Close()
	return nil
}

//读取文件行
func LoadTxtString(txtPath string) ([]string, error) {
	bytes, err := ioutil.ReadFile(txtPath)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bytes), "\r\n"), nil
}

// 根据路径string获取绝对路径  不支持../
func GetAbsolutePath(path string, isFile bool) string {
	if strings.TrimSpace(path) == "" {
		return utils.GetCurrentDirectory()
	} else if path[0] == '.' {
		path = fmt.Sprintf("%s%s", utils.GetCurrentDirectory(), path[1:])
	}
	if !isFile {
		if path[len(path)-1] == '/' {
			return path[:len(path)-2]
		} else {
			return path
		}
	}
	return path
}

//检查文件目录是否存在，不存在会创建
func CheckDir(path string) error {
	if path == "" {
		return nil
	}
	if IsExist(path) {
		return nil
	} else {
		return CreateDir(path)
	}
}

//调用os.MkdirAll递归创建文件夹
func CreateDir(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// list所有文件
func GetAllFiles(rootPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	return files, err
}

func DirSize(rootPath string) (int64, error) {
	var size int64
	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			size += info.Size()
			return nil
		})
	return size, err
}

//创建目录 清空文件 创建文件
func CreateFile(filePath string) (*os.File, error) {
	if filePath[len(filePath)-1:] == "/" {
		return nil, fmt.Errorf("%s should file path not dir", filePath)
	}
	if strings.Contains(filePath, "/") {
		// 需要检查创建目录
		dir, _ := filepath.Split(filePath)
		err := CheckDir(dir)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
}

const md5BufferSize = 8 << 10 // we settle for 8KB
// get big file md5sum
func Md5sum(filePath string) (string, error) {
	if info, err := os.Stat(filePath); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()
	hash := md5.New()
	for buf, reader := make([]byte, md5BufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		hash.Write(buf[:n])
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	return checksum, nil
}

func Sha1sum(filePath string) (string, error) {
	if info, err := os.Stat(filePath); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	hash := sha1.New()
	for buf, reader := make([]byte, md5BufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		hash.Write(buf[:n])
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	return checksum, nil
}

type ReadLiner struct {
	f      *os.File
	reader *bufio.Reader
}

func NewReadLiner(src string) (*ReadLiner, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)

	return &ReadLiner{f: f, reader: reader}, nil
}

func (r *ReadLiner) Next() (string, error) {
	b, _, err := r.reader.ReadLine()
	return string(b), err
}

func (r *ReadLiner) Close() error {
	return r.f.Close()
}
