
#### HashMap

Java 1.7 HashMap

* 问题：如果在resize的过程中，插入元素，有可能会导致元素丢失？所以HashMap不是线程安全的？

* table 定义

```Java
    // transient定义的变量不参与序列化
    transient Entry<K,V>[] table = (Entry<K,V>[]) EMPTY_TABLE;
```

* put方法

```Java
    /**
     * Associates the specified value with the specified key in this map.
     * If the map previously contained a mapping for the key, the old
     * value is replaced.
     *
     * @param key key with which the specified value is to be associated
     * @param value value to be associated with the specified key
     * @return the previous value associated with <tt>key</tt>, or
     *         <tt>null</tt> if there was no mapping for <tt>key</tt>.
     *         (A <tt>null</tt> return can also indicate that the map
     *         previously associated <tt>null</tt> with <tt>key</tt>.)
     */
    public V put(K key, V value) {
        // 如果table为空，则进行初始化
        if (table == EMPTY_TABLE) {
            inflateTable(threshold);
        }
        // 如果当前key == null 则放到table的第一个位置。因为null没有办法求hash值
        if (key == null)
            return putForNullKey(value);

        // 计算当前key的hash值
        int hash = hash(key);

        // 计算当前hash值在table中的位置
        int i = indexFor(hash, table.length);

        // 遍历table[i]位置的链表
        for (Entry<K,V> e = table[i]; e != null; e = e.next) {
            Object k;
            /*
             * 如果当前Entry的hash和输入key的hash相同（Hash冲突）且Entry的key和输入的key 相等，则用新的值覆盖老的值。
             */
            if (e.hash == hash && ((k = e.key) == key || key.equals(k))) {
                V oldValue = e.value;
                e.value = value;
                e.recordAccess(this);
                return oldValue;
            }
        }
        /**
         * The number of times this HashMap has been structurally modified
         * Structural modifications are those that change the number of mappings in
         * the HashMap or otherwise modify its internal structure (e.g.,
         * rehash).  This field is used to make iterators on Collection-views of
         * the HashMap fail-fast.  (See ConcurrentModificationException).
         */
        modCount++;
        // 将当前Entry放到table中的位置。
        addEntry(hash, key, value, i);
        return null;
    }
```

* get方法

```Java
    /**
     * Returns the value to which the specified key is mapped,
     * or {@code null} if this map contains no mapping for the key.
     *
     * <p>More formally, if this map contains a mapping from a key
     * {@code k} to a value {@code v} such that {@code (key==null ? k==null :
     * key.equals(k))}, then this method returns {@code v}; otherwise
     * it returns {@code null}.  (There can be at most one such mapping.)
     *
     * <p>A return value of {@code null} does not <i>necessarily</i>
     * indicate that the map contains no mapping for the key; it's also
     * possible that the map explicitly maps the key to {@code null}.
     * The {@link #containsKey containsKey} operation may be used to
     * distinguish these two cases.
     *
     * @see #put(Object, Object)
     */
    public V get(Object key) {
        // 如果Key为空，则返回Null Key的值
        if (key == null)
            return getForNullKey();
        Entry<K,V> entry = getEntry(key);

        return null == entry ? null : entry.getValue();
    }

    /**
     * Returns the entry associated with the specified key in the
     * HashMap.  Returns null if the HashMap contains no mapping
     * for the key.
     */
    final Entry<K,V> getEntry(Object key) {
        // 如果table的size是0，返回null
        if (size == 0) {
            return null;
        }
        // 计算key的hash值
        int hash = (key == null) ? 0 : hash(key);
        // 遍历table hash值对应位置上的链表
        for (Entry<K,V> e = table[indexFor(hash, table.length)];
             e != null;
             e = e.next) {
            Object k;
            // 找到hash相同且key相等的Entry将value进行返回。
            if (e.hash == hash && ((k = e.key) == key || (key != null && key.equals(k))))
                return e;
        }
        return null;
    }
```

* rehash方法

