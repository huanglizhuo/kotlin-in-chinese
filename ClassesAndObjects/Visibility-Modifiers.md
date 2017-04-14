## 可见性修饰词
类，对象，接口，构造函数，属性以及它们的 setter 方法都可以有可见性修饰词。( getter 方法作为属性时都是可见性)。在 Kotlin 中有四种修饰词：`private`,`protected`,`internal`,以及 `public` 。默认的修饰符是 `public`。
下面将解释不同类型的声明作用域。

### 包
函数，属性和类，对象和接口可以在 "top-level" 声明：

```kotlin
package foo
fun baz() {}
class bar {}
```

> 如果没有指明任何可见性修饰词，默认使用 `public` ,这意味着你的声明在任何地方都可见；

> 如果你声明为 `private` ，则只在包含声明的文件中可见；

> 如果用 `internal` 声明，则在同一模块中的任何地方可见；

> `protected` 在 "top-level" 中不可以使用

例子：

```kotlin
package foo

private fun foo() {} // visible inside example.kt

public var bar: Int = 5 // property is visible everywhere 

private set // setter is visible only in example.kt

internal val baz = 6 // visible inside the same module
```

### 类和接口
当在类中声明时：

> `private` 只在该类(以及它的成员)中可见

> `protected` 和 `private` 一样但在子类中也可见

> `internal` 在本模块的所有可以访问到声明区域的均可以访问该类的所有 `internal` 成员 ( internal — any client inside this module who sees the declaring class sees its internal members;)

> `public` 任何地方可见 (public — any client who sees the declaring class sees its public members.)

java 使用者注意：外部类不可以访问内部类的 private 成员。

例子：

```kotlin
open class Outer {
	private val a = 1
	protected val b = 2
	internal val c = 3
	val d = 4 // public by default
	protected class Nested { 
		public val e: Int = 5
	} 
}
class Subclass : Outer() {
	// a is not visible
	// b, c and d are visible // Nested and e are visible
}
class Unrelated(o: Outer) {
	// o.a, o.b are not visible
	// o.c and o.d are visible (same module)
	// Outer.Nested is not visible, and Nested::e is not visible either
}
```

### 构造函数
通过下面的语法来指定主构造函数(必须显示的使用 constructor 关键字)的可见性：

```kotlin
class C private constructor(a: Int) { ... }
```

这里构造函数是 private 。所有的构造函数默认是 `public` ,实际上只要类是可见的它们就是可见的
(注意 `internal` 类型的类中的 public 属性只能在同一个模块内才可以访问)

### 局部声明
局部变量，函数和类是不允许使用修饰词的

### 模块
`internal` 修饰符是指成员的可见性是只在同一个模块中才可见的。模块在 Kotlin 中就是一系列的 Kotlin 文件编译在一起：

— an IntelliJ IDEA module;

— a Maven or Gradle project;

— a set of files compiled with one invocation of the Ant task.
