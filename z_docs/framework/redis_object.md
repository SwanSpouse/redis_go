

![framework](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/framework/redis_object.png?raw=true)

## redis_object 结构设计

### TBase

base interface。 这里定义了作为redis object最基本方法。如：

* object的类型、编码和存储的值等。

* object的TTL，是否过期判断，LRU等等。

### 五种基本类型

TString TList THash TSet TZSet作为redis object的五种基本类型。

在TString TList THash TSet TZSet中分别定义了五种不同类型应该实现的功能。

### 类型的编码方式

encodings目录下的文件对应基本类型的实现方式。

例如TString类型有go string 和 int两种编码方式。string_raw.go和string_int.go就是为两种编码方式的不同实现。


### raw_type 原始类型

这里面定义了list 和 dict两种基本的数据结构。

* list 是双端链表实现的链表结构。list不是线程安全的。

* dict 是参考Java 1.7 ConcurrentHashMap实现的HashMap结构。

dict作为一种基本的数据结构。既是THash和TSet基本类型的底层编码方式。同样是redis database的实现方式。

所以为了提升系统的访问效率，dict是支持并发访问的。


