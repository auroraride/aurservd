package definition

type RiderPhoneDeviceReq struct {
	Imei         string `json:"imei"`         // IMEI
	DeviceSn     string `json:"deviceSn"`     // 设备编号
	Model        string `json:"model"`        // 手机型号
	Brand        string `json:"brand"`        // 手机品牌
	OsVersion    string `json:"osVersion"`    // 系统版本
	OsName       string `json:"osName"`       // 系统名称
	ScreenWidth  uint64 `json:"screenWidth"`  // 屏幕宽度
	ScreenHeight uint64 `json:"screenHeight"` // 屏幕高度
}
