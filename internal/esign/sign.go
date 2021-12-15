// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

import "fmt"

type CreateFlowOneStepRes struct {
    FlowId string `json:"flowId"`
}

// CreateFlowOneStep 一步发起签署
func (e *esign) CreateFlowOneStep() string {
    var (
        scene           = "签署测试"
        fileId          = "3f8a6de98af44f20b836985e81ce942c"
        personAccountId = "f8c1df05edb047f4869a4d24749d76df"
        psnSignBean     = PosBean{
            PosPage: "2",
            PosX:    376.15,
            PosY:    441.3,
        }
        entSignBean = PosBean{
            PosPage: "2",
            PosX:    155.65,
            PosY:    436.8,
        }
    )

    req := CreateFlowOneStepReq{
        FlowInfo: FlowInfo{
            AutoArchive:   true,
            AutoInitiate:  true,
            BusinessScene: scene,
            // ContractValidity: time.Now().Add(30*time.Minute).UnixNano() / 1e6, // 半小时有效期
            FlowConfigInfo: FlowConfigInfo{
                NoticeDeveloperUrl: e.config.Callback,
                RedirectUrl:        fmt.Sprintf("%s?path=%s", e.config.Redirect, "test"),
                // RedirectDelayTime: 0,
            },
        },
        Signers: []Signer{
            {
                SignOrder:     2,
                SignerAccount: SignerAccount{SignerAccountId: personAccountId},
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
