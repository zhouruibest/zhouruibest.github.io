
# CFS

cfs表示Completely Fair Scheduler完全公平调度器，是Linux内核的一部分，负责进程调度。

- cpu.cfs_period_us: 用来设置一个CFS调度时间周期长度，默认值是100000us(100ms)，一般cpu.cfs_period_us作为系统默认值我们不会去修改它。
- cpu.cfs_quota_us: 用来设置在一个CFS调度时间周期(cfs_period_us)内，允许此控制组执行的时间。默认值为-1表示不限制时间。
- cpu.shares: 用来设置cpu cgroup子系统对于控制组之间的cpu分配比例。默认值是1000。

Limit

使用cfs_quota_us/cfs_period_us，例如20000us/100000us=0.2，表示允许这个控制组使用的CPU最大是0.2个CPU，即限制使用20%CPU。 如果cfs_quota_us/cfs_period_us=2，就表示允许控制组使用的CPU资源配置是2个。

Request

对于cpu分配比例的使用，例如有两个cpu控制组foo和bar，foo的cpu.shares是1024，bar的cpu.shares是3072，它们的比例就是1:3。

在一台8个CPU的主机上，如果foo和bar设置的都是需要4个CPU的话(cfs_quota_us/cfs_period_us=4)，根据它们的CPU分配比例，控制组foo得到的是2个，而bar控制组得到的是6个。

*需要注意cpu.shares是在多个cpu控制组之间的分配比例，且只有到整个主机的所有CPU都打满时才会起作用*

例如刚才这个例子，如果是在一个4CPU的机器上，foo和bar的比例是1:3，如果foo控制组的程序满负载跑，而bar控制组程序没有运行或空负载，此时foo仍然能获得4个CPU，但如果foo和bar都满负载跑，因为它们都设置的需要4个CPU，而主机的CPU不够，只能按比例给它们分配1:3，则foo分配得到1个，而bar得到3个。

也就是说cpu.cfs_quota_us/cpu.cfs_period_us决定cpu控制组中所有进程所能使用CPU资源的最大值，而cpu.shares决定了cpu控制组间可用CPU的相对比例，这个比例只有当主机上的CPU完全被打满时才会起作用。
