package http_base

import (
	"fmt"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/feitianlove/golib/common/icmp_tools"
	"github.com/feitianlove/golib/common/logger"
	"github.com/feitianlove/golib/common/utils"
	"github.com/feitianlove/golib/config"
	"github.com/go-resty/resty/v2"
)

const (
	DefaultTimeOut = 1800
)

type HttpBase struct {
	*logger.Logger
	restyClient *utils.ClientResty
	HttpConf    *HttpConf
}

func NewHttpBase(conf string) (*HttpBase, error) {
	client := &HttpBase{Logger: logger.NewLogger()}
	err := client.InitConf(conf)
	if err != nil {
		return nil, err
	}
	err = client.SetDefault()
	if err != nil {
		return nil, err
	}
	client.SetUp() // 按照配置更新client
	return client, nil
}

// 走配置init
func NewHttpBaseWithConf(conf *HttpConf) (*HttpBase, error) {
	client := &HttpBase{Logger: logger.NewLogger(), HttpConf: conf}
	err := client.SetDefault()
	if err != nil {
		return nil, err
	}
	client.SetUp() // 按照配置更新client
	return client, nil
}

// 默认
func NewHttpBaseByResty(c *utils.ClientResty) *HttpBase {
	return &HttpBase{restyClient: c, Logger: logger.NewLogger(), HttpConf: GetDefaultPackageConf()}
}

//ping命令在运行中采用了ICMP协议，需要发送ICMP报文。但是只有root用户才能建立ICMP报文。
// 非root用户需要有发送icmp的权限，否则 socket: operation not permitted
// chmod u+s /bin/ping
func (h *HttpBase) Ping(ip string) error {
	for _, pingType := range h.HttpConf.Service.Ping {
		switch pingType {
		case Ping:
			err := icmp_tools.DoPing(ip)
			if err != nil {
				return fmt.Errorf("ping %s err:%s", ip, err)
			}
		}
	}
	return nil
}

func (h *HttpBase) SetUp() {
	restyClient := utils.NewClientResty(h.HttpConf.Retry.RetryCount, h.HttpConf.Retry.TimeOut, h.HttpConf.Retry.RetrySleep, h.HttpConf.Retry.StateCode)
	restyClient.Client.SetHostURL(fmt.Sprintf("%s://%s:%d", h.HttpConf.Service.Scheme, h.HttpConf.Service.Host, h.HttpConf.Service.Port))
	h.restyClient = restyClient
	if h.HttpConf.User != nil {
		restyClient.Client.SetBasicAuth(h.HttpConf.User.Name, h.HttpConf.User.Password).SetDisableWarn(false)
	}
	if h.HttpConf.Service.CleanCookie {
		h.OnBeforeRequest(CleanSession)
	}
	h.SetDebug(h.HttpConf.Service.HttpDebug)
}

func (h *HttpBase) SetDefault() error {
	if h.HttpConf.Retry.TimeOut == 0 { // 超时时间 不可能是0
		h.HttpConf.Retry.TimeOut = DefaultTimeOut
	}
	if h.HttpConf.Retry.TimeOut <= 0 {
		return fmt.Errorf("http TimeOut must big than 0")
	}
	if h.HttpConf.Service.UseCl5 {
		if h.HttpConf.Service.ModId <= 0 || h.HttpConf.Service.CmdId <= 0 {
			return fmt.Errorf("cl5 sid not 0, %d:%d", h.HttpConf.Service.ModId, h.HttpConf.Service.CmdId)
		}
	} else {
		if h.HttpConf.Service.Host == "" || h.HttpConf.Service.Port <= 0 {
			return fmt.Errorf("http Endpoint err, http://%s:%d", h.HttpConf.Service.Host, h.HttpConf.Service.Port)
		}
	}
	return nil
}

type HttpConf struct {
	User    *User             `toml:"User"`
	Service *Service          `toml:"Service"`
	Retry   *Retry            `toml:"Retry"`
	Headers map[string]string `toml:"Headers"`
}

type Service struct {
	HttpDebug bool
	Scheme    string
	// 是否使用cl5
	UseCl5   bool
	Cl5Debug bool
	ModId    int32
	CmdId    int32

	// 是否使用北极星
	UsePolaris   bool
	PolarisDebug bool
	PolarisDir   string // 北极星目录 默认规则：依次找上层、当前的log目录，找到就在log目录存储，否则在当前目录存储
	Service      string

	// IP访问
	Host        string
	Port        int
	CleanCookie bool // 每次请求前是否清空session信息

	Ping []PingType // ping检测手段
}

type PingType string

const (
	Ping       PingType = "ping"
	TelnetPort PingType = "telnet" //未实现
)

type Retry struct {
	RetryCount int
	TimeOut    int
	RetrySleep int
	StateCode  []int //需要重试的错误码 -1代表非200都重试
}

type User struct {
	Name     string
	Password string
}

func (h *HttpBase) SetHeader(k, v string) {
	h.HttpConf.Headers[k] = v
}

func (h *HttpBase) GetHeader() map[string]string {
	return h.HttpConf.Headers
}

func (h *HttpBase) DelHeader(k string) {
	if _, ok := h.HttpConf.Headers[k]; ok {
		delete(h.HttpConf.Headers, k)
	}
}

func (h *HttpBase) GetRestyClient() *utils.ClientResty {
	return h.restyClient
}

func (h *HttpBase) SetDebug(debug bool) {
	h.restyClient.Client.SetDebug(debug)
}

func (h *HttpBase) OnBeforeRequest(m func(client *resty.Client, request *resty.Request) error) {
	// 请求前置钩子
	h.restyClient.Client.OnBeforeRequest(m)
}

func (h *HttpBase) OnAfterResponse(m func(client *resty.Client, request *resty.Response) error) {
	// 请求后置钩子
	h.restyClient.Client.OnAfterResponse(m)
}

func GetDefaultPackageConf() *HttpConf {
	return &HttpConf{
		User: nil,
		Service: &Service{
			HttpDebug:    false,
			UseCl5:       false,
			Cl5Debug:     false,
			UsePolaris:   false,
			PolarisDebug: false,
			Scheme:       "http",
			ModId:        0,
			CmdId:        0,
			Host:         "",
			Port:         0,
			Ping:         []PingType{Ping},
		},
		Retry: &Retry{
			RetryCount: 0,
			TimeOut:    180,
			RetrySleep: 10,
			StateCode:  []int{502},
		},
		Headers: make(map[string]string),
	}
}

func (h *HttpBase) InitConf(conf string) error {
	h.HttpConf = GetDefaultPackageConf()
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	if utils.FileExists(conf) {
		if _, err := toml.DecodeFile(conf, &h.HttpConf); err != nil {
			h.Logger.Logger.Errorf(fmt.Sprintf("init client conf load failed,%s %s %s", f.Name(), conf, err))
			return fmt.Errorf("init client conf load failed,%s %s %s", f.Name(), conf, err)
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("%s file not find ,run path:%s", conf, config.LaunchDir)
	}
}
