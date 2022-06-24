## 2022-06-25

#### 接口新增

- [R5003 取消救援](http://localhost:5533/docs#tag/R/operation/RiderAssistanceCancel)





<br />

## 2022-06-23

#### 接口新增

- [R5001 获取救援原因]([极光出行API](http://localhost:5533/docs#tag/R/operation/RiderAssistanceBreakdown))
- [R5002 发起救援]([极光出行API](http://localhost:5533/docs#tag/R/operation/RiderAssistanceCreate))





<br />

## 2022-06-21

#### 接口调整

> 删除所有涉及到`voltage<number>(电压)`和`capacity<number>(容量)`的字段，统一改为`model<string>(电池型号)`

- 删除`R3001 电压型号(/rider/v1/battery/voltage)` 改用 [C4 获取生效中的电池型号](http://localhost:5533/docs#tag/C/operation/ManagerBatteryModel)
- [R3010 企业骑手获取可用电池](http://localhost:5533/docs#tag/R/operation/RiderEnterpriseListBattery)`(/rider/v1/enterprise/voltage)`URL改为`/rider/v1/enterprise/battery`
- 删除 [R3008 骑士卡购买历史](http://localhost:5533/docs#tag/R/operation/RiderOrderList) 中字段 `models`
- 删除 [R3009 订单详情](http://localhost:5533/docs#tag/R/operation/RiderOrderDetail) 中字段`models`





<br />

## 2022-06-19

#### 接口新增

- [R3010 订单支付状态](http://localhost:5533/docs#tag/R/operation/RiderOrderStatus)





<br />

## 2022-06-15

#### 接口调整

- 原`R3002 获取骑士卡`修改为`R3002 新购骑士卡`



#### 接口新增

- [R3003 续费骑士卡](http://localhost:5533/docs#tag/R/operation/RiderPlanRenewly)





<br />

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