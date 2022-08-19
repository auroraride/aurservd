// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-19
// Based on aurservd by liasica, magicrolan@qq.com.

package excel

type Row []any

type CurrentRow struct {
    rowi   Row
    index  int
    total  int
    column any
}

type Rows struct {
    rowsi []Row
    index int
    row   *CurrentRow
}

func (r *Rows) Next() bool {
    r.index += 1
    if r.index > len(r.rowsi)-1 {
        r.row = nil
        return false
    }
    r.row = &CurrentRow{
        rowi:  r.rowsi[r.index],
        index: -1,
        total: len(r.rowsi[r.index]),
    }
    return true
}

func (r *CurrentRow) Next() bool {
    r.index += 1
    if r.index > r.total-1 {
        r.column = nil
        return false
    }
    r.column = r.rowi[r.index]
    return true
}
