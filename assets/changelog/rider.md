## 2022-06-08

#### APP

- 新增团签用户激活逻辑



#### 接口新增

- [R30010 企业骑手获取可用电池](http://localhost:5533/docs#tag/R/operation/RiderEnterpriseListVoltage)
- [R30011 企业骑手选择电池](http://localhost:5533/docs#tag/R/operation/RiderEnterpriseSubscribe)
- [R30012 企业骑手订阅激活状态](http://localhost:5533/docs#tag/R/operation/RiderEnterpriseSubscribeStatus) (长连接轮询，30s轮询一次，返回结果为`TRUE`或出现错误时中断轮询)



<br />

## 2022-06-07

#### APP

- 骑手信息新增字段`enterprise`所属企业，仅团签用户有此字段，此字段可进行是否团签用户的判定，详情参考文档[R10006 获取个人信息](http://localhost:5533/docs#tag/R/operation/RiderRiderProfile)
- 团签用户隐藏我的骑士卡页面和购买历史页面
- `orderNotActived`字段调整为仅个签用户有此字段