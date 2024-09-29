当用户使用 kubectl，client-go 或者 REST API 请求 apiserver 时，都要经过认证、授权、准入控制的校验。

1. 认证解决的问题是识别用户的身份；
2. 授权是明确用户具有哪些权限；
3. 准入控制是作用于 kubernetes 中的资源对象。

# Kubernetes API Server 认证机制（Authentication）

一旦TLS连接建立，请求就进入到身份认证阶段，在这一阶段，请求由一个或多个认证器模块检查。

认证模块是管理员在集群创建过程中配置的，一个集群可能有多个认证模块配置，每个模块会依次尝试认证， 直到其中一个认证成功。

在主流的认证模块中会包括客户端证书、密码、plain tokens、bootstrap tokens以及JWT tokens（用于service account）。客户端证书的使用是默认的并且是最常见的方案。

1. X509 client certs
2. Static Token File
3. Bootstrap Tokens
4. Static Password File
5. Service Account Tokens
6. OpenId Connect Tokens
7. Webhook Token Authentication
8. Authticating Proxy
9. Anonymous requests
10. User impersonation
11. Client-go credential plugins

# Kubernetes 常用认证机制

## X509 client certs

X509 client certs 认证方式是用在一些客户端访问 apiserver 以及集群组件之间访问时使用，比如 kubectl 请求 apiserver 时。

适用对象：外部用户

X509是一种数字证书的格式标准，现在 HTTPS 依赖的 SSL 证书使用的就是使用的 X509 格式。X509 客户端证书认证方式是 kubernetes 所有认证中使用最多的一种，相对来说也是最安全的一种，kubernetes 的一些部署工具 kubeadm、minkube 等都是基于证书的认证方式。客户端证书认证叫作 TLS 双向认证，也就是服务器客户端互相验证证书的正确性，在都正确的情况下协调通信加密方案。目前最常用的 X509 证书制作工具有 openssl、cfssl 等

## Service Account Tokens

serviceaccounts 是用在 pod 中访问 apiserver 时进行认证的，比如使用自定义 controller 时。

适用对象：内部用户

有些情况下，我们希望在 pod 内部访问 apiserver，获取集群的信息，甚至对集群进行改动。针对这种情况，kubernetes 提供了一种特殊的认证方式：serviceaccounts。

serviceaccounts 是面向 namespace 的，每个 namespace 创建的时候，kubernetes 会自动在这个 namespace 下面创建一个默认的 serviceaccounts；并且这个 serviceaccounts 只能访问该 namespace 的资源。

serviceaccounts 和 pod、service、deployment 一样是 kubernetes 集群中的一种资源，用户也可以创建自己的 serviceaccounts。

serviceaccounts 主要包含了三个内容：namespace、token 和 ca，每个 serviceaccounts 中都对应一个 secrets，namespace、token 和 ca 信息都是保存在 secrets 中且都通过 base64 编码的。namespace 指定了 pod 所在的 namespace，ca 用于验证 apiserver 的证书，token 用作身份验证，它们都通过 mount 的方式保存在 pod 的文件系统中，其三者都是保存在 /var/run/secrets/kubernetes.io/serviceaccount/目录下。

## Kubernetes API Server 授权机制（Authorization）

请求经过认证之后，下一步就是确认这一操作是否被允许执行，即授权。

对于授权一个请求，Kubernetes主要关注三个方面：

- 请求者的用户名
- 请求动作
- 动作影响的对象

用户名从嵌入 token 的头部中提取，动作是映射到CRUD操作的HTTP动词之一（如 GET、POST、PUT、DELETE），对象是其中一个有效的 Kubernetes 资源对象。

Kubernetes基于一个白名单策略授权，默认没有访问权限。

kubernetes 目前支持如下四种授权机制：

- Node
- ABAC
- RBAC
- Webhook

# kubernetes 常用授权机制

## RBAC（基于角色的访问控制）

RBAC 中有三个比较重要的概念：

Role：角色，它其实是一组规则，定义了一组对 Kubernetes API 对象的操作权限；
Subject：被作用者，包括 user，group，serviceaccounts，通俗来讲就是认证机制中所识别的用户；
RoleBinding：定义了“Role”和“Subject”的绑定关系，也就是将用户以及操作权限进行绑定；
