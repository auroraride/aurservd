# 购车需求

### 项目架构

1. `./internal`目录禁止任何外部被引用
2. 与原项目公用数据库 `/internal/ent`
3. 数据表均以`purchase_`开头


### 表
- `purchase_order` 订单表
  - 用户支付成功后，修改字段`installment_index` 和 `next_date`

- `goods` 商品表
  - 增加字段 `installment` 分期方案配置，`[][]float64` 格式，保存记录时会通过 [Schema Hooks](https://entgo.io/docs/hooks/#schema-hooks) 校验和排序
