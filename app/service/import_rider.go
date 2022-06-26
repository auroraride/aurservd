// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import "context"

type importRiderService struct {
    ctx context.Context
}

func NewImportRider() *importRiderService {
    return &importRiderService{
        ctx: context.Background(),
    }
}

func (s *importRiderService) ParseCSV() {

}
