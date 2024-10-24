// Package rider Code generated by swaggo/swag. DO NOT EDIT
package rider

import "github.com/liasica/swag"

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
        "/rider/v2/purchase/contract/create": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contract - 合同"
                ],
                "summary": "创建合同",
                "operationId": "ContractCreate",
                "parameters": [
                    {
                        "type": "string",
                        "description": "骑手校验token",
                        "name": "X-Rider-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "desc",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ContractCreateReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.ContractCreateRes"
                        }
                    }
                }
            }
        },
        "/rider/v2/purchase/contract/sign": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contract - 合同"
                ],
                "summary": "签署合同",
                "operationId": "ContractSign",
                "parameters": [
                    {
                        "type": "string",
                        "description": "骑手校验token",
                        "name": "X-Rider-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "desc",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ContractSignNewReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.ContractSignNewRes"
                        }
                    }
                }
            }
        },
        "/rider/v2/purchase/contract/{docId}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contract - 合同"
                ],
                "summary": "查看合同",
                "operationId": "ContractDetail",
                "parameters": [
                    {
                        "type": "string",
                        "description": "骑手校验token",
                        "name": "X-Rider-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "合同ID",
                        "name": "docId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.ContractDetailRes"
                        }
                    }
                }
            }
        },
        "/rider/v2/purchase/order": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order - 购车订单"
                ],
                "summary": "订单列表",
                "operationId": "OrderList",
                "parameters": [
                    {
                        "type": "string",
                        "description": "骑手校验token",
                        "name": "X-Rider-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "enum": [
                            1,
                            2
                        ],
                        "type": "integer",
                        "x-enum-comments": {
                            "BillStatusNormal": "正常",
                            "BillStatusOverdue": "逾期"
                        },
                        "x-enum-varnames": [
                            "BillStatusNormal",
                            "BillStatusOverdue"
                        ],
                        "description": "还款状态",
                        "name": "billStatus",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "当前页, 从1开始, 默认1",
                        "name": "current",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "结束时间",
                        "name": "end",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "订单编号",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "关键字",
                        "name": "keyword",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页数据, 默认20",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "骑手ID",
                        "name": "riderId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "车架号",
                        "name": "sn",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "开始时间",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "pending",
                            "staging",
                            "ended",
                            "cancelled",
                            "refunded"
                        ],
                        "type": "string",
                        "x-enum-comments": {
                            "OrderStatusCancelled": "已取消",
                            "OrderStatusEnded": "已完成",
                            "OrderStatusPending": "待支付",
                            "OrderStatusRefunded": "已退款",
                            "OrderStatusStaging": "分期执行中"
                        },
                        "x-enum-varnames": [
                            "OrderStatusPending",
                            "OrderStatusStaging",
                            "OrderStatusEnded",
                            "OrderStatusCancelled",
                            "OrderStatusRefunded"
                        ],
                        "description": "订单状态",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "门店ID",
                        "name": "storeId",
                        "in": "query"
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
                                                "$ref": "#/definitions/model.PurchaseOrderListRes"
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
                    "Order - 购车订单"
                ],
                "summary": "创建订单",
                "operationId": "OrderCreate",
                "parameters": [
                    {
                        "type": "string",
                        "description": "骑手校验token",
                        "name": "X-Rider-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "请求参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PurchaseOrderCreateReq"
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
        "/rider/v2/purchase/order/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order - 购车订单"
                ],
                "summary": "订单详情",
                "operationId": "OrderDetail",
                "parameters": [
                    {
                        "type": "string",
                        "description": "管理员校验token",
                        "name": "X-Rider-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "订单ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.PurchaseOrderDetail"
                        }
                    }
                }
            }
        },
        "/rider/v2/purchase/pay": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order - 购车订单"
                ],
                "summary": "订单支付",
                "operationId": "PaymentPay",
                "parameters": [
                    {
                        "type": "string",
                        "description": "骑手校验token",
                        "name": "X-Rider-Token",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "请求参数",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PaymentReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "请求成功",
                        "schema": {
                            "$ref": "#/definitions/model.PurchasePayRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "definition.Goods": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "description": "创建时间",
                    "type": "string"
                },
                "headPic": {
                    "description": "商品头图",
                    "type": "string"
                },
                "id": {
                    "description": "商品ID",
                    "type": "integer"
                },
                "intro": {
                    "description": "商品介绍",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "lables": {
                    "description": "商品标签",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "description": "商品名称",
                    "type": "string"
                },
                "paymentPlans": {
                    "description": "付款方案",
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "number"
                        }
                    }
                },
                "photos": {
                    "description": "商品图片",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "price": {
                    "description": "商品价格",
                    "type": "number"
                },
                "remark": {
                    "description": "备注",
                    "type": "string"
                },
                "sn": {
                    "description": "商品编号",
                    "type": "string"
                },
                "status": {
                    "description": "商品状态 0-已下架 1-已上架",
                    "allOf": [
                        {
                            "$ref": "#/definitions/definition.GoodsStatus"
                        }
                    ]
                },
                "storeIds": {
                    "description": "门店Ids",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "stores": {
                    "description": "配置店铺信息",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Store"
                    }
                },
                "type": {
                    "description": "商品类型 1-电车",
                    "allOf": [
                        {
                            "$ref": "#/definitions/definition.GoodsType"
                        }
                    ]
                },
                "weight": {
                    "description": "商品权重",
                    "type": "integer"
                }
            }
        },
        "definition.GoodsStatus": {
            "type": "integer",
            "enum": [
                0,
                1
            ],
            "x-enum-comments": {
                "GoodsStatusOffline": "下架",
                "GoodsStatusOnline": "上架"
            },
            "x-enum-varnames": [
                "GoodsStatusOffline",
                "GoodsStatusOnline"
            ]
        },
        "definition.GoodsType": {
            "type": "integer",
            "enum": [
                1
            ],
            "x-enum-comments": {
                "GoodsTypeEbike": "电车"
            },
            "x-enum-varnames": [
                "GoodsTypeEbike"
            ]
        },
        "github_com_auroraride_aurservd_app_purchase_internal_model.Payway": {
            "type": "string",
            "enum": [
                "alipay",
                "wechat",
                "cash"
            ],
            "x-enum-varnames": [
                "Alipay",
                "Wechat",
                "Cash"
            ]
        },
        "model.BillStatus": {
            "type": "integer",
            "enum": [
                1,
                2
            ],
            "x-enum-comments": {
                "BillStatusNormal": "正常",
                "BillStatusOverdue": "逾期"
            },
            "x-enum-varnames": [
                "BillStatusNormal",
                "BillStatusOverdue"
            ]
        },
        "model.ContractCreateReq": {
            "type": "object",
            "required": [
                "orderId"
            ],
            "properties": {
                "orderId": {
                    "description": "订单ID",
                    "type": "integer"
                }
            }
        },
        "model.ContractCreateRes": {
            "type": "object",
            "properties": {
                "docId": {
                    "description": "合同ID",
                    "type": "string"
                },
                "effective": {
                    "description": "是否存在生效中的合同, 若返回值为true则代表无需签合同",
                    "type": "boolean"
                },
                "link": {
                    "description": "合同链接",
                    "type": "string"
                },
                "needRealName": {
                    "description": "是否需要实名认证   true:需要  false:不需要",
                    "type": "boolean"
                }
            }
        },
        "model.ContractDetailRes": {
            "type": "object",
            "properties": {
                "link": {
                    "description": "合同链接",
                    "type": "string"
                }
            }
        },
        "model.ContractSignNewReq": {
            "type": "object",
            "required": [
                "docId",
                "orderId",
                "seal"
            ],
            "properties": {
                "docId": {
                    "description": "合同ID",
                    "type": "string"
                },
                "orderId": {
                    "description": "订单ID",
                    "type": "integer"
                },
                "seal": {
                    "description": "签名Base64",
                    "type": "string"
                }
            }
        },
        "model.ContractSignNewRes": {
            "type": "object",
            "properties": {
                "link": {
                    "description": "合同链接",
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
        "model.OrderStatus": {
            "type": "string",
            "enum": [
                "pending",
                "staging",
                "ended",
                "cancelled",
                "refunded"
            ],
            "x-enum-comments": {
                "OrderStatusCancelled": "已取消",
                "OrderStatusEnded": "已完成",
                "OrderStatusPending": "待支付",
                "OrderStatusRefunded": "已退款",
                "OrderStatusStaging": "分期执行中"
            },
            "x-enum-varnames": [
                "OrderStatusPending",
                "OrderStatusStaging",
                "OrderStatusEnded",
                "OrderStatusCancelled",
                "OrderStatusRefunded"
            ]
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
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.Pagination"
                        }
                    ]
                }
            }
        },
        "model.PaymentDetail": {
            "type": "object",
            "properties": {
                "amount": {
                    "description": "订单金额",
                    "type": "number"
                },
                "billingDate": {
                    "description": "账单日期",
                    "type": "string"
                },
                "forfeit": {
                    "description": "逾期金额（滞纳金）",
                    "type": "number"
                },
                "id": {
                    "description": "分期订单ID",
                    "type": "integer"
                },
                "outTradeNo": {
                    "description": "交易订单号",
                    "type": "string"
                },
                "overdueDays": {
                    "description": "逾期天数",
                    "type": "integer"
                },
                "paymentTime": {
                    "description": "支付时间",
                    "type": "string"
                },
                "payway": {
                    "description": "支付方式 alipay-支付宝 wechat-微信支付 cash-现金",
                    "allOf": [
                        {
                            "$ref": "#/definitions/github_com_auroraride_aurservd_app_purchase_internal_model.Payway"
                        }
                    ]
                },
                "status": {
                    "description": "支付状态 obligation:待付款, paid:已支付, canceled:已取消, overdue:已逾期",
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.PaymentStatus"
                        }
                    ]
                },
                "total": {
                    "description": "支付金额（订单金额+逾期金额）",
                    "type": "number"
                }
            }
        },
        "model.PaymentReq": {
            "type": "object",
            "required": [
                "orderId",
                "payway",
                "planIndex"
            ],
            "properties": {
                "orderId": {
                    "description": "订单id",
                    "type": "integer"
                },
                "payway": {
                    "description": "支付方式",
                    "allOf": [
                        {
                            "$ref": "#/definitions/github_com_auroraride_aurservd_app_purchase_internal_model.Payway"
                        }
                    ]
                },
                "planIndex": {
                    "description": "付款计划索引",
                    "type": "integer"
                }
            }
        },
        "model.PaymentStatus": {
            "type": "string",
            "enum": [
                "obligation",
                "paid",
                "canceled",
                "overdue"
            ],
            "x-enum-comments": {
                "PaymentStatusCanceled": "已取消",
                "PaymentStatusObligation": "待付款",
                "PaymentStatusOverdue": "已逾期",
                "PaymentStatusPaid": "已支付"
            },
            "x-enum-varnames": [
                "PaymentStatusObligation",
                "PaymentStatusPaid",
                "PaymentStatusCanceled",
                "PaymentStatusOverdue"
            ]
        },
        "model.PurchaseOrderCreateReq": {
            "type": "object",
            "required": [
                "goodsId",
                "planIndex"
            ],
            "properties": {
                "goodsId": {
                    "description": "商品id",
                    "type": "integer"
                },
                "planIndex": {
                    "description": "付款计划索引",
                    "type": "integer"
                }
            }
        },
        "model.PurchaseOrderDetail": {
            "type": "object",
            "properties": {
                "activeName": {
                    "description": "激活人名称",
                    "type": "string"
                },
                "activePhone": {
                    "description": "激活人电话",
                    "type": "string"
                },
                "amount": {
                    "description": "订单金额",
                    "type": "number"
                },
                "billStatus": {
                    "description": "账单状态 // 1-正常 2-逾期",
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.BillStatus"
                        }
                    ]
                },
                "color": {
                    "description": "车辆颜色",
                    "type": "string"
                },
                "contractUrl": {
                    "description": "合同url",
                    "type": "string"
                },
                "createdAt": {
                    "description": "创建时间",
                    "type": "string"
                },
                "docId": {
                    "description": "合同ID",
                    "type": "string"
                },
                "follows": {
                    "description": "订单跟进数据",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.PurchaseOrderFollow"
                    }
                },
                "formula": {
                    "description": "违约说明",
                    "type": "string"
                },
                "goods": {
                    "description": "商品信息",
                    "allOf": [
                        {
                            "$ref": "#/definitions/definition.Goods"
                        }
                    ]
                },
                "id": {
                    "description": "订单编号",
                    "type": "integer"
                },
                "installmentPlan": {
                    "description": "分期方案",
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "installmentStage": {
                    "description": "当前分期阶段",
                    "type": "integer"
                },
                "installmentTotal": {
                    "description": "分期总数",
                    "type": "integer"
                },
                "paidAmount": {
                    "description": "已支付金额",
                    "type": "number"
                },
                "payments": {
                    "description": "分期订单数据（还款计划）",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.PaymentDetail"
                    }
                },
                "planIndex": {
                    "description": "付款计划索引",
                    "type": "integer"
                },
                "remark": {
                    "description": "备注",
                    "type": "string"
                },
                "riderName": {
                    "description": "骑手名称",
                    "type": "string"
                },
                "riderPhone": {
                    "description": "骑手电话",
                    "type": "string"
                },
                "signed": {
                    "description": "是否签约 true:已签约 false:未签约",
                    "type": "boolean"
                },
                "sn": {
                    "description": "车架号",
                    "type": "string"
                },
                "startDate": {
                    "description": "激活时间",
                    "type": "string"
                },
                "status": {
                    "description": "订单状态 pending: 待支付, staging: 分期执行中, ended: 已完成, cancelled: 已取消, refunded: 已退款",
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.OrderStatus"
                        }
                    ]
                },
                "storeId": {
                    "description": "门店ID",
                    "type": "integer"
                },
                "storeName": {
                    "description": "提车门店",
                    "type": "string"
                }
            }
        },
        "model.PurchaseOrderFollow": {
            "type": "object",
            "properties": {
                "content": {
                    "description": "跟进内容",
                    "type": "string"
                },
                "createdAt": {
                    "description": "跟进时间",
                    "type": "string"
                },
                "id": {
                    "description": "跟进ID",
                    "type": "integer"
                },
                "modifier": {
                    "description": "跟进人",
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.Modifier"
                        }
                    ]
                },
                "pics": {
                    "description": "跟进图片",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "model.PurchaseOrderListRes": {
            "type": "object",
            "properties": {
                "activeName": {
                    "description": "激活人名称",
                    "type": "string"
                },
                "activePhone": {
                    "description": "激活人电话",
                    "type": "string"
                },
                "amount": {
                    "description": "订单金额",
                    "type": "number"
                },
                "billStatus": {
                    "description": "账单状态 // 1-正常 2-逾期",
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.BillStatus"
                        }
                    ]
                },
                "color": {
                    "description": "车辆颜色",
                    "type": "string"
                },
                "contractUrl": {
                    "description": "合同url",
                    "type": "string"
                },
                "createdAt": {
                    "description": "创建时间",
                    "type": "string"
                },
                "docId": {
                    "description": "合同ID",
                    "type": "string"
                },
                "formula": {
                    "description": "违约说明",
                    "type": "string"
                },
                "goods": {
                    "description": "商品信息",
                    "allOf": [
                        {
                            "$ref": "#/definitions/definition.Goods"
                        }
                    ]
                },
                "id": {
                    "description": "订单编号",
                    "type": "integer"
                },
                "installmentPlan": {
                    "description": "分期方案",
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "installmentStage": {
                    "description": "当前分期阶段",
                    "type": "integer"
                },
                "installmentTotal": {
                    "description": "分期总数",
                    "type": "integer"
                },
                "paidAmount": {
                    "description": "已支付金额",
                    "type": "number"
                },
                "planIndex": {
                    "description": "付款计划索引",
                    "type": "integer"
                },
                "remark": {
                    "description": "备注",
                    "type": "string"
                },
                "riderName": {
                    "description": "骑手名称",
                    "type": "string"
                },
                "riderPhone": {
                    "description": "骑手电话",
                    "type": "string"
                },
                "signed": {
                    "description": "是否签约 true:已签约 false:未签约",
                    "type": "boolean"
                },
                "sn": {
                    "description": "车架号",
                    "type": "string"
                },
                "startDate": {
                    "description": "激活时间",
                    "type": "string"
                },
                "status": {
                    "description": "订单状态 pending: 待支付, staging: 分期执行中, ended: 已完成, cancelled: 已取消, refunded: 已退款",
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.OrderStatus"
                        }
                    ]
                },
                "storeId": {
                    "description": "门店ID",
                    "type": "integer"
                },
                "storeName": {
                    "description": "提车门店",
                    "type": "string"
                }
            }
        },
        "model.PurchasePayRes": {
            "type": "object",
            "properties": {
                "outTradeNo": {
                    "description": "交易编码",
                    "type": "string"
                },
                "prepay": {
                    "description": "预支付字符串",
                    "type": "string"
                }
            }
        },
        "model.StatusResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "boolean"
                }
            }
        },
        "model.Store": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "description": "门店名称",
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
	Title:            "极光出行API - 骑手端购车api",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
