package service

import (
	"context"
	"fmt"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/credential"
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

type MiniProgramConfig struct {
	AppID     string
	AppSecret string
}

type WechatService struct {
	MiniPrograms map[string]*miniProgramService
	redisClient  *cache.Redis
}

func NewWechatService() *WechatService {
	redisOpts := &cache.RedisOpts{
		Host:     ar.Config.Redis.Address,
		Password: ar.Config.Redis.Password,
		Database: 0,
	}

	redisClient := cache.NewRedis(context.Background(), redisOpts)

	return &WechatService{
		MiniPrograms: make(map[string]*miniProgramService),
		redisClient:  redisClient,
	}
}

func (ws *WechatService) AddMiniProgram(name string, cfg *MiniProgramConfig, params ...any) {
	if cfg.AppID == "" || cfg.AppSecret == "" {
		snag.Panic("微信小程序配置为空")
	}

	wc := wechat.NewWechat()
	wc.SetCache(ws.redisClient)
	miniProgram := wc.GetMiniProgram(&config.Config{
		AppID:     cfg.AppID,
		AppSecret: cfg.AppSecret,
	})

	// 更换获取access_token的方式
	miniProgram.SetAccessTokenHandle(credential.NewStableAccessToken(cfg.AppID, cfg.AppSecret, credential.CacheKeyMiniProgramPrefix, ws.redisClient))

	ws.MiniPrograms[name] = &miniProgramService{
		BaseService: newService(params...),
		MiniProgram: miniProgram,
	}
}

func (ws *WechatService) GetMiniProgram(name string) *miniProgramService {
	return ws.MiniPrograms[name]
}

// NewminiProgram 代理小程序
func NewminiProgram(params ...any) *miniProgramService {
	ws := NewWechatService()
	ws.AddMiniProgram("agent", &MiniProgramConfig{
		AppID:     ar.Config.WechatMiniprogram.Agent.AppID,
		AppSecret: ar.Config.WechatMiniprogram.Agent.AppSecret,
	}, params...)
	return ws.GetMiniProgram("agent")
}

// NewPromotionMiniProgram 推广小程序
func NewPromotionMiniProgram(params ...any) *miniProgramService {
	ws := NewWechatService()
	ws.AddMiniProgram("promotion", &MiniProgramConfig{
		AppID:     ar.Config.WechatMiniprogram.Promotion.AppID,
		AppSecret: ar.Config.WechatMiniprogram.Promotion.AppSecret,
	}, params...)
	return ws.GetMiniProgram("promotion")
}

// NewRiderMiniProgram 骑手小程序
func NewRiderMiniProgram(params ...any) *miniProgramService {
	ws := NewWechatService()
	ws.AddMiniProgram("promotion", &MiniProgramConfig{
		AppID:     ar.Config.WechatMiniprogram.Rider.AppID,
		AppSecret: ar.Config.WechatMiniprogram.Rider.AppSecret,
	}, params...)
	return ws.GetMiniProgram("rider")
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
func (s *miniProgramService) GetAuth(code string) model.OpenidRes {
	session, err := s.MiniProgram.GetAuth().Code2Session(code)
	if err != nil || session.ErrCode != 0 {
		zap.L().Error("获取openid失败", zap.Error(err))
		snag.Panic("获取openid失败")
	}
	return model.OpenidRes{Openid: session.OpenID}
}

// InviteQrcode 邀请骑手二维码
func (s *miniProgramService) InviteQrcode(enterprise *ent.Enterprise, req *model.EnterpriseRiderInviteReq) []byte {
	url := fmt.Sprintf("s=%d&e=%d", req.StationID, enterprise.ID)
	coderParam := qrcode.QRCoder{Scene: url, Page: "pages/rider-login/index"}
	code, err := s.MiniProgram.GetQRCode().GetWXACodeUnlimit(coderParam)
	if err != nil {
		zap.L().Error("生成二维码失败", zap.Error(err))
		snag.Panic("生成二维码失败")
	}
	return code
}

// PromotionQrcode 获取推广返佣二维码
func (s *miniProgramService) PromotionQrcode(id uint64) []byte {
	url := fmt.Sprintf("m=%d", id)

	envVersion := "release"
	if ar.Config.Environment == "environment" {
		envVersion = "trial"
	}

	coderParam := qrcode.QRCoder{Scene: url, EnvVersion: envVersion, Page: "pages/login/index"}
	code, err := s.MiniProgram.GetQRCode().GetWXACodeUnlimit(coderParam)
	if err != nil {
		zap.L().Error("生成二维码失败", zap.Error(err))
		snag.Panic("生成二维码失败")
	}
	return code
}
