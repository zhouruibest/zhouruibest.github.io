工具：kubehealthy

# 如何先于用户发现

## 可观测

LOG
Metric 根据Metric配置告警规则
Trace（主要是应用的）

# 黑盒测试

不去管指标看起来正不正常，而是亲自去用一下。
作为基础设施的维护者，模拟用户的各种使用平台的行为。
1）比如说发布一个应用，部署一份基线环境模板跑一下测试计划，这个可以定时跑
2）比如说公共组件升级了，部署专门的探测应用（自己编写的，测一下Log、Trace这些组件是否正常）
3）比如说刚扩容完成，在新扩容的节点上调度一下POD，测一下连通性等

# 巡检

对于可能的故障点定期地、持续的检查。有一些问题不是metric

（1）比如说内网CLB数量即将达到限额，可以调用云厂商接口提前预支
（2）比如证书有效期不足
（3）比如一级应用没有启用备灾

# 测试管理平台

通用的测试管理平台，能够作为rpc/http客户端发送请求。

发送请求本就是能够测试应用是否正常。
触发泳道管理平台执行巡检，平台毁掉测试管理平台，执行告警规则

## 用例
## 用例集管理
## 环境管理
## 测试计划管理
## 告警规则配置







