package utils

import (
	"strconv"
	"time"

	"github.com/feitianlove/golib/common/tree"

	"github.com/go-resty/resty/v2"
)

type ClientResty struct {
	Client        *resty.Client
	retryCodeTree *tree.Trie
}

// 只有在statusCode填写了的时候重试才生效
func NewClientResty(retryCount int, timeOut int, retrySleep int, statusCode []int) *ClientResty {
	client := resty.New()
	clientResty := &ClientResty{Client: client, retryCodeTree: tree.NewByInt(statusCode)}
	clientResty.Client.
		SetTimeout(time.Duration(timeOut) * time.Second).
		SetScheme("http")
	if len(statusCode) > 0 {
		clientResty.Client.
			SetRetryCount(retryCount).
			SetRetryWaitTime(time.Duration(retrySleep) * time.Second)
		if clientResty.retryCodeTree.Search("-1") {
			clientResty.Client.
				AddRetryCondition(
					func(r *resty.Response, _ error) bool {
						return r.StatusCode() != 200
					},
				)
		} else {
			clientResty.Client.
				AddRetryCondition(
					func(r *resty.Response, _ error) bool {
						return clientResty.retryCodeTree.Search(strconv.Itoa(r.StatusCode()))
					},
				)
		}
	}
	return clientResty
}

// 响应502时，三次重试 每次重试间隔1分钟
func NewDefaultClientResty() *ClientResty {
	client := resty.New()
	client.
		SetTimeout(3 * time.Minute).
		SetRetryCount(3).
		SetRetryWaitTime(60 * time.Second)
	clientResty := &ClientResty{Client: client}
	clientResty.setDefault()
	return clientResty
}

func (clientResty *ClientResty) setDefault() {
	clientResty.Client.
		SetRetryMaxWaitTime(60 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() == 502
			},
		)
}

func (clientResty *ClientResty) setDebug(debug bool) {
	clientResty.Client.SetDebug(debug)
}

func (clientResty *ClientResty) SetClientLog(l resty.Logger) {
	clientResty.Client.SetLogger(l)
}
