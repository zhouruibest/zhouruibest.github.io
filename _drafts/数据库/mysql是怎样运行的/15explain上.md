一条查询语句在经过MySQL查询优化器的各种基于成本和规则的优化会后生成一个所谓的执行计划，这个执行计划展示了接下来具体执行查询的方式，比如多表连接的顺序是什么，对于每个表采用什么访问方法来具体执行查询等等。设计MySQL的大佬贴心的为我们提供了EXPLAIN语句来帮助我们查看某个查询语句的具体执行计划


```
mysql> EXPLAIN SELECT 1;
+----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+----------------+
| id | select_type | table | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra          |
+----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+----------------+
|  1 | SIMPLE      | NULL  | NULL       | NULL | NULL          | NULL | NULL    | NULL | NULL |     NULL | No tables used |
+----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+----------------+
1 row in set, 1 warning (0.01 sec)

```

|列名|描述|
|-|-|
|id|在一个大的查询语句中每个SELECT关键字都对应一个唯一的id|
|select_type|SELECT关键字对应的那个查询的类型|
|table|表名|
|partitions|匹配的分区信息|
|type|针对单表的访问方法|
|possible_keys|可能用到的索引|
|key|实际上使用的索引|
|key_len|实际使用到的索引长度|
|ref|当使用索引列等值查询时，与索引列进行等值匹配的对象信息|
|rows|预估的需要读取的记录条数|
|filtered|某个表经过搜索条件过滤后剩余记录条数的百分比|
|Extra|一些额外的信息|

## talbe

EXPLAIN语句输出的每条记录都对应着某个单表的访问方法，该条记录的table列代表着该表的表名

## id

```sql
mysql> EXPLAIN SELECT * FROM s1 INNER JOIN s2;
+----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+---------------------------------------+
| id | select_type | table | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra                                 |
+----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+---------------------------------------+
|  1 | SIMPLE      | s1    | NULL       | ALL  | NULL          | NULL | NULL    | NULL | 9688 |   100.00 | NULL                                  |
|  1 | SIMPLE      | s2    | NULL       | ALL  | NULL          | NULL | NULL    | NULL | 9954 |   100.00 | Using join buffer (Block Nested Loop) |
+----+-------------+-------+------------+------+---------------+------+---------+------+------+----------+---------------------------------------+
2 rows in set, 1 warning (0.01 sec)
```

上述连接查询中参与连接的s1和s2表分别对应一条记录，但是这两条记录对应的id值都是1。这里需要大家记住的是，在连接查询的执行计划中，每个表都会对应一条记录，这些记录的id列的值是相同的，出现在前面的表表示驱动表，出现在后边的表表示被驱动表。所以从上面的EXPLAIN输出中我们可以看出，查询优化器准备让s1表作为驱动表，让s2表作为被驱动表来执行查询。

```sql
mysql> EXPLAIN SELECT * FROM s1 WHERE key1 IN (SELECT key3 FROM s2 WHERE common_field = 'a');
+----+-------------+-------+------------+------+---------------+----------+---------+-------------------+------+----------+------------------------------+
| id | select_type | table | partitions | type | possible_keys | key      | key_len | ref               | rows | filtered | Extra                        |
+----+-------------+-------+------------+------+---------------+----------+---------+-------------------+------+----------+------------------------------+
|  1 | SIMPLE      | s2    | NULL       | ALL  | idx_key3      | NULL     | NULL    | NULL              | 9954 |    10.00 | Using where; Start temporary |
|  1 | SIMPLE      | s1    | NULL       | ref  | idx_key1      | idx_key1 | 303     | xiaohaizi.s2.key3 |    1 |   100.00 | End temporary                |
+----+-------------+-------+------------+------+---------------+----------+---------+-------------------+------+----------+------------------------------+
2 rows in set, 1 warning (0.00 sec)
```
虽然查询语句是一个子查询，但是执行计划中s1和s2表对应的记录的id值全部是1，这就表明了查询优化器将子查询转换为了连接查询。

## select_type

一条大的查询语句里边可以包含若干个SELECT关键字，每个SELECT关键字代表着一个小的查询语句，而每个SELECT关键字的FROM子句中都可以包含若干张表（这些表用来做连接查询），每一张表都对应着执行计划输出中的一条记录，对于在同一个SELECT关键字中的表来说，它们的id值是相同的。

MySQL为每一个SELECT关键字代表的小查询都定义了一个称之为select_type的属性，意思是我们只要知道了某个小查询的select_type属性，就知道了这个小查询在整个大查询中扮演了一个什么角色.

