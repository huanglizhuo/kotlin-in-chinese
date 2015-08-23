##扩展

Kotlin 与 C# Gosu 类似，提供了不用从父类继承，或者使用像装饰模式这样的设计模式来给某个类进行扩展。这是通过叫扩展的特殊声明来达到的。现在， Kotlin 支持扩展函数和属性。

###扩展函数

声明一个扩展函数，我们需要添加一个接收者类型的的前缀。下面是给 `MutableList<Int>` 添加一个 `swap` 函数的例子：

```kotlin
fun MutableList<Int>.swap(x: Int, y: Int) {
	val temp = this[x] // this 对应 list
	this[x] = this[y]
	this[y] = tmp
}
```

在扩展函数中的 this 关键字对应接收者对象。现在我们可以在任何 `MutableList<Int>` 中使用这个函数了：

```kotlin
val l = mutableListOf(1, 2, 3)
l.swap(0, 2)
```

当然，这个函数对任何 `MutableList<Int>` 都是适用的，而且我们可以把它变的通用：

```kotlin
fun <T> MutableList<T>.swap(x: Int, y: Int) {
  val tmp = this[x] // 'this' corresponds to the list
  this[x] = this[y]
  this[y] = tmp
}
```

我们在函数名前声明了通用类型，从而使他可以接受任何参数。参看[通用函数](http://kotlinlang.org/docs/reference/generics.html)。

###扩展是**静态**解析的

扩展不需要修改它们扩展的类。定义一个扩展，你不需要在类中插入一个新的成员，而只需要添加一个可以通过 . 注解调用的函数就可以了。

必须注意的是扩展函数是静态触动的，它们的接受者类型不是虚拟的。如果有同名同参数的成员函数和扩展函数，总是触动成员函数，比如：

```kotlin
class C {
	fun foo() { Println("number") }
}

func C.foo { println("extention") }
```

###可空的接受者

值得注意的是扩展可以定义一个可空的接受者。这样的扩展可以算作对象的变量，即使它是空的，你可以在函数体内检查 `this == null` 。这样你就可以在 Kotlin 中不进行空检查就可以调用 toString() 方法：这样的检查是在扩展函数中做的。

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

注意，扩展并不是真正给类添加了成员，没有比给扩展添加备用字段更有效的办法了／这就是为什么初始化函数是不允许扩展属性。它们只能通过明确提供 getter setter 来作用。

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
MyClass.foo(0
```

###扩展的范围

大多数时候我们在 top level 定义扩展，就在包下面直接定义：

```kotlin
package foo.bar
fun Baz.goo() { ... }
```
在声明的包外面使用这样的扩展，我们需要在 import 时导入：

```kotlin
package com.example,usage

import foo.bar.goo//导入所有名字叫 "goo" 的扩展

import foo.bar.*

fun usage(baz: Baz) {
	baz.goo()
}
```

###动机

在 java 中，我经常给类命名为 "*Utils": `FileUtils`,`StringUtils`等等。很有名的 `java.util.Collections` 也是这样的，这些 Utils 类不方便的地方就是我们用起来总是像下面这样：

```java
//java
Collections.swap(list, Collections.binarySearch(list, Collections.max(otherList)), Collections.max(list))
```

这些类名总是通过一样的方式得到的。我我们可以使用静态导入并这样使用：

```java
swap(list, binarySearch(list, max(otherList)), max(list))
```

这样就好很多了，但这样我们就只能从 IDE 自动完成代码那里获得很少或得不到帮助。如果我们可以像下面这样那门就好多了

```kotlin
list.swap(list.binarySearch(otherList.max()), list.max())
```

但我们又不想在 List 类中实现所有可能的方法。这就是扩展带来的好处。