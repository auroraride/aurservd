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

