

### Dict 

#### 元素

DictEntry:
* Key-Value对，HashMap的最基本元素，同时里面存储着Key的Hash值，以免除在后续的处理流程中重复计算。

Segment:
* 等同于一个HashTMap。低层是一个Key-Value组成的数组，根据DictEntry的哈希值来计算DictEntry在Segment中的位置。
在发生Hash冲突时，采用的办法是链接地址法来解决Hash冲突。

Dict:
* 由多个Segment组成，为了解决并发处理的问题。在理想情况下，有多少个Segment就支持多少个线程的并发操作。

#### 整体架构
![dict](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/data_type/dict.png?raw=true)

#### 哈希值计算

Hash值的计算采用的github的第三方库(https://github.com/mitchellh/hashstructure)来根据计算DictEntry的Key计算哈希值。

#### 并发处理相关

参照 Java 1.7 ConcurrentHashMap的设计。

* 将所有的DictEntry利用hash值均匀分布到各个Segment中，不同的Segment之间相互独立，这样在理想情况下就可以支持和
Segment个数等同的并发量。

Segment

* 针对单独每一个Segment都有自己的读写锁，所以在同一个Segment之内，支持并发的读，但是对于写操作，还是串行的。

#### Hash索引相关 

计算Key在Dict中Segment的位置以及在Segment中对应的Index

* 这里有限有一个前提，就是保证Segment的size都是2的幂次，这样能够保证size -1的二进制表示所有位数上都是1。

* 定义sizeMask = size-1，在计算index的时候，用key对应的hashCode & sizeMask（hashCode与sizeMask进行与操作），这样就能够保证能到的
结果处于[0, size)区间内，不会溢出。同时在rehash的过程中，同样保证size的大小是2的幂次，这样，老元素在新Segment中的index，不是在原index，
就是在原index + 2^n 的位置上，方便计算。

* 仅用hashCode和sizeMask计算index会导致一个问题，就是在Dict中Segment的序号有可能和Segment中index的序号有可能重复，大大增加了Hash冲突
的概率，所以在计算Segment序号的时候多了一个SegmentShift，来避免Segment序号和Segment中的index不重复。

* DictEntry在Dict的哪个Segment中的计算方法：

```golang
	segmentIdx := (hashCode >> dict.segmentShift) & dict.segmentMask
```

* DictEntry在Segment中index的计算方法：

```golang
	idx := hashCode & dict.segments[segmentIdx].sizeMask
```

#### 哈希冲突 与 rehash

当Segment中的元素个数超过Segment容量的阈值后，再执行add操作，就会触发rehash操作。

* 在rehash的过程里，首先给当前的Segment加锁，同时创建一个新的Segment，将所有的DictEntry都从旧的Segment迁移到新的Segment，在
整体迁移完成之后用新的Segment替换老的Segment

#### TODO

参照Java1.7 ConcurrentHashMap这里有一个优化的点

* 这里采用了乐观锁的思想：在每次操作之前，不是先进行无脑的加锁。而是不考虑并发，直接来进行操作。同时记录segment的modCount， 在操作结束
之后判断一下modCount和操作之前是否一致，如果不一致，说明有线程在并发的进行操作。本次操作就会因为并发的原因置为失败。同时再重复本次操作，当
因为并发原因失败的次数达到一定的限制。再加锁来进行处理。

* 这样的做法虽然降低了并发条件下的操作效率，但是对于整体系统来说，并发操作的情况比较少，不是每次操作都无脑加锁，这样能够提高系统整体的效率。
