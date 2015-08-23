##可见性修饰词

类，对象，接口，构造函数，属性以及它们的 setter 方法都可以有可见性修饰词。( getter 方法作为属性时都是可见性)。在 Kotlin 中有四种修饰词：

> private -- 只在声明的范围和同一个模块的子范围可见；

> protected -- (只可以用在类或接口的成员上)和 privete 很像，但在子类中也可见；

> internal -- (默认使用) 在同一个模块中都可见；

> public -- 在任何地方均可见；

**注意**：函数如果有表达式并且所有属性均声明为 public 则必须有明确的返回值。(This is required so that we do not accidentally change a type that is a part of a public API by merely altering the implementation.)

```kotlin
public val foo: Int = 5 //明确的返回值

public fun bar(): Int = 5 //明确的返回值

public fun bar {} // 函数体为空，返回值是 Unit 不能随意改变，所以不需要指明

```

下面将解释不同类型的声明作用域。

###包

函数，属性和类，对象和接口可以在 "top-level" 声明：

```kotlin
package foo
fun baz() {}
class bar {}
```

> 如果没有指明任何可见性修饰词，默认使用 `internal` ,这就意味这在同一个模块内均是可见的

> 如果你声明为 `private` ，就只在本包或子包中可见，而且必须是同一个模块；

> 如果用 `public` 声明，则任何地方均可见

> `protected` 在 "top-level" 中不可以使用

例子：

```kotlin
package foo

private fun foo() {}//在本包及子包中可见

public var bar: Int = 5 // 任何地方均可见

private set // setter 仅在本包及子包中可见

internal val bax = 6 // 在同一个模块中可见，修饰词可省
```

###类和接口

当在类中声明时：

> `private` 只在该类(以及它的成员)中可见

> `protected` 和 `private` 一样但在子类中也可见

> `internal` 在本模块的所有可以访问该类的均可以访问该类的所有 `internal` 成员

> `public` 任何地方可见

java 使用者注意：外部类不可以访问内部类的 private 成员。

例子：

```kotlin
open class Outer {
	private val a = 1
	protected val b = 2
	val c = 3 //默认是 internal
	public val d: Int = 4 // 必须有返回值类型
	protected class Nested {
		public val e: Int = 5
	}
}

class Subclass : Outer() {
	//a 不可见
	//b c d 可见
	// 嵌套的 e 可见
}

class Unrelated(0: Outer) {
	//o.a , o.b 不可见
	//o.c , o.d 可见(必须是同一个模块)
	//Outer.Nested 不可见, Nested::e 也不可见
}
```

###构造函数

通过下面的语法(你必须显示的使用 constructor 关键字)来指定主构造函数的可见性：

```kotlin
class C private constructor(a: Int) { ... }
```

这里构造函数是公共的。不像其他的声明默认，所有的默认构造函数是 `public` ,实际上只要类是可见的它们就是可见的

###局部声明

局部变量，函数和类是不允许使用修饰词的