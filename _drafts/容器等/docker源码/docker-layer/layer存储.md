Docker镜像由多层只读镜像层组成，Docker容器实际上就是创建了一层读写层，然后将这层读写层与使用的镜像关联起来。所以实际上镜像和容器在Docker中都是通过“层”来进行存储，每一层对应一个目录，Docker层存储就是用来管理这些“层”以及这些层的元数据

Docker层存储主要包括存储驱动、层的元数据存储以及层的元数据映射表3个部分。在Docker源码中，关系如下如：

![](./层存储主要结构示意图.png)

层的内容存储于/var/lib/docker/driver_name目录中，由存储驱动管理

层的元素据存储于/var/lib/docker/image/driver_name/layerdb目录中，由层的元数据存储管理

# 主要结构

Docker层存储在源码中使用一个layerStore结构来表示。store字段表示层的元数据存储，driver字段表示具体的存储驱动

存储驱动用一个Driver结构来表示，它实现了graphdriver.Driver接口。Docker支持overlay2，overlay，aufs等多种驱动，我们这里使用overlay2作为存储驱动，因此这个Driver就是overlay2定义的Driver结构。存储驱动管理层的具体内容

层的元数据存储则使用了名为fileMetadataStore的结构来表示，它实现了MetadataStore接口。层的元数据存储管理层的元数据

容器层的元数据通过一个mountedLayer结构来记录，镜像层的元数据则是通过一个roLayer结构来记录。

> 目录说明
> - Docker目录：/var/lib/docker
> - 存储驱动目录：/var/lib/docker/overlay2
> - 层的元数据存储目录：/var/lib/docker/image/overlay2/layerdb

```go
func NewDaemon(ctx context.Context, config *config.Config, pluginStore *plugin.Store) (daemon *Daemon, err error) {
	setDefaultMtu(config) {

    // 对于linux， operatingSystem=linux
    for operatingSystem, gd := range d.graphDrivers {
		layerStores[operatingSystem], err = layer.NewStoreFromOptions(layer.StoreOptions{
			Root:                      config.Root, // Linux下一般是 /var/lib/docker
			MetadataStorePathTemplate: filepath.Join(config.Root, "image", "%s", "layerdb"),
			GraphDriver:               gd,
			GraphDriverOptions:        config.GraphOptions,
			IDMapping:                 idMapping,
			PluginGetter:              d.PluginStore,
			ExperimentalEnabled:       config.Experimental,
			OS:                        operatingSystem,
		})
		if err != nil {
			return nil, err
		}
    }
```

```go
type StoreOptions struct {
	Root                      string //存储路径：/var/lib/docker
	MetadataStorePathTemplate string //元数据存储路径模板：/var/lib/docker/image/%s/layerdb
	GraphDriver               string //存储驱动：overlay2、overlay、aufs、devicemmapper...
	GraphDriverOptions        []string
	IDMapping                 *idtools.IdentityMapping
	PluginGetter              plugingetter.PluginGetter
	ExperimentalEnabled       bool
	OS                        string
}
```