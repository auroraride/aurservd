// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-23
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/assets"
    "github.com/auroraride/aurservd/internal/amap"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/city"
    jsoniter "github.com/json-iterator/go"
    "github.com/spf13/cobra"
    "io/ioutil"
    "log"
    "strconv"
    "strings"
)

var cityCmd = &cobra.Command{
    Use:   "city",
    Short: "城市助手",
}

func cityCenterCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "center",
        Short: "从json中更新城市中心点",
        Run: func(cmd *cobra.Command, args []string) {
            type City struct {
                Adcode   uint64  `json:"adcode"`
                Name     string  `json:"name"`
                Code     string  `json:"code"`
                Lng      float64 `json:"lng,omitempty"`
                Lat      float64 `json:"lat,omitempty"`
                Children []City  `json:"children,omitempty"`
            }
            var items []City
            err := jsoniter.Unmarshal(assets.City, &items)
            if err != nil {
                log.Fatal(err)
            }
            m := &model.Modifier{
                ID:    1,
                Name:  "初始管理员",
                Phone: "18888888888",
            }
            orm := ent.Database.City
            for _, item := range items {
                for _, child := range item.Children {
                    orm.UpdateOneID(child.Adcode).
                        SetLastModifier(nil).
                        SetLng(child.Lng).
                        SetLat(child.Lat).
                        SetLastModifier(m).
                        SaveX(context.Background())
                }
            }
        },
    }
}

func cityAmapCenterCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "centerx",
        Short: "获取城市中心点",
        Run: func(cmd *cobra.Command, args []string) {
            // 城市列表
            orm := ent.Database.City
            items := orm.Query().
                Where(city.ParentIDNotNil()).
                Where(city.Or(
                    city.LngIsNil(),
                    city.LatIsNil(),
                )).
                AllX(context.Background())
            for _, item := range items {
                res, err := amap.New().Geo(item.Name)
                if err != nil {
                    log.Fatal(err)
                }
                location := strings.Split(res.Location, ",")
                lng, _ := strconv.ParseFloat(location[0], 10)
                lat, _ := strconv.ParseFloat(location[1], 10)
                _, err = orm.UpdateOne(item).SetLng(lng).SetLat(lat).Save(context.Background())
                if err != nil {
                    log.Fatal(err)
                }
            }
        },
    }
    return cmd
}

func cityJsonCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "json",
        Short: "更新json文件",
        Run: func(cmd *cobra.Command, args []string) {
            type City struct {
                Adcode   uint64  `json:"adcode"`
                Name     string  `json:"name"`
                Code     string  `json:"code"`
                Lng      float64 `json:"lng,omitempty"`
                Lat      float64 `json:"lat,omitempty"`
                Children []City  `json:"children,omitempty"`
            }
            models := ent.Database.City.Query().WithChildren(func(cq *ent.CityQuery) {
                cq.Order(ent.Asc(city.FieldID))
            }).Order(ent.Asc(city.FieldID)).Where(city.ParentIDIsNil()).AllX(context.Background())
            items := make([]City, len(models))
            for i, model := range models {
                item := City{
                    Children: make([]City, len(model.Edges.Children)),
                    Adcode:   model.ID,
                    Name:     model.Name,
                    Code:     model.Code,
                }
                for n, child := range model.Edges.Children {
                    item.Children[n] = City{
                        Adcode: child.ID,
                        Name:   child.Name,
                        Code:   child.Code,
                        Lng:    child.Lng,
                        Lat:    child.Lat,
                    }
                }
                items[i] = item
            }
            b, _ := jsoniter.MarshalIndent(items, "", "  ")
            _ = ioutil.WriteFile("assets/city.json", b, 0755)
        },
    }
}

func init() {
    cityCmd.AddCommand(cityAmapCenterCmd())
    cityCmd.AddCommand(cityCenterCmd())
    cityCmd.AddCommand(cityJsonCmd())
}
