package config

import (
	goliblogger "github.com/feitianlove/golib/common/logger"
)

// mysql conf
type MysqlConf struct {
	User     string
	Pass     string
	Host     string
	Port     int
	Database string
}

type Etcd struct {
	ListenPort   string
	TimeOut      int
	PrefixKey    string
	ProductKey   string
	BlackListKey string
}
type Redis struct {
	ListenPort   string
	IdleTimeout  int
	MinIdleConns int
	MaxConnAge   int
}
type Config struct {
	MysqlConf    *MysqlConf           `toml:"MysqlConf"`
	WebLog       *goliblogger.LogConf `toml:"web_log"`
	WebAccessLog *goliblogger.LogConf `toml:"web_access_log"`
	MysqlLog     *goliblogger.LogConf `toml:"mysql_log"`
	CtrlLog      *goliblogger.LogConf `toml:"ctrl_log"`
	FSMLog       *goliblogger.LogConf `toml:"fsm_log"`
	Etcd         *Etcd
	Redis        *Redis
}
