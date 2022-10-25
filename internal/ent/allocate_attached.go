// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-25
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
    "context"
    "errors"
)

func (u *AllocateUpsertOne) Save(ctx context.Context) (*Allocate, error) {
    if len(u.create.conflict) == 0 {
        return nil, errors.New("ent: missing options for AllocateCreate.OnConflict")
    }
    return u.create.Save(ctx)
}

func (u *AllocateUpsertOne) SaveX(ctx context.Context) *Allocate {
    allo, err := u.create.Save(ctx)
    if err != nil {
        panic(err)
    }
    return allo
}
