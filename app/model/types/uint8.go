package types

import (
    "database/sql/driver"
    "fmt"
    "strconv"
)

type Uint8Type uint8

func (p *Uint8Type) Scan(v interface{}) (err error) {
    switch x := v.(type) {
    case nil:
    case uint8:
        *p = Uint8Type(x)
    case int64:
        *p = Uint8Type(x)
    case int, int8, int16, int32, uint, uint16, uint32, uint64, uintptr, float32, float64:
        var n uint64
        n, err = strconv.ParseUint(fmt.Sprintf("%x", x), 10, 64)
        if err != nil {
            return
        }
        *p = Uint8Type(n)
    default:
        err = fmt.Errorf("unexpected type %T", v)
    }
    return
}

func (p Uint8Type) Value() (driver.Value, error) {
    return p.Uint8(), nil
}

func (p Uint8Type) Uint8() uint8 {
    return uint8(p)
}
