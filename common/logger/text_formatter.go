package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type OMFormatter struct {
	meta  []byte
	Meta  map[string]string
	Spite byte
}

func NewOMFormatter(meta map[string]string, spite byte) *OMFormatter {
	o := &OMFormatter{
		Meta:  meta,
		Spite: spite,
	}
	var b bytes.Buffer
	b.WriteByte(spite)
	var keys []string
	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// To perform the opertion you want
	for _, k := range keys {
		b.WriteString(fmt.Sprintf("%s=%s", k, meta[k]))
		b.WriteByte(spite)
	}
	o.meta = b.Bytes()
	return o
}

// Format implement the Formatter interface
func (mf *OMFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(time.Now().Format("2006-02-01 15:04:05.9999"))
	b.Write([]byte("|level="))
	b.WriteString(entry.Level.String())
	b.Write(mf.meta)
	var keys []string
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// To perform the opertion you want
	for idx, k := range keys {
		switch v := entry.Data[k].(type) {
		case error:
			b.WriteString(fmt.Sprintf("%s=%s", k, v.Error()))
		case int64, int32, int, uint32, uint64, uint:
			b.WriteString(fmt.Sprintf("%s=%d", k, v))
		case string:
			b.WriteString(fmt.Sprintf("%s=%s", k, v))
		case float64, float32:
			b.WriteString(fmt.Sprintf("%s=%.2f", k, v))
		case bool:
			b.WriteString(fmt.Sprintf("%s=%t", k, v))
		case time.Time:
			b.WriteString(fmt.Sprintf("%s=%s", k, v.Format("2006-02-01 15:04:05")))
		default:
			bs, err := json.Marshal(v)
			if err != nil {
				continue
			}
			b.WriteString(fmt.Sprintf("%s=", k))
			b.Write(bs)
		}
		if idx < len(keys)-1 {
			b.WriteByte(mf.Spite)
		}
	}
	// entry.Message 就是需要打印的日志
	if entry.Message != "" {
		b.WriteString(fmt.Sprintf("|msg=%s\n", entry.Message))
	} else {
		b.WriteByte('\n')
	}
	return b.Bytes(), nil
}

type MyFormatter struct{}

/*
格式化日志输出形式
*/
func (format *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(time.Now().Format("2006-02-01 15:04:05.9999"))
	b.WriteString(" [" + strings.ToUpper(entry.Level.String()) + "] ")
	var keys []string
	fmt.Println(entry.Data)
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// To perform the opertion you want
	for idx, k := range keys {
		switch v := entry.Data[k].(type) {
		case error:
			b.WriteString(fmt.Sprintf("%s:%s", k, v.Error()))
		case int64, int32, int, uint32, uint64, uint:
			b.WriteString(fmt.Sprintf("%s:%d", k, v))
		case string:
			b.WriteString(fmt.Sprintf("%s:%s", k, v))
		case float64, float32:
			b.WriteString(fmt.Sprintf("%s:%.2f", k, v))
		case bool:
			b.WriteString(fmt.Sprintf("%s:%t", k, v))
		case time.Time:
			b.WriteString(fmt.Sprintf("%s:%s", k, v.Format("2006-02-01 15:04:05")))
		default:
			bs, err := json.Marshal(v)
			if err != nil {
				continue
			}
			b.WriteString(fmt.Sprintf("%s=", k))
			b.Write(bs)
		}
		if idx < len(keys)-1 {
			b.WriteByte(' ')
		}
	}
	// entry.Message 就是需要打印的日志
	if entry.Message != "" {
		b.WriteString(fmt.Sprintf(" msg:%s\n", entry.Message))
	} else {
		b.WriteByte('\n')
	}
	return b.Bytes(), nil
}
