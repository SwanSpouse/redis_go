

#### SQL语句执行顺序

查询中用到的关键词主要包含六个，并且他们的顺序依次为: 

* select --> from --> where --> group by --> having --> order by 

SQL Select语句完整的执行顺序【从DBMS使用者角度】:

1. from:需要从哪个数据表检索数据 

2. where:过滤表中数据的条件 

3. group by:如何将上面过滤出的数据分组 

4. having:对上面已经分组的数据进行过滤的条件  

5. select:查看结果集中的哪个列，或列的计算结果 

6. order by :按照什么样的顺序来查看返回的数据 