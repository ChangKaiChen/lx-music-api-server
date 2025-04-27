## 概述

适用于 LX Music 的解析接口服务器，采用了分布式架构。它使用了Kitex、Hertz、Redis、Etcd、Elasticsearch、Fluentd、Kibana、Prometheus、Cadvisor、Grafana和APMPlus等技术。
## 特点

### 功能
- [x] 白名单key鉴权
- [x] wy
  - [x] url获取
  - [x] cookie保活
- [x] tx
  - [x] url获取
  - [x] musicKey刷新
- [ ] kg（抓不了包，暂无法实现）
- [ ] kw（抓不了包，暂无法实现）
- [ ] mg（抓不了包，暂无法实现）
### 特性

- 云原生：采用原生 Go 语言分布式架构设计，基于字节跳动的最佳实践。
- 高性能：支持异步 RPC、非阻塞 I/O。
- 可扩展性：基于模块化和分层结构设计，代码清晰易读。
- 可观测性：采用火山引擎的应用性能监控全链路版（简称APMPlus，**收费的**）、EFK日志管理方案（Elasticsearch、Fluentd、Kibana）和Grafana可视化监控。
## 项目结构

整体结构
```text
.
├─app                  # 各服务的实现
├─cmd                  # 各服务启动入口
├─global               # 全局配置
├─idl                  # 接口定义
├─kitex_gen            # kitex生成的代码
└─pkg
    ├─cache            # 简易redis缓存
    ├─consts           # 一些常量
    ├─limiter          # 简易key限流器
    ├─logger           # 简易日志器
    ├─response         # 封装的响应
    └─utils           
```
网关
```text
./app/gateway
├─config               # 配置文件
├─handler              # 请求处理器
├─middleware           # 中间件
├─router               # 路由
└─rpc                  # rpc调用
```
微服务（wy模块）
```text
./app/wy
├─config               # 配置文件
├─crypto               
├─refresh              # cookie保活
└─rpc                  # rpc接口具体实现
```
## 中间件

```text
middleware.RecoveryMW()
cache.NewCacheByRequestURI()    # 缓存HTTP响应
middleware.SentinelMW()         # 限流，默认100QPS
middleware.LimiterMW()          # 针对单key限流，默认5s内允许10次请求，超出则惩罚1h
```
## 可视化示例

![Kibana](https://github.com/ChangKaiChen/lx-music-api-server/tree/main/images/Kibana.png)
![APMPlus](https://github.com/ChangKaiChen/lx-music-api-server/tree/main/images/APMPlus.png)
![Grafana](https://github.com/ChangKaiChen/lx-music-api-server/tree/main/images/Grafana.png)
## 部署

---
先修改全局配置文件global.yaml和各服务的配置文件

etcd默认prefix为kitex/registry-etcd，仅在authEnable为true时需要注意
### 直接启动
```text
go mod tidy
go run main.go
```
### docker
global中isLocal要设为false
```text
docker-compose up -d
```
## 备注

---
本项目仅为练习项目，代码可能存在不规范之处（野路子），欢迎各位大佬提出宝贵意见并给予修正建议。
## 参考项目

---
https://github.com/MeoProject/lx-music-api-server

https://github.com/west2-online/DomTok/tree/main