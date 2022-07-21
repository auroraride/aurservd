// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-23
// Based on aurservd by liasica, magicrolan@qq.com.

package amap

import (
    "encoding/json"
    "errors"
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

type Regeocode struct {
    AddressComponent struct {
        City     any    `json:"city,omitempty"`
        Province string `json:"province,omitempty"`
        Adcode   string `json:"adcode,omitempty"`
        District string `json:"district,omitempty"`
        Citycode string `json:"citycode,omitempty"`
    } `json:"addressComponent,omitempty"`
    FormattedAddress string `json:"formatted_address,omitempty"`
}

type ReGeoRes struct {
    Status    string    `json:"status,omitempty"`
    Regeocode Regeocode `json:"regeocode,omitempty"`
    Info      string    `json:"info,omitempty"`
    Infocode  string    `json:"infocode,omitempty"`
}

func (r *ReGeoRes) String() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (a *amap) ReGeo(lng, lat float64) (*Regeocode, error) {
    res := new(ReGeoRes)
    _, err := resty.New().R().SetResult(res).Get(fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=%s&location=%f,%f&poitype=&radius=&extensions=base&batch=false&roadlevel=0", a.Key, lng, lat))
    if err != nil {
        log.Error(err)
        return nil, err
    }
    if res.Infocode != "10000" {
        log.Error(res)
        return nil, errors.New(res.Info)
    }
    out := res.Regeocode
    return &out, nil
}
