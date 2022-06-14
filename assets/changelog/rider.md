## 2022-06-14

#### 接口调整

- 骑手订阅信息新增字段`business`是否可办理业务





<br />

## 2022-06-08

#### APP

- 新增团签用户激活逻辑



#### 接口调整

- 接口编号调整
- 骑手登录或资料新增字段
  - `phone`电话
  - `name`骑手姓名，仅骑手认证后才会有



#### 接口新增

- [R3010 企业骑手获取可用电池](http://localhost:5533/docs#tag/R/operation/RiderEnterpriseListVoltage)
- [R3011 企业骑手选择电池](http://localhost:5533/docs#tag/R/operation/RiderEnterpriseSubscribe)
- [R3012 企业骑手订阅激活状态](http://localhost:5533/docs#tag/R/operation/RiderEnterpriseSubscribeStatus) (长连接轮询，30s轮询一次，返回结果为`TRUE`或出现错误时中断轮询)



<br />

## 2022-06-07

#### APP

- 骑手信息新增字段`enterprise`所属企业，仅团签用户有此字段，此字段可进行是否团签用户的判定，详情参考文档[R1006 获取个人信息](http://localhost:5533/docs#tag/R/operation/RiderRiderProfile)
- 团签用户隐藏我的骑士卡页面和购买历史页面
- `orderNotActived`字段调整为仅个签用户有此字段