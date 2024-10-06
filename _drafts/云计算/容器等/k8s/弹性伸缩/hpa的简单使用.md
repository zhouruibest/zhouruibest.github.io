## 前提（基础）metric-server

metric-server: 列出集群中POD使用的资源
```sh
kubectl -n kube-system top pod # 可以查看指定命名空间下POD使用的资源（CPU和Memory）
```
## 1 给POD分配request资源和limit资源

## 2 创建hpa

```sh 
# 为deployment资源创建hpa，pod上限3个，最低1个， 在POD的平均CPU达到50%后开始扩容
kubectl autoscale deployment web --max=3 --min=1 --cpu-percent=50
```
## 完成

这样，一旦CPU使用率上来了，那么它会自动扩容；流量降下来之后，他会缩容