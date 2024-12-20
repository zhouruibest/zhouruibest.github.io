OpenTelemetry 可以用于从应用程序收集数据。它是一组工具、API 和 SDK 集合，我们可以使用它们来检测、生成、收集和导出遥测数据（指标、日志和追踪），以帮助分析应用的性能和行为。

![](./opentelemetry全景图.png)

OpenTelemetry 是：

- 开源的
- 受到可观测领域行业领导者的采用和支持
- 一个 CNCF 项目
- 与供应商无关的

OpenTelemetry 包括可观测性的三个支柱：追踪、指标和日志。（本文将重点关注追踪）

- 分布式追踪是一种跟踪服务请求在分布式系统中从开始到结束的方法。
- 指标是对一段时间内活动的测量，以便了解系统或应用程序的性能。
- 日志是系统或应用程序在特定时间点发生的事件的文本记录。

分布式链路追踪数据流

![](./opentelemetry分布式链路追踪数据流.png)

# 传播器

可以将 W3C tracecontext、baggage 和b3 等传播器（Propagators）添加到配置中。

> 不同的传播器定义特定的行为规范，以便跨进程边界传播带上上下文数据。

- Trace Context：用于在 HTTP headers 中编码 trace 数据，以便在不同的服务间传递这些数据。
- Baggage：用于在 span 之间传递键值对数据，例如用户 ID、请求 ID 等。
- B3：用于在 HTTP headers 中编码 trace 数据，以便在不同的服务间传递这些数据（主要用于 Zipkin 或其兼容的系统）。

# 采样

采样是一种通过减少收集和发送到后端的追踪样本数量来控制 OpenTelemetry 引入的噪声和开销的机制。

可以告诉 OpenTelemetry 根据要发送的追踪/流量的数量执行采样。（比如只采样 10% 的追踪数据）。

![](./数据采样.png)

两种常见的采样技术是头采样和尾采样。

# OpenTelemetry 协议（OTLP）
OpenTelemetry 协议（OTLP）规范描述了遥测数据在遥测源、收集器和遥测后端之间的编码、传输和传递机制。

每种语言的 SDK 都提供了一个 OTLP 导出器，可以配置该导出器来通过 OTLP 导出数据。然后，OpenTelemetry SDK 会将事件转换为 OTLP 数据。

OTLP 是代理（配置为导出器）和收集器（配置为接收器）之间的通信。

# OpenTelemetry Collectors

![](./otel-collector.png)

收集器是 OpenTelemetry 的一个组件，它接收遥测数据（span、metrics、logs 等），处理（预处理数据）并导出数据（将其发送到想要的通信后端）。

otel-collector中分三个阶段处理数据

![](./otel-collector中分三个阶段处理数据.png)

## 阶段1 Receivers

接收器 Receivers 是数据进入收集器的方式，可以是推送或拉取。OpenTelemetry 收集器可以以多种格式接收遥测数据。

![](./otel-collector-receivers.png)

以下是接收器在端口 4317(gRPC) 和 4318(http) 上接受 OTLP 数据的配置示例：

``` js
otlp:
  protocols:
    http:
    grpc:
      endpoint: "0.0.0.0:4317"
```

同样下面的示例，它可以以 Jaeger Thrift HTTP 协议方式接收遥测数据。

```js
jaeger: # Jaeger 协议接收器
  protocols: # 定义接收器支持的协议
    thrift_http: # 通过 Jaeger Thrift HTTP 协议接收数据
      endpoint: "0.0.0.0:14278"
```

## 阶段2 Processors

一旦接收到数据，收集器就可以处理数据。处理器在接收和导出之间处理数据。处理器是可选的，但有些是推荐的。

比如 batch 处理器是非常推荐的。批处理器接收跨度、指标或日志，并将它们放入批次中。批处理有助于更好地压缩数据，减少传输数据所需的传出连接数量。该处理器支持基于大小和时间的批处理。

```js
processors:
  batch:
```

需要注意的是配置处理器并不会启用它。需要通过 service 部分的 pipelines 启用。

```js
service:
  pipelines:
    traces:
      receivers: [jaeger]
      processors: [batch]
      exporters: [zipkin]
```

## 阶段3 Exporters

为了可视化和分析遥测数据，我们还需要使用导出器。导出器是 OpenTelemetry 的一个组件，也是数据发送到不同系统/后端的方式。

比如 console exporter 是一种常见的导出器，对于开发和调试任务非常有用，他会将数据打印到控制台。

在 exporters 部分，可以添加更多目的地。例如，如果想将追踪数据发送到 Grafana Tempo，只需添加如下所示的配置：

```js
exporters:
  logging:
  otlp:
    endpoint: "<tempo_endpoint>"
    headers:
      authorization: Basic <api_token>
```

当然最终要生效也需要在 service 部分的 pipelines 中启用。

```js
jsservice:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [logging, otlp]
```

OpenTelemetry 附带了各种导出器，在 OpenTelemetry 收集器 Contrib 存储库中可以找到。

## Extensions

扩展主要适用于不涉及处理遥测数据的任务。扩展的示例包括健康监控、服务发现和数据转发。扩展是可选的。

```js
extensions:
  health_check:
  pprof:
  zpages:
  memory_ballast:
    size_mib: 512
```

# OpenTelemetry Collector 部署模式/策略

## Agent 模式

在这种情况下，OpenTelemetry 检测的应用程序将数据发送到与应用程序一起驻留的（收集器）代理。然后，该代理程序将接管并处理所有来自应用程序的追踪数据。

收集器可以通过 sidecar 方式部署为代理，sidecar 可以配置为直接将数据发送到存储后端。

![](./简单的agent模式部署.png)

## Gateway 模式

还可以决定将数据发送到另一个 OpenTelemetry 收集器，然后从（中心）收集器进一步将数据发送到存储后端。在这种配置中，我们有一个中心的 OpenTelemetry 收集器，它使用 deployment 模式部署，具有许多优势，如自动扩展。

![](./带collector的gateway模式.png)

使用中心收集器的一些优点是：

- 消除对团队的依赖
- 强制执行批处理、重试、加密、压缩的配置/策略
- 在中心位置进行身份验证
- 丰富的元数据信息
- 进行抽样决策
- 通过 HPA 进行扩展

# 在 Kubernetes 上部署 OpenTelemetry Collector 时可以使用的模式

## 在k8s中的部署模式-作为sidecar

![](./在k8s中的部署模式-作为sidecar.png)

## daemonset 模式：

Agent 作为 DaemonSet，这样我们每个 Kubernetes 节点就有一个代理 pod。

![](./在k8s中的部署模式-作为daemonset.png)


# 参考
https://cloud.tencent.com/developer/article/2327988