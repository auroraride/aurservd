// Package v2 Code generated by swaggo/swag. DO NOT EDIT
package v2

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
        "/v2/certification/face": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Person - 实人"
                ],
                "summary": "获取人脸核身参数",
                "operationId": "CertificationFace",
                "parameters": [
                    {
                        "type": "string",
                        "description": "身份信息 // TODO: 加密传输",
                        "name": "identity",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "订单号，用户使用OCR识别时不为空",
                        "name": "orderNo",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/definition.PersonCertificationFaceRes"
                        }
                    }
                }
            }
        },
        "/v2/certification/face/result": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Person - 实人"
                ],
                "summary": "获取实人核身结果",
                "operationId": "CertificationFaceResult",
                "parameters": [
                    {
                        "type": "string",
                        "description": "订单号",
                        "name": "orderNo",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/definition.PersonCertificationFaceResultRes"
                        }
                    }
                }
            }
        },
        "/v2/certification/ocr": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Person - 实人"
                ],
                "summary": "获取人脸核身OCR参数",
                "operationId": "CertificationOcr",
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/definition.PersonCertificationOcrRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "definition.PersonCertificationFaceRes": {
            "type": "object",
            "properties": {
                "appId": {
                    "description": "WBAppid",
                    "type": "string"
                },
                "faceId": {
                    "type": "string"
                },
                "licence": {
                    "type": "string"
                },
                "nonce": {
                    "description": "随机字符串",
                    "type": "string"
                },
                "orderNo": {
                    "description": "订单号",
                    "type": "string"
                },
                "sign": {
                    "description": "签名",
                    "type": "string"
                },
                "userId": {
                    "description": "用户唯一标识",
                    "type": "string"
                },
                "version": {
                    "description": "版本号",
                    "type": "string"
                }
            }
        },
        "definition.PersonCertificationFaceResultRes": {
            "type": "object",
            "properties": {
                "success": {
                    "description": "是否成功",
                    "type": "boolean"
                }
            }
        },
        "definition.PersonCertificationOcrRes": {
            "type": "object",
            "properties": {
                "appId": {
                    "description": "WBAppid",
                    "type": "string"
                },
                "nonce": {
                    "description": "随机字符串",
                    "type": "string"
                },
                "orderNo": {
                    "description": "订单号",
                    "type": "string"
                },
                "sign": {
                    "description": "签名",
                    "type": "string"
                },
                "userId": {
                    "description": "用户唯一标识",
                    "type": "string"
                },
                "version": {
                    "description": "版本号",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "2.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "极光出行API - 骑手端api",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}