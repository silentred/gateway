# Roadmap

## Key Feature

### Day 0

- [x] 补充 Uint Test
- [x] 配置文件解析; Unit Test
- [x] 路由文本 Parser; Unit Test
- [x] 路由热更新

### Day 1

- [x] 签名验证 Guard 实现; Unit Test
- [x] 时间戳验证 Guard 实现; Unit Test

### Day 2

- [x] 后端异常熔断器 Reactor 实现; Unit Test

### Day 3

- [x] 基于权重的负载均衡; Unit Test

### Day 4

- [x] 监控接口实现, 包含: 接口维度的访问量, QPS;

## Other Feature

### Day 5

- [x] Reactor 接口添加 ObserveError(), 需要观察 RoundTripper 返回的 error
- [x] 支持 Consul 路由表

### Day 6

- [x] 防重放 Reactor 实现; Unit Test

### Day 8

- [x] 限流器 Reactor 实现, UID 维度; Unit Test

### Day 9

- [x] 支持 websocket 协议; Unit Test

### Day 10

- [x] 运维接口实现
- [x] Web管理界面

## Advanced Features

- [ ] 自定义 Guard Group
- [ ] Guard Group 管理接口
- [ ] 自定义 Reactor Group

- [x] 去除静态资源依赖
- [x] 使用Trie Tree判断路由前缀
- [x] 处理信号