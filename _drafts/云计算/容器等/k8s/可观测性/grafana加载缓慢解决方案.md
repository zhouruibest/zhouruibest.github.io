目前随着数据和图表的逐渐增多，[Grafana](https://so.csdn.net/so/search?q=Grafana&spm=1001.2101.3001.7020) 页面加载速度明显变慢，严重影响了用户体验。

对于速度优化，我们可以从以下方面进行入手：

1. 优化前端资源加载

**使用反向代理和压缩**：通过 NGINX 或其他[反向代理](https://so.csdn.net/so/search?q=%E5%8F%8D%E5%90%91%E4%BB%A3%E7%90%86&spm=1001.2101.3001.7020)服务器启用 `gzip` 压缩和缓存静态资源，减少页面加载的时间。

**静态资源缓存**：确保静态资源如 CSS、JS 等在浏览器中缓存，以避免每次加载都重新获取这些资源。

2. 优化grafana **server**端服务器资源

**增加服务器性能**：检查服务器的 CPU、内存和 I/O 是否有瓶颈，适当增加服务器资源配置。

**调整 Grafana 服务器配置**：增加 Grafana 的 `concurrent_requests_limit` 设置，允许更多的并发请求。

3. 数据库优化

如果你使用的是 `Grafana` 自己的 `sqlite` 或 `MySQL`/`PostgreSQL`，请确保这些数据库被适当优化，数据库性能问题也可能导致慢加载。

4. 分离数据源

**多实例 Prometheus**：如果你使用的是 Prometheus 数据源，可以考虑使用多个 Prometheus 实例来分担负载，特别是如果你的查询数据量很大。

5. 优化数据源查询

减少查询时间范围：设置默认时间范围为较短的时间段（如过去 5 分钟或 15 分钟），以减少加载时的数据量。

使用高效的数据源：检查你的数据源（如 Prometheus 或 Elasticsearch）是否有性能瓶颈。数据源的响应速度慢会直接影响 Grafana 的加载速度。

优化查询：确保你的查询尽可能高效，避免不必要的复杂计算或过滤条件。使用 rate() 或 avg_over_time() 等函数优化大数据量的查询。

# 前端访问优化

缓存哪部分内容

grafana静态文件的目录是在/usr/share/grafana/public

我们需要将这部分内容放到缓存中或者CDN或者OSS中

**由于CDN和OSS都需要花钱，我们暂时就用nginx来做一个类似缓存的功能**

> 注意事项：这里我建议直接在正在运行的grafana容器中把这个public包下载下来（想办法下载，比如压缩后下载，docker cp也行），因为grafana运行了很长时间，里面或多或少会新增一些js文件，而这些文件，是在官方纯净版public是没有的。官方纯净版：Download Grafana | Grafana Labs 选择对应的版本，下载window版本，解压后里面有public文件

#### nginx操作

部署一个nginx，将public静态文件夹放到nginx下面

default.conf如下：

```conf
server {
    listen       80;

    listen  [::]:80;
    server_name  xx.xx.xx;

    location /grafana-oss/9.5.7/public/ {
        alias   /tmp/public/;

        # 跨域配置
        add_header 'Access-Control-Allow-Origin' '*';
        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
        add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization';

        autoindex on;
        autoindex_exact_size on;
        autoindex_localtime on;
    }
    error_page   500 502 503 504  /50x.html;

    location = /50x.html {
        root   /usr/share/nginx/html;
    }
  }


```

内容解释

location 必须是/grafana-oss/你的grafana版本/public/，因为grafana后台请求时，这些是它源码里面携带的地址，你自己改不了

alias /tmp/public/ 这个路径是 下载下来的public 放在 nginx容器里面的某个位置（我这里是/tmp/public/），可以自定义，当然你也可以是用root，两者的区别：nginx 中location中root和alias的区别


接下来验证一下：

访问地址：域名/grafana-oss/9.5.7/public
