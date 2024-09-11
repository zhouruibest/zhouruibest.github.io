5.1版本开始引入show profile剖析单条语句功能，支持show profiles和show profile语句，参数have_profiling;控制是否开启：

查看是否支持这个功能（查询为yes表示支持）：

```
mysql > show variables like 'have_profiling';
+----------------+-------+
| Variable_name  | Value |
+----------------+-------+
| have_profiling | YES   |
+----------------+-------+
1 row in set (0.00 sec)
```

需要临时使用时直接sql命令行中输入：set profiling=1;来开启

```
mysql> set profiling=1;

Query OK, 0 rows affected, 1 warning (0.00 sec)
```
 

然后在服务器上执行SQL语句，都会被测量其消耗的时间和其他一些查询执行状态变更相关的数据

```
mysql> select count(*) from xx;

+----------+

| count(*) |

+----------+

|   262144 |

+----------+

1 row in set (0.05 sec)
```
 

然后再执行：show prifiles;命令，所有的查询SQL都会被列出来

```
mysql> show profiles;

+----------+------------+-------------------------+

| Query_ID | Duration   | Query                   |

+----------+------------+-------------------------+

|        1 | 0.05645950 | select count(*) from xx |

+----------+------------+-------------------------+

1 row in set, 1 warning (0.00 sec)
```
 

然后根据编号查询具体SQL的执行过程，这里演示只执行了一句，那就选项query id为1

```
mysql> show profile for query 1;

+----------------------+----------+

| Status               | Duration |

+----------------------+----------+

| starting             | 0.000041 |

| checking permissions | 0.000004 |

| Opening tables       | 0.000017 |

| init                 | 0.000010 |

| System lock          | 0.000006 |

| optimizing           | 0.000004 |

| statistics           | 0.000009 |

| preparing            | 0.000008 |

| executing            | 0.000001 |

| Sending data         | 0.056110 |

| end                  | 0.000009 |

| query end            | 0.000007 |

| closing tables       | 0.000011 |

| freeing items        | 0.000121 |

| logging slow query   | 0.000001 |

| logging slow query   | 0.000093 |

| cleaning up          | 0.000010 |

+----------------------+----------+

17 rows in set, 1 warning (0.00 sec)
```
 

当查到最耗时的线程状态时，可以进一步选择all或者cpu,block io,page faults等明细类型来查看mysql在每个线程状态中使用什么资源上耗费了过高的时间：

```
show profile cpu for query 2;
```
  

上面的输出中可以以很高的精度显示了查询的响应时间，列出了查询执行的每个步骤花费的时间，其结果很难确定哪个步骤花费的时间太多，因为输出是按照执行顺序排序，而不是按照花费大小来排序的，如果要按照花费大小排序，就不能使用show prifile命令，而是直接使用information_schema.profiling表。如：

```
mysql> set profiling=1;

Query OK, 0 rows affected, 1 warning (0.00 sec)

 

mysql> select count(*) from xx;

+----------+

| count(*) |

+----------+

|   262144 |

+----------+

1 row in set (0.05 sec)

 

mysql> show profiles;

+----------+------------+-------------------------+

| Query_ID | Duration   | Query                   |

+----------+------------+-------------------------+

|        1 | 0.05509950 | select count(*) from xx |

+----------+------------+-------------------------+

1 row in set, 1 warning (0.00 sec)


mysql> set @query_id=1;

Query OK, 0 rows affected (0.00 sec)

mysql> select state,sum(duration) as total_r,round(100*sum(duration)/(select sum(duration) from information_schema.profiling where query_id=@query_id),2) as pct_r,count(*) as calls,sum(duration)/count(*) as "r/call" from information_schema.profiling where query_id=@query_id group by state order by total_r desc;

+----------------------+----------+-------+-------+--------------+

| state                | total_r  | pct_r | calls | r/call       |

+----------------------+----------+-------+-------+--------------+

| Sending data         | 0.054629 | 99.14 |     1 | 0.0546290000 |

| freeing items        | 0.000267 |  0.48 |     1 | 0.0002670000 |

| logging slow query   | 0.000070 |  0.13 |     2 | 0.0000350000 |

| starting             | 0.000040 |  0.07 |     1 | 0.0000400000 |

| Opening tables       | 0.000016 |  0.03 |     1 | 0.0000160000 |

| closing tables       | 0.000011 |  0.02 |     1 | 0.0000110000 |

| init                 | 0.000010 |  0.02 |     1 | 0.0000100000 |

| cleaning up          | 0.000010 |  0.02 |     1 | 0.0000100000 |

| end                  | 0.000009 |  0.02 |     1 | 0.0000090000 |

| statistics           | 0.000009 |  0.02 |     1 | 0.0000090000 |

| preparing            | 0.000008 |  0.01 |     1 | 0.0000080000 |

| query end            | 0.000007 |  0.01 |     1 | 0.0000070000 |

| System lock          | 0.000006 |  0.01 |     1 | 0.0000060000 |

| checking permissions | 0.000005 |  0.01 |     1 | 0.0000050000 |

| optimizing           | 0.000004 |  0.01 |     1 | 0.0000040000 |

| executing            | 0.000001 |  0.00 |     1 | 0.0000010000 |

+----------------------+----------+-------+-------+--------------+

16 rows in set (0.01 sec)
```
 

从上面的结果中可以看到，第一个是sending data（如果产生了临时表，第一就不是它了，那么临时表也是优先要解决的优化问题），另外还有sorting result（结果排序）也要注意，如果占比比较高，也要想办法优化，一般不建议在tuning sort buffer（优化排序缓冲区）或者类似的活动上花时间去优化。

 

如果要查询query id为1的Sending data状态的详细信息，可以使用如下SQL查询：

select * from information_schema.profiling where query_id=1 and state='Sending data'\G;

 

最后，做完剖析测试别忘记断开你的连接或者set profiling=0关闭这个功能