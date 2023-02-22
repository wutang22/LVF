# simple-demo

## 抖音项目服务端简单示例

工程无其他依赖，直接编译运行即可

```shell
air
```

### 功能说明

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

### 代码结构

```
LVF
├─ .air.toml
├─ controller
│  ├─ comment.go
│  ├─ common.go
│  ├─ favorite.go
│  ├─ feed.go
│  ├─ message.go
│  ├─ publish.go
│  ├─ relation.go
│  └─ user.go
├─ go.mod
├─ go.sum
├─ main.go
├─ mydb
│  └─ mydb.go
├─ public
├─ README.md
├─ router.go
├─ service
│  └─ message.go
└─ test
   ├─ base_api_test.go
   ├─ common.go
   ├─ interact_api_test.go
   ├─ message_server_test.go
   └─ social_api_test.go

```
