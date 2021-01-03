package http_base

import "github.com/go-resty/resty/v2"

var CleanSession = func(client *resty.Client, request *resty.Request) error {
	client.SetCookieJar(nil)
	return nil
}
