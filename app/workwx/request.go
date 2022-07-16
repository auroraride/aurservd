// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package workwx

import (
    "errors"
    "fmt"
    "github.com/go-resty/resty/v2"
    "strings"
)

const (
    // MethodGet HTTP method
    MethodGet = "GET"

    // MethodPost HTTP method
    MethodPost = "POST"

    // MethodPut HTTP method
    MethodPut = "PUT"

    // MethodDelete HTTP method
    MethodDelete = "DELETE"

    // MethodPatch HTTP method
    MethodPatch = "PATCH"

    // MethodHead HTTP method
    MethodHead = "HEAD"

    // MethodOptions HTTP method
    MethodOptions = "OPTIONS"
)

type Request struct {
    *Client
    *resty.Request
}

func (w *Client) R() *Request {
    return &Request{
        w,
        resty.New().SetBaseURL(baseURL).R(),
    }
}

// RequestExecute 请求API
func (r *Request) RequestExecute(url string, method string, res response, params ...bool) (err error) {
    q := "?"
    if strings.Contains(url, "?") {
        q = "&"
    }

    _, err = r.SetResult(res).Execute(method, fmt.Sprintf(`%s%saccess_token=%s`, url, q, r.getAccessToken(params...)))
    if err != nil {
        return
    }

    if res.IsSuccess() {
        return nil
    }

    if res.TokenInValid() && (len(params) == 0 || !params[0]) {
        // 强制获取access_token
        return r.RequestExecute(url, method, res, true)
    }

    return errors.New(res.Message())
}

func (w *Client) RequestGet(url string, res response) error {
    return w.R().RequestExecute(url, MethodGet, res)
}

func (w *Client) RequestPost(url string, body any, res response) error {
    r := w.R()
    r.SetBody(body)
    return r.RequestExecute(url, MethodPost, res)
}
