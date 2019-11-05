# eHttp

本包使用 gin 框架编写 http 服务.

## 运行服务
```
cd goms/eHttp/cmd

go run . & 

```

## 测试API
http
```
# 使用 http 方法 /call/ping
curl  localhost:8080/call/ping

# 使用 http 方法 /call/ping, 带参数 message=xxx
curl  localhost:8080/call/ping?message=xxx
```