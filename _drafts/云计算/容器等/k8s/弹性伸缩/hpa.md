# Kubernetes 默认支持根据容器的 CPU 和内存的使用率扩缩容

举例：
```yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: web-app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: web-app
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      targetAverageUtilization: 80
```
这里 scaleTargetRef 表示要操作的 Kind 对象 Deployment，其中副本最小数量为1，最大为10。而对条件的定义就是 metrics字段，这个示例里只指定了CPU，表示当 CPU 平均使用率达到80% 的时候就开始扩容。

# 也可以对扩缩容行为进行定义

举例
```yaml
behavior:
  scaleDown:
    policies:
    - type: Pods
      value: 4
      periodSeconds: 60
    - type: Percent
      value: 10
      periodSeconds: 20
```

这里只对 scaleDown 缩容行为进行了定义，一共有两个策略。第一个策略（Pods）允许在一分钟内最多缩容 4 个副本，第二个策略（Percent） 允许在一分钟内最多缩容当前副本个数的百分之十

# 现在 HPA 支持从其他的 API 中获取指标来进行扩容

HPA 控制器会从 apiservices 管理的 API 中获取一些指标，然后根据定义好指标的阈值来触发一些扩缩容。

```
$ kubectl get apiservices.apiregistration.k8s.io
NAME                                   SERVICE                                         AVAILABLE   AGE
v1beta1.custom.metrics.k8s.io          monitoring/prometheus-adapter                   True        30h
v1beta1.external.metrics.k8s.io        addons-system/keda-operator-metrics-apiserver   True        5h37m
v1beta1.metrics.k8s.io                 kube-system/metrics-server
```

- 对于资源指标，使用 metrics.k8s.io API， 一般由 metrics-server 提供。 它可以作为集群插件启动。
- 对于自定义指标，使用 custom.metrics.k8s.io API。 它由其他“适配器（Adapter）” API 服务器提供。从上面的命令输出可以看出，这里是由 prometheus-adapter 来提供的。
- 对于外部指标，将使用 external.metrics.k8s.io API。

其中 Prometheus-Adapter 是一款用于将 Prometheus 指标转换为 Kubernetes 自定义指标的工具。它可以用于将 Prometheus 监控的应用程序的指标暴露给 Kubernetes Horizontal Pod Autoscaler (HPA) 等工具。

工作经历三个阶段：

Prometheus-Adapter 会定期（默认每 30 秒）从 Prometheus 服务器拉取指标数据。
Prometheus-Adapter 会将拉取到的指标数据转换为 Kubernetes API Server 可以理解的形式。
Prometheus-Adapter 会将转换后的指标数据暴露给 Kubernetes API Server。
从上面这三个阶段可以看出HPA方案有存在一个很大的缺点，那就是K8S集群无法实时根据负载情况动态扩缩容，存在一定的延时（默认30秒）。