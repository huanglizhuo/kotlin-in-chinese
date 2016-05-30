##对象表达式和声明

有时候我们需要创建一个对当前类做轻微修改的对象，而不用重新声明一个子类。java 用匿名内部类来解决这个问题。Kotlin 更希望推广用对象表达式和声明来解决这个问题。

###对象表达式

我们通过下面这样的方式创建继承自某种(或某些)匿名类的对象：

```kotlin
window.addMouseListener(object: MouseAdapter () {
	override fun mouseClicked(e: MouseEvent) {
		//...
	}
})
```

如果父类有构造函数，则必须传递相应的构造函数。多个父类可以用逗号隔开，跟在冒号后面：

```kotlin
open class A(x: Int) {
	public open val y: Int = x
}

interface B { ... }

val ab = object : A(1), B {
	override val y = 14
}
```

有时候我们只是需要一个没有父类的对象，我们可以这样写：

```kotlin
val adHoc = object {
	var x: Int = 0
	var y: Int = 0
}

print(adHoc.x + adHoc.y)
```

像 java 的匿名内部类一样，对象表达式可以访问闭合范围内的变量 (和 java 不一样的是，这些不用声明为 final)

```kotlin
fun countClicks(windows: JComponent) {
	var clickCount = 0
	var enterCount = 0
	window.addMouseListener(object : MouseAdapter() {
		override fun mouseClicked(e: MouseEvent) {
			clickCount++
		}
		override fun mouseEntered(e: MouseEvent){
			enterCount++
		}
	})
}
```

###对象声明

[单例模式](http://en.wikipedia.org/wiki/Singleton_pattern)是一种很有用的模式，Kotln 中声明它很方便：

```kotlin
object DataProviderManager {
	fun registerDataProvider(provider: Dataprovider) {
		//...
	}
	val allDataProviders : Collection<DataProvider>
		get() = //...
}
```

这叫做对象声明。如果在 object 关键字后面有个名字，我们不能把它当做表达式了。虽然不能把它赋值给变量，但是我们可以直接通过名字来使用这个类。这样的对象可以有父类：

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

**注意**：对象声明不可以是局部的(比如不可以直接在函数内部声明)，但可以在其它对象的声明或非内部类中使用

###伴随对象

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
inerface Factory<T> {
	fun create(): T
}

class MyClass {
	companion object : Factory<MyClass> {
		override fun create(): MyClass = MyClass()
	}
}
```

当然你可以通过 `@platfoemStatic` 注解使 JVM 将伴随对象生成为静态方法和字段。参看 [java interoperabillity](http://kotlinlang.org/docs/reference/java-interop.html#static-methods-and-fields)

###对象表达式和声明的区别

他俩之间只有一个特别重要的区别：

>　对象声明是 lazily 初始化的，我们只能访问一次

>　对象表达式在我们使用的地方立即初始化并执行的
