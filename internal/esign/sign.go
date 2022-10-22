// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/pkg/snag"
    "time"
)

type CreateFlowOneStepRes struct {
    FlowId string `json:"flowId"`
}

type CreateFlowReq struct {
    Scene           string
    FileId          string
    PersonAccountId string
    PsnSignBean     PosBean
    EntSignBean     PosBean
}

// CreateFlowOneStep 一步发起签署
func (e *Esign) CreateFlowOneStep(data CreateFlowReq) string {
    if e.sn == "" {
        snag.Panic("合同编号为空")
    }
    // 定义跳转url/app
    e.redirect = e.Config.Redirect + "/" + e.sn
    var (
        scene           = data.Scene
        fileId          = data.FileId
        personAccountId = data.PersonAccountId
        psnSignBean     = data.PsnSignBean
        entSignBean     = data.EntSignBean
    )

    req := CreateFlowOneStepReq{
        FlowInfo: FlowInfo{
            AutoArchive:   true,
            AutoInitiate:  true,
            BusinessScene: scene,
            SignValidity:  time.Now().Add(model.ContractExpiration * time.Minute).UnixMilli(),
            FlowConfigInfo: FlowConfigInfo{
                RedirectDelayTime:        0,
                NoticeDeveloperUrl:       e.Config.Callback,
                RedirectUrl:              e.redirect,
                SignPlatform:             "1",
                NoticeType:               "",
                PersonAvailableAuthTypes: []string{"PSN_FACEAUTH_BYURL"},
                WillTypes:                []string{"FACE_ZHIMA_XY"},
            },
        },
        Signers: []Signer{
            {
                SignOrder: 2,
                SignerAccount: SignerAccount{
                    SignerAccountId: personAccountId,
                    NoticeType:      "",
                },
                Signfields: []Signfield{
                    {
                        FileId:   fileId,
                        PosBean:  psnSignBean,
                        SealType: "0",
                    },
                },
            },
            {
                PlatformSign: true,
                SignOrder:    1,
                Signfields: []Signfield{
                    {
                        AutoExecute:        true,
                        ActorIndentityType: 2,
                        FileId:             fileId,
                        PosBean:            entSignBean,
                    },
                },
            },
        },
        Docs: []Doc{
            {
                FileId: fileId,
            },
        },
    }

    res := new(CreateFlowOneStepRes)
    e.request(createFlowOneStepUrl, methodPost, req, res)
    return res.FlowId
}

// executeUrlRes 获取签署地址返回
type executeUrlRes struct {
    Url      string `json:"url"`
    ShortUrl string `json:"shortUrl"`
}

// ExecuteUrl 获取签署地址
// @doc https://open.esign.cn/doc/detail?id=opendoc%2Fsaas_api%2Ffdtfqf&namespace=opendoc%2Fsaas_api
// @doc https://open.esign.cn/doc/detail?id=opendoc%2Fpaas_api%2Fevvsef&namespace=opendoc%2Fpaas_api
func (e *Esign) ExecuteUrl(flowId, accountId string, params ...string) string {
    res := new(executeUrlRes)
    redirect := e.redirect
    if redirect == "" && len(params) > 0 {
        redirect = params[0]
    }
    url := fmt.Sprintf(executeUrl, flowId, accountId, redirect)
    e.request(url, methodGet, nil, res)
    return res.ShortUrl
}

type signResult struct {
    ProcessId                    string      `json:"processId,omitempty"`
    ContractNo                   interface{} `json:"contractNo,omitempty"`
    FlowId                       string      `json:"flowId,omitempty"`
    AppId                        string      `json:"appId,omitempty"`
    AppName                      interface{} `json:"appName,omitempty"`
    AutoArchive                  bool        `json:"autoArchive,omitempty"`
    FlowStatus                   int         `json:"flowStatus,omitempty"`
    FlowDesc                     string      `json:"flowDesc,omitempty"`
    FlowStartTime                int64       `json:"flowStartTime,omitempty"`
    FlowEndTime                  int64       `json:"flowEndTime,omitempty"`
    BusinessScene                string      `json:"businessScene,omitempty"`
    InitiatorClient              string      `json:"initiatorClient,omitempty"`
    InitiatorAccountId           string      `json:"initiatorAccountId,omitempty"`
    InitiatorAuthorizedAccountId string      `json:"initiatorAuthorizedAccountId,omitempty"`
    PayerAccountId               string      `json:"payerAccountId,omitempty"`
    SignValidity                 interface{} `json:"signValidity,omitempty"`
    ContractValidity             interface{} `json:"contractValidity,omitempty"`
    ContractEffective            interface{} `json:"contractEffective,omitempty"`
    ContractRemind               interface{} `json:"contractRemind,omitempty"`
    ConfigInfo                   struct {
        SignPlatform       string      `json:"signPlatform,omitempty"`
        NoticeType         string      `json:"noticeType,omitempty"`
        NoticeDeveloperUrl interface{} `json:"noticeDeveloperUrl,omitempty"`
        RedirectUrl        string      `json:"redirectUrl,omitempty"`
        ArchiveLock        bool        `json:"archiveLock,omitempty"`
    } `json:"configInfo,omitempty"`
}

// Result 签署结果查询
// @doc https://open.esign.cn/doc/detail?id=opendoc%2Fpaas_api%2Fghywlg&namespace=opendoc%2Fpaas_api
func (e *Esign) Result(flowId string) model.ContractStatus {
    var res = new(signResult)
    e.request(fmt.Sprintf(signResultUrl, flowId), methodGet, nil, res)
    return model.ContractStatus(res.FlowStatus)
}
