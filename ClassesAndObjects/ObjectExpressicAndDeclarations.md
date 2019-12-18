## 对象表达式和对象声明
有时候我们想要创建一个对某个类有一点小修改的对象，而不是声明一个新子类。Kotlin 中用*对象表达式*和*对象声明*来实现。

### 对象表达式
以下方式可以创建继承自某种(或某些)匿名类的对象：

```kotlin
window.addMouseListener(object : MouseAdapter() {
    override fun mouseClicked(e: MouseEvent) { /*...*/ }

    override fun mouseEntered(e: MouseEvent) { /*...*/ }
})
```

如果父类有构造函数，则必须传递相应的构造参数。多个父类可以用逗号隔开，跟在冒号后面：

```kotlin
open class A(x: Int) {
	public open val y: Int = x
}

interface B { ... }

val ab = object : A(1), B {
	override val y = 14
}
```

有时候我们只是需要一个的对象，没有任何父类,我们可以这样写：

```kotlin
fun foo() {
    val adHoc = object {
        var x: Int = 0
        var y: Int = 0
    }
    print(adHoc.x + adHoc.y)
}
```

这里需要注意匿名对象作为类型只能出现在本地或者私有声明中. 如果把匿名对象作为公有函数返回类型或者公有属性时, 真正的类型将会是匿名函数的超类,如果声明没有超类则会使 `Any` .作为成员变量的匿名对象是不可访问的.


```kotlin
class C {
    // Private function, so the return type is the anonymous object type
    private fun foo() = object {
        val x: String = "x"
    }

    // Public function, so the return type is Any
    fun publicFoo() = object {
        val x: String = "x"
    }

    fun bar() {
        val x1 = foo().x        // Works
        val x2 = publicFoo().x  // ERROR: Unresolved reference 'x'
    }
}
```

在对象表达式中可以访问来自包含它的作用域的变量.

```kotlinl 
fun countClicks(window: JComponent) {
    var clickCount = 0
    var enterCount = 0

    window.addMouseListener(object : MouseAdapter() {
        override fun mouseClicked(e: MouseEvent) {
            clickCount++
        }

        override fun mouseEntered(e: MouseEvent) {
            enterCount++
        }
    })
    // ...
}
```

### 对象声明
[单例模式](http://en.wikipedia.org/wiki/Singleton_pattern)在很多情形中很实用，Kotln(在 Scala 之后)大大简化了声明方式：

```kotlin
object DataProviderManager {
    fun registerDataProvider(provider: DataProvider) {
        // ...
    }

    val allDataProviders: Collection<DataProvider>
        get() = // ...
}
```

这叫做对象声明，跟在 object 关键字后面是对象名。和变量声明一样，对象声明并不是表达式，而且不能作为右值用在赋值语句。

对象声明的初始化

想要访问这个类，直接通过名字来使用这个类：

```kotlin
DataProviderManager.registerDataProvider(...)
```

这样类型的对象可以有父类型：

```kotlin
object DefaultListener : MouseAdapter() {
    override fun mouseClicked(e: MouseEvent) {
        // ...
    }

    override fun mouseEntered(e: MouseEvent) {
        // ...
    }
}
```

**注意**：对象声明不可以是局部的(比如不可以直接在函数内部声明)，但可以在其它对象的声明或非内部类中进行内嵌入

### 伴随对象
在类声明内部可以用 companion 关键字标记对象声明：

```kotln
class MyClass {
	companion object Factory {
		fun create(): MyClass = MyClass()
	}
}
```

伴随对象的成员可以通过类名做限定词直接使用：

```kotlin
val instance = MyClass.create()
```

在使用了 `companion` 关键字时，伴随对象的名字可以省略：

```kotlin
class MyClass {
	companion object {

	}
}
```

注意，尽管伴随对象的成员很像其它语言中的静态成员，但在运行时它们任然是真正对象的成员实例，比如可以实现接口：

```kotlin
interface Factory<T> {
	fun create(): T
}

class MyClass {
	companion object : Factory<MyClass> {
		override fun create(): MyClass = MyClass()
	}
}
```

如果你在 JVM 上使用 `@JvmStatic` 注解，你可以有多个伴随对象生成为真实的静态方法和属性。参看 [java interoperabillity](https://kotlinlang.org/docs/reference/java-interop.html#static-methods-and-fields)。

### 对象表达式和声明的区别
他俩之间只有一个特别重要的区别：

>　对象表达式在我们使用的地方立即初始化并执行的
>
>　对象声明是懒加载的，是在我们第一次访问时初始化的。
>
>​    伴随对象是在对应的类加载时初始化的，和 Java 的静态初始是对应的。
