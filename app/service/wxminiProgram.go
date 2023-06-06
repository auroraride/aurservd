package service

import (
	"context"
	"fmt"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
)

type miniProgramService struct {
	*BaseService
	MiniProgram *miniprogram.MiniProgram
}

func NewminiProgram(params ...any) *miniProgramService {
	redisOpts := &cache.RedisOpts{
		Host:     ar.Config.Redis.Address,
		Password: ar.Config.Redis.Password,
		Database: ar.Config.Redis.DB,
	}
	wc := wechat.NewWechat()
	redisCache := cache.NewRedis(context.Background(), redisOpts)
	wc.SetCache(redisCache)

	if ar.Config.WechatMiniprogram.Agent.AppID == "" || ar.Config.WechatMiniprogram.Agent.AppSecret == "" {
		snag.Panic("微信小程序配置为空")
	}
	miniProgram := wc.GetMiniProgram(&config.Config{
		AppID:     ar.Config.WechatMiniprogram.Agent.AppID,
		AppSecret: ar.Config.WechatMiniprogram.Agent.AppSecret,
	})
	return &miniProgramService{
		BaseService: newService(params...),
		MiniProgram: miniProgram,
	}
}

// GetPhoneNumber 通过code换取手机号码
func (s *miniProgramService) GetPhoneNumber(code string) string {
	resultPhone, err := s.MiniProgram.GetAuth().GetPhoneNumber(code)
	if err != nil || resultPhone.ErrCode != 0 {
		zap.L().Error("获取手机号码失败"+ar.Config.WechatMiniprogram.Agent.AppID+" "+"ar.Config.WechatMiniprogram.Agent.AppSecret", zap.Error(err))
		snag.Panic("获取手机号码失败")
	}
	phoneNumber := resultPhone.PhoneInfo.PhoneNumber
	return phoneNumber
}

// GetAuth 通过code换取openid
func (s *miniProgramService) GetAuth(code string) string {
	session, err := s.MiniProgram.GetAuth().Code2Session(code)
	if err != nil || session.ErrCode != 0 {
		snag.Panic("获取openid失败")
	}
	return session.OpenID

}

// Invite 邀请骑手二维码
func (s *miniProgramService) Invite(enterprise *ent.Enterprise, req *model.EnterpriseRiderInviteReq) []byte {
	url := fmt.Sprintf("s=%d&e=%d", req.StationID, enterprise.ID)
	coderParam := qrcode.QRCoder{Scene: url, Path: "pages/rider-login/index"}
	code, err := s.MiniProgram.GetQRCode().GetWXACodeUnlimit(coderParam)
	if err != nil {
		snag.Panic("生成二维码失败")
	}
	return code
}
