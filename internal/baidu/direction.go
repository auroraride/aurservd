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

	"github.com/auroraride/aurservd/internal/ar"
)

const (
	host string = "https://api.map.baidu.com"
	uri  string = "/direction/v2/riding"
)

// DirectionRes 骑行路线规划响应
type DirectionRes struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Info    struct {
		Copyright struct {
			Text     string `json:"text"`
			ImageUrl string `json:"imageUrl"`
		} `json:"copyright"`
	} `json:"info"`
	Type   int `json:"type"`
	Result struct {
		Routes      []Routes    `json:"routes"`      // 方案列表
		Origin      Origin      `json:"origin"`      // 起点信息
		Destination Destination `json:"destination"` // 终点信息
	} `json:"result"`
}

// Routes 骑行路线方案
type Routes struct {
	RestrictionsStatus int    `json:"restrictions_status"` // 限行类型	0x01表示禁行；0x02表示逆行
	RestrictionsInfo   string `json:"restrictions_info"`   // 限行信息	如 "包含禁行路段|包含逆行路段"
	Distance           int    `json:"distance"`            // 方案距离 单位：米
	Duration           int    `json:"duration"`            // 单位：秒
	Steps              []struct {
		LegIndex           int           `json:"leg_index"`
		Area               int           `json:"area"`              // 文档未标注
		Direction          int           `json:"direction"`         // 当前道路方向角
		Distance           int           `json:"distance"`          // 路段距离	单位：米
		Duration           int           `json:"duration"`          // 路段耗时	单位：秒
		Instructions       string        `json:"instructions"`      // 路段描述	如“骑行50米“
		Name               string        `json:"name"`              // 该路段道路名称	如“信息路“ 若道路未命名或百度地图未采集到该道路名称，则返回"无名路"
		Path               string        `json:"path"`              // 路段位置坐标描述
		Pois               []interface{} `json:"pois"`              // 文档未标注  该路段途径的POI列表？
		Type               int           `json:"type"`              // 文档未标注
		TurnType           string        `json:"turn_type"`         // 行驶转向方向	如“直行”、“左前方转弯”
		RestrictionsInfo   string        `json:"restrictions_info"` // 限行信息	如 "包含禁行路段|包含逆行路段"
		StepOriginLocation struct {
			Lng float64 `json:"lng"` // 路段起点经度
			Lat float64 `json:"lat"` // 路段起点纬度
		} `json:"stepOriginLocation"` // 路段起点坐标
		StepDestinationLocation struct {
			Lng float64 `json:"lng"` // 路段终点经度
			Lat float64 `json:"lat"` // 路段终点经度
		} `json:"stepDestinationLocation"` // 路段终点坐标
		StepOriginInstruction      string `json:"stepOriginInstruction"`      // 路段起点描述？
		StepDestinationInstruction string `json:"stepDestinationInstruction"` // 路段终点描述？
		RestrictionsStatus         int    `json:"restrictions_status"`        // 限行类型	0x01表示禁行；0x02表示逆行
		Links                      []struct {
			Length int `json:"length"` // link长度	单位：米
			Attr   int `json:"attr"`   // link属性	0x01表示禁行；0x02表示逆行
		} `json:"links"` // link信息
	} `json:"steps"`
	OriginLocation struct {
		Lng float64 `json:"lng"` // 路线终点经度
		Lat float64 `json:"lat"` // 路线终点纬度
	} `json:"originLocation"` // 路线起点坐标
	DestinationLocation struct {
		Lng float64 `json:"lng"` // 路线终点经度
		Lat float64 `json:"lat"` // 路线终点纬度
	} `json:"destinationLocation"` // 路线终点坐标
}

// Origin 起点信息
type Origin struct {
	AreaId   int    `json:"area_id"` // 起点区域ID
	Cname    string `json:"cname"`   // 起点城市名称
	Uid      string `json:"uid"`     // 起点UID
	Wd       string `json:"wd"`      // 起点名称
	OriginPt struct {
		Lng float64 `json:"lng"` // 起点经度
		Lat float64 `json:"lat"` // 起点纬度
	} `json:"originPt"` // 起点坐标
}

// Destination 终点信息
type Destination struct {
	AreaId        int    `json:"area_id"` // 终点区域ID
	Cname         string `json:"cname"`   // 终点城市名称
	Uid           string `json:"uid"`     // 终点UID
	Wd            string `json:"wd"`      // 终点名称
	DestinationPt struct {
		Lng float64 `json:"lng"` // 终点经度
		Lat float64 `json:"lat"` // 终点纬度
	} `json:"destinationPt"` // 终点坐标
}

type direction struct {
}

func NewDirection() *direction {
	return &direction{}
}

func (d *direction) GetDirection(origin, destination string) (res *DirectionRes, err error) {
	res = &DirectionRes{}
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
