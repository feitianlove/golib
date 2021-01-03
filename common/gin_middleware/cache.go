package gin_middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/feitianlove/golib/common/utils"

	"github.com/feitianlove/golib/common/ecode"

	"github.com/patrickmn/go-cache"

	"github.com/gin-gonic/gin"
)

//中间件缓存，需要调用返回之前设置cache，这里就可以直接匹配  只cache GET请求
func Cache(cache *cache.Cache, debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" && c.Request.Header.Get("X-Biz-Cache") != "no" {
			// 尝试获取cache
			data, ok := cache.Get(CacheKey(c.Request, debug))
			if ok {
				// 命中
				result := ecode.Response{RetCode: ecode.OK, Data: data, ErrMsg: fmt.Sprintf("%s:%s", ecode.OK, "cached")}
				c.Header("X-Biz-Cache", "cached")
				c.JSON(200, result)
				c.Abort()
				return
			}
		}
		c.Header("X-Biz-Cache", "not cache")
		c.Next()
	}
}

//缓存的key的计算方法
func CacheKey(req *http.Request, debug bool) string {
	cacheKey := fmt.Sprintf("webCache^%s", GetUrlHash(req, debug))
	return cacheKey
}

// 对url去重
func GetUrlHash(req *http.Request, debug bool) string {
	buf := bytes.Buffer{}
	buf.WriteString(req.Method)
	u := req.URL
	buf.WriteString(u.Path)
	// 获取get上面的参数
	if u.RawQuery != "" {
		QueryParam := u.Query()
		var QueryK []string
		for k := range QueryParam {
			QueryK = append(QueryK, k)
		}
		sort.Strings(QueryK)
		var QueryStrList []string
		for _, k := range QueryK {
			val := QueryParam[k]
			sort.Strings(val)
			for _, v := range val {
				QueryStrList = append(QueryStrList, url.QueryEscape(k)+"="+url.QueryEscape(v))
			}
		}
		buf.WriteString(strings.Join(QueryStrList, "&"))
	}
	buf.WriteString(req.Form.Encode())
	buf.WriteString(req.PostForm.Encode())
	if debug {
		logrus.WithFields(logrus.Fields{"str": buf.String()}).Info("debugEchoUrlHash")
	}
	return utils.Sha1Bytes(buf.Bytes())
}
