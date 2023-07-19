package model

import "path/filepath"

type Legal string

const (
	LegalAppPolicy        Legal = "policy"            // APP隐私政策
	LegalAppAgreement     Legal = "agreement"         // APP服务协议
	LegalAgentPolicy      Legal = "agent-policy"      // 代理端小程序隐私政策
	LegalAgentAgreement   Legal = "agent-agreement"   // 代理端小程序服务协议
	LegalPromotePolicy    Legal = "promote-policy"    // 推广端小程序隐私政策
	LegalPromoteAgreement Legal = "promote-agreement" // 推广端小程序服务协议
)

func (l Legal) Filepath() string {
	return filepath.Join("public", "pages", string(l)+".html")
}

func (l Legal) Title() string {
	switch l {
	case LegalAppPolicy:
		return "极光出行骑手APP隐私政策"
	case LegalAppAgreement:
		return "极光出行骑手APP服务协议"
	case LegalAgentPolicy:
		return "极光出行代理小程序隐私政策"
	case LegalAgentAgreement:
		return "极光出行代理小程序服务协议"
	case LegalPromotePolicy:
		return "极光出行推广小程序隐私政策"
	case LegalPromoteAgreement:
		return "极光出行推广小程序服务协议"
	}
	return " - "
}

var (
	Legals = []Legal{
		LegalAppPolicy,
		LegalAppAgreement,
		LegalAgentPolicy,
		LegalAgentAgreement,
		LegalPromotePolicy,
		LegalPromoteAgreement,
	}
)

type LegalName struct {
	// policy: APP隐私政策; agreement: APP服务协议; agent-policy: 代理端小程序隐私政策; agent-agreement: 代理端小程序服务协议; promote-policy: 推广端小程序隐私政策; promote-agreement: 推广端小程序服务协议
	Name Legal `json:"name" param:"name" query:"name" validate:"required" enums:"policy,agreement,agent-policy,agent-agreement,promote-policy,promote-agreement" trans:"名称"`
}

type LegalRes struct {
	Title   string `json:"title"`   // 法规项目名称
	Content string `json:"content"` // 内文
	Url     string `json:"url"`     // url <前端需要加入API前缀，例如`https://api.auroraride.com/`>
}

type LegalSaveReq struct {
	LegalName
	Content string `json:"content" validate:"required" trans:"政策内容"`
}
