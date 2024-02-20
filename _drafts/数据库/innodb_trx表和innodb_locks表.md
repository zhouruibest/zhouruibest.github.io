# innodb_trx

1. innodb_trx表提供了当前innodb引擎内每个事务的信息（只读事务除外），包括当一个事务启动，事务是否在等待一个锁，以及交易正在执行的语句（如果有的话）。查询语句：

```sql
select * from information_schema.innodb_trx;
select * from information_schema.innodb_trx\G
```

2. innodb_trx表列信息详解：

```
trx_id：唯一事务id号，只读事务和非锁事务是不会创建id的。
TRX_WEIGHT：事务的高度，代表修改的行数（不一定准确）和被事务锁住的行数。为了解决死锁，innodb会选择一个高度最小的事务来当做牺牲品进行回滚。已经被更改的非交易型表的事务权重比其他事务高，即使改变的行和锁住的行比其他事务低。
TRX_STATE：事务的执行状态，值一般分为：RUNNING, LOCK WAIT, ROLLING BACK, and COMMITTING.
TRX_STARTED：事务的开始时间
TRX_REQUESTED_LOCK_ID:如果trx_state是lockwait,显示事务当前等待锁的id，不是则为空。想要获取锁的信息，根据该lock_id，以innodb_locks表中lock_id列匹配条件进行查询，获取相关信息。
TRX_WAIT_STARTED：如果trx_state是lockwait,该值代表事务开始等待锁的时间；否则为空。
TRX_MYSQL_THREAD_ID：mysql线程id。想要获取该线程的信息，根据该thread_id，以INFORMATION_SCHEMA.PROCESSLIST表的id列为匹配条件进行查询。
TRX_QUERY：事务正在执行的sql语句。
TRX_OPERATION_STATE：事务当前的操作状态，没有则为空。
TRX_TABLES_IN_USE：事务在处理当前sql语句使用innodb引擎表的数量。
TRX_TABLES_LOCKED：当前sql语句有行锁的innodb表的数量。（因为只是行锁，不是表锁，表仍然可以被多个事务读和写）
TRX_LOCK_STRUCTS：事务保留锁的数量。
TRX_LOCK_MEMORY_BYTES：在内存中事务索结构占得空间大小。
TRX_ROWS_LOCKED：事务行锁最准确的数量。这个值可能包括对于事务在物理上存在，实际不可见的删除标记的行。
TRX_ROWS_MODIFIED：事务修改和插入的行数
TRX_CONCURRENCY_TICKETS：该值代表当前事务在被清掉之前可以多少工作，由 innodb_concurrency_tickets系统变量值指定。
TRX_ISOLATION_LEVEL：事务隔离等级。
TRX_UNIQUE_CHECKS：当前事务唯一性检查启用还是禁用。当批量数据导入时，这个参数是关闭的。
TRX_FOREIGN_KEY_CHECKS：当前事务的外键坚持是启用还是禁用。当批量数据导入时，这个参数是关闭的。
TRX_LAST_FOREIGN_KEY_ERROR：最新一个外键错误信息，没有则为空。
TRX_ADAPTIVE_HASH_LATCHED：自适应哈希索引是否被当前事务阻塞。当自适应哈希索引查找系统分区，一个单独的事务不会阻塞全部的自适应hash索引。自适应hash索引分区通过 innodb_adaptive_hash_index_parts参数控制，默认值为8。
TRX_ADAPTIVE_HASH_TIMEOUT：是否为了自适应hash索引立即放弃查询锁，或者通过调用mysql函数保留它。当没有自适应hash索引冲突，该值为0并且语句保持锁直到结束。在冲突过程中，该值被计数为0，每句查询完之后立即释放门闩。当自适应hash索引查询系统被分区（由 innodb_adaptive_hash_index_parts参数控制），值保持为0。
TRX_IS_READ_ONLY：值为1表示事务是read only。
TRX_AUTOCOMMIT_NON_LOCKING：值为1表示事务是一个select语句，该语句没有使用for update或者shared mode锁，并且执行开启了autocommit，因此事务只包含一个语句。当TRX_AUTOCOMMIT_NON_LOCKING和TRX_IS_READ_ONLY同时为1，innodb通过降低事务开销和改变表数据库来优化事务。
```

# innodb_locks详解
1. INFORMATION_SCHEMA INNODB_LOCKS 提供innodb事务去请求但没有获取到的锁信息和事务阻塞其他事务的锁信息。执行命令如下：

```sql
select * from information_schema.innodb_locks\G
```

2. innodb_locks各列参数详解：

```sql
lock_id:innodb唯一lock id。把他当做一个不透明的字符串。虽然lock_id当前包含trx_id，lock_id的数据格式在任何时间都肯能改变。不要写用于解析lock_id值得应用程序。
lock_trx_id：持有锁的事务id。查询事务信息，与innodb_trx表中trx_id列匹配。
lock_mode:锁请求。该值包括： S, X, IS, IX, GAP, AUTO_INC, and UNKNOWN。锁模式标识符可以组合用于识别特定的锁模式。查看更多信息，点击[此处]((https://dev.mysql.com/doc/refman/8.0/en/innodb-locking.html))
lock_type:锁类型。行锁为record，表锁为table。
lock_table:被锁的表名，或者包含锁记录的表名。
lock_index:lock_type为行锁时，该值为索引名，否则为空。
lock_space:lock_type为行锁时，该值为锁记录的表空间的id，否则为空。
lock_page：lock_type为行锁时，该值为锁记录页数量，否则为空。
lock_rec:lock_type为行锁时，页内锁记录的堆数，否则为空。
lock_data：与锁相关的数据。如果lock_type为行锁时，该值是锁记录的主键值，否则为空。这列包含锁定行的主键列的值，转化为一个有效的字符串，如果没有主键，lock_data是唯一innodb内部行id号。如果是键值或者范围大于索引的最大值会使用间隙锁，lock_data表示为supremum pseudo-record。当包含锁记录的页不在buffer pool内，innodb不去从磁盘获取页，为了避免不必要的磁盘操作，lock_data为空。
```