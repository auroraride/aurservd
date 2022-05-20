// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
)

type personService struct {
    ctx context.Context
}

func NewPerson() *personService {
    return &personService{
        ctx: context.Background(),
    }
}

func (s *personService) Ban() {

}