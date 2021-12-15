// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/15
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

import "fmt"

// docTemplateRes 模板文件返回结构体
type docTemplateRes struct {
    TemplateId       string `json:"templateId"`
    FileSize         int    `json:"fileSize"`
    TemplateName     string `json:"templateName"`
    TemplateType     int    `json:"templateType"`
    CreateTime       int64  `json:"createTime"`
    UpdateTime       int64  `json:"updateTime"`
    StructComponents []struct {
        Id      string `json:"id"`
        Key     string `json:"key"`
        Type    int    `json:"type"`
        Context struct {
            Label    string `json:"label"`
            Limit    string `json:"limit"`
            Required bool   `json:"required"`
            Style    struct {
                Font      int     `json:"font"`
                FontSize  int     `json:"fontSize"`
                TextColor string  `json:"textColor"`
                Width     float64 `json:"width"`
                Height    float64 `json:"height"`
            } `json:"style"`
            Pos struct {
                X    float64 `json:"x"`
                Y    float64 `json:"y"`
                Page int     `json:"page"`
            } `json:"pos"`
        } `json:"context"`
    } `json:"structComponents"`
}

// DocTemplate 查询模板文件详情
func (e *esign) DocTemplate() *docTemplateRes {
    res := new(docTemplateRes)
    e.request(fmt.Sprintf(docTemplateUrl, e.config.TemplateId), "GET", nil, res)
    return res
}