```Java
    /**
     * Adds a new entry with the specified key, value and hash code to
     * the specified bucket.  It is the responsibility of this
     * method to resize the table if appropriate.
     *
     * Subclass overrides this to alter the behavior of put method.
     */
    void addEntry(int hash, K key, V value, int bucketIndex) {
        // 在插入新节点的时候，会对size是否超过阈值进行判断，如果超过会触发resize
        if ((size >= threshold) && (null != table[bucketIndex])) {
            resize(2 * table.length);
            hash = (null != key) ? hash(key) : 0;
            bucketIndex = indexFor(hash, table.length);
        }

        createEntry(hash, key, value, bucketIndex);
    }

    /**
     * Rehashes the contents of this map into a new array with a
     * larger capacity.  This method is called automatically when the
     * number of keys in this map reaches its threshold.
     *
     * If current capacity is MAXIMUM_CAPACITY, this method does not
     * resize the map, but sets threshold to Integer.MAX_VALUE.
     * This has the effect of preventing future calls.
     *
     * @param newCapacity the new capacity, MUST be a power of two;
     *        must be greater than current capacity unless current
     *        capacity is MAXIMUM_CAPACITY (in which case value
     *        is irrelevant).
     */
    void resize(int newCapacity) {
        Entry[] oldTable = table;
        int oldCapacity = oldTable.length;
        // 如果旧的table的容量已经达到MAXIMUM_CAPACITY 2ˆ30的大小，则阈值就被设置为Integer.MAX_VALUE
        if (oldCapacity == MAXIMUM_CAPACITY) {
            threshold = Integer.MAX_VALUE;
            return;
        }
        // 创建一个newCapacity这么大的新的table
        Entry[] newTable = new Entry[newCapacity];
        // 将数据进行转移
        transfer(newTable, initHashSeedAsNeeded(newCapacity));
        // table从老table指向新table
        table = newTable;
        // 重新计算容量阈值
        threshold = (int)Math.min(newCapacity * loadFactor, MAXIMUM_CAPACITY + 1);
    }

    /**
     * Transfers all entries from current table to newTable.
     */
    void transfer(Entry[] newTable, boolean rehash) {
        int newCapacity = newTable.length;
        // 遍历老table里面的所有Entry
        for (Entry<K,V> e : table) {
            // 遍历所有Key Hash值相同的Entry
            while(null != e) {
                Entry<K,V> next = e.next;
                if (rehash) {
                    e.hash = null == e.key ? 0 : hash(e.key);
                }
                int i = indexFor(e.hash, newCapacity);
                // 头插法，先将e.next指向newTable i位置的元素，然后将e赋值给newTable的位置元素
                // 在并发的条件下，HashMap容易形成回路。造成死循环。
                /*
                假如有两个线程P1、P2，以及链表 a=>b=>null
                    1、P1先运行，运行完"Entry<K,V> next = e.next;"代码后发生堵塞，或者其它情况不再运行下去，此时e=a。next=b
                    2、而P2已经运行完整段代码，于是当前的新链表newTable[i]为b=>a=>null
                    3、P1又继续运行"Entry<K,V> next = e.next;"之后的代码，则运行完"e=next;"后，newTable[i]为a<=>b。则造成回路，while(e!=null)一直死循环
                */
                e.next = newTable[i];
                newTable[i] = e;
                e = next;
            }
        }
    }
```

Java 1.8 HashMap

* table 定义

```Java
    transient Node<K,V>[] table;
```

* put 方法
```Java
 /**
     * Associates the specified value with the specified key in this map.
     * If the map previously contained a mapping for the key, the old
     * value is replaced.
     *
     * @param key key with which the specified value is to be associated
     * @param value value to be associated with the specified key
     * @return the previous value associated with <tt>key</tt>, or
     *         <tt>null</tt> if there was no mapping for <tt>key</tt>.
     *         (A <tt>null</tt> return can also indicate that the map
     *         previously associated <tt>null</tt> with <tt>key</tt>.)
     */
    public V put(K key, V value) {
        return putVal(hash(key), key, value, false, true);
    }

    /**
     * Implements Map.put and related methods
     *
     * @param hash hash for key
     * @param key the key
     * @param value the value to put
     * @param onlyIfAbsent if true, don't change existing value
     * @param evict if false, the table is in creation mode.
     * @return previous value, or null if none
     */
    final V putVal(int hash, K key, V value, boolean onlyIfAbsent, boolean evict) {
        Node<K,V>[] tab; Node<K,V> p; int n, i;
        if ((tab = table) == null || (n = tab.length) == 0)
            // 如果table == null 或者 length = 0 的话，就需要进行初始化，resize函数也有初始化的功能。
            n = (tab = resize()).length;
        if ((p = tab[i = (n - 1) & hash]) == null)
            // 找到hash对应数组中的位置，如果那个位置为null，则newNode出一个元素来。放到相应的位置上去。
            tab[i] = newNode(hash, key, value, null);
        else {
             // 找到hash对应数组中的位置，如果那个位置不为空有元素在，
            Node<K,V> e; K k;
            if (p.hash == hash && ((k = p.key) == key || (key != null && key.equals(k))))
            // 如果p e的hash值相同并且key值相等。则用e = p
                e = p;
            else if (p instanceof TreeNode)
                // 如果p当前的结构是红黑树，则进行putTreeVal操作
                e = ((TreeNode<K,V>)p).putTreeVal(this, tab, hash, key, value);
            else {
                // 如果p当前的结构不是红黑树，则把元素插入到链表的末端。
                for (int binCount = 0; ; ++binCount) {
                    if ((e = p.next) == null) {
                        p.next = newNode(hash, key, value, null);
                        if (binCount >= TREEIFY_THRESHOLD - 1) // -1 for 1st
                            treeifyBin(tab, hash);
                        break;
                    }
                    // 如果在遍历的过程中，发现有key相同的，则退出
                    if (e.hash == hash && ((k = e.key) == key || (key != null && key.equals(k))))
                        break;
                    p = e;
                }
            }
            // 如果有key相同的，用新的值覆盖老的值。并返回老的值。
            if (e != null) { // existing mapping for key
                V oldValue = e.value;
                if (!onlyIfAbsent || oldValue == null)
                    e.value = value;
                afterNodeAccess(e);
                return oldValue;
            }
        }
        // modCount++ 保证fast-fail
        ++modCount;
        // 如果size > threshold会进行resize
        if (++size > threshold)
            resize();
        // 这是一个callback
        afterNodeInsertion(evict);
        return null;
    }
```


