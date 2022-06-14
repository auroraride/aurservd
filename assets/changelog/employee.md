## 2022-06-14

#### 接口新增

- [E2003 骑手业务详情](http://localhost:5533/docs#tag/E/operation/EmployeeBusinessRider)
  - 通过判定`enterpriseName`是否存在来判定骑手是否团签用户，团签用户只能办理退租
  - 当`business = false`的时候显示当前状态无法办理业务
  - 当`status=2`寄存中的时候，办理寄存调整为`结束寄存`
- [E2004 寄存电池](http://localhost:5533/docs#tag/E/operation/EmployeeBusinessPause)
- [E2005 结束寄存电池](http://localhost:5533/docs#tag/E/operation/EmployeeBusinessContinue)
- [E1005 店员资料](http://localhost:5533/docs#tag/E/operation/EmployeeEmployeeProfile)
- [E2006 退租](http://localhost:5533/docs#tag/E/operation/EmployeeBusinessUnSubscribe)





<br />

## 2022-06-13

#### 接口新增

- [E1003 打卡预检](http://localhost:5533/docs#tag/E/operation/EmployeeAttendancePrecheck)
- [E1004 考勤打卡](http://localhost:5533/docs#tag/E/operation/EmployeeAttendanceCreate)
- [E2001 待激活骑士卡详情](http://localhost:5533/docs#tag/E/operation/EmployeeSubscribeInactive)



#### 接口调整

- [E1001 登录](http://localhost:5533/docs#tag/E/operation/EmployeeEmployeeSignin)
  - 新增字段`onduty`和`store`



<br />

## 2022-06-08

#### 接口新增

- [E1001 登录](http://localhost:5533/docs#tag/E/operation/EmployeeEmployeeSignin)
- [E1002 更新二维码](http://localhost:5533/docs#tag/E/operation/EmployeeEmployeeQrcode)