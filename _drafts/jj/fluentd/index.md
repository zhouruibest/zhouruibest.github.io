# 背景

公司的测试环境使用了fluentd收集Kubernetes中容器的日志。各个业务的日志都打到标准输出中了。日志的格式比较多，有单行的，有多行的，有json的，还有一些Java堆栈的。。。观察发现fluentd出现解析错误，很多日志不能有效推送到ES。

# 目的

解决fluentd错误配置的问题，正确解析并推送日志

# 安装

以root用户身份在ubuntu上安装

```
curl -L https://toolbelt.treasuredata.com/sh/install-ubuntu-bionic-td-agent3.sh | sh
```

安装好之后可以以Daemon方式启动

systemctl start td-agent.service
systemctl status td-agent.service

```
$ systemctl status td-agent.service
● td-agent.service - td-agent: Fluentd based data collector for Treasure Data
     Loaded: loaded (/lib/systemd/system/td-agent.service; disabled; vendor preset: enabled)
     Active: active (running) since Tue 2024-12-17 09:32:39 CST; 10s ago
       Docs: https://docs.treasuredata.com/articles/td-agent
    Process: 622610 ExecStart=/opt/td-agent/embedded/bin/fluentd --log $TD_AGENT_LOG_FILE --daemon /var/run/td-agent/td-agent.pid $TD_AGENT_OPTIONS (code=exited, status=0/SUCCESS)
   Main PID: 622630 (fluentd)
      Tasks: 11 (limit: 19105)
     Memory: 58.7M
     CGroup: /system.slice/td-agent.service
             ├─622630 /opt/td-agent/embedded/bin/ruby /opt/td-agent/embedded/bin/fluentd --log /var/log/td-agent/td-agent.log --daemon /var/run/td-agent/td-agent.pid
             └─622635 /opt/td-agent/embedded/bin/ruby -Eascii-8bit:ascii-8bit /opt/td-agent/embedded/bin/fluentd --log /var/log/td-agent/td-agent.log --daemon /var/run/td-agent/td-agent.pid -->

Dec 17 09:32:38 ubuntu20.04.5-template systemd[1]: Starting td-agent: Fluentd based data collector for Treasure Data...
Dec 17 09:32:39 ubuntu20.04.5-template systemd[1]: Started td-agent: Fluentd based data collector for Treasure Data.
```

安装没有问题。`systemctl stop td-agent.service`停掉服务。使用 /opt/td-agent/embedded/bin/fluentd和当前文件夹下的配置文件试验。

# Fluentd 的配置文件

source 指令确定输入源。
match 指令确定输出目的地。
filter 指令确定事件处理管道。
system 指令设置系统范围的配置。
label 指令对内部路由的输出和过滤器进行分组。
@include 指令包含其他文件。

## source
通过使用 source 指令选择和配置所需的输入插件来启用 Fluentd 输入源，Fluentd 标准输入插件包括 http 和 forward。http 提供了一个 HTTP 端点来接受传入的 HTTP 消息，而 forward 提供了一个 TCP 端点来接受 TCP 数据包。当然，也可以同时是两者。例如：

```
# test.conf
<source>      # 输入源
  @type http  # 打开http端口
  port 9880   #  端口号
</source>

<match *.*>
  @type stdout # 收到消息之后直接输出到stdout
</match>
```

启动服务

opt/td-agent/embedded/bin/fluentd -c test.conf

可以在fluentd的标准输出看到:

```
....
2024-12-17 10:01:14.119904126 +0800 my.tag: {"event":"data"} 

```



输入源可以一次指定多个，@type 参数用来指定输入插件，输入插件扩展了 Fluentd，以检索和提取来自外部的日志事件，一个输入插件通常创建一个线程、套接字和一个监听套接字，它也可以被写成定期从数据源中提取数据。Fluentd 支持非常多种输入插件，包括：

in_tail
in_forward
in_udp
in_tcp
in_unix
in_http
in_syslog
in_exec
in_sample
in_windows_eventlog

tail 插件应该是平时我们使用得最多的输入插件了，in_tail 输入插件允许 Fluentd 从文本文件的尾部读取事件，其行为类似于 tail -F 命令.

# 测试

## 日志样例1：

```log
2024-12-17 09:54:51.358 [INFO ] [TID: N/A] [Thread-19] [com.xxl.job.core.thread.ExecutorRegistryThread] [run] [60] ->>>>>>>>>>>> xxl-job registry success, registryParam:RegistryParam{registGroup='EXECUTOR', registryKey='liquidation-checking-executor', registryValue='172.10.29.100:9235'}, registryResult:ReturnT [code=200, msg=null, content=null]
```

配置文件如下

```
# test.conf
<source>      # 输入源
  @type tail
  path /root/fluentd-test/*.log
  pos_file /root/fluentd-test/test.log.pos
  tag tag123
  <parse>
    @type multiline
    format_firstline /^(?<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) \[(?<level>[^\]]+)\]/
    format1 /^(?<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) \[(?<level>[^\]]+)\] (?<log>.+)/
    time_format %Y-%m-%d %H:%M:%S.%L
  </parse>
</source>

<match tag123>
  @type stdout # 收到消息之后直接输出到stdout
</match>
```

运行情况如下，能正确解析

```log
...

2024-12-17 09:54:51.358000000 +0800 tag123: {"level":"INFO ","log":"[TID: N/A] [Thread-19] [com.xxl.job.core.thread.ExecutorRegistryThread] [run] [60] ->>>>>>>>>>>> xxl-job registry success, registryParam:RegistryParam{registGroup='EXECUTOR', registryKey='liquidation-checking-executor', registryValue='172.10.29.100:9235'}, registryResult:ReturnT [code=200, msg=null, content=null]"}
```

## 日志样例2

```log
{"@timestamp":"2024-12-17T14:01:28.263+08:00","level":"info","content":"创建Merge PR成功, 合并分支：bak-auto/20241217140127-feature/以旧换新v1.0, 目标分支：tmp-auto/20241217140127-40910, prID：95, version：145623","trace":"73d241b66a1a0487f230ebfbc8e8e7d5","span":"938553de06697eb5"}
```

配置文件如下：

```
<source>      # 输入源
  @type tail
  path /root/fluentd-test/*.log
  pos_file /root/fluentd-test/test.log.pos
  tag tag123
  <parse>
    @type json
    time_key @timestamp
    time_format %Y-%m-%dT%H:%M:%S.%N%:z
  </parse>
</source>

<match tag123>
  @type stdout # 收到消息之后直接输出到stdout
</match>
```

解释： %Y-%m-%dT%H:%M:%S.%N%:z

•%Y：四位数的年份（例如：2024）
•%m：两位数的月份（01 到 12）
•%d：两位数的日期（01 到 31）
•%T 或 %H:%M:%S：24小时制的时间，格式为 HH:MM:SS
•%N：纳秒部分（可以是9位数字）
•%:z：带有冒号的时区偏移量（例如：+08:00）

## multi_parser

TODO

在同一个文件中解析多种格式的日志.

需要安装 fluent-plugin-multi-format-parser 这个插件。

# 其他

## 正则表达式匹配日志中的日期

格式字符串 "%Y-%m-%dT%H:%M:%S.%N%:z"

解释：
•%Y：四位数的年份（例如：2024）
•%m：两位数的月份（01 到 12）
•%d：两位数的日期（01 到 31）
•%T 或 %H:%M:%S：24小时制的时间，格式为 HH:MM:SS
•%N：纳秒部分（可以是9位数字）
•%:z：带有冒号的时区偏移量（例如：+08:00）

