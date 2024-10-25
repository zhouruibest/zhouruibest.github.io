2019年的博客
https://cloud.tencent.com/developer/article/1397084
https://blog.tianfeiyu.com/2019/06/09/node_status/

2023年的博客，跟上面的结论不一样
https://www.cnblogs.com/WJQ2017/p/17066768.html   


# Kubelet进程异常，Pod状态变化

一个节点上运行着pod前提下，这个时候把kubelet进程停掉。里面的pod会被干掉吗？会在其他节点recreate吗？

（1）Node状态变为NotReady
（2）Pod 5分钟之内状态无变化，5分钟之后的状态变化：Daemonset的Pod状态变为Nodelost，Deployment、Statefulset和Static Pod的状态先变为NodeLost，然后马上变为Unknown/Terminating。Deployment的pod会recreate，但是Deployment如果是node selector停掉kubelet的node，则recreate的pod会一直处于Pending的状态。Static Pod和Statefulset的Pod会一直处于Unknown状态。

如果某节点死掉或者与集群中其他节点失联，Kubernetes 会实施一种策略，将失去的节点上运行的所有 Pod 的 **phase** 设置为 Failed。

> 注意：以上，不要混淆状态和phase！！！！

# Kubelet恢复，Pod行为

Kubelet恢复，Pod行为

结论：

（1）Node状态变为Ready。

（2）Daemonset的pod不会recreate，旧pod状态直接变为Running。

（3）Deployment的则是将kubelet进程停止的Node删除（原因可能是因为旧Pod状态在集群中有变化，但是Pod状态在变化时发现集群中Deployment的Pod实例数已经够了，所以对旧Pod做了删除处理）

（4）Statefulset的Pod会重新recreate。

（5）Staic Pod没有重启，但是Pod的运行时间会在kubelet起来的时候置为0。

在kubelet停止后，statefulset的pod会变成nodelost，接着就变成unknown，但是不会重启，然后等kubelet起来后，statefulset的pod才会recreate。

> 还有一个就是Static Pod在kubelet重启以后应该没有重启，但是集群中查询Static Pod的状态时，Static Pod的运行时间变了

# StatefulSet Pod为何在Node异常时没有Recreate

Node down后，StatefulSet Pods並沒有重建，為什麼？

我们在node controller中发现，除了daemonset pods外，都会调用delete pod api删除pod。

但并不是调用了delete pod api就会从apiserver/etcd中删除pod object，仅仅是设置pod 的deletionTimestamp，标记该pod要被删除。真正删除Pod的行为是kubelet，kubelet grace terminate该pod后去真正删除pod object。这个时候statefulset controller 发现某个replica缺失就会去recreate这个pod。

但此时由于kubelet挂了，无法与master通信，导致Pod Object一直无法从etcd中删除。如果能成功删除Pod Object，就可以在其他Node重建Pod。

另外，要注意，statefulset只会针对isFailed Pod，（但现在Pods是Unkown状态）才会去delete Pod。