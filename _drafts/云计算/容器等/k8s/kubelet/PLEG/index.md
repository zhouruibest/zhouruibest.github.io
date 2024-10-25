# PLEG

PLEG轮询检测运行容器的状态

kubelet的两项重要的配置
```sh
–node-monitor-grace-period=40s（node驱逐时间）
–node-monitor-period=5s（轮询间隔时间）
```

上面两项参数表示每隔 5 秒 kubelet 去检测 Pod 的健康状态，如果在 40 秒后依然没有检测到 Pod 的健康状态便将其置为 NotReady 状态，5 分钟后就将节点下所有的 Pod 进行驱逐。