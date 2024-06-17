// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-22
// Based on aurservd by liasica, magicrolan@qq.com.

package baidu

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/auroraride/adapter/log"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
)

const (
	mapBaseUrl = "https://api.map.baidu.com"
	ridingUrl  = `/direction/v2/riding?ak=%s&destination=%s&origin=%s`
	rectifyUrl = `/rectify/v1/track`
)

type mapClient struct {
	ak     string
	sk     string
	apiKey string
}

func NewMap() *mapClient {
	cfg := ar.Config.Baidu.Map
	return &mapClient{
		ak:     cfg.Ak,
		sk:     cfg.Sk,
		apiKey: cfg.ApiKey,
	}
}

const (
	MapCoordType = "gcj02"
)

// var keys = map[string]string{
// 	"%21": "!",
// 	"%23": "#",
// 	"%24": "$",
// 	"%25": "%",
// 	"%26": "&",
// 	"%27": "'",
// 	"%28": "(",
// 	"%29": ")",
// 	"%2A": "*",
// 	"%2B": "+",
// 	"%2C": ",",
// 	"%2F": "/",
// 	"%3A": ":",
// 	"%3B": ";",
// 	"%3D": "=",
// 	"%3F": "?",
// 	"%40": "@",
// 	"%5B": "[",
// 	"%5D": "]",
// }

func (c *mapClient) getSignedUrl(u string) string {
	x, _ := url.Parse(u)
	v := x.Query()
	str := x.Path + "?" + v.Encode() + c.sk
	sn := utils.Md5String(url.QueryEscape(str))
	v.Add("sn", sn)
	result := fmt.Sprintf("%s%s?%s", mapBaseUrl, x.Path, v.Encode())
	return result
}

func (c *mapClient) getSignedSn(u string) string {
	x, _ := url.Parse(u)
	v := x.Query()
	str := x.Path + "?" + v.Encode() + c.sk
	sn := utils.Md5String(url.QueryEscape(str))
	return sn
}

type MapRes[T any] struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Info    struct {
		Copyright struct {
			Text     string `json:"text,omitempty"`
			ImageUrl string `json:"imageUrl,omitempty"`
		} `json:"copyright,omitempty"`
	} `json:"info,omitempty"`
	Type   int `json:"type,omitempty"`
	Result T   `json:"result,omitempty"`
}

type RidingRoute struct {
	RestrictionsStatus int    `json:"restrictions_status,omitempty"`
	RestrictionsInfo   string `json:"restrictions_info,omitempty"`
	Distance           int    `json:"distance,omitempty"`
	Duration           int    `json:"duration,omitempty"`
	Steps              []struct {
		LegIndex           int           `json:"leg_index,omitempty"`
		Area               int           `json:"area,omitempty"`
		Direction          int           `json:"direction,omitempty"`
		Distance           int           `json:"distance,omitempty"`
		Duration           int           `json:"duration,omitempty"`
		Instructions       string        `json:"instructions,omitempty"`
		Name               string        `json:"name,omitempty"`
		Path               string        `json:"path,omitempty"`
		Pois               []interface{} `json:"pois,omitempty"`
		Type               int           `json:"type,omitempty"`
		TurnType           string        `json:"turn_type,omitempty"`
		RestrictionsInfo   string        `json:"restrictions_info,omitempty"`
		StepOriginLocation struct {
			Lng float64 `json:"lng,omitempty"`
			Lat float64 `json:"lat,omitempty"`
		} `json:"stepOriginLocation,omitempty"`
		StepDestinationLocation struct {
			Lng float64 `json:"lng,omitempty"`
			Lat float64 `json:"lat,omitempty"`
		} `json:"stepDestinationLocation,omitempty"`
		StepOriginInstruction      string `json:"stepOriginInstruction,omitempty"`
		StepDestinationInstruction string `json:"stepDestinationInstruction,omitempty"`
		RestrictionsStatus         int    `json:"restrictions_status,omitempty"`
		Links                      []struct {
			Length int `json:"length,omitempty"`
			Attr   int `json:"attr,omitempty"`
		} `json:"links,omitempty"`
	} `json:"steps,omitempty"`
	OriginLocation struct {
		Lng float64 `json:"lng,omitempty"`
		Lat float64 `json:"lat,omitempty"`
	} `json:"originLocation,omitempty"`
	DestinationLocation struct {
		Lng float64 `json:"lng,omitempty"`
		Lat float64 `json:"lat,omitempty"`
	} `json:"destinationLocation,omitempty"`
}

