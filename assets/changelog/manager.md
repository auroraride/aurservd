## 2022-06-02

#### 用户属性

原型单一用户状态拆分为

- 业务状态 `subscribe.status`  若`subscribe`对象字段不存在则代表用户从未使用过
- 用户状态 `status` 
- 认证状态 `authStatus`

详细见接口文档 [M70001列举骑手](/docs#tag/M/operation/RiderList)解释，筛选亦如是。