|名称|描述|
|--|--|
|SIMPLE|Simple SELECT (not using UNION or subqueries)|
|PRIMARY|Outermost SELECT|
|UNION|Second or later SELECT statement in a UNION|
|UNION RESULT|Result of a UNION, MySQL选择使用临时表来完成UNION查询的去重工作，针对该临时表的查询的select_type就是UNION RESULT|
|SUBQUERY|First SELECT in subquery|
|DEPENDENT SUBQUERY|First SELECT in subquery, dependent on outer query|
|DEPENDENT UNION|Second or later SELECT statement in a UNION, dependent on outer query|
|DERIVED|Derived table|
|MATERIALIZED|Materialized subquery|
|UNCACHEABLE SUBQUERY|A subquery for which the result cannot be cached and must be re-evaluated for each row of the outer query|
|UNCACHEABLE UNION|The second or later select in a UNION that belongs to an uncacheable subquery (see UNCACHEABLE SUBQUERY)|


## type

执行计划的一条记录就代表着MySQL对某个表的执行查询时的访问方法，其中的type列就表明了这个访问方法是什么.

完整的访问方法包含了：system，const，eq_ref，ref，fulltext，ref_or_null，index_merge，unique_subquery，index_subquery，range，index，ALL

- system 当表中只有一条记录并且该表使用的存储引擎的统计数据是精确的，比如MyISAM、Memory，那么对该表的访问方法就是system

- const 根据主键或者唯一二级索引列与常数进行等值匹配时，对单表的访问方法就是const

- eq_ref 如果被驱动表是通过主键或者唯一二级索引列等值匹配的方式进行访问的（如果该主键或者唯一二级索引是联合索引的话，所有的索引列都必须进行等值比较），则对该被驱动表的访问方法就是eq_ref

- ref 当通过普通的二级索引列与常量进行等值匹配时来查询某个表，那么对该表的访问方法就**可能**是ref

— fulltext 全文索引

- ref_or_null 当对普通二级索引进行等值匹配查询，该索引列的值也可以是NULL值时，那么对该表的访问方法就**可能**是ref_or_null。

```sql
EXPLAIN SELECT * FROM s1 WHERE key1 = 'a' OR key1 IS NULL;
```

- index_merge 一般情况下对于某个表的查询只能使用到一个索引，但我们介绍单表访问方法时特意强调了在某些场景下可以使用Intersection、Union、Sort-Union这三种索引合并的方式来执行查询

- unique_subquery 类似于两表连接中被驱动表的eq_ref访问方法，unique_subquery是针对在一些包含IN子查询的查询语句中，如果查询优化器决定将IN子查询转换为EXISTS子查询，而且子查询可以使用到主键进行等值匹配的话，那么该子查询执行计划的type列的值就是unique_subquery

- index_subquery index_subquery与unique_subquery类似，只不过访问子查询中的表时使用的是普通的索引

- range 如果使用索引获取某些范围区间的记录，那么就可能使用到range访问方法

- index 当我们可以使用索引覆盖，但需要扫描全部的索引记录时，该表的访问方法就是index

- ALL 全表扫描

## possible_keys和key

possible_keys列表示在某个查询语句中，对某个表执行单表查询时可能用到的索引有哪些，key列表示实际用到的索引有哪些

不过有一点比较特别，就是在使用index访问方法来查询某个表时，possible_keys列是空的，而key列展示的是实际使用到的索引

## key_len

表示当优化器决定使用某个索引执行查询时，该索引记录的最大长度

## ref

当使用索引列等值匹配的条件去执行查询时，也就是在访问方法是const、eq_ref、ref、ref_or_null、unique_subquery、index_subquery其中之一时，ref列展示的就是与索引列作等值匹配的东东是什么，比如只是一个常数或者是某个列。大家看下面这个查询：

1. 可以看到ref列的值是const，表明在使用idx_key1索引执行查询时，与key1列作等值匹配的对象是一个常数，当然有时候更复杂一点：
```sql
mysql> EXPLAIN SELECT * FROM s1 WHERE key1 = 'a';
+----+-------------+-------+------------+------+---------------+----------+---------+-------+------+----------+-------+
| id | select_type | table | partitions | type | possible_keys | key      | key_len | ref   | rows | filtered | Extra |
+----+-------------+-------+------------+------+---------------+----------+---------+-------+------+----------+-------+
|  1 | SIMPLE      | s1    | NULL       | ref  | idx_key1      | idx_key1 | 303     | const |    8 |   100.00 | NULL  |
+----+-------------+-------+------------+------+---------------+----------+---------+-------+------+----------+-------+
1 row in set, 1 warning (0.01 sec)

```

