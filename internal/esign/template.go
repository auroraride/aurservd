// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/15
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

import (
    "fmt"
)

// DocTemplateRes 模板文件返回结构体
type DocTemplateRes struct {
    TemplateId       string  `json:"templateId"`
    FileSize         float64 `json:"fileSize"`
    TemplateName     string  `json:"templateName"`
    TemplateType     float64 `json:"templateType"`
    CreateTime       int64   `json:"createTime"`
    UpdateTime       int64   `json:"updateTime"`
    StructComponents []struct {
        Id      string  `json:"id"`
        Key     string  `json:"key"`
        Type    float64 `json:"type"`
        Context struct {
            Label    string `json:"label"`
            Limit    string `json:"limit"`
            Required bool   `json:"required"`
            Style    struct {
                Font      float64 `json:"font"`
                FontSize  float64 `json:"fontSize"`
                TextColor string  `json:"textColor"`
                Width     float64 `json:"width"`
                Height    float64 `json:"height"`
            } `json:"style"`
            Pos struct {
                X    float64 `json:"x"`
                Y    float64 `json:"y"`
                Page float64 `json:"page"`
            } `json:"pos"`
        } `json:"context"`
    } `json:"structComponents"`
}

// CreateByTemplateReq 填充内容生成PDF请求体
type CreateByTemplateReq struct {
    Name             string                 `json:"name"`
    SimpleFormFields map[string]interface{} `json:"simpleFormFields"`
    TemplateId       string                 `json:"templateId"`
}

// CreateByTemplateRes 填充内容生成PDF返回体
type CreateByTemplateRes struct {
    DownloadUrl string `json:"downloadUrl"`
    FileId      string `json:"fileId"`
    FileName    string `json:"fileName"`
}

// DocTemplate 查询模板文件详情
func (e *Esign) DocTemplate(templateId string) *DocTemplateRes {
    res := new(DocTemplateRes)
    e.request(fmt.Sprintf(docTemplateUrl, templateId), methodGet, nil, res)
    return res
}

// CreateByTemplate 填充内容生成PDF
func (e *Esign) CreateByTemplate(req CreateByTemplateReq) *CreateByTemplateRes {
    res := new(CreateByTemplateRes)
    e.request(createByTemplateUrl, methodPost, req, res)
    return res
}
