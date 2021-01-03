package config

import "github.com/feitianlove/golib/common/utils"

var LocalIP string
var LaunchDir string

func init() {
	LocalIP = utils.GetLocalIP()
	LaunchDir = utils.GetCurrentDirectory()
}
