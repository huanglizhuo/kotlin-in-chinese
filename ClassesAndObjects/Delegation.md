## 委托

### 委托属性

委托属性在单独页面描述：[委托属性](DelegationProperties.md) 

### 委托实现

[委托模式](https://zh.wikipedia.org/wiki/%E5%A7%94%E6%89%98%E6%A8%A1%E5%BC%8F)已被证实是继承的一个很好的替代方案，而且 kotlin 原生支持该模式并不需要任何模板代码。`Derived` 类可以通过将其所有公有成员都委托给指定对象来实现一个接口 `Base`：

```kotlin
interface Base {
    fun print()
}

class BaseImpl(val x: Int) : Base {
    override fun print() { print(x) }
}

class Derived(b: Base) : Base by b

fun main() {
    val b = BaseImpl(10)
    Derived(b).print()
}
```

`Derived` 的超类列表中的 *by*子句表示 `b` 将会存储在 `Derived` 内部，并且编译器将生成所有委托给 `b` 的 `Base` 的方法。

### 覆写由委托实现的接口成员


覆写将和你预期的一样工作：编译器会使用你 `override` 后的实现而不是委托对象中的实现。如果将 `override fun printMessage() { print("abc") }` 添加给 `Derived`，那么当调用 `printMessage` 时程序会输出“abc”而不是“10”：

```kotlin
interface Base {
    fun printMessage()
    fun printMessageLine()
}

class BaseImpl(val x: Int) : Base {
    override fun printMessage() { print(x) }
    override fun printMessageLine() { println(x) }
}

class Derived(b: Base) : Base by b {
    override fun printMessage() { print("abc") }
}

fun main() {
    val b = BaseImpl(10)
    Derived(b).printMessage()
    Derived(b).printMessageLine()
}
```

但请注意，以这种方式覆写的成员不会从委托对象的成员调用，该成员只能访问其自己的接口成员实现：

```kotlin
interface Base {
    val message: String
    fun print()
}

class BaseImpl(val x: Int) : Base {
    override val message = "BaseImpl: x = $x"
    override fun print() { println(message) }
}

class Derived(b: Base) : Base by b {
    // 在 b 的 `print` 实现中不会访问到这个属性
    override val message = "Message of Derived"
}

fun main() {
    val b = BaseImpl(10)
    val derived = Derived(b)
    derived.print()
    println(derived.message)
}
```