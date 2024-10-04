# k8s informer 是如何保证事件不丢失的?

## informer的list-watch机制

client-go中的reflector模块首先会list apiserver获取某个资源的全量信息，然后根据list到的resourceversion来watch资源的增量信息。且希望使用client-go编写的控制器组件在与apiserver发生连接异常时，尽量的re-watch资源而不是re-list

> watch是使用了http/2的Server-Push特性：K8s 为了充分利用 HTTP/2 在 Server-Push、Multiplexing 上的高性能 Stream 特性，在实现 RESTful Watch 时，提供了 HTTP1.1/HTTP2 的协议协商(ALPN, Application-Layer Protocol Negotiation) 机制，在服务端优先选中([参考](HTTP2https://www.cnblogs.com/tencent-cloud-native/p/16206606.html))...

## re-list的场景

场景一：very short watch

reflector与api建立watch连接，但apiserver关闭了连接，则会重新re-list

这意味着 apiserver 接受了监视请求，但立即终止了连接，如果您偶尔看到它，则表明存在暂时性错误，并不值得警惕。如果您反复看到它，则意味着 apiserver（或 etcd）有问题。

如下只返回了0个item，因此做了一下re-list

```sh
I0728 11:32:06.170821 67483 streamwatcher.go:114] Unexpected EOF during watch stream event decoding: unexpected EOF I0728 11:32:06.171062 67483 reflector.go:391] k8s.io/client-go/informers/factory.go:134: Watch close - *v1.Deployment total 0 items received W0728 11:32:06.187394 67483 reflector.go:302] k8s.io/client-go/informers/factory.go:134: watch of *v1.Deployment ended with: very short watch: k8s.io/client-go/informers/factory.go:134: Unexpected watch close - watch lasted less than a second and no items received
```

场景二：401 Gone

为什么跟etcd不会一直记录历史版本有关 [参考：bookmark机制](https://blog.csdn.net/qq_43684922/article/details/131869680?spm=1001.2014.3001.5501)

reflector与api建立watch连接，但是出现watch的相关事件丢失时（etcd不会一直记录历史版本），api返回401 Gone，reflector提示too old resource version并重新re-list

```sh
I0728 14:40:58.807670 71423 reflector.go:300] k8s.io/client-go/informers/factory.go:134: watch of *v1.Deployment ended with: too old resource version: 332167941 (332223202) I0728 14:40:59.808153 71423 reflector.go:159] Listing and watching *v1.Deployment from k8s.io/client-go/informers/factory.go:134 I0728 14:41:00.300695 71423 reflector.go:312] reflector list resourceVersion: 332226582
```

## resync场景

1. resync不是re-list，resync不需要访问apiserver

2. resync 是重放 informer 中的 obj 到 DeltaFIFO 队列中，触发 handler 再次处理 obj。目的是防止有些 handler 处理失败了而缺乏重试的机会。特别是，**需要修改外部系统的状态的时候，需要做一些补偿的时候**。

比如说，根据 networkpolicy刷新 node 上的 iptables。iptables 有可能会被其他进程或者管理员意外修改，有 resync 的话，才有机会定期修正。这也说明，**回调函数的实现需要保证幂等性。对于 OnUpdate 函数而言，有可能会拿到完全一样的两个 Obj，实现 OnUpdate 时要考虑到**。

3. re-list 是指 reflector 重新调用 kube-apiserver 全量同步所有 obj。list 的时机一般是在程序第一次启动，或者 watch 有错误，才会 re-list。

5. resync 是一个水平触发的模式

水平触发是只要处于某个状态，就会一直通知。比如在这里，对象已经在缓存里，会触发不止一次回调函数。

5. process 函数怎么区分从 DeltaFIFO 里拿到的 obj 是新的还是重放的呢？

根据 obj 的 key（namespace/name）从 index 里拿到旧的 obj，和新出队的 obj 比较 resource revision，这两个 resource revision 如果一样，就是重放的，如果不一样，就是从 kube-apiserver 拿到的新的。因为 resource revision 只有在 etcd 才能更新。 index 作为客户端缓存，这个值是不变的。

## 6、resync要注意的问题

1. 如何配置resync的周期？

func NewSharedIndexInformer(lw ListerWatcher, exampleObject runtime.Object, defaultEventHandlerResyncPeriod time.Duration, indexers Indexers)

第三个参数 defaultEventHandlerResyncPeriod 就指定多久 resync 一次。如果为0，则不 resync。

AddEventHandlerWithResyncPeriod也可以给单独的 handler 定义 resync period，否则默认和 informer 的是一样的。

2. 配置resync周期间隔太小会有什么问题

此时会以比较高的频率促使事件重新入队进行reconcile，造成controller的压力过大

3.resync用于解决什么问题，resync 多久一次比较合适？或者需不需要 resync？

根据具体业务场景来，根据外部状态是不是稳定的、是否需要做这个补偿来决定的，

举例：假设controller是一个LB controller

（1）当watch到了service创建，然后调用云平台接口去创建一个对应的LB （k8s系统外部的）
（2）然后如果此时这个对应的LB由于某种bug被删除了，此时service就不通了，那么此时状态不一致了，集群里有这个service，云平台那边没有对应的LB，并且由于云平台的bug被删除了，而不是删除service而触发LB删除的，此时service是没有变化的，也就不会出发reconcile了。
（3）假设我们reconcile里有逻辑是判断如果service没有对应的LB就创建，那么此时reconcile不会被出发，那也就没有被执行了。此时如果有resync，定时将indexer里的对象，也就是缓存的对象来一次update事件的入队，进行后续出队触发reconcile，那我们就会发现service对应的LB没了，进而进行创建。
（4）也就是说，resync是防止业务层的bug。
（5）且resync将indexer的对象重入队，里面的service不是所有service，而是创建了LB的service，

