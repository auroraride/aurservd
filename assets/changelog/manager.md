## 2022-06-11

#### 接口调整

- [M8001 订单列表](http://localhost:5533/docs#tag/M/operation/ManagerOrderList)
  - 筛选新增`riderId`(骑手ID) `refund`筛选退款
  - 返回新增`refund`退款详情



#### 接口新增

- [M8002 退款审核](http://localhost:5533/docs#tag/M/operation/ManagerOrderRefundAudit)



#### 后台调整

- 删除骑手详情下的订单、退款tabs
- 订单管理新增退款筛选，若有退款申请（判定条件: `status == 2`），则新增操作处理退款同意或拒绝，若拒绝则必填理由





<br />

## 2022-06-10

#### 接口新增

- [M7006 暂停计费](http://localhost:5533/docs#tag/M/operation/ManagerRiderPause)
- [M7007 继续计费](http://localhost:5533/docs#tag/M/operation/ManagerRiderContinue)
- [M7008 强制退租](http://localhost:5533/docs#tag/M/operation/ManagerSubscribeHalt)
- [M7009 修改押金](http://localhost:5533/docs#tag/M/operation/ManagerSubscribeDeposit)
- [M7010 修改骑手资料](http://localhost:5533/docs#tag/M/operation/ManagerSubscribeModify)



#### 接口调整

- [M7005 查看骑手操作日志](http://localhost:5533/docs#tag/M/operation/ManagerRiderLog)
  - 调整入参和返回为标准分页模型
  - 入参新增参数`type`详情参见文档
- [M7001列举骑手](http://localhost:5533/docs#tag/M/operation/RiderList)
  - 删除顶级字段`idCardNumber`
  - 新增字段`person`用户认证信息
  - 新增字段`contract`合同





<br />

## 2022-06-09

#### 接口新增

- [M1003 列举管理员](http://localhost:5533/docs#tag/M/operation/ManagerManagerList)
- [M1004 删除管理员](http://localhost:5533/docs#tag/M/operation/ManagerManagerDelete)



#### 接口调整

- [M8001 订单列表](http://localhost:5533/docs#tag/M/operation/ManagerOrderList) 筛选项和返回值调整



<br />

## 2022-06-08

#### 接口调整

- 接口编号调整



#### 接口新增

- [MA010 新增店员](http://localhost:5533/docs#tag/M/operation/ManagerEmployeeCreate)
- [MA011 修改店员](http://localhost:5533/docs#tag/M/operation/ManagerEmployeeModify)
- [MA012 列举店员](http://localhost:5533/docs#tag/M/operation/ManagerEmployeeList)
- [MA013 删除店员](http://localhost:5533/docs#tag/M/operation/ManagerEmployeeDelete)



<br />

## 2022-06-07

#### 接口新增

- [M9004 企业详情](http://localhost:5533/docs#tag/M/operation/ManagerEnterpriseDetail)
- [M9006 创建站点](http://localhost:5533/docs#tag/M/operation/ManagerEnterpriseCreateStation)
- [M9007 编辑站点](http://localhost:5533/docs#tag/M/operation/ManagerEnterpriseModifyStation)
- [M9008 列举站点](http://localhost:5533/docs#tag/M/operation/ManagerEnterpriseListStation)
- [M9009 添加骑手](http://localhost:5533/docs#tag/M/operation/ManagerEnterpriseCreateRider)
- [M9010 列举骑手](http://localhost:5533/docs#tag/M/operation/ManagerEnterpriseListRider)



<br />

## 2022-06-06

#### 接口新增

- [M9003 列举企业](/docs#tag/M/operation/ManagerEnterpriseList)
- [M9005 企业预付费](/docs#tag/M/operation/ManagerEnterprisePrepayment)



<br />

## 2022-06-05

#### 接口调整

- 骑手列表新增字段
  - `deletedAt`账户删除时间（已删除账户会有此字段，团签骑手退租后会自动软删除账户，下次使用骑手会重新注册）
  - `remark`备注



#### 接口新增

- [M9001 创建企业](/docs#tag/M/operation/ManagerEnterpriseCreate)
- [M9002 更新企业](/docs#tag/M/operation/ManagerEnterpriseModify)



<br />

## 2022-06-03

#### 接口调整

- [M7001 列举骑手](/docs#tag/M/operation/RiderList) 新增字段 `address`(户籍地址)



#### 接口新增

- [M8001 订单列表](/docs#tag/M/operation/ManagerOrderList)



<br />

## 2022-06-02

#### 接口新增

- [M7004 修改订阅时间](/docs#tag/M/operation/ManagerSubscribeAlter)
- [M7005 查看骑手操作日志](/docs#tag/M/operation/ManagerRiderLog)



#### 骑手属性

原型单一骑手状态拆分为

- 业务状态 `subscribe.status`  若`subscribe`对象字段不存在则代表用户从未使用过
- 骑手状态 `status` 
- 认证状态 `authStatus`

详细见接口文档 [M7001列举骑手](/docs#tag/M/operation/RiderList)解释，筛选亦如是。

