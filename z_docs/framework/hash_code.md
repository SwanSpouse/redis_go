
## HashCode

### HashCode作用

数据结构Set可以保证其中的元素是无序且不重复的。如何保证其中元素不重复，有两种方法。

* 在插入新元素的时候，依次和Set中元素调用equals方法进行比较。不过这种方法在数据量较小的情况。

* 在插入新元素的时候，计算新元素的HashCode值。通过HashCode的值寻找其对应的位置，如果对应的位置有元素，则再次调用equals方法来进行判断。

简言之:HashCode的作用就是用来`寻址`

* Java中的hashCode方法就是根据一定的规则将与对象相关的信息（比如对象的存储地址，对象的 字段等）映射成一个数值，这个数值称作为散列值

### HashCode设计中的矛盾：

一个对象势必会存在若干个属性，

* 如果我们将所有属性进行散列，这必定会是一个糟糕的设计，因为对象的hashCode方法无时无刻不是在被调用，如果太多的属性参与散列，那么需要的操作数时间将会大大增加，这将严重影响程序的性能。

* 如果较少属相参与散列，散列的多样性会削弱，会产生大量的散列"冲突"，除了不能够很好的利用空间外，在某种程度也会影响对象的查询效率。

### HashCode 和 equals: 

equals所需要遵循的规则:

* `对称性：` 如果x.equals(y)返回是"true"，那么y.equals(x)也应该返回是"true"。

* `反射性：` x.equals(x)必须返回是"true"。
  
* `类推性：` 如果x.equals(y)返回是"true"，而且y.equals(z)返回是"true"，那么z.equals(x)也应该返回是"true"。
  
* `一致性：` 如果x.equals(y)返回是"true"，只要x和y内容一直不变，不管你重复x.equals(y)多少次，返回都是"true"。

* 任何情况下，x.equals(null)，永远返回是"false"；x.equals(和x不同类型的对象)永远返回是"false"。

hashCode 所需要遵循的规则:

* 在一个应用程序执行期间，如果一个对象的equals方法做比较所用到的信息没有被修改的话，则对该对象调用hashCode方法多次，它必须始终如一地返回同一个整数。

* 如果两个对象根据equals(Object o)方法是相等的，则调用这两个对象中任一对象的hashCode方法必须产生相同的整数结果。

* 如果两个对象根据equals(Object o)方法是不相等的，则调用这两个对象中任一个对象的hashCode方法，不要求产生不同的整数结果。但如果能不同，则可能提高散列表的性能。

equals和hashCode联系:

* 如果x.equals(y)返回“true”，那么x和y的hashCode()必须相等。

* 如果x.equals(y)返回“false”，那么x和y的hashCode()有可能相等，也有可能不等。

### HashCode的实现:

Java String类的HashCode实现:

```java

public int hashCode() {  
    int h = hash;  
    if (h == 0) {  
        int off = offset;  
        char val[] = value;  
        int len = count;  
  
            for (int i = 0; i < len; i++) {  
                h = 31*h + val[off++];  
            }  
            hash = h;  
        }  
        return h;  
    }
} 
```
* 上述方法可以总结为: s[0]*31^(n-1) + s[1]*31^(n-2) + ... + s[n-1]  

* 为什么这里用31，而不是其它数呢?《Effective Java》是这样说的：之所以选择31，是因为它是个奇素数，如果乘数是偶数，并且乘法溢出的话，信息就会丢失，因为与2相乘等价于移位运算。使用素数的好处并不是很明显，但是习惯上都使用素数来计算散列结果。31有个很好的特性，就是用移位和减法来代替乘法，可以得到更好的性能：31*i==(i<<5)-i。现在的JVM可以自动完成这种优化。

Java Object类的HashCode实现：

* 这个相对来说比较复杂。以后有兴趣可以看看: https://www.zhihu.com/question/29976202


### HashCode In redis_go

采用三方库来实现: 

* 项目: https://github.com/mitchellh/hashstructure
* 项目文档:  https://godoc.org/github.com/mitchellh/hashstructure

Code Sample: 

```java

type ComplexStruct struct {
    Name     string
    Age      uint
    Metadata map[string]interface{}
}

v := ComplexStruct{
    Name: "mitchellh",
    Age:  64,
    Metadata: map[string]interface{}{
        "car":      true,
        "location": "California",
        "siblings": []string{"Bob", "John"},
    },
}

hash, err := Hash(v, nil)
if err != nil {
    panic(err)
}

fmt.Printf("%d", hash)

//Output: 6691276962590150517
 
```

### reference

* http://www.importnew.com/20381.html
* https://www.cnblogs.com/lchzls/p/6714146.html
* https://www.zhihu.com/question/29976202
* https://github.com/mitchellh/hashstructure
* https://godoc.org/github.com/mitchellh/hashstructure