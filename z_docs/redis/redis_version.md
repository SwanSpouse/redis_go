### redis版本

#### 2.6  

Lua脚本支持

新增PEXIRE、PTTL、PSETEX过期设置命令，key过期时间可以设置为毫秒级

新增位操作命令：BITCOUNT、BITOP

新增命令：dump、restore，即序列化与反序列化操作

新增命令：INCRBYFLOAT、HINCRBYFLOAT，用于对值进行浮点数的加减操作

新增命令：MIGRATE，用于将key原子性地从当前实例传送到目标实例的指定数据库上

放开了对客户端的连接数限制

hash函数种子随机化，有效防止碰撞

SHUTDOWN命令添加SAVE和NOSAVE两个参数，分别用于指定SHUTDOWN时用不用执行写RDB的操作

虚拟内存Virtual Memory相关代码全部去掉

sort命令会拒绝无法转换成数字的数据模型元素进行排序

不同客户端输出缓冲区分级，比如普通客户端、slave机器、pubsub客户端，可以分别控制对它们的输出缓冲区大小

更多的增量过期(减少阻塞)的过期key收集算法 ,当非常多的key在同一时间失效的时候,意味着redis可以提高响应的速度

底层数据结构优化，提高存储大数据时的性能

#### 2.8

引入PSYNC，主从可以增量同步，这样当主从链接短时间中断恢复后，无需做完整的RDB完全同步

从显式ping主，主可以扫描到可能超时的从

新增命令：SCAN、SSCAN、HSCAN和ZSCAN

crash的时候自动内存检查

新增键空间通知功能，客户端可以通过订阅/发布机制，接收改动了redis指定数据集的事件

可绑定多个IP地址

可通过CONFIGSET设置客户端最大连接数

新增CONFIGREWRITE命令，可以直接把CONFIGSET的配置修改到redis.conf里

新增pubsub命令，可查看pub/sub相关状态

支持引用字符串，如set 'foo bar' "hello world\n"

新增redis master-slave集群高可用解决方案（Redis-Sentinel）

当使用SLAVEOF命令时日志会记录下新的主机

#### 3.0

实现了分布式的Redis即Redis Cluster，从而做到了对集群的支持

全新的"embedded string"对象编码方式，从而实现了更少的缓存丢失和性能的提升

大幅优化LRU近似算法的性能

新增CLIENT PAUSE命令，可以在指定时间内停止处理客户端请求

新增WAIT命令，可以阻塞当前客户端，直到所有以前的写命令都成功传输并和指定的slaves确认

AOF重写过程中的"last write"操作降低了AOF child -> parent数据传输的延迟

实现了对MIGRATE连接缓存的支持，从而大幅提升key迁移的性能

为MIGRATE命令新增参数：copy和replace，copy不移除源实例上的key，replace替换目标实例上已存在的key

提高了BITCOUNT、INCR操作的性能

调整Redis日志格式

#### 3.2

新增对GEO（地理位置）功能的支持

新增Lua脚本远程debug功能

SDS相关的优化，提高redis性能

修改Jemalloc相关代码，提高redis内存使用效率

提高了主从redis之间关于过期key的一致性

支持利用upstart和systemd管理redis进程

将list底层数据结构类型修改为quicklist，在内存占用和RDB文件大小方面有极大的提升

SPOP命令新增count参数，可控制随机删除元素的个数

支持为RDB文件增加辅助字段，比如创建日期，版本号等，新版本可以兼容老版本RDB文件，反之不行

通过调整哈希表大小的操作码RDB_OPCODE_RESIZEDB，redis可以更快得读RDB文件

新增HSTRLEN命令，返回hash数据类型的value长度

提供了一个基于流水线的MIGRATE命令，极大提升了命令执行速度

redis-trib.rb中实现将slot进行负载均衡的功能

改进了从机迁移的功能

改进redis sentine高可用方案，使之可以更方便地监控多个redis主从集群

#### 4.0

加入模块系统，用户可以自己编写代码来扩展和实现redis本身不具备的功能，它与redis内核完全分离，互不干扰

优化了PSYNC主从复制策略，使之效率更高

为DEL、FLUSHDB、FLUSHALL命令提供非阻塞选项，可以将这些删除操作放在单独线程中执行，从而尽可能地避免服务器阻塞

新增SWAPDB命令，可以将同一redis实例指定的两个数据库互换

新增RDB-AOF持久化格式，开启后，AOF重写产生的文件将同时包含RDB格式的内容和AOF格式的内容，其中 RDB格式的内容用于记录已有的数据，而AOF格式的内存则用于记录最近发生了变化的数据

新增MEMORY内存命令，可以用于查看某个key的内存使用、查看整体内存使用细节、申请释放内存、深入查看内存分配器内部状态等功能

兼容NAT和Docker