##扩展

与 C# 和 Gosu 类似, Kotlin 也提供了一种渠道,可以在不继承父类，也不使用类似装饰器这样的设计模式的情况下对指定类进行扩展。我们可以通过一种叫做扩展的特殊声明来实现他。现在， Kotlin 支持扩展函数和扩展属性。

###扩展函数

为了声明一个扩展函数，我们需要在函数名使用接收者类型作为前缀。下面我们会为 `MutableList<Int>` 添加一个 `swap` 函数：

```kotlin
fun MutableList<Int>.swap(x: Int, y: Int) {
	val temp = this[x] // this 对应 list
	this[x] = this[y]
	this[y] = tmp
}
```

在扩展函数中的 this 关键字对应接收者对象。现在我们可以在任何 `MutableList<Int>` 实例中使用这个函数了：

```kotlin
val l = mutableListOf(1, 2, 3)
l.swap(0, 2)
```

当然，这个函数对任意的 `MutableList<T>` 都是适用的，而且我们可以把它变的通用：

```kotlin
fun <T> MutableList<T>.swap(x: Int, y: Int) {
  val tmp = this[x] // 'this' corresponds to the list
  this[x] = this[y]
  this[y] = tmp
}
```

我们在函数名前声明了通用类型，从而使他可以接受任何参数。参看[通用函数](http://kotlinlang.org/docs/reference/generics.html)。

###扩展是**静态**解析的

扩展实际上并没有修改它所扩展的类。定义一个扩展，你并没有在类中插入一个新的成员，只是让这个类的实例对象能够通过`.`调用新的函数。

需要强调的是扩展函数是静态分发的，举个例子,它们并不是接受者类型的虚拟方法。如果有同名同参数的成员函数和扩展函数，调用的时候必然会使用成员函数，比如：

```kotlin
class C {
	fun foo() { Println("member") }
}

func C.foo { println("extension") }
```

当我们对C的实力c调用`c.foo()`的时候,他会输出"member",而不是"extension"

###空接受者

注意扩展可以使用空接受者进行定义。那样的话,扩展可以在一个值为空的对象变量被调用，并在函数体内检查 `this == null` 。这样你就可以在 Kotlin 中任意调用 toString() 方法而不进行空指针检查：空指针检查延后到扩展函数中完成。

```kotlin
fun Any?.toString(): String {
	if (this == null) return "null"
	return toString()
}
```

###扩展属性

和函数类似， Kotlin 也支持属性扩展：

```kotlin
val <T> List<T>.lastIndex:  Int
	get() = size-1
```

注意，由于扩展并不会真正给类添加了成员属性，因此也没有办法让扩展属性拥有一个备份字段.这也是为什么初始化函数是不允许扩展属性。扩展属性只能够通过直接提供 getter 和 setter方法来进行定义.

例子：
```kotlin
val Foo.bar = 1 //error: initializers are not allowed for extension properties
```

###伴随对象扩展

如果一个对象定义了伴随对象，你也可以给伴随对象添加扩展函数或扩展属性：

```kotlin
class MyClass {
	companion object {} 
}

fun Myclass.Companion,foo() {

}
```

和普通伴随对象的成员一样，它们可以只用类的名字就调用：

```kotlin
MyClass.foo()
```

###扩展的域

大多数时候我们在 top level 定义扩展，就在包下面直接定义：

```kotlin
package foo.bar
fun Baz.goo() { ... }
```

为了在除声明的包外使用这个扩展，我们需要在 import 时导入：

```kotlin
package com.example,usage

import foo.bar.goo // 导入所有名字叫 "goo" 的扩展
				
					// 或者

import foo.bar.* // 导入foo.bar包下得所有数据

fun usage(baz: Baz) {
	baz.goo()
}
```

###动机

在 java 中，我们通常使用一系列名字为 "*Utils" 的类: `FileUtils`,`StringUtils`等等。很有名的 `java.util.Collections` 也是其中一员的，但我们不得不像下面这样使用他们：

```java
//java
Collections.swap(list, Collections.binarySearch(list, Collections.max(otherList)), Collections.max(list))
```

由于这些类名总是不变的。我们可以使用静态导入并这样使用：

```java
swap(list, binarySearch(list, max(otherList)), max(list))
```

这样就好很多了，但这样我们就只能从 IDE 自动完成代码那里获得很少或得不到帮助信息。如果我们可以像下面这样那么就好多了

```kotlin
list.swap(list.binarySearch(otherList.max()), list.max())
```

但我们又不想在 List 类中实现所有可能的方法。这就是扩展带来的好处。