# Entree

API gateway in golang

![](https://img.shields.io/badge/language-golang-blue.svg)
![](https://img.shields.io/badge/license-MIT-000000.svg)

## Roadmap

- [x] Protocol (http, websocket)
- [x] Guard Iterface (sign, timestamp)
- [x] Metrics/Reactor (backend failover, circuit breaker, rate limiter, playback)
- [x] Route (match targets URL, pick a target)
- [x] config (file, consul, k8s)

## Usage

```
# make build
# ./entree -log_dir=/tmp -vv=3
```

```
Route not found, HTTP Code: 404
Backend failure, HTTP Code: 502
Guard reject, HTTP Code: 403
```

## Deployment

### build prerequisite

- node, npm
- go-bindata

### prerequisite

- Consul
- config.yaml

## TCP tuning

```
# Protection from SYN flood attack.
net.ipv4.tcp_syncookies = 1

# --------------------------------------------------------------------
# The following allow the server to handle lots of connection requests
# --------------------------------------------------------------------

# Increase number of incoming connections that can queue up
# before dropping
net.core.somaxconn = 50000

# Handle SYN floods and large numbers of valid HTTPS connections
net.ipv4.tcp_max_syn_backlog = 30000

# Increase the length of the network device input queue
net.core.netdev_max_backlog = 5000

# Increase system file descriptor limit so we will (probably)
# never run out under lots of concurrent requests.
# (Per-process limit is set in /etc/security/limits.conf)
fs.file-max = 200000

# Widen the port range used for outgoing connections
net.ipv4.ip_local_port_range = 10000 65000

# Disconnect dead TCP connections after 1 minute
net.ipv4.tcp_keepalive_time = 60

# Wait a maximum of 5 * 2 = 10 seconds in the TIME_WAIT state after a FIN, to handle
# any remaining packets in the network. 
net.ipv4.netfilter.ip_conntrack_tcp_timeout_time_wait = 5

# Allow a high number of timewait sockets
net.ipv4.tcp_max_tw_buckets = 2000000

# Timeout broken connections faster (amount of time to wait for FIN)
net.ipv4.tcp_fin_timeout = 10

# Let the networking stack reuse TIME_WAIT connections when it thinks it's safe to do so
net.ipv4.tcp_tw_reuse = 1

# Determines the wait time between isAlive interval probes (reduce from 75 sec to 15)
net.ipv4.tcp_keepalive_intvl = 15

# Determines the number of probes before timing out (reduce from 9 sec to 5 sec)
net.ipv4.tcp_keepalive_probes = 5
```

