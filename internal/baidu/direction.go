package baidu

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/internal/ar"
)

const (
	host string = "https://api.map.baidu.com"
	uri  string = "/direction/v2/riding"
)

type direction struct {
}

func NewDirection() *direction {
	return &direction{}
}

func (d *direction) GetDirection(origin, destination string) (res *definition.DirectionRes, err error) {
	res = &definition.DirectionRes{}
	ak := ar.Config.Baidu.Map.Ak
	sk := ar.Config.Baidu.Map.Sk

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	// 设置请求参数
	params := [][]string{
		{"origin", origin},
		{"destination", destination},
		{"ak", ak},
		{"timestamp", timestamp},
		{"riding_type", "1"}, // 电动车
	}

	paramsArr := make([]string, 0)
	for _, v := range params {
		kv := v[0] + "=" + (v[1])
		paramsArr = append(paramsArr, kv)
	}
	paramsStr := strings.Join(paramsArr, "&")

	// 计算sn
	queryStr := url.QueryEscape(uri + "?" + paramsStr)
	sn := calculateSN(queryStr, sk)

	requestURL := fmt.Sprintf("%s%s?%s&sn=%s", host, uri, paramsStr, sn)

	r, err := resty.New().R().SetResult(res).Get(requestURL)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(r.Body(), res)
	if err != nil {
		return nil, err
	}

	if res.Status != 0 {
		zap.L().Error("获取骑行路线失败", zap.Error(err))
		return nil, errors.New("获取骑行路线失败")
	}
	return
}

func calculateSN(queryStr string, sk string) string {
	str := queryStr + sk
	key := md5.Sum([]byte(str))
	sn := fmt.Sprintf("%x", key)
	return sn
}
