// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

import "fmt"

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
            // ContractValidity: time.Now().Add(30*time.Minute).UnixNano() / 1e6, // 半小时有效期
            FlowConfigInfo: FlowConfigInfo{
                NoticeDeveloperUrl: e.Config.Callback,
                RedirectUrl:        fmt.Sprintf("%s?path=%s", e.Config.Redirect, "test"),
                // RedirectDelayTime: 0,
            },
        },
        Signers: []Signer{
            {
                SignOrder: 2,
                SignerAccount: SignerAccount{
                    SignerAccountId: personAccountId,
                    NoticeType:      "1",
                },
                Signfields: []Signfield{
                    {
                        FileId:  fileId,
                        PosBean: psnSignBean,
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

type executeUrlRes struct {
    Url      string `json:"url"`
    ShortUrl string `json:"shortUrl"`
}

// ExecuteUrl 获取签署地址
// @doc https://open.esign.cn/doc/detail?id=opendoc%2Fsaas_api%2Ffdtfqf&namespace=opendoc%2Fsaas_api
func (e *Esign) ExecuteUrl(flowId, accountId string) string {
    res := new(executeUrlRes)
    e.request(fmt.Sprintf(executeUrl, flowId, accountId), methodGet, nil, res)
    return res.Url
}
