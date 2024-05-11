# OCI image-spec
OCI 规范中的镜像规范 image-spec 决定了我们的镜像按照什么标准来构建，以及构建完镜像之后如何存放，接着下文提到的 Dockerfile 则决定了镜像的 layer 内容以及镜像的一些元数据信息。一个镜像规范 image-spec 和一个 Dockerfile 就指导着我们构建一个镜像.根据官方文档的描述，OCI 镜像规范的主要由以下几个 markdown 文件组成：
```sh
├── annotations.md         # 注解规范
├── config.md              # image config 文件规范
├── considerations.md      # 注意事项
├── conversion.md          # 转换为 OCI 运行时
├── descriptor.md          # OCI Content Descriptors 内容描述
├── image-index.md         # manifest list 文件
├── image-layout.md        # 镜像的布局
├── implementations.md     # 使用 OCI 规范的项目
├── layer.md               # 镜像层 layer 规范
├── manifest.md            # manifest 规范
├── media-types.md         # 文件类型
├── README.md              # README 文档
├── spec.md                # OCI 镜像规范的概览
```

总结以上几个 markdown 文件， OCI 容器镜像规范主要包括以下几块内容：

## layer
文件系统：以 layer （镜像层）保存的文件系统，每个 layer 保存了和上层之间变化的部分，layer 应该保存哪些文件，怎么表示增加、修改和删除的文件等。

## image config
image config 文件：保存了文件系统的层级信息（每个层级的 hash 值，以及历史信息），以及容器运行时需要的一些信息（比如环境变量、工作目录、命令参数、mount 列表），指定了镜像在某个特定平台和系统的配置，比较接近我们使用 docker inspect  看到的内容

## manifest
manifest 文件 ：**镜像的 config 文件索引，有哪些 layer，额外的 annotation 信息**，manifest 文件中保存了很多和当前平台有关的信息。另外 manifest 中的 layer 和 config 中的 layer 表达的虽然都是镜像的 layer ，但二者代表的意义不太一样，稍后会讲到。**manifest 文件是存放在 registry 中**，当我们拉取镜像的时候，会根据该文件拉取相应的 layer 。

manifest 也分好几个版本，目前主流的版本是 Manifest Version 2, Schema 2 
 
```

image manifest
```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
  "config": {
    "mediaType": "application/vnd.docker.container.image.v1+json",
    "size": 1509,
    "digest": "sha256:a24bb4013296f61e89ba57005a7b3e52274d8edd3ae2077d04395f806b63d83e"
  },
  "layers": [
    {
      "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
      "size": 5844992,
      "digest": "sha256:50644c29ef5a27c9a40c393a73ece2479de78325cae7d762ef3cdc19bf42dd0a"
    }
  ]
}
```

### 解释

Manifest 是一个 JSON 文件，其定义包括两个部分，分别是 Config 和 Layers。Config 是一个 JSON 对象，Layers 是一个由 JSON 对象组成的数组。可以看到，Config 与 Layers 中的每一个对象的结构相同，都包括三个字段，分别是 digest、mediaType 和 size。其中 digest 可以理解为是这一对象的 ID。mediaType 表明了这一内容的类型。size 是这一内容的大小。

容器镜像的 Config 有着固定的 mediaType application/vnd.oci.image.config.v1+json。一个 Config 的示例配置如下，它记录了关于容器镜像的配置，可以理解为是镜像的元数据。通常它会被镜像仓库用来在 UI 中展示信息，以及区分不同操作系统的构建等。

而容器镜像的 Layers 是由多层 mediaType 为 application/vnd.oci.image.layer.v1.*（其中最常见的是 application/vnd.oci.image.layer.v1.tar+gzip) 的内容组成的。众所周知，容器镜像是分层构建的，每一层就对应着 Layers 中的一个对象。

容器镜像的 Config，和 Layers 中的每一层，都是以 Blob 的方式存储在镜像仓库中的，它们的 digest 作为 Key 存在。因此，在请求到镜像的 Manifest 后，Docker 会利用 digest 并行下载所有的 Blobs，其中就包括 Config 和所有的 Layers。

# image manifest index

(存在于registry)

在 docker 的 distribution 中称之为 Manifest List 在 OCI 中就叫 OCI Image Index Specification 。其实两者是指的同一个文件，甚至两者 GitHub 上文档给的 example 都一一模样，应该是 OCI 复制粘贴 Docker 的文档。index 文件是个可选的文件，包含着一个列表为**同一个镜像不同的处理器 arch 指向不同平台的 manifest 文件，这个文件能保证一个镜像可以跨平台使用**，每个处理器 arch 平台拥有不同的 manifest 文件，使用 index 作为索引。当我们使用 arm 架构的处理器时要额外注意，在拉取镜像的时候要拉取 arm 架构的镜像，一般处理器的架构都接在镜像的 tag 后面，默认 latest tag 的镜像是 x86 的，在 arm 处理器的机器这些镜像上是跑不起来的。


example manifest list

```json 
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
  "manifests": [
    {
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "size": 7143,
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
      "platform": {
        "architecture": "ppc64le",
        "os": "linux",
      }
    },
    {
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "size": 7682,
      "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
      "platform": {
        "architecture": "amd64",
        "os": "linux",
        "features": [
          "sse4"
        ]
      }
    }
  ]
}



