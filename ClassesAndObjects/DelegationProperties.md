## 代理属性
很多常用属性，虽然我们可以在需要的时候手动实现它们，但更好的办法是一次实现多次使用，并放到库。比如：

> 延迟属性：只在第一次访问是计算它的值
>观察属性：监听者从这获取这个属性更新的通知
>在 map 中存储的属性，而不是单独存在分开的字段

为了满足这些情形，Kotllin 支持代理属性：

```kotlin
class Example {
	var p: String by Delegate()
}
```

语法结构是： `val/var <property name>: <Type> by <expression>` 在 by 后面的属性就是代理，这样这个属性的 get() 和 set() 方法就代理给了它。

属性代理不需要任何接口的实现，但必须要提供 `get()` 方法(如果是变量还需要 `set()` 方法)。像这样：

```kotlin
class Delegate {
	fun get(thisRef: Any?, prop: PropertyMetadata): String {
		return "$thisRef, thank you for delegating '${prop.name}' to me !"
	}
	
	fun set(thisRef: Any?, prop: PropertyMatada, value: String) {
		println("$value has been assigned to '${prop.name} in $thisRef.'")
	}
}
```

当我们从 `p` 也就是 `Delegate` 的代理，中读东西时，会调用 `Delegate` 的 `get()` 函数，因此第一个参数是我们从 `p` 中读取的，第二个参数是 `p` 自己的一个描述。比如：

```kotlin
val e = Example()
println(e.p)
```

打印结果：　

>Example@33a17727, thank you for delegating ‘p’ to me!

同样当我们分配 `p` 时 `set()` 函数就会调动。前俩个参数所以一样的，第三个持有分配的值：

```kotlin
e.p = "NEW"
```

打印结果：　

>NEW has been assigned to ‘p’ in Example@33a17727.

### 代理属性的要求
这里总结一些代理对象的要求。

只读属性 (val)，代理必须提供一个名字叫 `get` 的方法并接受如下参数：

> 接收者--必须是相同的，或者是属性拥有者的子类型

> 元数据--必须是 `PropertyMetadata` 或这它的子类型

这个函数必须返回同样的类型作为属性。

可变属性 (var)，代理必须添加一个叫 `set` 的函数并接受如下参数：

> 接受者--与 `get()` 一样
> 元数据--与 `get()` 一样
>新值--必须和属性类型一致或是它的字类型

### 标准代理
`kotlin.properties.Delegates` 对象是标准库提供的一个工厂方法并提供了很多有用的代理

#### 延迟
`Delegate.lazy()` 是一个接受 lamdba 并返回一个实现延迟属性的代理：第一次调用 `get()` 执行 lamdba 并传递 `lazy()` 并记下结果，随后调用 `get()` 并简单返回之前记下的值。

```kotlin
import kotlin.properties.Delegates

val lazy: String by Delegates.lazy {
    println("computed!")
    "Hello"
}

fun main(args: Array<String>) {
    println(lazy)
    println(lazy)
}
```

如果你想要线程安全，使用 `blockingLazy()`: 它还是按照同样的方式工作，但保证了它的值只会在一个线程中计算，并且所有的线程都获取的同一个值。

#### 观察者
`Delegates.observable()` 需要俩个参数：一个初始值和一个修改者的 handler 。每次我们分配属性时都会调用handler (在分配前执行)。它有三个参数：一个分配的属性，旧值，新值：

```kotlin
class User {
	var name: String by Delegates.observable("<no name>") {
		d.old,new -> println("$old -> $new")
	}
}

fun main(args: Array<String>) {
	val user = User()
	user.name = "first"
	user.name = "second"
}
```
打印结果

><no name> -> first
first -> second

如果你想能够截取它的分配并取消它，用 `vetoable()`代替  `observable()`

#### 非空
有时我们有一个非空的 var ，但我们在构造函数中没有一个合适的值，比如它必须稍后再分配。问题是你不能持有一个未初始化并且是非抽象的属性：

```kotlin
class Foo {
	var bar: Bat //错误必须初始化
}
```

我们可以用 null 初始化它，但我们不用每次访问时都检查它。

`Delegates.notNull()`可以解决这个问题

```kotlin
class Foo {
	var bar: Bar by Delegates.notNull()
}
```

如果这个属性在第一次写之前读，它就会抛出一个异常，只有分配之后才会正常。

#### 在 Map 中存储属性
`Delegates.mapVal()` 拥有一个 map 实例并返回一个可以从 map 中读其中属性的代理。在应用中有很多这样的例子，比如解析 JSON 或者做其它的一些 "动态"的事情：

```kotlin
class User(val map: Map<String, Any?>) {
	val name: String by Delegates.mapVal(map)
	val age: Int     by Delegates.mapVal(map)
}
```

在这个例子中，构造函数持有一个 map :

```kotlin
val user = User(mapOf (
	"name" to "John Doe",
	"age" to 25
))
```

代理从这个 map 中取指(通过属性的名字)：

```kotlin
println(user.name) // Prints "John Doe"
println(user.age)  // Prints 25
```

var 可以用 `mapVar` 
