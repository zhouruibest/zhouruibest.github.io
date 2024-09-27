1. 配置文件默认是没有的，通过命令生成

 containerd config default > /etc/containerd/config.toml

 2. 修改配置文件，在镜像库里面做一个映射

 /etc/containerd/certs.d/<镜像库, e.g. gcr.io>/hosts.toml

 ![](./镜像仓库的映射.png)

 3. 查看配置文件的位置

 ![](./containerd中关于镜像仓库的配置.png)

 4. 添加新的映射

![](./添加新的映射.png)