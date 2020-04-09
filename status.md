# 2017-05-26

- [x] 运维接口
- [x] 后台 Web UI

# 2017-05-17

- [x] Reactor 接口添加 ObserveError(), 需要观察 RoundTripper 返回的 error
- [x] 支持 Consul 路由表
- [x] 防重放 Reactor 实现; Unit Test
- [x] 限流器 Reactor 实现, UID 维度; Unit Test
- [x] 支持 websocket 协议

# 2017-05-11

- [x] 监控接口实现, 包含: 接口维度的访问量, QPS;
- [x] 关键位置日志补全; 整体流程检查;

目前功能包含 签名检查，时间戳检查，后端异常熔断，后端负载均衡(round-robin, random, weight-based)，
可以部署测试了。

# 2017-05-10

进展比较顺利，目前提前完成了一天的任务.

- [x] 补充 Uint Test
- [x] 配置文件解析; Unit Test
- [x] 路由文本 Parser; Unit Test
- [x] 路由热更新
- [x] 签名验证 Guard 实现; Unit Test
- [x] 时间戳验证 Guard 实现; Unit Test
- [x] 后端异常熔断器 Reactor 实现; Unit Test
- [x] 基于权重的负载均衡; Unit Test

## Test Coverage

```
go test $(go list ./...| grep -vE 'vendor') -cover
?       github.com/silentred/gateway     [no test files]
ok      github.com/silentred/gateway/config      1.023s  coverage: 69.9% of statements
ok      github.com/silentred/gateway/guard       0.014s  coverage: 92.6% of statements
ok      github.com/silentred/gateway/proxy       0.025s  coverage: 66.2% of statements
ok      github.com/silentred/gateway/reactor     0.016s  coverage: 95.8% of statements
ok      github.com/silentred/gateway/reactor/circuit_breaker     1.525s  coverage: 98.1% of statements
ok      github.com/silentred/gateway/route       0.012s  coverage: 78.6% of statements
ok      github.com/silentred/gateway/util        0.042s  coverage: 25.6% of statements
```

## Benchmark

本地测试，10 线程, 10s 压测, 网关没加业务逻辑，结果和 直接压测后端服务的结果 几乎一致。

```
# k6 run --vus 10 --duration 10s - <node.js

✗ 0.00% - is status 200

checks................: 0.00%
data_received.........: 9.6 MB (956 kB/s)
data_sent.............: 5.9 MB (587 kB/s)
http_req_blocked......: avg=920ns max=5.78ms med=0s min=0s p90=0s p95=0s
http_req_connecting...: avg=907ns max=5.73ms med=0s min=0s p90=0s p95=0s
http_req_duration.....: avg=1.54ms max=58.08ms med=1.22ms min=203.11µs p90=2.4ms p95=3.04ms
http_req_receiving....: avg=175.01µs max=57.13ms med=21.79µs min=6.74µs p90=116.13µs p95=182.27µs
http_req_sending......: avg=23.56µs max=17.4ms med=11.57µs min=5.64µs p90=36.12µs p95=65.21µs
http_req_waiting......: avg=1.34ms max=25.37ms med=1.15ms min=175.63µs p90=2.27ms p95=2.85ms
http_reqs.............: 58674 (5867.4/s)
iterations............: 58674 (5867.4/s)
vus...................: 10
vus_max...............: 10
```
