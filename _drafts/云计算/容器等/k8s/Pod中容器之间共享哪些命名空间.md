ipc, net and uts namespace are shared.

pid: By default its not shared. However the namespace sharing can be enabled in the pod [spec](https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace/).

mnt: Not shared. Different containers are built from different images. From each container only its own root directory will be visible.

user: This is not supported by k8s and some form of user namespace sharing could be implemented in the future as mentioned [here](https://kinvolk.io/blog/2020/12/improving-kubernetes-and-container-security-with-user-namespaces/#bringing-user-namespaces-to-kubernetes).