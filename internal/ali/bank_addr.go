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
	bankaddressurl = "https://bankaddress.shumaidata.com/bankaddress?bankcard=%s"
)

type bankAddr struct {
	AppCode string
}

func NewBankAddr() *bankAddr {
	cfg := ar.Config.Aliyun.BankAddr
	return &bankAddr{
		AppCode: cfg.AppCode,
	}
}

type BankDataRes struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    *Data  `json:"data"`
}

type Data struct {
	OrderNo      string `json:"order_no"`
	Bank         string `json:"bank"`
	Province     string `json:"province"`
	City         string `json:"city"`
	CardName     string `json:"card_name"`
	Tel          string `json:"tel"`
	Type         string `json:"type"`
	Logo         string `json:"logo"`
	Abbreviation string `json:"abbreviation"`
	CardBin      string `json:"card_bin"`
	BinDigits    int    `json:"bin_digits"`
	CardDigits   int    `json:"card_digits"`
	IsLuhn       bool   `json:"isLuhn"`
	WebURL       string `json:"weburl"`
}

// GetBankAddr 获取银行卡信息
func (i *bankAddr) GetBankAddr(bankcard string) (res *BankDataRes) {
	res = new(BankDataRes)
	r, err := resty.New().R().
		SetHeader("Authorization", "APPCODE "+i.AppCode).
		SetResult(res).
		Get(fmt.Sprintf(bankaddressurl, bankcard))
	if err != nil {
		snag.Panic(err)
	}

	if !res.Success {
		zap.L().Info("获取银行卡URL", log.ResponseBody(r.Body()))
		snag.Panic("获取银行卡信息失败")
	}
	return
}
