据了解，在Kubernetes 1.22版本开始，默认开启了pod-deletion-cost特性，允许用户设置Pod的删除成本，它是一个整数值，可以为正数、零或负数，分值越低在缩容时的优先级越高。

我们在网上翻过很多资料，大多都是通过手动修改，亦或编写脚本定时进行批量修改，都不是很理想。

Openkruise这个组件，可以通过自定义探针PodProbeMarker自动给Pods注入pod-deletion-cost的分值，将CPU使用率较低的删除成本设置为5，将CPU使用率较高的设置为10

# Pod Readiness Gate
https://www.cnblogs.com/jiangbo4444/p/14589407.html
https://www.alibabacloud.com/help/zh/ack/ack-managed-and-ack-dedicated/user-guide/use-readinessgate-to-check-pod-readiness-before-associating-pods-with-alb-ingresses-during-rolling-updates