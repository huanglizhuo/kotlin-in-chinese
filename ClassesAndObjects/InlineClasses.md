## 内联类

> 内联类在 kotlin1.3 开始支持，并且目前还是*实验性*的。参看下文

有时候会为了某些业务逻辑而对某些类型进行包装。然而由于额外的堆分配操作，这会给运行时带来性能损耗。除此之外，如果包装类型是原始类型，性能损耗尤为可怕，因为原始类型通常会经由运行时进行大幅度优化，然而这些包装器并不会享受任何特殊待遇。

为了解决这些问题，Kotlin引入了一种称为内联类的特殊类，它通过在类的名称前面放置一个 inline 关键字来声明：

```kotlin
inline class Password(val value: String)
```

内联类必须在主构造函数中初始化唯一属性。在运行时，将使用此单个属性表示内联类的实例（请参阅下面有关运行时表示的详细信息）：

```kotlin
// 并没有 Password 实例
// 在运行时，'securePassword' 只包含'String'
val securePassword = Password("Don't try this in production")
```

这是内联类的主要特性。之所以有命名为内联，是因为类的数据内联在它使用处（与内联函数概念类似）

###成员

内联类也支持一些普通类的功能。尤其是可以声明属性和函数：

```kotlin
inline class Name(val s: String) {
    val length: Int
        get() = s.length

    fun greet() {
        println("Hello, $s")
    }
}    

fun main() {
    val name = Name("Kotlin")
    name.greet() // method `greet` is called as a static method
    println(name.length) // property getter is called as a static method
}
```

当然内联类成员也有些其它限制：

* 内联类不能有 init 块
* 内联类不能有 inner 类
* 内联类属性不能有备用字段(backing fields)
    * 内联类只允许有简单的可计算属性（不能含有延迟初始化/代理属性）

### 继承

内联类允许继承接口：

```kotlin
interface Printable {
    fun prettyPrint(): String
}

inline class Name(val s: String) : Printable {
    override fun prettyPrint(): String = "Let's $s!"
}    

fun main() {
    val name = Name("Kotlin")
    println(name.prettyPrint()) // Still called as a static method
}
```

禁止内联类进行类关系继承。这意味着内联函数不能继承其它类而且必须是 final

### 表示

在生成的代码中，kotlin 编译器会为每个内联类保留一个包装器。内联类实例在运行时既可以表示为包装器也可以表示为基础类型。这和 `Int` 既可以表示为基础类型`int`也可以表示为包装器 `Integer`。

kotlin 编译器会更加倾向于使用基础类型而不是包装器，这样可以提高性能并优化代码。然而有时保留包装器也是很有必要的。一般来说，只要将内联类用作另一种类型，它们就会被装箱。

```kotlin
interface I

inline class Foo(val i: Int) : I

fun asInline(f: Foo) {}
fun <T> asGeneric(x: T) {}
fun asInterface(i: I) {}
fun asNullable(i: Foo?) {}

fun <T> id(x: T): T = x

fun main() {
    val f = Foo(42) 
    
    asInline(f)    // unboxed: used as Foo itself
    asGeneric(f)   // boxed: used as generic type T
    asInterface(f) // boxed: used as type I
    asNullable(f)  // boxed: used as Foo?, which is different from Foo
    
    // below, 'f' first is boxed (while being passed to 'id') and then unboxed (when returned from 'id') 
    // In the end, 'c' contains unboxed representation (just '42'), as 'f' 
    val c = id(f)  
}
```

因为内联类既可以表示为基础类型有可以表示为包装器，引用相等对于内联类而言毫无意义，因而也禁止此项操作。

### 类名重排

由于内联类被编译为其基础类型，因此可能会带来各种模糊的错误，例如意想不到的平台签名冲突：

```kotlin
inline class UInt(val x: Int)

// Represented as 'public final void compute(int x)' on the JVM
fun compute(x: Int) { }

// Also represented as 'public final void compute(int x)' on the JVM!
fun compute(x: UInt) { }
```

为了缓解这种问题，一般会通过在函数名拼接一段哈希值重命名函数。 `fun compute(x: UInt)` 将会被表示为 `public final void compute-<hashcode>(int x)`，以此来解决冲突的问题。

> 请注意在 Java 中 `-` 是一个 *无效的* 符号，也就是说在 Java 中不能调用使用内联类作为形参的函数。

### 内联类 vs 类型别名

乍一看，内联类似乎与类型别名非常相似，两者似乎都引入了一种新的类型，并且都在运行时表示为基础类型。

然而，关键的区别在于类型别名与其基础类型(以及具有相同基础类型的其他类型别名)是 *赋值兼容* 的，而内联类却不是这样。

换句话说，内联类真正引入了新类型，而类型别名仅仅是为现有的类型取了个新的替代名称(别名)：

```kotlin
typealias NameTypeAlias = String
inline class NameInlineClass(val s: String)

fun acceptString(s: String) {}
fun acceptNameTypeAlias(n: NameTypeAlias) {}
fun acceptNameInlineClass(p: NameInlineClass) {}

fun main() {
    val nameAlias: NameTypeAlias = ""
    val nameInlineClass: NameInlineClass = NameInlineClass("")
    val string: String = ""

    acceptString(nameAlias) // 正确: 传递别名类型的实参替代函数中基础类型的形参
    acceptString(nameInlineClass) // 错误: 不能传递内联类的实参替代函数中基础类型的形参

    // And vice versa:
    acceptNameTypeAlias("") // 正确: 传递基础类型的实参替代函数中别名类型的形参
    acceptNameInlineClass("") // 错误: 不能传递基础类型的实参替代函数中内联类类型的形参
}
```

</div>


### 内联类的实验性状态

内联类的设计目前是实验性的，这就是说此特性是正在 *快速变化*的，并且不保证其兼容性。在 Kotlin 1.3+ 中使用内联类时，将会得到一个警告，来表明此特性还是实验性的。

要想移除警告，可以对 `kotlinc` 指定 `-XXLanguage:+InlineClasses`参数来选择使用该实验性的特性。

### 在 Gradle 中启用内联类：

``` groovy
compileKotlin {
    kotlinOptions.freeCompilerArgs += ["-XXLanguage:+InlineClasses"]
}
```

### 在 Maven 中启用内联类

```xml
<configuration>
    <args>
        <arg>-XXLanguage:+InlineClasses</arg> 
    </args>
</configuration>
```

## 进一步讨论

关于其他技术详细信息和讨论，请参见[内联类的语言提议](https://github.com/Kotlin/KEEP/blob/master/proposals/inline-classes.md)