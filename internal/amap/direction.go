// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-23
// Based on aurservd by liasica, magicrolan@qq.com.

package amap

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "gopkg.in/resty.v1"
    "strconv"
)

type DirectionRidingRes struct {
    Status   string `json:"status,omitempty"`
    Info     string `json:"info,omitempty"`
    Infocode string `json:"infocode,omitempty"`
    Count    string `json:"count,omitempty"`
    Route    struct {
        Origin      string `json:"origin,omitempty"`
        Destination string `json:"destination,omitempty"`
        Paths       []struct {
            Distance string `json:"distance,omitempty"`
            Duration string `json:"duration,omitempty"`
            Steps    []struct {
                Instruction  string `json:"instruction,omitempty"`
                Orientation  string `json:"orientation,omitempty"`
                RoadName     string `json:"road_name,omitempty"`
                StepDistance int    `json:"step_distance,omitempty"`
                Cost         struct {
                    Duration string `json:"duration,omitempty"`
                } `json:"cost,omitempty"`
            } `json:"steps,omitempty"`
        } `json:"paths,omitempty"`
    } `json:"route,omitempty"`
}

func (a *amap) DirectionRiding(origin, destination LngLat) (res *DirectionRidingRes) {
    r, err := resty.R().SetBody(res).Get(fmt.Sprintf(
        `https://restapi.amap.com/v5/direction/electrobike?key=%s&origin=%f,%f&destination=%f,%f&show_fields=cost`,
        a.Key,
        origin.Lng,
        origin.Lat,
        destination.Lng,
        destination.Lat,
    ))
    if err != nil {
        log.Error(err)
    }
    log.Info(string(r.Body()))
    return
}

// DirectionRidingMinutes 骑行规划
func (a *amap) DirectionRidingMinutes(origin, destination LngLat) (total int) {
    res := a.DirectionRiding(origin, destination)
    if res == nil {
        return
    }
    for _, path := range res.Route.Paths {
        for _, step := range path.Steps {
            cost, _ := strconv.Atoi(step.Cost.Duration)
            total += cost
        }
    }
    return
}
