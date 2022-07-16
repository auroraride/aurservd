// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package workwx

type Client struct {
    agentID    int64
    corpID     string
    corpSecret string
}

const (
    baseURL = `https://qyapi.weixin.qq.com/cgi-bin`
)

type response interface {
    IsSuccess() bool
    TokenInValid() bool
    Message() string
}

type baseResponse struct {
    Errcode int    `json:"errcode"`
    Errmsg  string `json:"errmsg"`
}

func New(agentID int64, corpID, corpSecret string) (w *Client) {
    w = &Client{
        agentID:    agentID,
        corpID:     corpID,
        corpSecret: corpSecret,
    }
    return
}

func (b *baseResponse) IsSuccess() bool {
    return b.Errcode == 0 && b.Errmsg == "ok"
}

func (b *baseResponse) TokenInValid() bool {
    return b.Errcode == 40014
}

func (b *baseResponse) Message() string {
    return b.Errmsg
}
