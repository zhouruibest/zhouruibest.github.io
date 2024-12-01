# ubuntu 安装 docker

1. 检查卸载老版本docker
ubuntu下自带了docker的库，不需要添加新的源。
但是ubuntu自带的docker版本太低，需要先卸载旧的再安装新的。

注：docker的旧版本不一定被称为docker，docker.io 或 docker-engine也有可能，所以我们卸载的命令为

apt-get remove docker docker-engine docker.io containerd runc

2. 更新软件包
sudo apt update
sudo apt upgrade

3. 安装docker依赖

apt-get install ca-certificates curl gnupg lsb-release

4. 添加Docker官方GPG密钥
curl -fsSL http://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -

5. 添加Docker软件源
sudo add-apt-repository "deb [arch=amd64] http://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"

6. 安装docker

apt-get install docker-ce docker-ce-cli containerd.io

7. 运行docker

systemctl start docker

8. 安装工具

apt-get -y install apt-transport-https ca-certificates curl software-properties-common

# 安装 Docker-Compose
curl -L https://github.com/docker/compose/releases/download/1.26.2/docker-compose-$(uname -s)-$(uname -m) > /usr/bin/docker-compose
chmod +x /usr/bin/docker-compose
docker-compose --version

# 构建镜像

1. 配置代理
/etc/docker/daemon.json, 国内的镜像经常被封掉或者自己停了, 需要时长翻新
```json
{
    "registry-mirrors": [
    	"https://docker.unsee.tech",
        "https://dockerpull.org",
        "https://dockerhub.icu"
    ]
}
```

改完之后，systemctl daemon-reload && systemctl restart docker

2. 拉取镜像
docker pull nginx:latest
docker pull tsund/php:7.2.3-fpm # 这是其中一个博主的
docker pull mysql:5.7

其中 nginx 为官方最新镜像，mysql:5.7 为官方 5.7 镜像，tsund/php:7.2.3-fpm 的 Dockerfile 如下：

```dockerfile
FROM php:7.2.3-fpm
LABEL maintainer="tsund" \
      email="tsund@qq.com" \
      version="7.2.3"

RUN apt-get update \
    && docker-php-ext-install pdo_mysql \
    && echo "output_buffering = 4096" > /usr/local/etc/php/conf.d/php.ini

```

在官方镜像的基础上，添加了 PDO_MYSQL（如果使用 MySQL 作为 Typecho 的数据库，则需安装此扩展），并设置 buffer 为 4kb，即一个内存页

# 本地配置

新建 blog 文件夹，其目录结构如下：

.
├── docker-compose.yml      Docker Compose 配置文件
├── mysql                   mysql 持久化目录
├── mysql.env               mysql 配置信息
├── nginx                   nginx 配置文件的持久化目录
├── ssl                     ssl 证书目录
└── typecho                 站点根目录

1. 配置 docker-compose.yml

```yml
version: "3"

services:
  nginx:
    image: nginx
    ports:
      - "80:80"
      - "443:443"
    restart: always
    volumes:
      - ./typecho:/var/www/html
      - ./ssl:/var/www/ssl
      - ./nginx:/etc/nginx/conf.d
    depends_on:
      - php
    networks:
      - web

  php:
    image: tsund/php:7.2.3-fpm
    restart: always
    ports:
      - "9000:9000"
    volumes:
      - ./typecho:/var/www/html
    environment:
      - TZ=Asia/Shanghai
    depends_on:
      - mysql
    networks:
      - web

  mysql:
    image: mysql:5.7
    restart: always
    ports:
      - "3306:3306"
    volumes:
      - ./mysql/data:/var/lib/mysql
      - ./mysql/logs:/var/log/mysql
      - ./mysql/conf:/etc/mysql/conf.d
    env_file:
      - mysql.env
    networks:
      - web

networks:
  web:
```

2. 配置nginx
在 ./nginx 目录下新建 default.conf 文件，参考内容如下：

shengheblog.com  买的便宜域名

``` conf
server {
    listen       80;
    server_name  shengheblog.com;
    rewrite ^(.*) https://shengheblog.com$1 permanent;
}

server {
    listen 443 ssl http2 reuseport;
    server_name shengheblog.com;

    root /var/www/html;
    index index.php;

    access_log /var/log/nginx/typecho_access.log main;

    ssl_certificate /var/www/ssl/tsund_cn.crt;
    ssl_certificate_key /var/www/ssl/tsund_cn.key;
    ssl_session_cache shared:SSL:1m;
    ssl_session_timeout 5m;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
    ssl_prefer_server_ciphers on;

    if (!-e $request_filename) {
        rewrite ^(.*)$ /index.php$1 last;
    }

    location ~ .*\.php(\/.*)*$ {
        fastcgi_pass   php:9200;
        fastcgi_index  index.php;
        fastcgi_split_path_info ^(.+?.php)(/.*)$;
        fastcgi_param  PATH_INFO $fastcgi_path_info;
        fastcgi_param  PATH_TRANSLATED $document_root$fastcgi_path_info;
        fastcgi_param  SCRIPT_NAME $fastcgi_script_name;
        fastcgi_param  SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include        fastcgi_params;
    }
}
```

3. 配置 mysql

mysql.env 参考内容如下：

```commandline
# MySQL的root用户默认密码，这里自行更改
MYSQL_ROOT_PASSWORD=root1234
# MySQL镜像创建时自动创建的数据库名称
MYSQL_DATABASE=blog
# MySQL镜像创建时自动创建的用户名
MYSQL_USER=zhourui
# MySQL镜像创建时自动创建的用户密码
MYSQL_PASSWORD=zhourui1234
# 时区
TZ=Asia/Shanghai
```

# 安装

1. 编排容器

在 blog 目录下, docker-compose up -d && docker-compose ps


2. 安装 typecho

进入typecho目录

wget https://github.com/typecho/typecho/releases/download/v1.2.0/typecho.zip

apt install zip -y

unzip typecho.zip

# 配置 typecho



