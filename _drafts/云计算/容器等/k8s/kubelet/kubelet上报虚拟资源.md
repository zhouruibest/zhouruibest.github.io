# 概述

某些场景下需要给节点增加一种资源，如果开源或者社区没有这种资源的 Operator 或者 Device Plugin 做资源的上报，或者说部署这些组件的成本比较高，在测试和验证的阶段如果还不需要做到这种程度，可以考虑通过 PATCH 的方式，给节点添加一种类型的资源。

# 操作

```sh
# 节点上操作
kubectl proxy
curl -X PATCH \
  -H "Content-Type: application/json-patch+json" \
  -d '[{"op":"add","path":"/status/capacity/example.com~1dongle","value":"4"}]' \
  http://localhost:8001/api/v1/nodes/node1/status
```

# 结果

这样之后，Node 的 Status 字段就会被增加了上述的资源。

```yaml
status:
  allocatable:
    cpu: "6"
    ephemeral-storage: "66051905018"
    example.com/dongle: "4"
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 65212372Ki
    pods: "110"
```

# 参考

[Advertise Extended Resources for a Node](https://kubernetes.io/docs/tasks/administer-cluster/extended-resource-node/)