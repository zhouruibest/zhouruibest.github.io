在MySQL 5.6以及之后的版本中，MySQL提出了一个optimizer trace的功能，这个功能可以让用户方便的查看优化器生成执行计划的整个过程，这个功能的开启与关闭由系统变量optimizer_trace决定

```sql
mysql> SHOW VARIABLES LIKE 'optimizer_trace';
+-----------------+--------------------------+
| Variable_name   | Value                    |
+-----------------+--------------------------+
| optimizer_trace | enabled=off,one_line=off |
+-----------------+--------------------------+
1 row in set (0.02 sec)

mysql> SET optimizer_trace="enabled=on";
```

打开之后，可以到information_schema数据库下的OPTIMIZER_TRACE表中查看完整的优化过程。OPTIMIZER_TRACE表有4个列：

1. QUERY：表示我们的查询语句。
2. TRACE：表示优化过程的JSON格式文本。
3. MISSING_BYTES_BEYOND_MAX_MEM_SIZE：由于优化过程可能会输出很多，如果超过某个限制时，多余的文本将不会被显示，这个字段展示了被忽略的文本字节数。
4. INSUFFICIENT_PRIVILEGES：表示是否没有权限查看优化过程，默认值是0，只有某些特殊情况下才会是1。

