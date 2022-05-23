// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-23
// Based on aurservd by liasica, magicrolan@qq.com.

package amap

import (
    "fmt"
    "github.com/go-resty/resty/v2"
    log "github.com/sirupsen/logrus"
)

type Geocode struct {
    FormattedAddress string        `json:"formatted_address"`
    Country          string        `json:"country"`
    Province         string        `json:"province"`
    Citycode         string        `json:"citycode"`
    City             string        `json:"city"`
    District         []interface{} `json:"district"`
    Township         []interface{} `json:"township"`
    Neighborhood     struct {
        Name []interface{} `json:"name"`
        Type []interface{} `json:"type"`
    } `json:"neighborhood"`
    Building struct {
        Name []interface{} `json:"name"`
        Type []interface{} `json:"type"`
    } `json:"building"`
    Adcode   string        `json:"adcode"`
    Street   []interface{} `json:"street"`
    Number   []interface{} `json:"number"`
    Location string        `json:"location"`
    Level    string        `json:"level"`
}

type GeoRes struct {
    Status   string    `json:"status"`
    Info     string    `json:"info"`
    Infocode string    `json:"infocode"`
    Count    string    `json:"count"`
    Geocodes []Geocode `json:"geocodes"`
}

func (a *amap) Geo(name string) (*Geocode, error) {
    res := new(GeoRes)
    _, err := resty.New().R().SetResult(res).Get(fmt.Sprintf("https://restapi.amap.com/v3/geocode/geo?key=%s&address=%s", a.Key, name))
    if err != nil {
        log.Error(err)
        return nil, err
    }
    out := res.Geocodes[0]
    return &out, nil
}
