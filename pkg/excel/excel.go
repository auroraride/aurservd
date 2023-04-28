// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-19
// Based on aurservd by liasica, magicrolan@qq.com.

package excel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

type excel struct {
	rows     *Rows
	excelize *excelize.File
}

// GetCell 获取表格 cell 信息
func GetCell(row, column int) string {
	name := fmt.Sprintf("%s%d", string(rune(65+column)), row+1)
	return name
}

func (e *excel) Save(path string) {
	d := filepath.Dir(path)
	_, err := os.Stat(d)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(d, 0755)
	}
}

func New(rows []Row, args ...any) *excel {
	sheet := "Sheet1"
	if len(args) > 0 {
		sheet = args[0].(string)
	}
	fe := excelize.NewFile()
	fe.SetSheetName("Sheet1", sheet)
	return &excel{
		rows: &Rows{
			rowsi: rows,
			index: -1,
		},
		excelize: fe,
	}
}

func (e *excel) InsertRows() {
	for e.rows.Next() {
		println("")
		for e.rows.row.Next() {

		}
	}
}
