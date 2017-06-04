## 泛型
像 java 一样，Kotlin 中的类可以拥有类型参数：

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
因此 java 禁止了这样的事情来保证运行时安全。但这有些其它影响。比如，`Collection` 接口的 `addAll()` 方法。这个方法的签名在哪呢？直觉告诉我们应该是这样的：

```java
//java
interface Collection<E> ... {
	void addAdd(Collection<E> items);
}
```
但接下来我们就不能做下面这些操作了(虽然这些操作都是安全的)：

```java
// Java
void copyAll(Collection<Object> to, Collection<String> from) {
  to.addAll(from); // !!! Would not compile with the naive declaration of addAll:
                   //       Collection<String> is not a subtype of Collection<Object>
}
```

这就是为什么 `addAll()` 的签名是下面这样的：

```java
//java
interface Collection<E> ... {
	void addAll(Colletion<? extend E> items);
}
```

这个通配符参数 `? extends T` 意味着这个方法接受一些 T 类型的子类而非 T 类型本身。这就是说我们可以安全的读 `T's`(这里表示 T 子类元素的集合)，但不能写，因为我们不知道 T 的子类究竟是什么样的，针对这样的限制，我们很想要这样的行为：`Collection<String>` 是 `Collection<? extens Object>`的子类。换句话讲，带 **extends** 限定（**上界**）的通配符类型使得类型是**协变的（covariant）**。

这个技巧其实很简单：如果你只能从集合中读数据，那么使用`String` 集合并从中读取 `Objects` 是安全的，如果你只能给 `存入` 集合 ，那么给 `Objects` 集合存入 `String` 也是可以的：在 Java 中`List<? super String>` 是 `List<Object>`的超类。

后者称为**逆变性（contravariance）**，并且对于 `List <? super String>` 你只能调用接受 String 作为参数的方法 （例如，你可以调用 `add(String)` 或者 `set(int, String)`），当然 如果调用函数返回 `List<T>` 中的 `T`，你得到的并非一个 `String` 而是一个 `Object`。

Joshua Bloch 称只能**读取**的对象为**生产者**，只能**写入**的对象为**消费者**。他建议：“*为了灵活性最大化，在表示生产者或消费者的输入参数上使用通配符类型*”，并提出了以下助记符：

*PECS 代表生产者-Extens，消费者-Super（Producer-Extends, Consumer-Super）。*

*注意*：如果你使用一个生产者对象，如 `List<? extends Foo>`，在该对象上不允许调用 `add()` 或 `set()`。但这并不意味着 该对象是**不可变的**：例如，没有什么阻止你调用 `clear()`从列表中删除所有项目，因为 `clear()` 根本无需任何参数。通配符（或其他类型的型变）保证的唯一的事情是**类型安全**。不可变性完全是另一回事。