// Package v1 Code generated by swaggo/swag. DO NOT EDIT
package v1

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/common/captcha": {
            "get": {
                "description": "生成的图片验证码有效时间为10分钟",
                "consumes": [
                    "image/png"
                ],
                "produces": [
                    "image/png"
                ],
                "tags": [
                    "Communal - 公共接口"
                ],
                "summary": "C1 生成图片验证码",
                "operationId": "CaptchaGenerate",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        },
                        "headers": {
                            "X-Captcha-Id\ttrue": {
                                "type": "string",
                                "description": "Captcha验证码ID"
                            }
                        }
                    }
                }
            }
        },
        "/common/oss/token": {
            "get": {
                "description": "上传文件必须，单次获取有效时间为1个小时",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Communal - 公共接口"
                ],
                "summary": "C3 获取阿里云oss临时凭证",
                "operationId": "AliyunOssToken",
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.AliyunOssStsRes"
                        }
                    }
                }
            }
        },
        "/common/sms": {
            "post": {
                "description": "上传文件必须，单次获取有效时间为1个小时",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Communal - 公共接口"
                ],
                "summary": "C2 发送短信验证码",
                "operationId": "SendSmsCode",
                "parameters": [
                    {
                        "description": "请求参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.SmsReq"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Captcha验证码ID",
                        "name": "X-Captcha-Id",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.SmsRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.AliyunOssStsRes": {
            "type": "object",
            "properties": {
                "accessKeyId": {
                    "type": "string"
                },
                "accessKeySecret": {
                    "type": "string"
                },
                "bucket": {
                    "type": "string"
                },
                "expiration": {
                    "type": "string"
                },
                "region": {
                    "type": "string"
                },
                "stsToken": {
                    "type": "string"
                }
            }
        },
        "model.SmsReq": {
            "type": "object",
            "required": [
                "captchaCode",
                "phone"
            ],
            "properties": {
                "captchaCode": {
                    "description": "captcha 验证码",
                    "type": "string"
                },
                "phone": {
                    "description": "手机号",
                    "type": "string"
                }
            }
        },
        "model.SmsRes": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "任务ID",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "极光出行API - 公共api",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}