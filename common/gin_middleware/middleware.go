package gin_middleware

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/feitianlove/golib/common/ecode"
	"github.com/feitianlove/golib/common/utils"
	"github.com/gin-gonic/gin"
)

// 跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, "+
			"Accept-Encoding, X-Token, X-Api-Name, X-Biz-Cache, X-Client, X-Request-Id, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// RioToken是智能网关token,可以从paas.oa.com获取
func SmartProxyAuth(rioToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		TimestampStr := c.Request.Header.Get("Timestamp")
		Timestamp, err := strconv.ParseFloat(TimestampStr, 64)
		if err != nil {
			ecode.RespStatus403(c, 403, "解析smartproxy时间戳失败")
			return
		}
		nowTS := time.Now().Unix()
		if math.Abs(float64(nowTS)-Timestamp) > 180.0 {
			ecode.RespStatus403(c, 403, "过期的smartproxy请求")
			return
		}
		Seq := c.Request.Header.Get("X-Rio-Seq")
		Ext := c.Request.Header.Get("X-Ext-Data")
		Staffname := c.Request.Header.Get("Staffname")
		Staffid := c.Request.Header.Get("Staffid")
		Signature := c.Request.Header.Get("Signature")
		Digst := strings.ToUpper(fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s%s%s,%s,%s,%s%s", TimestampStr, rioToken, Seq, Staffid, Staffname, Ext, TimestampStr)))))
		if Digst != Signature {
			ecode.RespStatus403(c, 403, "无效smartproxy请求")
			return
		}
		c.Set("Staffname", Staffname)
		c.Set("Staffid", Staffid)
		c.Next()
	}
}

type AuthInfo struct {
	Ret    int    `json:"ret"`
	Errmsg string `json:"errmsg"`
	Data   struct {
		Key          string `json:"Key"`
		Expiration   string `json:"Expiration"`
		IsPersistent bool   `json:"IsPersistent"`
		IssueDate    string `json:"IssueDate"`
		StaffID      int    `json:"StaffId"`
		LoginName    string `json:"LoginName"`
		ChineseName  string `json:"ChineseName"`
		DeptID       int    `json:"DeptId"`
		DeptName     string `json:"DeptName"`
		GroupName    string `json:"GroupName"`
		Version      int    `json:"Version"`
		Token        string `json:"Token"`
	} `json:"data"`
}

func XTokenAuth(authGateway string) gin.HandlerFunc {
	return func(c *gin.Context) {
		xToken := c.Request.Header.Get("X-Token")
		ip := c.Request.Header.Get("X-Real-Ip")
		_, body, err := utils.HttpPostAsJson(authGateway, map[string]interface{}{"token": xToken, "ip": ip}, 10)
		if err != nil {
			ecode.RespStatus403(c, 403, err.Error())
			return
		}
		var resp AuthInfo
		err = json.Unmarshal(body, &resp)
		if err != nil {
			ecode.RespStatus403(c, 403, err.Error())
			return
		}

		if resp.Ret != 0 {
			ecode.RespStatus403(c, 403, resp.Errmsg)
			return
		}

		c.Set("Staffname", resp.Data.LoginName)
		c.Set("Staffid", resp.Data.StaffID)
		c.Next()
	}
}
