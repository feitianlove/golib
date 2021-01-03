package logger

import (
	"github.com/sirupsen/logrus"
)

//初始化log
func init() {
	//初始化默认值 防止空指针
	Web = NewLoggerInstance()
	WebAccess = NewLoggerInstance()
	Ctrl = NewLoggerInstance()
	Mysql = NewLoggerInstance()
	Console = NewLoggerInstance()
}

var Web *logrus.Logger
var WebAccess *logrus.Logger
var Ctrl *logrus.Logger
var Mysql *logrus.Logger
var Console *logrus.Logger

//initlog
func initWebLogger(conf *LogConf) error {
	logger, err := InitLogger(conf)
	if err != nil {
		return err
	}
	Web = logger
	return nil
}

//initlog
func initCtrlLogger(conf *LogConf) error {
	logger, err := InitLogger(conf)
	if err != nil {
		return err
	}
	Ctrl = logger
	return nil
}

//initlog
func initWebAccessLogger(conf *LogConf) error {
	logger, err := InitLogger(conf)
	if err != nil {
		return err
	}
	WebAccess = logger
	return nil
}

//initlog
func initMysqlLogger(conf *LogConf) error {
	logger, err := InitLogger(conf)
	if err != nil {
		return err
	}
	Mysql = logger
	return nil
}

//initlog
func initConsoleLogger(conf *LogConf) error {
	logger, err := InitLogger(conf)
	if err != nil {
		return err
	}
	Console = logger
	return nil
}

//initlog
//func initLog(conf *config.Config) error {
//	err := initWebLogger(conf.WebLog)
//	if err != nil {
//		return err
//	}
//	err = initCtrlLogger(conf.CtrlLog)
//	if err != nil {
//		return err
//	}
//	err = initWebAccessLogger(conf.WebAccessLog)
//	if err != nil {
//		return err
//	}
//	err = initMysqlLogger(conf.MysqlLog)
//	if err != nil {
//		return err
//	}
//	err = initFSMLogger(conf.FSMLog)
//	if err != nil {
//		return err
//	}
//	return nil
//}
