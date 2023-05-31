// Created at 2023-05-31

package model

import "database/sql/driver"

type Payway uint8

const (
	PaywayUnknown            Payway = iota // 未知
	PaywayCash                             // 现金支付
	PaywayAgentWxMiniprogram               // 代理商小程序
)

func (p *Payway) Scan(src interface{}) error {
	*p = Payway(src.(int64))
	return nil
}

func (p Payway) Value() (driver.Value, error) {
	return uint8(p), nil
}
