package scikits

import (
	"crypto/tls"
	"net/http"
	"time"
)

// `这里请注意，使用 InsecureSkipVerify: true 来跳过证书验证`
func GetHttpClient() *http.Client {
	client := &http.Client{
		Timeout: 20 * time.Second,
		Transport: &http.Transport{
			// DisableKeepAlives 参数很关键，业务场景是服务A调用B,每次调用都会创建新的http client对象，这样每次请求都会创建一个连接，而我们的接口请求量很大，这样就会创建大量的连接，所以不能够启用连接池进行连接复用。
			DisableKeepAlives: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}
	return client
}
