// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

type agentFeedbackService struct {
	*BaseService
}

func NewAgentFeedback(params ...any) *agentFeedbackService {
	return &agentFeedbackService{
		BaseService: newService(params...),
	}
}
