## 使用方法

```
下载源码
go get gitlab.luojilab.com/igetserver/entree

编译
cd $GOPATH/src/gitlab.luojilab.com/igetserver/entree
make build

配置文件: 目录下 config.yaml, route.cfg
./entree -log_dir=/tmp -vv=3

在 /tmp 目录下会生成 entree.INFO 日志文件
```

config.yaml

```
listen: :8088
recover:
  stackSize: 4096
  disableStackAll: false
  disablePrintStack: false
proxy:
  strategy: rnd
  matcher: prefix
  noRoute: 404
  maxConn: 10000
  bufSize: 128
  shutdownWait: 5s
  dialTimeout: 30s
  flushInterval: 1s
  respHeaderTimeout: 5s
  keepAliveTimeout: 30s
route:
  backend: file
  filePath: ./route.cfg
guard:
  # for sign
  secret: test123
  # for timestamp range
  timeRange: 600
reactor:
  circuitBreaker:
    thresholdCount: 60
    thresholdDuration: 1m
    blockDuration: 1m
```

route.cfg

```
Route:www.luojilab.com Prefix:/v1/A/hello Service:serviceA Strip:/v1/A Targets(127.0.0.1:8080,100)
Route:s2.luojilab.com Prefix:/v1/B/world Service:serviceB Strip:/v1/B Targets(127.0.0.1:8081,100)
Route:s3.luojilab.com Prefix:/v1/C/hello Service:serviceC Strip:/v1/C Targets(127.0.0.1:8082,100)
Route:192.168.16.230:8088 Prefix:/v1/hello Service:localSvc Strip:/v1 Targets(127.0.0.1:8080,20)
```

## 功能点

### 转发请求
route.cfg 是路由表配置，一行一条，负责把前端的请求转发到后端， 例如
```
Route:192.168.16.230:8088 Prefix:/v1/hello Service:localSvc Strip:/v1 Targets(127.0.0.1:8080,20)
```
作用是：把请求地址为 192.168.16.230:8088, URI 前缀为 /v1/hello 的请求，转发到 127.0.0.1:8080，并把前缀 /v1 去掉。
service的名称为 localSvc, 作为 service的ID使用。Targets的第二个参数是权重。

例如:

```
curl http://192.168.16.230:8088/v1/hello/sdf 
的请求会转发到 http://127.0.0.1/hello/sdf 
```

如果路由没有匹配到后端服务，那么返回 HTTP Code = 404

### 签名

前端的请求需要有签名，签名方案在文档中, http://wiki.luojilab.com/display/IG/Gateway+Signing+Method 。
签名不合法的返回 HTTP Code = 403.
Secret的值在 config.yaml 中配置， guard.secret 

### 时间戳

签名规则中有提到 请求的 Header 中包含一个 X-Timestamp 字段，该字段为 Unix timestamp, 
必须在服务器时间的前后10分钟内, 超过范围则返回 HTTP Code = 403
时间容忍范围在 config.yaml 中， guard.timeRange

### 后端异常熔断

条件: 一定时间段 x 内，后端返回的 HTTP code >= 500 的个数超过了 y, 
行为: 则在 z 时间段内按一定比例 r 直接拒绝前端的请求。
过了 z 时间段后，比例r 会减小。

x, y, z 可以配置
x = reactor.circuitBreaker.thresholdDuration
y = reactor.circuitBreaker.thresholdCount
z = reactor.circuitBreaker.blockDuration

r 的具体规则为：
第一次条件达成时, r = 1/2， 会有 1/2 请求直接被拒绝；
假设在这段时间 z 内， 条件再一次被达成， 那么 r = ( 1 - 1/4 ) = 3/4 ， 有 3/4 请求 会被拒绝， 1/4 被放行.
依次类推。
如果在时间段 z 内，条件没有达成，那么 r 会恢复, 100% 的请求会被放行.


## Cases

### 签名

- [x] Go client 封装, 并发访问两个service, 观察是否有路由错误

### 时间戳

- [x] 针对时间戳范围内三个边界时间点的请求：并发3个请求，不同时间戳。

### 熔断

- [x] 到达熔断条件后，是否生效 1/2 拦截; 再次达到熔断条件后，是否生效 3/4 拦截；依次类推： 并发访问，后端返回错误，统计拦截率。
- [x] 熔断恢复后是否正常处理请求，测试多个接口(被熔断和未被熔断接口)：并发请求多个接口，观察返回。

### 速率

- [x] UID访问次数到达上限后的拦截行为：并发请求，判断是否拦截
- [x] 判断是否恢复正常：并发少量请求，观察是否拦截(被拦截UID和其他UID)

### Reply

- [x] 判断是否拦截：并发请求，超过数量后是否拦截
- [x] 判断是否恢复正常：并发请求，同一个X-Nonce是否被拦截。