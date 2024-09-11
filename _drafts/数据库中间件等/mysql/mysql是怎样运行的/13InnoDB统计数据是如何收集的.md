InnoDB统计数据的收集和存储方式

# 存储位置

innodb_stats_persistent为ON时，统计数据存储到磁盘，否则存储到内存、重启后待合适时机重新收集。

InnoDB以表为单位收集和统计数据。

```SQL
CREATE TABLE 表名 (...) Engine=InnoDB, STATS_PERSISTENT = (1|0)
ALTER TABLE 表名 (...) Engine=InnoDB, STATS_PERSISTENT = (1|0)

-- STATS_PERSISTENT 默认取值为系统变量innodb_stats_persistent
```

# 基于磁盘的存储

当我们选择把某个表以及该表索引的统计数据存放到磁盘上时，实际上是把这些统计数据存储到了两个表里：

```sql
mysql> SHOW TABLES FROM mysql LIKE 'innodb%';
+---------------------------+
| Tables_in_mysql (innodb%) |
+---------------------------+
| innodb_index_stats        |
| innodb_table_stats        |
+---------------------------+
2 rows in set (0.01 sec)
```

- innodb_table_stats 存储了关于表的统计数据，每一条记录对应着一个表的统计数据。
- innodb_index_stats 存储了关于索引的统计数据，每一条记录对应着一个索引的一个统计项的统计数据。

## innodb_table_stats

|字段名|描述|
|-|-|
|database_name|数据库名|
|table_name|表名|
|last_update|本条记录最后更新时间|
|n_rows|表中记录的条数|
|clustered_index_size|表的聚簇索引占用的页面数量|
|sum_of_other_index_sizes|表的其他索引占用的页面数量|

- n_rows 采样均值 * 页面数
- clustered_index_size、sum_of_other_index_sizes 根据页、区段、数据字典等结构估算

## innodb_index_stats

innodb_index_stats表的每条记录代表着一个索引的一个统计项。

|字段名|描述|
|--|--|
|database_name|数据库名|
|table_name|表名|
|index_name|索引名|
|last_update|本条记录最后更新时间|
|stat_name|统计项的名称|
|stat_value|对应的统计项的值|
|sample_size|为生成统计数据而采样的页面数量|
|stat_description|对应的统计项的描述|

## 定期更新统计数据

- 自动重新计算统计数据(innodb_stats_auto_recalc默认值ON)。如果发生变动的记录数量超过了表大小的10%，并且自动重新计算统计数据的功能是打开的，那么服务器会重新进行一次统计数据的计算，并且更新innodb_table_stats和innodb_index_stats表。

```sql
-- 单独为某个表指定STATS_AUTO_RECALC。1表示自动重新计算统计数据，每个表区系统变量innodb_stats_auto_recalc
CREATE TABLE 表名 (...) Engine=InnoDB, STATS_AUTO_RECALC = (1|0);
ALTER TABLE 表名 Engine=InnoDB, STATS_AUTO_RECALC = (1|0);

```

- 手动调用ANALYZE TABLE语句来更新统计信息

ANALYZE TABLE语句会立即重新计算统计数据，也就是这个过程是同步的，在表中索引多或者采样页面特别多时这个过程可能会特别慢.

## 手动更新innodb_table_stats和innodb_index_stats表

```sql
UPDATE innodb_table_stats 
    SET n_rows = 1
    WHERE table_name = 'single_table';
FLUSH TABLE single_table; --更新完innodb_table_stats只是单纯的修改了一个表的数据，需要让MySQL查询优化器重新加载更改过的数据
```

# 基于内存的非永久性统计数据