* get 方法

```Java
    /**
     * Returns the value to which the specified key is mapped,
     * or {@code null} if this map contains no mapping for the key.
     *
     * <p>More formally, if this map contains a mapping from a key
     * {@code k} to a value {@code v} such that {@code (key==null ? k==null :
     * key.equals(k))}, then this method returns {@code v}; otherwise
     * it returns {@code null}.  (There can be at most one such mapping.)
     *
     * <p>A return value of {@code null} does not <i>necessarily</i>
     * indicate that the map contains no mapping for the key; it's also
     * possible that the map explicitly maps the key to {@code null}.
     * The {@link #containsKey containsKey} operation may be used to
     * distinguish these two cases.
     *
     * @see #put(Object, Object)
     */
    public V get(Object key) {
        Node<K,V> e;
        return (e = getNode(hash(key), key)) == null ? null : e.value;
    }

    /**
     * Implements Map.get and related methods
     *
     * @param hash hash for key
     * @param key the key
     * @return the node, or null if none
     */
    final Node<K,V> getNode(int hash, Object key) {
        Node<K,V>[] tab; Node<K,V> first, e; int n; K k;
        // table != null 且长度不为0， 对应hash在数组中index的位置上第一个元素不为null
        if ((tab = table) != null && (n = tab.length) > 0 && (first = tab[(n - 1) & hash]) != null) {
            // always check first node 首先校验hash值是否相等
            // 再判断是否和第一个元素的key相等。
            if (first.hash == hash && ((k = first.key) == key || (key != null && key.equals(k))))
                return first;
            if ((e = first.next) != null) {
                // 如果first是一个TreeNode元素的话，那么就强转成TreeNode，然后用getTreeNode的方法。
                // 数组+红黑树的方式进行存储
                if (first instanceof TreeNode)
                    return ((TreeNode<K,V>)first).getTreeNode(hash, key);
                do {
                // 如果不是一个TreeNode元素，那么就说明first是一个链表，那么就遍历这个链表寻找key相同的元素就好了。
                // 数组+链表的方式进行存储，这个和Java 1.7很像。
                    if (e.hash == hash &&
                        ((k = e.key) == key || (key != null && key.equals(k))))
                        return e;
                } while ((e = e.next) != null);
            }
        }
        return null;
    }
```

* resize

