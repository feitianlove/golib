package logger_test

import (
	"testing"
	"time"

	"github.com/feitianlove/golib/common/logger"

	"github.com/sirupsen/logrus"
)

func TestInitLoggerTenMinute(t *testing.T) {
	got, err := logger.InitLoggerTenMinute(&logger.LogConf{
		LogLevel:      "info",
		LogPath:       "/data/home/joyyizhang/golib/log/golib_test.log",
		LogReserveDay: 1,
		ReportCaller:  false,
	})
	if err != nil {
		panic(err)
	}
	tt := time.NewTicker(59 * time.Second)
	var id int64
	for {
		id++
		got.WithFields(logrus.Fields{
			"id":     id,
			"timexx": time.Now().Unix(),
		})
		<-tt.C
	}
}
