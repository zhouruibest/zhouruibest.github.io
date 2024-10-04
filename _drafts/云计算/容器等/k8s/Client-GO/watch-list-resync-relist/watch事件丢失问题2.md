# 现象
v1.27 的 K8s，在 kube-apiserver 的日志中会看到 “etcd event received with PrevKv=nil” 的字样，资源对象被删除后在 Etcd 中已经不存在了但在 Reflector store 中仍然存在，可以在 Informer 或者 watchCache 中看到对应的对象，依赖 Informer 的组件也不会感知到资源对象被删除，通过 List API 设置 RV=“0” 去 kube-apiserver 的 watchCache 中获取的话也可以看到已经被删除的对象仍然存在。

# 回顾
出现 PrevKV=nil 的话，肯定就是 Etcd 返回的数据有问题。前一篇中讲到

Etcd compaction 导致出现 PrevKV=nil 的 delete event，而这个问题已经在 Etcd 中修复，最终返回 ErrCompacted；
K8s 侧做了兜底，任何 PrevKV=nil 的非 Create 事件都会导致 Reflector 收到 InternalServerError 报错，进而触发 ListAndWatch() 的重新执行，规避 Etcd 返回异常数据带来的影响；
既然都已经在 Etcd 与 K8s 侧都进行了功能完善，那么理论上就不会再出现 PrevKV=nil 导致的丢事件的问题了，但为何又出现了呢？

# 新问题
新的问题是在名为 “APIServer watchcache lost events” 的 issue 中提出的，使用了 v1.27 的 K8s 版本。目前基本已经知道问题原因，但尚未完全修复，也就是说如果你在使用 v1.27 或者更新的 K8s 版本，就有可能会遇到这个问题，但也并不是说使用低版本的 K8s 不会遇到这个问题，只是概率不同而已，另外也和具体的使用方式有关，下面会介绍。

这次的问题主要是在 Etcd，当然 K8s 侧也有一点小问题。本篇重点介绍 K8s 侧的相关逻辑，后续再详细介绍 Etcd 的相关逻辑问题。截止目前，这些问题都已经被定位到，有相关 PR，但尚未全部合入 master。

# 直接原因
PR cacher: Fix watch behaviour for unset RV 随着 v1.27 发布，触发了 Etcd 的问题，导致最终问题的出现。
...

原文：https://www.likakuli.com/posts/k8seventlost2/