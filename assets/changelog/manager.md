## 2022-07-05

#### 接口新增

- [MD001 权限列表](http://localhost:5533/docs#tag/M/operation/ManagerPermissionList)







<br />

## 2022-07-04

#### 接口新增

- [M7012 创建骑手跟进](http://localhost:5533/docs#tag/M/operation/ManagerRiderFollowUpCreate)
- [M7013 获取骑手跟进](http://localhost:5533/docs#tag/M/operation/ManagerRiderFollowUpList)





<br />

## 2022-07-01

#### 接口新增

- [MB007 筛选企业](http://localhost:5533/docs#tag/M/operation/ManagerSelectionEnterprise)
- [M8003 换电列表](http://localhost:5533/docs#tag/M/operation/ManagerExchangeList)
- [M5010 电柜数据表](http://localhost:5533/docs#tag/M/operation/ManagerCabinetData)





<br />

## 2022-06-29

#### 接口新增

- [M7011 删除骑手](http://localhost:5533/docs#tag/M/operation/ManagerRiderDelete)







<br />

## 2022-06-24

#### 接口新增

- [MC006 拒绝救援](http://localhost:5533/docs#tag/M/operation/ManagerAssistanceRefuse)





<br />

## 2022-06-23

#### 接口新增

- [MC001 救援列表](http://localhost:5533/docs#tag/M/operation/ManagerAssistanceList)
- [MC002 救援详情](http://localhost:5533/docs#tag/M/operation/ManagerAssistanceDetail)
- [MC003 附近门店](http://localhost:5533/docs#tag/M/operation/ManagerAssistanceNearby)
- [MC004 分配救援任务](http://localhost:5533/docs#tag/M/operation/ManagerAssistanceAllocate)
- [MC005 救援免费](http://localhost:5533/docs#tag/M/operation/ManagerAssistanceFree)





<br />

## 2022-06-21

#### 接口调整

> 删除所有涉及到`voltage<number>(电压)`和`capacity<number>(容量)`的字段，统一改为`model<string>(电池型号)`

- 删除`M4003 列举电池电压型号(/manager/v1/battery/voltages) `改用 [C4 获取生效中的电池型号](http://localhost:5533/docs#tag/C/operation/ManagerBatteryModel)

- 删除 [M8001 订单列表](http://localhost:5533/docs#tag/M/operation/ManagerOrderList) 中字段`models`

  





<br />

## 2022-06-20

#### 接口新增

- [M9012 结账](http://localhost:5533/docs#tag/M/operation/ManagerStatementBill)
- [MA014 店员业绩](http://localhost:5533/docs#tag/M/operation/ManagerEmployeeActivity)
- [MA015 启用/禁用店员](http://localhost:5533/docs#tag/M/operation/ManagerEmployeeEnable)
- [M3010 修改合同底单](http://localhost:5533/docs#tag/M/operation/ManagerBranchSheet)
- [MA016 工作流](http://localhost:5533/docs#tag/M/operation/ManagerAttendanceList)
- [MB006 筛选网点](http://localhost:5533/docs#tag/M/operation/ManagerSelectionBranch)



#### 接口调整

- M5002 查询电柜
  - 新增电池型号筛选
- MA012 列举店员
  - 修改返回数据





<br />

## 2022-06-19

#### 接口调整

- `M9003 列举企业` / `M9004 企业详情` 新增字段 `unsettlement` (未结算天数, 预付费企业此字段强制为0)



#### 接口新增

- [M9011 获取账单](http://localhost:5533/docs#tag/M/operation/ManagerStatementGetBill)





<br />

## 2022-06-18

#### 接口新增

- [MB001 筛选骑士卡](http://localhost:5533/docs#tag/M/operation/ManagerSelectionPlan)
- [MB002 筛选骑手](http://localhost:5533/docs#tag/M/operation/ManagerSelectionRider)
- [MB003 筛选门店](http://localhost:5533/docs#tag/M/operation/ManagerSelectionStore)
- [MB004 筛选店员](http://localhost:5533/docs#tag/M/operation/ManagerSelectionEmployee)
- [MB005 筛选启用的城市](http://localhost:5533/docs#tag/M/operation/ManagerSelectionCity)





<br />

## 2022-06-17

#### 接口调整

- 设置项新增`EXCEPTION`(物资异常)





<br />

## 2022-06-13

#### 接口新增

- [M1016 物资管理概览](http://localhost:5533/docs#tag/M/operation/ManagerStockOverview)
- [M1018 可调拨物资清单](http://localhost:5533/docs#tag/M/operation/ManagerInventoryTransferable) *调拨物资时使用*



#### 接口调整

- 门店相关接口返回门店二维码属性`qrcode`





<br />

## 2022-06-12

#### 接口调整

- 骑士卡相关，骑士卡多日期调整，详情参考接口文档，涉及接口如下：
  - [M6001 创建骑士卡](http://localhost:5533/docs#tag/M/operation/PlanCreate)
  - [M6004 列举骑士卡](http://localhost:5533/docs#tag/M/operation/PlanList)



#### 接口新增

- [M4003 列举电池电压型号](http://localhost:5533/docs#tag/M/operation/ManagerBatteryListVoltages)
- [M1012 物资设定创建或更新](http://localhost:5533/docs#tag/M/operation/ManagerInventoryCreateOrModify)
- [M1013 列举物资设定](http://localhost:5533/docs#tag/M/operation/ManagerInventoryList)
- [M1014 删除物资设定](http://localhost:5533/docs#tag/M/operation/ManagerInventoryDelete)
- [M1015 调拨物资](http://localhost:5533/docs#tag/M/operation/ManagerStockCreate)
- [M1017 门店物资详细](http://localhost:5533/docs#tag/M/operation/ManagerStockList)



<br />

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

