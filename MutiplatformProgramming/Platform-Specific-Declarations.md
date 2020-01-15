## 平台相关声明

**跨平台项目是 Kotlin 1.2 和 1.3 中的实验性特性。本文档中描述的所有语言和工具功能都可能在将来的Kotlin版本中发生变更**

Kotlin的多平台代码的主要功能之一是使通用代码依赖于特定于平台的声明的方式。 在其他语言中，通常可以通过在通用代码中构建一组接口并在特定于平台的模块中实现这些接口来实现。 但是，如果在其中一个平台上所需功已经有一个库，并且希望直接使用该库的API而无需额外的包装器，则这种方法并不适合。 另外，它要求将公共声明表示为接口，但并不能涵盖所有可能的情况。

作为替代，Kotlin提供了预期和实际声明的机制。 通过这种机制，公共模块可以定义期望的声明，而平台模块可以提供与期望的声明相对应的实际声明。 为了了解它是如何工作的，我们首先来看一个示例。 此代码是通用模块的一部分：


```Kotlin
package org.jetbrains.foo

expect class Foo(bar: String) {
    fun frob()
}

fun main() {
    Foo("Hello").frob()
}
```

这是对应 JVM 的模块：

```kotlin
package org.jetbrains.foo

actual class Foo actual constructor(val bar: String) {
    actual fun frob() {
        println("Frobbing the $bar")
    }
}
```

这里演示了一个重要的点：

- 通用模块中的预期声明及其实际对应项始终具有完全相同的限定名称。

- 预期的声明用 expect 关键字标记； 实际的声明中用 actual 关键字标记。

- 与预期声明的任何部分匹配的所有实际声明都需要标记为 actual 。

- 预期的声明绝不包含任何实现代码。

请注意，预期的声明不限于接口和接口成员。 在此示例中，期望的类具有构造函数，可以直接从通用代码创建。 还可以将 expect 修饰符应用于其他声明，包括顶级声明和注解：

```Kotlin

// Common
expect fun formatString(source: String, vararg args: Any): String

expect annotation class Test

// JVM
actual fun formatString(source: String, vararg args: Any) =
    String.format(source, *args)
    
actual typealias Test = org.junit.Test

```

编译器确保每个期望的声明在实现相应公共模块的所有平台模块中都有实际的声明，并在缺少任何实际的声明时报告错误。 IDE提供了可帮助创建缺少实际声明的工具。

如果您要在通用代码中使用特定平台的库，同时为另一个平台提供自己的实现，则可以为现有类提供类型别名作为实际的声明：

```Kotlin
expect class AtomicRef<V>(value: V) {
  fun get(): V
  fun set(value: V)
  fun getAndSet(value: V): V
  fun compareAndSet(expect: V, update: V): Boolean
}

actual typealias AtomicRef<V> = java.util.concurrent.atomic.AtomicReference<V>
```