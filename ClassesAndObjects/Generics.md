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

### 变型
java 类型系统最棘手的一部分就是通配符类型。但 kotlin 没有，代替它的是两种其它的东西：声明变型和类型投影(declaration-site variance and type projections)。

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
	void addAll(Collection<E> items);
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



### 声明处变型

假如有个范型接口`Source<T>`，没有任何接收 `T` 作为参数的方法，唯一的方法就是返回  `T`:

```Kotlin 
// Java
interface Source<T> {
  T nextT();
}
```

存储一个`Source<String>`的实例引用给一个类型为 `Source<Object>` 是十分安全的。但 Java并不知道，而且依然禁止这么做：

```Kotlin 
// Java
void demo(Source<String> strs) {
  Source<Object> objects = strs; // !!! Not allowed in Java
  // ...
}
```

为次，我们不得不声明对象类型为 `Source<? extends Object>`，这样做并没有太大的意义，因为我们可以像以前一样调用所有方法，因此并没有通过复杂的类型添加什么值。但编译器不知道。

在 Kotlin 中，有种可以将这些东西解释给编译器的办法，叫做声明处变型：通过注解**类型参数** `T` 的来源，来确保它仅从 `Source<T>` 成员中**返回**（生产），并从不被消费。 为此，我们提供 **out** 修饰符：

```Kotlin 
abstract class Source<out T> {
    abstract fun nextT(): T
}

fun demo(strs: Source<String>) {
    val objects: Source<Any> = strs // This is OK, since T is an out-parameter
    // ...
}
```

一般原则是：当一个类 `C` 的类型参数 `T` 被声明为 **out** 时，它就只能出现在 `C` 的成员的**输出**-位置，结果是 `C<Base>` 可以安全地作为 `C<Derived>`的超类。

更聪明的说法就是，当类 C 在类型参数 T 之下是协变的，或者 T 是一个协变类型。可以把 C 想象成 T 的生产者，而不是 T 的消费者。

`out` 修饰符本来被称之为变型注解，但由于同处与类型参数声明处，我们称之为声明处变型。这与 Java 中的使用处变型相反。

另外除了 **out**，Kotlin 又补充了一个变型注释：**in**。它接受一个类型参数**逆变**：只可以被消费而不可以 被生产。非变型类的一个很好的例子是 `Comparable`：

```Kotlin 
abstract class Comparable<in T> {
    abstract fun compareTo(other: T): Int
}

fun demo(x: Comparable<Number>) {
    x.compareTo(1.0) // 1.0 has type Double, which is a subtype of Number
    // Thus, we can assign x to a variable of type Comparable<Double>
    val y: Comparable<Double> = x // OK!
}
```

我们相信 **in** 和 **out** 两词是自解释的（因为它们已经在 C# 中成功使用很长时间了）， 因此上面提到的助记符不是真正需要的，并且可以将其改写为更高的目标：