2. 可以看到对被驱动表s2的访问方法是eq_ref，而对应的ref列的值是xiaohaizi.s1.id，这说明在对被驱动表进行访问时会用到PRIMARY索引，也就是聚簇索引与一个列进行等值匹配的条件，于s2表的id作等值匹配的对象就是xiaohaizi.s1.id列（注意这里把数据库名也写出来了）。
```sql
mysql> EXPLAIN SELECT * FROM s1 INNER JOIN s2 ON s1.id = s2.id;
+----+-------------+-------+------------+--------+---------------+---------+---------+-----------------+------+----------+-------+
| id | select_type | table | partitions | type   | possible_keys | key     | key_len | ref             | rows | filtered | Extra |
+----+-------------+-------+------------+--------+---------------+---------+---------+-----------------+------+----------+-------+
|  1 | SIMPLE      | s1    | NULL       | ALL    | PRIMARY       | NULL    | NULL    | NULL            | 9688 |   100.00 | NULL  |
|  1 | SIMPLE      | s2    | NULL       | eq_ref | PRIMARY       | PRIMARY | 4       | xiaohaizi.s1.id |    1 |   100.00 | NULL  |
+----+-------------+-------+------------+--------+---------------+---------+---------+-----------------+------+----------+-------+
2 rows in set, 1 warning (0.00 sec)
```

## rows

如果查询优化器决定使用全表扫描的方式对某个表执行查询时，执行计划的rows列就代表预计需要扫描的行数，如果使用索引来执行查询时，执行计划的rows列就代表预计扫描的索引记录行数。比如下面这个查询：


## filtered

之前在分析连接查询的成本时提出过一个condition filtering的概念，就是MySQL在计算驱动表扇出时采用的一个策略：

- 如果使用的是全表扫描的方式执行的单表查询，那么计算驱动表扇出时需要估计出满足搜索条件的记录到底有多少条。
- 如果使用的是索引执行的单表扫描，那么计算驱动表扇出的时候需要估计出满足除使用到对应索引的搜索条件外的其他搜索条件的记录有多少条。

例子1
```sql
mysql> EXPLAIN SELECT * FROM s1 WHERE key1 > 'z' AND common_field = 'a';
+----+-------------+-------+------------+-------+---------------+----------+---------+------+------+----------+------------------------------------+
| id | select_type | table | partitions | type  | possible_keys | key      | key_len | ref  | rows | filtered | Extra                              |
+----+-------------+-------+------------+-------+---------------+----------+---------+------+------+----------+------------------------------------+
|  1 | SIMPLE      | s1    | NULL       | range | idx_key1      | idx_key1 | 303     | NULL |  266 |    10.00 | Using index condition; Using where |
+----+-------------+-------+------------+-------+---------------+----------+---------+------+------+----------+------------------------------------+
1 row in set, 1 warning (0.00 sec)
```

从执行计划的key列中可以看出来，该查询使用idx_key1索引来执行查询，从rows列可以看出满足key1 > 'z'的记录有266条。执行计划的filtered列就代表查询优化器预测在这266条记录中，有多少条记录满足其余的搜索条件，也就是common_field = 'a'这个条件的百分比。此处filtered列的值是10.00，说明查询优化器预测在266条记录中有10.00%的记录满足common_field = 'a'这个条件。

对于单表查询来说，这个filtered列的值没什么意义，**我们更关注在连接查询中驱动表对应的执行计划记录的filtered值**.

例子二

```sql
mysql> EXPLAIN SELECT * FROM s1 INNER JOIN s2 ON s1.key1 = s2.key1 WHERE s1.common_field = 'a';
+----+-------------+-------+------------+------+---------------+----------+---------+-------------------+------+----------+-------------+
| id | select_type | table | partitions | type | possible_keys | key      | key_len | ref               | rows | filtered | Extra       |
+----+-------------+-------+------------+------+---------------+----------+---------+-------------------+------+----------+-------------+
|  1 | SIMPLE      | s1    | NULL       | ALL  | idx_key1      | NULL     | NULL    | NULL              | 9688 |    10.00 | Using where |
|  1 | SIMPLE      | s2    | NULL       | ref  | idx_key1      | idx_key1 | 303     | xiaohaizi.s1.key1 |    1 |   100.00 | NULL        |
+----+-------------+-------+------------+------+---------------+----------+---------+-------------------+------+----------+-------------+
2 rows in set, 1 warning (0.00 sec)

```
从执行计划中可以看出来，查询优化器打算把s1当作驱动表，s2当作被驱动表。我们可以看到驱动表s1表的执行计划的rows列为9688， filtered列为10.00，这意味着驱动表s1的扇出值就是9688 × 10.00% = 968.8，这说明还要对被驱动表执行大约968次查询。

























