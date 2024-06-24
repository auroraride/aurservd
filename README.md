# 极光出行


## 代码查重
```shell
go install github.com/mibk/dupl@master
dupl -html $(find . -name '*.go' -not -path "./internal/ent/*"  -not -path "./assets/docs/*") > ~/Downloads/aurservd-dupl.html
dupl -html -t 50 $(find ./app/service ./app/model ./app/biz -name '*.go') > ~/Downloads/aurservd-dupl.html
```


## 项目依赖第三方

#### redis hash过期 (未使用)

- [TairHash](https://github.com/alibaba/TairHash/blob/develop/README-CN.md)
- [tair-go](https://github.com/alibaba/tair-go/blob/main/README-CN.md)

编译指令

```shell
mkdir -p /var/lib/redis/modules
rm -rf build
mkdir build
cd build
cmake -DSLAB_MODE=yes ../
make -j
mv ../lib/tairhash_module.so /var/lib/redis/modules
```

配置文件中配置: `loadmodule /var/lib/redis/modules/tairhash_module.so`



## 注意

- 积分消耗: 若不支付, 预消耗20分钟后释放



## 调试

```shell
dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient attach $PID
```



## 待办

- 从`struct`创建或更新
    - https://github.com/ent/ent/issues/761
    - https://entgo.io/docs/templates/#examples



## 常用命令

- 服务端debug: https://juejin.cn/post/7035910722382987271
    - `GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/release/aurservd cmd/aurservd/main.go`
    - `./dlv --listen=:3333 --headless=true --api-version=2 --continue --accept-multiclient exec ./aurservd -- server`
    - `./dlv --listen=:3333 --headless=true --api-version=2 --accept-multiclient exec ./aurservd -- server`



## swag

> 新增了trans作为备注
> https://github.com/liasica/swag

```shell
git clone git@github.com:liasica/swag.git
cd swag/cmd/swag
go install
```

```bash
go run ./cmd/ent init TABLE
go run ./cmd/ent generate
```

> 转换为OpenAPI3:

- https://github.com/getkin/kin-openapi
- https://github.com/swaggo/swag/issues/386



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
