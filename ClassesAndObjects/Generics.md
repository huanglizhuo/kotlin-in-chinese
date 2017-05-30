## 泛型
像 java 一样，Kotlin 中可以拥有类型参数：

```kotlin
class Box<T>(t: T){
    var value = t
}
```

通常来说，创建一个这样类的实例，我们需要提供类型参数：

```kotlin
val box: Box<Int> = Box<Int>(1)
```

但如果类型有可能是推断的，比如来自构造函数的参数或者通过其它的一些方式，一个可以忽略类型的参数：

```kotin
val box = Box(1)//1是 Int 型，因此编译器会推导出我们调用的是 Box<Int>
```

### 变化
java 类型系统最棘手的一部分就是通配符类型。但 kotlin 没有，代替它的是两种其它的东西：声明变化和类型投影(declaration-site variance and type projections)。

首先，我们想想为什么 java 需要这些神秘的通配符。这个问题在[Effective Java](http://www.oracle.com/technetwork/java/effectivejava-136174.html),条目18中是这样解释的：使用界限通配符增加 API 的灵活性。首先 java 中的泛型是不变的，这就意味着 `List<String>` 不是 `List<Object>` 的子类型。为什么呢，如果 List 不是不变的，就会引发下面的问题：

```java
// Java
List<String> strs = new ArrayList<String>();
List<Object> objs = strs; // !!! The cause of the upcoming problem sits here. Java prohibits this!
objs.add(1); // Here we put an Integer into a list of Strings
String s = strs.get(0); // !!! ClassCastException: Cannot cast Integer to String
```
因此 java 禁止了这样的事情来保证运行时安全。但这有些其它影响。比如，`Collection` 接口的 `addAll()` 方法。这个方法的签名在哪呢？直观来讲是这样的：

```java
//java
interface Collection<E> ... {
	void addAdd(Collection<E> items);
}
```
但接下来我们就不能做下面这些简单事情了：

```java
//java
void copyAll(Collection<Object> to, Collection<String> from){
	to.addAll(from);
}
```

这就是为什么 `addAll()` 的签名是下面这样的：

```java
//java
interface Collection<E> ... {
	void addAll(Colletion<? extend E> items);
}
```

这个通配符参数 `? extends T` 意味着这个方法接受一些 T 类型的子类而非 T 类型本身。这就是说我们可以安全的读 `T's`(这里表示 T 子类元素的集合)，但不能写，因为我们不知道 T 的子类究竟是什么样的，这个限制的补偿，是我们很想要的行为：`Collection<String>` 是 `Collection<? extends Object>`的子类。用更高级的话说，使用 extends 界限（上界）的通配符使类型成为协变类型。

要理解这为什么可行是非常简单的：如果你只能从一个集合中获取元素，那么从一个 String 集合中读取 Object 是可以的。相反，如果你只能把元素放入集合中，往 Object 集合中放入 String 也是可以的：在 Java 中 List<? super String> 是 List<Object> 的父类。

