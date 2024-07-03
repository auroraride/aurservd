package biz

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/riderphonedevice"
)

type riderPhoneDeviceBiz struct {
	orm *ent.RiderPhoneDeviceClient
}

func NewRiderPhoneDeviceBiz() *riderPhoneDeviceBiz {
	return &riderPhoneDeviceBiz{
		orm: ent.Database.RiderPhoneDevice,
	}
}

// ReportPhoneDevice  上报骑手手机设备信息
func (b *riderPhoneDeviceBiz) ReportPhoneDevice(ctx echo.Context, r *ent.Rider, req *definition.RiderPhoneDeviceReq) error {
	err := ent.Database.RiderPhoneDevice.
		Create().
		SetRiderID(r.ID).
		SetDeviceSn(req.DeviceSn).
		SetModel(req.Model).
		SetBrand(req.Brand).
		SetImei(req.Imei).
		SetOsVersion(req.OsVersion).
		SetOsName(req.OsName).
		SetScreenWidth(req.ScreenWidth).
		SetScreenHeight(req.ScreenHeight).
		OnConflictColumns(riderphonedevice.FieldRiderID).
		UpdateNewValues().
		SetDeviceSn(req.DeviceSn).
		SetModel(req.Model).
		SetBrand(req.Brand).
		SetImei(req.Imei).
		SetOsVersion(req.OsVersion).
		SetOsName(req.OsName).
		SetScreenWidth(req.ScreenWidth).
		SetScreenHeight(req.ScreenHeight).
		Exec(ctx.Request().Context())
	if err != nil {
		return err
	}
	return nil
}