```Java
    /**
     * Initializes or doubles table size.  If null, allocates in
     * accord with initial capacity target held in field threshold.
     * Otherwise, because we are using power-of-two expansion, the
     * elements from each bin must either stay at same index, or move
     * with a power of two offset in the new table.
     *
     * @return the table
     */
    final Node<K,V>[] resize() {
        Node<K,V>[] oldTab = table;
        int oldCap = (oldTab == null) ? 0 : oldTab.length;
        int oldThr = threshold;
        int newCap, newThr = 0;
        if (oldCap > 0) {
            // 如果原来的大小就已经超过MAXIMUM_CAPACITY了，则直接把阈值设置为Integer.MAX_VALUE了。
            if (oldCap >= MAXIMUM_CAPACITY) {
                threshold = Integer.MAX_VALUE;
                return oldTab;
            }
            // 否则newCap是老的的2倍
            else if ((newCap = oldCap << 1) < MAXIMUM_CAPACITY && oldCap >= DEFAULT_INITIAL_CAPACITY)
                newThr = oldThr << 1; // double threshold
        }
        else if (oldThr > 0) // initial capacity was placed in threshold
            newCap = oldThr;
        else {               // zero initial threshold signifies using defaults
            newCap = DEFAULT_INITIAL_CAPACITY;
            newThr = (int)(DEFAULT_LOAD_FACTOR * DEFAULT_INITIAL_CAPACITY);
        }
        // 设置新的阈值
        if (newThr == 0) {
            float ft = (float)newCap * loadFactor;
            newThr = (newCap < MAXIMUM_CAPACITY && ft < (float)MAXIMUM_CAPACITY ?
                      (int)ft : Integer.MAX_VALUE);
        }
        threshold = newThr;
        @SuppressWarnings({"rawtypes","unchecked"})
        // 创建一个新的table
        Node<K,V>[] newTab = (Node<K,V>[])new Node[newCap];
        table = newTab;
        // 将老table中的数据都迁移到新的table中。
        if (oldTab != null) {
            // 一次遍历oldtable的每个位置
            for (int j = 0; j < oldCap; ++j) {
                Node<K,V> e;
                if ((e = oldTab[j]) != null) {
                    // 获取到之后就把老表那个位置置为空
                    oldTab[j] = null;
                    // 如果e仅仅有一个元素，不是链表也不是红黑树节点，那么就直接放到新table的对应位置上。
                    if (e.next == null)
                        newTab[e.hash & (newCap - 1)] = e;
                    else if (e instanceof TreeNode)
                        // 如果e是红黑树的一个节点。那么xxxx
                        ((TreeNode<K,V>)e).split(this, newTab, j, oldCap);
                    else { // preserve order
                        // 如果e是链表的头结点，依次遍历老table中对应位置上的所有节点
                        // 保证e这个list在新表中和老表中的顺序一致。
                        Node<K,V> loHead = null, loTail = null;
                        Node<K,V> hiHead = null, hiTail = null;
                        Node<K,V> next;
                        do {
                            next = e.next;
                            if ((e.hash & oldCap) == 0) {
                                if (loTail == null)
                                    loHead = e;
                                else
                                    loTail.next = e;
                                loTail = e;
                            } else {
                                if (hiTail == null)
                                    hiHead = e;
                                else
                                    hiTail.next = e;
                                hiTail = e;
                            }
                        } while ((e = next) != null);
                        if (loTail != null) {
                            loTail.next = null;
                            newTab[j] = loHead;
                        }
                        if (hiTail != null) {
                            hiTail.next = null;
                            newTab[j + oldCap] = hiHead;
                        }
                    }
                }
            }
        }
        // 返回新的table
        return newTab;
    }
```

Java ConcurrentHashMap

ConcurrentHashMap为了提高本身的并发能力，在内部采用了一个叫做Segment的结构，一个Segment其实就是一个类Hash Table的结构，Segment内部维护了一个链表数组。

ConcurrentHashMap定位一个元素的过程需要进行两次Hash操作，第一次Hash定位到Segment，第二次Hash定位到元素所在的链表的头部，因此，这一种结构的带来的副作用是Hash的过程要比普通的HashMap要长，
但是带来的好处是写操作的时候可以只对元素所在的Segment进行加锁即可，不会影响到其他的Segment，这样，在最理想的情况下，ConcurrentHashMap可以最高同时支持Segment数量大小的写操作
（刚好这些写操作都非常平均地分布在所有的Segment上），所以，通过这一种结构，ConcurrentHashMap的并发能力可以大大的提高。

结合在实际中的应用，感觉concurrentHashMap更加适用于redis。所以打算采用这个对hashMap进行实现。


```Java
    static final class HashEntry<K,V> {
        final K key;
        final int hash;
        volatile V value;
        final HashEntry<K,V> next;
    }
```


上述代码中，除了value其他所有的变量都是final的。这意味着不能从hash链的中间或尾部添加或删除节点，因为这需要修改next引用值，所有的节点的修改只能从头部开始。

对于put操作，可以一律添加到Hash链的头部。但是对于remove操作，可能需要从中间删除一个节点，这就需要将要删除节点的前面所有节点整个复制一遍，最后一个节点指向要删除结点的下一个结点。

为了确保读操作能够看到最新的值，将value设置成volatile，这避免了加锁。 remove操作要注意一个问题：如果某个读操作在删除时已经定位到了旧的链表上，那么此操作仍然将能读到数据，只不过读取到的是旧数据而已，这在多线程里面是没有问题的。

HashEntry 类的 value 域被声明为 Volatile 型，Java 的内存模型可以保证：某个写线程对 value 域的写入马上可以被后续的某个读线程“看”到。

#### reference
* https://www.cnblogs.com/slwenyi/p/6393829.html
* https://www.ibm.com/developerworks/cn/java/java-lo-concurrenthashmap/















