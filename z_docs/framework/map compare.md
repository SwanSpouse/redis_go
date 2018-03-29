
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
                e.next = newTable[i];
                newTable[i] = e;
                e = next;
            }
        }
    }
```