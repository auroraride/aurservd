// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

// CreateFlowOneStepReq 一步发起签署 请求体
// @doc https://open.esign.cn/doc/detail?id=opendoc%2Fpaas_api%2Fpwd6l4&namespace=opendoc%2Fpaas_api&page=9eb4dd5a6a4573e9fb4d22fcbb536452_h3_4
type CreateFlowOneStepReq struct {
    Attachments []Attachment `json:"attachments,omitempty"`
    Copiers     []Copier     `json:"copiers,omitempty"`
    Docs        []Doc        `json:"docs,omitempty"`
    FlowInfo    FlowInfo     `json:"flowInfo,omitempty"`
    Signers     []Signer     `json:"signers,omitempty"`
}

type FlowInfo struct {
    AutoArchive                  bool           `json:"autoArchive,omitempty"`
    AutoInitiate                 bool           `json:"autoInitiate,omitempty"`
    BusinessScene                string         `json:"businessScene,omitempty"`
    ContractRemind               int64          `json:"contractRemind,omitempty"`
    ContractValidity             int64          `json:"contractValidity,omitempty"`
    FlowConfigInfo               FlowConfigInfo `json:"flowConfigInfo,omitempty"`
    InitiatorAccountId           string         `json:"initiatorAccountId,omitempty"`
    InitiatorAuthorizedAccountId string         `json:"initiatorAuthorizedAccountId,omitempty"`
    Remark                       string         `json:"remark,omitempty"`
    SignValidity                 string         `json:"signValidity,omitempty"`
}

type FlowConfigInfo struct {
    NoticeDeveloperUrl string `json:"noticeDeveloperUrl,omitempty"`
    NoticeType         string `json:"noticeType,omitempty"`
    RedirectUrl        string `json:"redirectUrl,omitempty"`
    SignPlatform       string `json:"signPlatform,omitempty"`
    RedirectDelayTime  int64  `json:"redirectDelayTime,omitempty"`
}

type Signer struct {
    PlatformSign  bool          `json:"platformSign,omitempty"`
    SignOrder     int64         `json:"signOrder,omitempty"`
    SignerAccount SignerAccount `json:"signerAccount,omitempty"`
    Signfields    []Signfield   `json:"signfields,omitempty"`
    ThirdOrderNo  string        `json:"thirdOrderNo,omitempty"`
}

type SignerAccount struct {
    SignerAccountId     string   `json:"signerAccountId,omitempty"`
    AuthorizedAccountId string   `json:"authorizedAccountId,omitempty"`
    NoticeType          string   `json:"noticeType,omitempty"`
    WillTypes           []string `json:"willTypes,omitempty"`
}

type Signfield struct {
    AutoExecute         bool         `json:"autoExecute,omitempty"`
    ActorIndentityType  int          `json:"actorIndentityType,omitempty"`
    FileId              string       `json:"fileId,omitempty"`
    Order               string       `json:"order,omitempty"`
    SealId              string       `json:"sealId,omitempty"`
    SealType            string       `json:"sealType,omitempty"`
    SignType            int64        `json:"signType,omitempty"`
    PosBean             PosBean      `json:"posBean,omitempty"`
    Width               int64        `json:"width,omitempty"`
    SignDateBeanType    int64        `json:"signDateBeanType,omitempty"`
    SignDateBean        SignDateBean `json:"signDateBean,omitempty"`
    AuthorizedAccountId string       `json:"authorizedAccountId,omitempty"`
    SignerAccountId     string       `json:"signerAccountId,omitempty"`
}

type SignDateBean struct {
    FontSize int64   `json:"fontSize,omitempty"`
    Format   string  `json:"format,omitempty"`
    PosPage  int64   `json:"posPage,omitempty"`
    PosX     float32 `json:"posX,omitempty"`
    PosY     float32 `json:"posY,omitempty"`
}

type PosBean struct {
    PosPage        string  `json:"posPage,omitempty"`
    PosX           float64 `json:"posX,omitempty"`
    PosY           float64 `json:"posY,omitempty"`
    AddSignTime    bool    `json:"addSignTime,omitempty"`
    SignTimeFormat string  `json:"signTimeFormat,omitempty"`
}

type Attachment struct {
    FileId         string `json:"fileId,omitempty"`
    AttachmentName string `json:"attachmentName,omitempty"`
}

type Copier struct {
    CopierAccountId                string `json:"copierAccountId,omitempty"`
    CopierIdentityAccountType      int64  `json:"copierIdentityAccountType,omitempty"`
    PrcopierIdentityAccountIdivate string `json:"copierIdentityAccountId,omitempty"`
}

type Doc struct {
    FileId       string `json:"fileId,omitempty"`
    FileName     string `json:"fileName,omitempty"`
    Encryption   int64  `json:"encryption,omitempty"`
    FilePassword string `json:"filePassword,omitempty"`
}
