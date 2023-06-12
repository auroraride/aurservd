// Created at 2023-06-12

package model

type AgentSubscribeActiveReq struct {
	ID           uint64 `json:"id" param:"id" validate:"required"` // 骑手ID
	BatteryID    uint64 `json:"batteryId"`                         // 电池ID
	EnterpriseID uint64 `json:"enterpriseId"`                      // 团签id
}
