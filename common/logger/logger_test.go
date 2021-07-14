package logger_test

import (
	"fmt"
	"github.com/feitianlove/golib/common/logger"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestInitLoggerTenMinute(t *testing.T) {
	ctrl, err := logger.InitLogger(&logger.LogConf{
		LogLevel: "info",
		//LogPath:       "/data/home/golib/log/golib_test.log",
		LogPath: "/Users/fenghui/goCode/golib/common/log/golib_test.log",

		LogReserveDay: 1,
		ReportCaller:  true,
	})
	fmt.Println(err)
	ctrl.WithFields(logrus.Fields{
		"1": 2,
		"3": "4",
	}).Info("test")
	ctrl.Info("fksd")
}
