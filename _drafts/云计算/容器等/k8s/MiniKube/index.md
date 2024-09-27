# MiniKube结构

1. 要有一台物理机
2. 要有容器环境比如PodMan、Docker（或者虚拟机环境）
3. 起一个容器，作为MiniKube的单节点
4. 在容器里面， MiniKube会SetUp一个容器运行时，例如Containerd。当这个单节点跑起来的时候，MiniKube会调用这个容器运行时会拉镜像、跑容器。因此，是容器套容器的架构

# Podman Driver（可以平替Docker Destop）

但是比较坑，暂时不要用（2024/9/26）

# Docker Driver (推荐)

1. 安装Docker Engine

https://docs.docker.com/engine/install/ubuntu/

2. 下载MiniKube

https://minikube.sigs.k8s.io/docs/start/?arch=%2Flinux%2Fx86-64%2Fstable%2Fbinary+download

minikube start --driver=docker --container-runtime=containerd 

3. 启动MiniKube

4. 配通gcr.io