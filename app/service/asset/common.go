package asset

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/h2non/filetype"
	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"

	"github.com/auroraride/aurservd/app/model/asset"
)

// ParseExcel 解析excel文件
func ParseExcel(c echo.Context) (list [][]string, err error) {
	source, err := c.FormFile("file")
	if err != nil {
		return nil, errors.New("上传文件失败:" + err.Error())
	}
	if source == nil {
		return nil, errors.New("未获取到上传的文件")
	}

	// 打开文件
	open, err := source.Open()
	if err != nil {
		return nil, errors.New("打开文件失败:" + err.Error())
	}
	defer func() {
		_ = open.Close()
	}()

	// 限制文件类型
	kind, err := filetype.MatchReader(open)
	if err != nil {
		return nil, errors.New("获取文件类型错误:" + err.Error())
	}
	if kind != filetype.GetType("xlsx") {
		return nil, errors.New("文件格式错误，必须为标准xlsx格式，当前为：" + kind.Extension)
	}
	_, _ = open.Seek(0, 0)

	// 获取文件内容
	r, err := excelize.OpenReader(open)
	if err != nil {
		return nil, errors.New("打开文件失败:" + err.Error())
	}
	defer func() {
		_ = r.Close()
	}()

	sheet := r.GetSheetName(0)

	rows, err := r.GetRows(sheet)
	if err != nil {
		return nil, errors.New("获取文件内容失败:" + err.Error())
	}

	for i, columns := range rows {
		// 第一行为标题
		if i == 0 {
			continue
		}
		list = append(list, columns)
	}
	return list, nil
}

// ParseBatterySN 解析电池编号
func ParseBatterySN(sn string) (bat asset.Battery, err error) {
	if len(sn) < 16 {
		return bat, errors.New("电池编号长度不足")
	}

	b := make([]byte, len(sn))
	for i := range b {
		c := sn[i]
		switch {
		case c >= 'a' && c <= 'z':
			c -= 'a' - 'A'
		case c < '0', c > 'z', c > '9' && c < 'A', c > 'Z' && c < 'a':
			return asset.Battery{}, errors.New("电池编号包含非法字符")
		}
		b[i] = c
	}
	sn = ConvertBytes2String(b)

	switch len(sn) {
	case asset.BatteryTbLength:
		bat.Brand = asset.BatteryBrandXC
		bat.Model = asset.BatteryModelXC[sn[3:5]]
	case asset.BatteryXcLength:
		bat.Brand = asset.BatteryBrandTB
		bat.Model = sn[4:6] + "V" + sn[7:9] + "AH"
	default:
		return asset.Battery{}, errors.New("电池编号解析失败")
	}
	bat.SN = sn
	return
}

func ConvertBytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// GenerateTemplate 生成模版excel文件
func (s *batteryService) GenerateTemplate(name string, data [][]string) error {
	f := excelize.NewFile()
	// 创建一个新的工作表
	index, err := f.NewSheet(name)
	if err != nil {
		return fmt.Errorf(": %v", err)
	}
	// 填充数据到工作表中
	for i, row := range data {
		for j, cell := range row {
			cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			err = f.SetCellValue(name, cellName, cell)
			if err != nil {
				return err
			}
		}
	}

	// 设置活动工作表
	f.SetActiveSheet(index)

	// 保存Excel文件
	if err := f.SaveAs(name + ".xlsx"); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}
	return nil
}
