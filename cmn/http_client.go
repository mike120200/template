package cmn

import (
	neturl "net/url"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// SendHttpRequest 使用 fasthttp 发送 HTTP 请求（可指定方法），默认 Content-Type 为 application/json
// 步骤：0) 校验 URL 1) 获取请求/响应对象 2) 设置方法/URL 3) 设置基础请求头 3.1) 追加自定义请求头 4) 写入请求体 5) 发送请求(支持超时) 6) 读取响应 7) 非2xx返回错误 8) 返回响应体
// 成功返回响应体（2xx 视为成功），否则返回 CodeError，Code 为 HTTP 状态码
// headers 参数用于附加自定义请求头（若与基础头冲突，将覆盖基础头）
// timeout 为超时时间，若 <= 0 则默认使用 5s
func SendHttpRequest(method, url string, body []byte, headers map[string]string, timeout time.Duration) ([]byte, error) {
	// 0) 校验 URL 前缀与合法性
	u, err := neturl.ParseRequestURI(url)
	if err != nil {
		Logger().Error("Invalid URL", zap.Error(err), zap.String("url", url))
		return nil, NewAppError(CommonError, "invalid url: "+err.Error())
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		Logger().Error("Unsupported URL scheme", zap.String("url", url))
		return nil, NewAppError(CommonError, "unsupported url scheme: "+u.Scheme)
	}
	Logger().Info("Send HTTP request", zap.String("method", method), zap.String("url", url))
	// 1) 从对象池获取 Request/Response
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	// 1.1) 函数结束时释放对象，归还到池中
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// 2) 设置 HTTP 方法与 URL
	req.Header.SetMethod(method)
	req.SetRequestURI(url)

	// 3) 设置基础请求头（User-Agent / Content-Type）
	req.Header.SetContentType("application/json; charset=utf-8")

	// 3.1) 追加自定义请求头（会覆盖基础头）
	for k, v := range headers {
		if k == "" {
			continue
		}
		req.Header.Set(k, v)
	}

	// 4) 写入请求体（可为空）
	if body != nil {
		req.SetBodyRaw(body)
	} else {
		req.SetBodyRaw([]byte{})
	}

	// 5) 发送请求（支持超时）
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	if err := fasthttp.DoTimeout(req, resp, timeout); err != nil {
		Logger().Error("HTTP request error", zap.Error(err), zap.String("url", url))
		return nil, err
	}

	// 6) 读取状态码并复制响应体（复制后可安全释放 resp）
	status := resp.StatusCode()
	respBody := append([]byte(nil), resp.Body()...)

	// 7) 非 2xx 返回 CodeError（将响应体作为错误消息）
	if status < 200 || status >= 300 {
		Logger().Error("HTTP request error", zap.Int("status", status), zap.String("url", url))
		return nil, NewAppError(status, string(respBody))
	}
	Logger().Info("HTTP request success", zap.Int("status", status), zap.String("url", url), zap.String("body", string(respBody)))
	// 8) 返回响应体
	return respBody, nil
}