**[存在性（The Existential）](https://en.wikipedia.org/wiki/Existentialism) 转变：消费者 in, 生产者 out!** :-)

### 类型投影

#### 使用处变型：类型投影

声明类型参数 T 为 *out* 很方便，而且可以避免在使用出子类型的麻烦，但有些类 **不能** 限制它只返回 `T` ，Array 就是一个例子：

```kotlin
class Array<T>(val size: Int) {
    fun get(index: Int): T { /* ... */ }
    fun set(index: Int, value: T) { /* ... */ }
}
```

这个类既不能是协变的也不能是逆变的，这会在一定程度上降低灵活性。考虑下面的函数：

```kotlin
fun copy(from: Array<Any>, to: Array<Any>) {
    assert(from.size == to.size)
    for (i in from.indices)
        to[i] = from[i]
}
```

该函数作用是复制 array ，让我们来实际应用一下：

```kotlin
val ints: Array<Int> = arrayOf(1, 2, 3)
val any = Array<Any>(3) { "" } 
copy(ints, any) // Error: expects (Array<Any>, Array<Any>)
```

这里我们又遇到了同样的问题 `Array<T>` 中的`T` 是不可变型的，因此 `Array<Int>` 和 `Array<Any>` 互不为对方的子类，导致复制失败。为什么呢？应为复制可能会有不合适的操作，比如尝试写入，当我们尝试将 Int 写入 String 类型的 array 时候将会导致 `ClassCastException` 异常。

我们想做的就是确保 `copy()` 不会做类似的不合适的操作，为阻止向`from`写入，我们可以这样：

```kotlin 
fun copy(from: Array<out Any>, to: Array<Any>) {
 // ...
}
```

这就是类型投影：这里的`from`不是一个简单的 array， 而是一个投影，我们只能调用那些返回类型参数 `T` 的方法，在这里意味着我们只能调用`get()`。这是我们处理调用处变型的方法，类似 Java 中`Array<? extends Object>`，但更简单。

当然也可以用`in`做投影：

```kotin
fun fill(dest: Array<in String>, value: String) {
    // ...
}
```

`Array<in String>` 对应 Java 中的 `Array<? super String>`，`fill()`函数可以接受任何`CharSequence` 类型或 `Object`类型的 array 。

#### 星投影

有时你对类型参数一无所知，但任然想安全的使用它。保险的方法就是定一个该范型的投影，每个该范型的正确实例都将是该投影的子类。

Kotlin 提供了一种星投影语法：

- For `Foo<out T>`, where `T` is a covariant type parameter with the upper bound `TUpper`, `Foo<*>` is equivalent to `Foo<out TUpper>`. It means that when the `T` is unknown you can safely *read* values of `TUpper` from `Foo<*>`.
- For `Foo<in T>`, where `T` is a contravariant type parameter, `Foo<*>` is equivalent to `Foo<in Nothing>`. It means there is nothing you can *write* to `Foo<*>` in a safe way when `T` is unknown.
- For `Foo<T>`, where `T` is an invariant type parameter with the upper bound `TUpper`, `Foo<*>` is equivalent to `Foo<out TUpper>` for reading values and to `Foo<in Nothing>` for writing values.

If a generic type has several type parameters each of them can be projected independently. For example, if the type is declared as `interface Function<in T, out U>` we can imagine the following star-projections:

- `Function<*, String>` means `Function<in Nothing, String>`;
- `Function<Int, *>` means `Function<Int, out Any?>`;
- `Function<*, *>` means `Function<in Nothing, out Any?>`.

*Note*: star-projections are very much like Java's raw types, but safe.   (这部分暂未翻译)

### 范型函数

函数也可以像类一样有类型参数。类型参数在函数名之前：

```kotlin
fun <T> singletonList(item: T): List<T> {
    // ...
}

fun <T> T.basicToString() : String {  // extension function
    // ...
}
```

调用范型函数需要在函数名后面制定类型参数：

```kotlin
val l = singletonList<Int>(1)
```

### 范型约束

指定类型参数代替的类型集合可以用通过范型约束进行限制。

#### 上界(**upper bound**)

最常用的类型约束是上界，在 Java 中对应 `extends`关键字：

```kotlin
fun <T : Comparable<T>> sort(list: List<T>) {
    // ...
}	
```

冒号后面指定的类型就是上界：只有 `Comparable<T>`的子类型才可以取代 `T` 比如：

```kotlin
sort(listOf(1, 2, 3)) // OK. Int is a subtype of Comparable<Int>
sort(listOf(HashMap<Int, String>())) // Error: HashMap<Int, String> is not a subtype of Comparable<HashMap<Int, String>>
```

默认的上界是 `Any?`。在尖括号内只能指定一个上界。如果要指定多种上界，需要用 **where** 语句指定：

```kotlin
fun <T> cloneWhenGreater(list: List<T>, threshold: T): List<T>
    where T : Comparable,
          T : Cloneable {
  return list.filter { it > threshold }.map { it.clone() }
}
```

