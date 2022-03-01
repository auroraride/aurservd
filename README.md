# 极光出行

## 待办
- 从`struct`创建或更新
  - https://github.com/ent/ent/issues/761
  - https://entgo.io/docs/templates/#examples

## 常用命令

```bash
go run ./cmd/ent init TABLE
go run ./cmd/ent generate
```

## 第三方服务列表
- 阿里云
  - SLS: 存储日志
  - SMS: 短信服务
  - 语音通知
- E签宝
- 百度人脸
  - 实名认证
  - 人脸比对

## 测试

### 推送

```go
mob.NewPush().SendMessage(mob.Req{
    RegId:    "65kzlib1miggt8g",
    Platform: mob.PlatformAndroid,
    Content:  "测试推送内容",
    Title:    "测试推送",
    MessageData: []mob.MessageData{
        {Key: "key", Value: "val"},
    },
    Channel: mob.ChannelSystem,
})
```