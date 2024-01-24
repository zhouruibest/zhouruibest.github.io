MySQLdump对于MySQL数据库备份是有一个很好用的命令，并且是MySQL自带的。

-d：只备份表结构，备份文件是SQL语句形式；只备份创建表的语句，插入的数据不备份。

-t：只备份数据，数据是文本形式；表结构不备份

-T [--tab]：表结构与数据分离，表结构为sql文件，数据为普通文件

-A：导出所有数据库

-B：导出指定数据库

-x, --lock-all-tables： 锁表
锁表原理：从执行定时备份脚本起（带-x参数），不能往表里更新，但是缺点，锁表后无法更新，如果单库一般在低谷，比如凌晨后半夜里；多库，就从从库里锁表备份（并且从库不对外，只做备份）
Locks all tables across all databases. This is achieved by taking a global read lock for the duration of the whole dump.
Automatically turns --single-transaction and --lock-tables off. 启用该选项，会自动关闭 --single-transaction 和 --lock-tables.

-l, --lock-tables： 只读锁表
Lock all tables before dumping them
Lock all tables for read.
(Defaults to on; use --skip-lock-tables to disable.)
该选项默认打开的，上面已经说到了。它的作用是在导出过程中锁定所有表。--single-transaction 和 --lock-all-tables 都会将该选项关闭。
在用LOCK TABLES给表显式加表锁时，必须同时取得所有涉及到表的锁，也就是说，在执行LOCK TABLES后，只能访问显式加锁的这些表，不能访问未加锁的表；同时，如果加的是读锁，那么只能执行锁表的查询操作，MyISAM总是一次获得SQL语句所需要的全部锁。这也正是MyISAM表不会出现死锁（Deadlock Free）的原因。

--single-transaction
--single-transaction 可以得到一致性的导出结果。他是通过将导出行为放入一个事务中达到目的的。
它有一些要求：只能是 innodb 引擎；导出的过程中，不能有任何人执行 alter table, drop table, rename table, truncate table等DDL语句。
实际上DDL会被事务所阻塞，因为事务持有表的metadata lock 的共享锁，而DDL会申请metadata lock的互斥锁，所以阻塞了。
--single-transaction 会自动关闭 --lock-tables 选项；上面我们说到mysqldump默认会打开了--lock-tables，它会在导出过程中锁定所有表。
因为 --single-transaction 会自动关闭--lock-tables，所以单独使用--single-transaction是不会使用锁的。与 --master-data 合用才有锁。

-q： 不做缓冲查询，直接导到标准输出

-R：导出存储过程和函数

-E,--events：导出调度事件

--add-drop-database
在CREATE DATABASE语句前增加DROP DATABASE语句，一般配合--all-databases 或 --databases使用，因为只有使用了这二者其一，才会记录CREATE DATABASE语句。

--add-drop-table 
在CREATE TABLE语句前增加DROP TABLE语句。

--add-drop-trigger
在CREATE TRIGGER语句前增加DROP TRIGGER语句

--all-tablespaces, -Y
导出全部表空间。该参数目前仅用在MySQL Cluster表上（NDB引擎）

--add-locks
在每个表导出之前增加LOCK TABLES并且之后UNLOCK TABLE。(默认为打开状态，使用--skip-add-locks取消选项)

--opt
等同于--add-drop-table, --add-locks, --create-options, --quick, --extended-insert, --lock-tables, --set-charset, --disable-keys 该选项默认开启, 可以用--skip-opt禁用.

-F,--flush-logs：刷新binlog日志

--master-data
mysqldump导出数据时，当这个参数的值为1的时候，mysqldump出来的备份文件就会包括CHANGE MASTER TO这个语句，CHANGE MASTER TO后面紧接着就是file和position的记录，在slave上导入数据时会执行该语句，salve就会根据CHANGE MASTER TO后面指定的binlog文件和binlog日志文件里的位置点，从master端复制binlog。默认情况下这个值是1 。运维经常使用到该参数，主从复制时，该参数是一个很好的功能，同时也可以做增量恢复。
当这个参数的值为2的时候mysqldump导出来的备份文件也会包含CHANGE MASTER TO语句，但是该语句被注释掉，不会生效，只是提供一个信息。

--master-data也会刷新binlog日志。





