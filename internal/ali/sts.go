// Copyright (C) liasica. 2024-present.
//
// Created at 2024-02-28
// Based on aurservd by liasica, magicrolan@qq.com.

package ali

import (
	"errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
)

// 获取临时访问token
// @doc https://help.aliyun.com/document_detail/442329.html
// @doc https://help.aliyun.com/document_detail/469176.html
// @doc https://help.aliyun.com/document_detail/442255.html
// @doc https://help.aliyun.com/zh/ram/user-guide/create-a-ram-role-for-a-trusted-alibaba-cloud-account
// @doc https://help.aliyun.com/zh/ram/developer-reference/use-the-sts-openapi-example?spm=a2c4g.11186623.0.0.27194118OoyjXx
func stsToken(regionId, accessKeyId, accessKeySecret, arn, ramRole string) (credentials *sts.Credentials, err error) {
	var client *sts.Client
	client, err = sts.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		return
	}

	// 构建请求对象
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	request.RoleArn = arn
	request.RoleSessionName = ramRole
	request.DurationSeconds = "900"

	// 发起请求，并得到响应
	var response *sts.AssumeRoleResponse
	response, err = client.AssumeRole(request)
	if err != nil {
		return
	}

	if !response.IsSuccess() {
		return nil, errors.New("临时凭证申请失败")
	}

	return &response.Credentials, nil
}
