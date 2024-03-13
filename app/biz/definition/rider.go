package definition

import "github.com/auroraride/aurservd/internal/baidu"

// RiderDirectionReq 骑手路径规划请求
type RiderDirectionReq struct {
	Origin      string `json:"origin" query:"origin" validate:"required"`           // 起点
	Destination string `json:"destination" query:"destination" validate:"required"` // 终点
}

// RiderDirectionRes 骑手路径规划响应
type RiderDirectionRes struct {
	Routes      []baidu.Routes    `json:"routes"`      // 方案列表
	Origin      baidu.Origin      `json:"origin"`      // 起点信息
	Destination baidu.Destination `json:"destination"` // 终点信息
}
