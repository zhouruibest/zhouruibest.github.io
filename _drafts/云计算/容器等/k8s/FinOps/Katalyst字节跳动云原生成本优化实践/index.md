# 资源治理方案

字节内部尝试过若干不同类型的资源治理方案，包括

1. 资源运营：定期帮助业务跑资源利用情况并推动资源申请治理，问题是运维负担重且无法根治利用率问题
> 
2. 动态超售：在系统侧评估业务资源量并主动缩减配额，问题是超售策略不一定准确且可能导致挤兑风险
> 
3. 动态扩缩：问题是如果只针对在线服务扩缩，由于在线服务的流量波峰波谷类似，无法充分实现全天利用率提升

所以最终字节采用**混合部署**，将在线和离线同时运行在相同节点，充分利用在线和离线资源之间的互补特性，实现更好的资源利用；最终我们期望达到如下图效果，即**二次销售在线未使用的资源**，利用离线工作负载能够很好地填补这部分超售资源，实现资源利用效率在全天保持在较高水平。

# 字节混部发展历程

## 在离线分时混部

第一个阶段主要进行在线和离线的分时混合部署。

对在线：在该阶段我们构建了在线服务弹性平台，用户可以根据业务指标配置横向伸缩规则；例如，凌晨时业务流量减少，业务主动缩减部分实例，系统将在实例缩容基础上进行资源 bing packing 从而腾出整机；

对离线：在该阶段离线服务可获取到大量 spot 类型资源，由于其供应不稳定所以成本上享受一定折扣；同时对于在线来说，将未使用的资源卖给离线，可以在成本上获得一定返利。

该方案优势在于不需要采取复杂的单机侧隔离机制，技术实现难度较低；但同样存在一些问题，例如

1. 转化效率不高，bing packing 过程中会出现碎片等问题；
2. 离线使用体验可能也不好，当在线偶尔发生流量波动时离线可能会被强制杀死，导致资源波动较强烈；
3. 对业务会造成实例变化，**实际操作过程中业务通常会配置比较保守的弹性策略，导致资源提升上限较低**。

## Kubernetes/YARN 联合混部

为解决上述问题我们进入了第二个阶段，尝试将离线和在线真正跑在一台节点上。

由于在线部分早先已经基于 Kubernetes 进行了原生化改造，但大多数离线作业仍然基于 YARN 进行运行。为推进混合部署，我们在单机上引入第三方组件负责确定协调给在线和离线的资源量，并与 Kubelet 或 Node Manager 等单机组件打通；同时当在线和离线工作负载调度到节点上后，也由该协调组件异步更新这两种工作负载的资源分配。

该方案使得我们完成混部能力的储备积累，并验证可行性，但仍然存在一些问题:

1. 两套系统异步执行，使得在离线容器只能旁路管控，存在 race；且中间环节资源损耗过多
2. 对在离线负载的抽象简单，使得我们无法描述复杂 QoS 要求
3. 在离线元数据割裂，使得极致的优化困难，无法实现全局调度优化

![](./在离线统一调度混部.awebp)

# Katalyst 

## 系统概览

Katalyst 系统大致分为四层，从上到下依次包括

- 最上层的标准 API，为用户抽象不同的 QoS 级别，提供丰富的资源表达能力；
- 中心层则负责统一调度、资源推荐以及构建服务画像等基础能力；
- 单机层包括自研的数据监控体系，以及负责资源实时分配和动态调整的资源分配器；
- 最底层是字节定制的内核，通过增强内核的 patch 和底层隔离机制解决在离线跑时单机性能问题。

![](./katalyst架构图.awebp)

## 抽象标准化：QoS Class

Katalyst QoS 可以从宏观和微观两个视角进行解读

- 宏观上，Katalyst 以 CPU 为主维度定义标准了的 QoS 级别；具体来说我们将 QoS 分为四类：独占型、共享型、回收型和为系统关键组件预留的系统型；
- 微观上，Katalyst 最终期望状态无论什么样的 workload，都能实现在相同节点上的并池运行，不需要通过硬切集群来隔离，实现更好的资源流量效率和资源利用效率。

![](./抽象标准化QoS.awebp)

在 QoS 的基础上，Katalyst 同时也提供了丰富的扩展 Enhancement 来表达除 CPU 核心外其他的资源需求

QoS Enhancement：扩展表达业务对于 NUMA /网卡绑定、网卡带宽分配、IO Weight 等多维度的资源诉求
Pod Enhancement：扩展表达业务对于各类系统指标的敏感程度，比如 CPU 调度延迟对业务性能的影响
Node Enhancement：通过扩展原生的 TopologyPolicy 表示多个资源维度间微拓扑的组合诉求

## 管控同步化：QoS Resource Manager

为在 K8s 体系下实现同步管控的能力，我们有三种 hook 方式：CRI 层、OCI 层、Kubelet 层；最终 Katalyst 选择在 Kubelet 侧实现管控，即实现和原生的 Device Manager 同层级的 QoS Resource Manager，该方案的优势包括

- 在 admit 阶段实现拦截，无需在后续步骤靠兜底措施来实现管控
- 与 Kubelet 进行元数据对接，将单机微观拓扑信息通过标准接口报告到节点 CRD，实现与调度器的对接
- 在此框架上，可以灵活实现可插拔的 plugin，满足定制化的管控需求


# 参考

1. 字节跳动云原生成本优化实践：https://juejin.cn/post/7265128702940512293

2. Katalyst 社区完成了 0.4.0 版本发布。除了持续优化 QoS 能力之外，我们还在新版本中提供了可以独立在原生 Kubernetes 上使用的潮汐混部和资源超售能力。https://juejin.cn/post/7327686118887473164

2.1 潮汐混部
在潮汐混部中引入了潮汐节点池的概念，并且将集群中的节点划分为“在线”和“离线”两种类型。潮汐混部主要分为两个部分：

实例数管理：通过 HPA、CronHPA 等各种横向扩缩能力来管理在线业务的实例数，在夜间可以腾出资源给离线业务使用

潮汐节点池管理：Tidal Controller 基于设定好的策略对潮汐节点池中的节点做 binpacking，将腾出的资源折合成整机出让给离线业务

2.2 在线超分

Over-commit Webhook：劫持 kubelet 上报心跳的请求，并对 Allocatable 资源量进行放大
Over-commit Controller：超分配置管理
Katalyst Agent：通过干扰检测和驱逐，保障超分后节点的性能和稳定性；根据指标数据，计算并上报动态的超分比
Katalyst Scheduler：对需要绑核的 Pod 进行准入，避免超分导致实际无法绑核而启动失败
