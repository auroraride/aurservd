// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

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
        "/commom/sms": {
            "post": {
                "description": "上传文件必须，单次获取有效时间为1个小时",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "[C]公共接口"
                ],
                "summary": "C2.发送短信验证码",
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
        },
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
                    "[C]公共接口"
                ],
                "summary": "C1.生成图片验证码",
                "operationId": "CaptchaGenerate",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        },
                        "headers": {
                            "X-Captcha-Id": {
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
                    "[C]公共接口"
                ],
                "summary": "C3.获取阿里云oss临时凭证",
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
        "/manager/v1/branch": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "[M]管理接口"
                ],
                "summary": "M3.1 网点列表",
                "operationId": "BranchList",
                "parameters": [
                    {
                        "type": "string",
                        "description": "管理员校验token",
                        "name": "X-Manager-Token",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.PaginationRes"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "items": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/model.Branch"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "[M]管理接口"
                ],
                "summary": "M3.2 新增网点",
                "operationId": "BranchAdd",
                "parameters": [
                    {
                        "type": "string",
                        "description": "管理员校验token",
                        "name": "X-Manager-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "网点数据",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Branch"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.StatusResponse"
                        }
                    }
                }
            }
        },
        "/manager/v1/branch/{id}": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "[M]管理接口"
                ],
                "summary": "M3.3 编辑网点",
                "operationId": "BranchModify",
                "parameters": [
                    {
                        "type": "string",
                        "description": "管理员校验token",
                        "name": "X-Manager-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "网点数据",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.Branch"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.StatusResponse"
                        }
                    }
                }
            }
        },
        "/manager/v1/city": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "[M]管理接口"
                ],
                "summary": "M2.1 城市列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "管理员校验token",
                        "name": "X-Manager-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "maximum": 2,
                        "minimum": 0,
                        "type": "integer",
                        "description": "启用状态 0:全部 1:未启用 2:已启用",
                        "name": "status",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.CityItem"
                            }
                        }
                    }
                }
            }
        },
        "/manager/v1/city/{id}": {
            "put": {
                "description": "desc",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "[M]管理接口"
                ],
                "summary": "M2.2 修改城市",
                "operationId": "CityModify",
                "parameters": [
                    {
                        "type": "string",
                        "description": "管理员校验token",
                        "name": "X-Manager-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "城市ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "城市数据",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.CityModifyReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.StatusResponse"
                        }
                    }
                }
            }
        },
        "/manager/v1/user/signin": {
            "post": {
                "description": "管理员登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "[M]管理接口"
                ],
                "summary": "M1.1 用户登录",
                "operationId": "ManagerSignin",
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.ManagerSigninRes"
                        }
                    }
                }
            }
        },
        "/manager/v1/{id}/contract": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "[M]管理接口"
                ],
                "summary": "M3.4 新增合同",
                "operationId": "BranchAddContract",
                "parameters": [
                    {
                        "type": "string",
                        "description": "管理员校验token",
                        "name": "X-Manager-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "合同数据",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.BranchContract"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.StatusResponse"
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
        "model.Branch": {
            "type": "object",
            "required": [
                "address",
                "cityId",
                "lat",
                "lng",
                "name",
                "photos"
            ],
            "properties": {
                "address": {
                    "description": "详细地址",
                    "type": "string"
                },
                "cityId": {
                    "description": "城市",
                    "type": "integer"
                },
                "contracts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.BranchContract"
                    }
                },
                "creator": {
                    "$ref": "#/definitions/model.Modifier"
                },
                "lastModifier": {
                    "$ref": "#/definitions/model.Modifier"
                },
                "lat": {
                    "description": "纬度",
                    "type": "number"
                },
                "lng": {
                    "description": "经度",
                    "type": "number"
                },
                "name": {
                    "description": "网点名称",
                    "type": "string"
                },
                "photos": {
                    "description": "网点照片",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "model.BranchContract": {
            "type": "object",
            "required": [
                "area",
                "bankNumber",
                "electricity",
                "electricityPledge",
                "endTime",
                "file",
                "idCardNumber",
                "landlordName",
                "lease",
                "phone",
                "pledge",
                "rent",
                "sheets",
                "startTime"
            ],
            "properties": {
                "area": {
                    "description": "网点面积",
                    "type": "number"
                },
                "bankNumber": {
                    "description": "房东银行卡号",
                    "type": "string"
                },
                "electricity": {
                    "description": "电费单价",
                    "type": "number"
                },
                "electricityPledge": {
                    "description": "电费押金",
                    "type": "number"
                },
                "endTime": {
                    "description": "租期结束时间",
                    "type": "string"
                },
                "file": {
                    "description": "合同文件",
                    "type": "string"
                },
                "idCardNumber": {
                    "description": "房东身份证",
                    "type": "string"
                },
                "landlordName": {
                    "description": "房东姓名",
                    "type": "string"
                },
                "lease": {
                    "description": "租期月数",
                    "type": "integer"
                },
                "phone": {
                    "description": "房东手机号",
                    "type": "string"
                },
                "pledge": {
                    "description": "押金",
                    "type": "number"
                },
                "rent": {
                    "description": "租金",
                    "type": "number"
                },
                "sheets": {
                    "description": "底单",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "startTime": {
                    "description": "租期开始时间",
                    "type": "string"
                }
            }
        },
        "model.CityItem": {
            "type": "object",
            "properties": {
                "children": {
                    "description": "城市列表",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.CityItem"
                    }
                },
                "id": {
                    "description": "城市或省份ID",
                    "type": "integer"
                },
                "name": {
                    "description": "城市/省份",
                    "type": "string"
                },
                "open": {
                    "description": "是否启用",
                    "type": "boolean"
                }
            }
        },
        "model.CityModifyReq": {
            "type": "object",
            "required": [
                "open"
            ],
            "properties": {
                "open": {
                    "description": "状态",
                    "type": "boolean"
                }
            }
        },
        "model.ManagerSigninRes": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "model.Modifier": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "model.Pagination": {
            "type": "object",
            "properties": {
                "current": {
                    "description": "当前页",
                    "type": "integer"
                },
                "pages": {
                    "description": "总页数",
                    "type": "integer"
                },
                "total": {
                    "description": "总条数",
                    "type": "integer"
                }
            }
        },
        "model.PaginationRes": {
            "type": "object",
            "properties": {
                "items": {
                    "description": "返回数据"
                },
                "pagination": {
                    "description": "分页属性",
                    "$ref": "#/definitions/model.Pagination"
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
        },
        "model.StatusResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "description": "默认接口成功返回",
                    "type": "boolean"
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
	Title:            "极光出行API",
	Description:      "### 说明\n> 接口采用非标准Restful API，所有http返回代码均为`200`，当返回为非`200`时应为network错误，需要及时排查。\n> 接口返回说明查看 **[返回](#返回)**\n\n### 认证\n项目接口使用简单认证，认证方式为`header`中添加对应的认证`token`\n|  header   |  类型  |  接口  |\n| :-----: | :----: | :--: |\n|  X-Rider-Token   |  string   |  骑手API  |\n| X-Manager-Token | string |  后台API  |\n|  X-Employee-Token   | string |  员工API  |\n\n### 返回\n一个标准的返回应包含以下结构\n\n|  字段   |  类型  |  必填  |  说明  |\n| :-----: | :----: | :--: | :--: |\n|  code   |  int   |  是  |  返回代码  |\n| message | string |  是  |  返回消息  |\n|  data   | object |  是  |  返回数据  |\n\n`code`代码取值说明\n\n| 二进制 | 十六进制 | 说明 |\n| :----: | :------: | :--: |\n| 0  |  0x000  | 请求成功 |\n| 256 |  0x100  | 请求失败 |\n| 512 |  0x200  | 需要认证 (需要登录) |\n| 768 |  0x300  | 没有权限 |\n| 1024 |  0x400  | 资源未获 |\n| 1280 |  0x500  | 未知错误 |\n| 1536 |  0x600  | 需要实名 |\n| 1792 |  0x700  | 需要验证 (更换设备需要人脸验证) |\n| 2048 |  0x800  | 需要联系人 |\n| 2304 |  0x900  | 请求过期 |\n\n\n比如：\n> 默认成功返回\n```json\n{\n  \"code\": 0,\n  \"message\": \"OK\",\n  \"data\": {\n    \"status\": true\n  }\n}\n```",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}