### 说明
> 接口采用非标准Restful API，所有http返回代码均为`200`，当返回为非`200`时应为network错误，需要及时排查。
> 接口返回说明查看 **[返回](#返回)**

### 认证
项目接口使用简单认证，认证方式为`header`中添加对应的认证`token`
|  header   |  类型  |  接口  |
| :-----: | :----: | :--: |
|  X-Rider-Token   |  string   |  骑手API  |
| X-Manager-Token | string |  后台API  |
|  X-Employee-Token   | string |  员工API  |

### 返回
一个标准的返回应包含以下结构

|  字段   |  类型  |  必填  |  说明  |
| :-----: | :----: | :--: | :--: |
|  code   |  int   |  是  |  返回代码  |
| message | string |  是  |  返回消息  |
|  data   | object |  是  |  返回数据  |

`code`代码取值说明

| 二进制 | 十六进制 | 说明 |
| :----: | :------: | :--: |
| 0  |  0x000  | 请求成功 |
| 256 |  0x100  | 请求失败 |
| 512 |  0x200  | 需要认证 (需要登录) |
| 768 |  0x300  | 没有权限 |
| 1024 |  0x400  | 资源未获 |
| 1280 |  0x500  | 未知错误 |
| 1536 |  0x600  | 需要实名 |
| 1792 |  0x700  | 需要验证 (更换设备需要人脸验证) |
| 2048 |  0x800  | 需要联系人 |
| 2304 |  0x900  | 请求过期 |


比如：
> 默认成功返回
```json
{
  "code": 0,
  "message": "OK",
  "data": {
    "status": true
  }
}
```