type Riding struct {
	Routes []RidingRoute `json:"routes,omitempty"`
	Origin struct {
		AreaId   int    `json:"area_id,omitempty"`
		Cname    string `json:"cname,omitempty"`
		Uid      string `json:"uid,omitempty"`
		Wd       string `json:"wd,omitempty"`
		OriginPt struct {
			Lng float64 `json:"lng,omitempty"`
			Lat float64 `json:"lat,omitempty"`
		} `json:"originPt,omitempty"`
	} `json:"origin,omitempty"`
	Destination struct {
		AreaId        int    `json:"area_id,omitempty"`
		Cname         string `json:"cname,omitempty"`
		Uid           string `json:"uid,omitempty"`
		Wd            string `json:"wd,omitempty"`
		DestinationPt struct {
			Lng float64 `json:"lng,omitempty"`
			Lat float64 `json:"lat,omitempty"`
		} `json:"destinationPt,omitempty"`
	} `json:"destination,omitempty"`
}

type Rectify struct {
	Status           int      `json:"status"`
	Message          string   `json:"message"`
	Total            int      `json:"total"`
	Distance         float64  `json:"distance"`
	TollDistance     int      `json:"toll_distance"`
	LowSpeedDistance int      `json:"low_speed_distance"`
	Points           []Points `json:"points"`
}
type Points struct {
	LocTime   int     `json:"loc_time"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
	Direction float64 `json:"direction"`
}

type Point struct {
	LocTime        int32   `json:"loc_time"`
	Longitude      float64 `json:"longitude"`
	Latitude       float64 `json:"latitude"`
	CoordTypeInput string  `json:"coord_type_input"`
}

type Track []Point

func (c *mapClient) Riding(origin, destination string) (*Riding, error) {
	u := fmt.Sprintf(ridingUrl+"&riding_type=1&timestamp=%d", c.ak, destination, origin, time.Now().Unix())
	var res MapRes[*Riding]
	r, err := resty.New().R().Get(c.getSignedUrl(u))
	if err != nil {
		return nil, err
	}
	_ = jsoniter.Unmarshal(r.Body(), &res)
	if res.Status != 0 {
		return nil, errors.New(res.Message)
	}
	return res.Result, nil
}

func (c *mapClient) RidingX(origin, destination string) *Riding {
	r, err := c.Riding(origin, destination)
	if err != nil {
		snag.Panic(err)
	}
	return r
}

func (c *mapClient) RidingPlan(origin, destination string) (seconds int, distance int, polylines []string, err error) {
	var r *Riding
	r, err = c.Riding(origin, destination)
	if err != nil {
		return
	}
	for _, route := range r.Routes {
		seconds += route.Duration
		distance += route.Distance
		for _, step := range route.Steps {
			polylines = append(polylines, step.Path)
		}
	}
	return
}

func (c *mapClient) RidingPlanX(origin, destination string) (seconds int, distance int, polylines []string) {
	var err error
	seconds, distance, polylines, err = c.RidingPlan(origin, destination)
	if err != nil {
		snag.Panic(err)
	}
	return
}

// TrackRectify 轨迹纠偏
func (c *mapClient) TrackRectify(req Track) (res *model.BatteryTrackRes, err error) {
	uri := fmt.Sprintf(mapBaseUrl + rectifyUrl)
	rectify := new(Rectify)
	res = new(model.BatteryTrackRes)
	// 纠偏参数
	option := "need_mapmatch:1|transport_mode:riding|denoise_grade:1|vacuate_grade:1"
	pointList, err := jsoniter.Marshal(req)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf(uri+"?point_list=%s&rectify_option=%s&coord_type_output=%s&ak=%s", string(pointList), option, MapCoordType, c.ak)
	var r *resty.Response
	r, err = resty.New().R().
		SetFormData(map[string]string{
			"point_list":        string(pointList),
			"rectify_option":    option,
			"ak":                c.ak,
			"sn":                c.getSignedSn(u),
			"coord_type_output": MapCoordType,
		}).
		Post(uri)
	if err != nil {
		return nil, err
	}
	zap.L().Info("获取轨迹纠偏", log.ResponseBody(r.Body()))
	_ = jsoniter.Unmarshal(r.Body(), &rectify)
	if rectify == nil || rectify != nil && rectify.Status != 0 {
		return nil, fmt.Errorf("获取轨迹纠偏失败: %s", rectify.Message)
	}

	points := make([]model.BatteryPoint, 0, len(rectify.Points))
	for _, l := range rectify.Points {
		points = append(points,
			model.BatteryPoint{
				LocTime:   int32(l.LocTime),
				Longitude: l.Longitude,
				Latitude:  l.Latitude,
			})
	}
	res.Points = points
	return res, nil
}
