package ali

import (
	"fmt"

	"github.com/auroraride/adapter/log"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/snag"
)

const (
	idverifyurl = "https://eid.shumaidata.com/eid/check?idcard=%s&name=%s"
)

type idVerify struct {
	AppCode string
}

func NewIdVerify() *idVerify {
	cfg := ar.Config.Aliyun.IdVerify
	return &idVerify{
		AppCode: cfg.AppCode,
	}
}

type IdVerifyRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Result  Result `json:"result"`
}

type Result struct {
	Name        string `json:"name"`
	IDCard      string `json:"idcard"`
	Res         string `json:"res"`
	Description string `json:"description"`
	Sex         string `json:"sex"`
	Birthday    string `json:"birthday"`
	Address     string `json:"address"`
}

func (i *idVerify) Verify(name, idCard string) (res *IdVerifyRes) {
	res = new(IdVerifyRes)
	r, err := resty.New().R().
		SetHeader("Authorization", "APPCODE "+i.AppCode).
		SetResult(res).
		Post(fmt.Sprintf(idverifyurl, idCard, name))
	if err != nil {
		snag.Panic(err)
	}

	if res.Code != "0" {
		zap.L().Info("校验身份证信息URL", log.ResponseBody(r.Body()))
		return
	}
	return
}
