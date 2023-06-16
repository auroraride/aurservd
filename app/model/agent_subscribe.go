// Created at 2023-06-12

package model

type AgentSubscribeActiveReq struct {
	ID        *uint64 `json:"id" validate:"required" trans:"订阅ID"`
	BatteryID *uint64 `json:"batteryId"` // 电池ID
	EbikeID   *uint64 `json:"ebikeId"`   // 电车ID
}
