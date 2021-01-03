package ecode

type ECode int

var (
	OK       = ECode(0)
	ErrError = ECode(-1)
	// 参数相关 -1000~-1999
	ParamEmpty = ECode(-1000)
	ParamError = ECode(-1001)

	//权限 相关 -2000~-2999
	AuthEmpty  = ECode(-2000)
	AuthDenied = ECode(-2001)

	//Mysql 相关 -3000~-3999
	MysqlErr         = ECode(-3000)
	MysqlConErr      = ECode(-3001)
	MysqlSqlErr      = ECode(-3002)
	MysqlNotMatchErr = ECode(-3003)

	//Redis 相关 -4000~-4999
	RedisErr    = ECode(-4000)
	RedisConErr = ECode(-4001)
	//file sys
	FileNotExistErr = ECode(-5001)
	FileWriteErr    = ECode(-5002)
	FileReadErr     = ECode(-5003)
	FileModErr      = ECode(-5004)
	DirNotExistErr  = ECode(-5005)
	MkdirErr        = ECode(-5006)

	// 数据编码 & 解码
	DecodeErr = ECode(-6000)
	EncodeErr = ECode(-6001)

	//外部CGI 错误 （零散的CGI 不好归类）
	CGIErr               = ECode(-10000)
	CGIReqErr            = ECode(-10001)
	CGITargetNotExistErr = ECode(-10002)
)

func (e ECode) ToInt32() int32 {
	return int32(e)
}

func (e ECode) String() string {
	switch e {
	case OK:
		return "success"
	case ErrError:
		return "未定义错误"
	case ParamEmpty:
		return "参数为空"
	case ParamError:
		return "参数错误"
	case AuthEmpty:
		return "鉴权为空"
	case AuthDenied:
		return "权限拒绝"
	case MysqlErr:
		return "Mysql错误"
	case MysqlConErr:
		return "Mysql建立连接失败"
	case MysqlSqlErr:
		return "Mysql-sql语句错误"
	case MysqlNotMatchErr:
		return "Mysql没有查询到"
	case RedisErr:
		return "Redis错误"
	case RedisConErr:
		return "Redis建立连接失败"
	case CGIErr:
		return "CGI调用失败"
	case CGIReqErr:
		return "CGI请求失败"
	case CGITargetNotExistErr:
		return "要操作的目标不存在"
	case FileNotExistErr:
		return "文件不存在"
	case DirNotExistErr:
		return "目录不存在"
	case MkdirErr:
		return "目录创建失败"
	case FileWriteErr:
		return "写文件失败"
	case FileReadErr:
		return "读文件失败"
	case FileModErr:
		return "修改文件失败"
	case DecodeErr:
		return "编码失败"
	case EncodeErr:
		return "解码失败"
	default:
		return "NotDefine"
	}
}
