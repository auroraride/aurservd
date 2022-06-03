## 2022-06-03

#### 字段调整

- [M70001 列举骑手](/docs#tag/M/operation/RiderList) 新增字段 `address`(户籍地址)



#### 接口新增

- [M80001 订单列表](/docs#tag/M/operation/ManagerOrderList)



## 2022-06-02

#### 接口新增

- [M70004 修改订阅时间](/docs#tag/M/operation/ManagerSubscribeAlter)
- [M70005 查看骑手操作日志](/docs#tag/M/operation/ManagerRiderLog)



#### 用户属性

原型单一用户状态拆分为

- 业务状态 `subscribe.status`  若`subscribe`对象字段不存在则代表用户从未使用过
- 用户状态 `status` 
- 认证状态 `authStatus`

详细见接口文档 [M70001列举骑手](/docs#tag/M/operation/RiderList)解释，筛选亦如是。

