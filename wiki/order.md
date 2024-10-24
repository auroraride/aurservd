## 问题列表

- 计算天数规则: 现行进一法

## 订单逻辑

> 当前骑士卡在`rider`表中进行更新

- 骑手购买(可能携带押金)支付保存订单`order`表同时计算是否需要保存提成`commission`表
- 骑手未激活订单判定条件
    - 已支付
    - 未退款
    - 开始时间为空
    - 类型(或)
        - 新签
        - 重签
        - 更改电池
- 当前骑士卡几种状态
    - 计费中 (start不为空, end为空)
    - 暂停中 (pause_at不为空)
    - 已逾期 (时间到但未归还)
    - 已退租 (end为空, 且children为空)
- 订单结束日期
    - 空值
        - 未开始
        - 计费中
    - 非空
        - 更改电池
        - 退款
        - 暂停计费

### 暂停 / 续签

> 需要更新父订单ID`parent_id`

### 退款

- 提交后, 需要财务审核(2022-05-30)
- 申请退款修改状态为 退款申请中, 退款成功后订单状态改为已退款

申请退款需满足条件:
1. 骑手骑士卡订单
2. 未使用
3. 订单状态为已支付

### 查找使用中的骑士卡

1. 订单为新签 / 重签
2. 已经激活且未退还(`start_at`不为空且`end_at`为空)
3. 未欠费(可使用时间不小于0)

### 逾期

条件:

1. 开始计费(start不为空)
2. 未归还电池(end为空)
3. 骑士卡剩余时间小于等于0

## 提成逻辑

1. 骑手新购/重构保存`commission`表
2. 当骑手激活的时候`commission`添加当前激活的员工ID, 此时提成生效
3. 当骑手申请退款时(只有未使用订单才可以退款), 提成订单需要删除