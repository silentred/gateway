# 文档

## 支持协议

HTTP, websocket

## 签名认证

采用类似 aliyun 方式, sign = Hash(HTTP query string in order + timestamp + random string + appID)

## 时间戳范围

15分钟以内, 否则返回错误

## 防重放

1 分钟内 random string 不能重复

## 负载均衡

API 和 后端服务 支持 一对多 配置，负载均衡支持 round-robin, random

## 后端异常

后端返回 HTTP Code > 500 则统计为错误数据, 超过某个阈值时, 屏蔽此接口 3 分钟

## 限流

根据 UID 限流: 单个UID每秒限制
根据 接口URI 限流: 单个接口每秒限制

## 服务降级

接口配置一个降级有损服务，当接口被降级时，返回有损服务内容

## 监控

支持 Prometheus 格式监控接口

## 路由发现

路由存储支持 file, consul, k8s

