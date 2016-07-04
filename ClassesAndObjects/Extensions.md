##扩展

与 C# 和 Gosu 类似, Kotlin 也提供了一种,可以在不继承父类，也不使用类似装饰器这样的设计模式的情况下对指定类进行扩展。我们可以通过一种叫做扩展的特殊声明来实现他。Kotlin 支持函数扩展和属性扩展。

###函数扩展

为了声明一个函数扩展，我们需要在函数前加一个接收者类型作为前缀。下面我们会为 `MutableList<Int>` 添加一个 `swap` 函数：

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
l.swap(0, 2)// 在 `swap()` 函数中 `this` 持有的值是 `l`
```

当然，这个函数对任意的 `MutableList<T>` 都是适用的，而且我们可以把它变的通用：

```kotlin
fun <T> MutableList<T>.swap(x: Int, y: Int) {
  val tmp = this[x] // 'this' corresponds to the list
  this[x] = this[y]
  this[y] = tmp
}
```

我们在函数名前声明了通用类型，从而使它可以接受任何参数。参看[泛型函数](http://kotlinlang.org/docs/reference/generics.html)。

###扩展是被**静态**解析的

扩展实际上并没有修改它所扩展的类。定义一个扩展，你并没有在类中插入一个新的成员，只是让这个类的实例对象能够通过`.`调用新的函数。

需要强调的是扩展函数是静态分发的，举个例子,它们并不是接受者类型的虚拟方法。这意味着扩展函数的调用时由发起函数调用的表达式的类型决定的，而不是在运行时动态获得的表达式的类型决定。比如

```Kotlin
open class C 

class D: C()fun C.foo() = "c" 
fun D.foo() = "d"fun printFoo(c: C) { 	println(c.foo())} 
printFoo(D())
```

这个例子会输出 `c`，因为这里扩展函数的调用决定于声明的参数 `c` 的类型，也就是 `C`。

如果有同名同参数的成员函数和扩展函数，调用的时候必然会使用成员函数，比如：

```kotlin

class C {
	fun foo() { println("member") }

}
fun C.foo() { println("extension") }
```

当我们对C的实例c调用`c.foo()`的时候,他会输出"member",而不是"extension"

但你可以用不同的函数签名通过扩展函数的方式重载函数的成员函数，比如下面这样：

```Kotlin
class C {
	fun foo() { println("number") }
}

fun C.foo(i:Int) { println("extention") }
```

`C().foo(1)` 的调用会打印 “extentions”。

###可空的接收者

注意扩展可以使用空接收者类型进行定义。这样的扩展使得，即使是一个空对象仍然可以调用该扩展，然后在扩展的内部进行 `this == null` 的判断。这样你就可以在 Kotlin 中任意调用 toString() 方法而不进行空指针检查：空指针检查延后到扩展函数中完成。

```kotlin
fun Any?.toString(): String {
	if (this == null) return "null"
	// 在空检查之后，`this` 被自动转为非空类型，因此 toString() 可以被解析到任何类的成员函数中
	return toString()
}
```

###属性扩展

和函数类似， Kotlin 也支持属性扩展：

```kotlin
val <T> List<T>.lastIndex:  Int
	get() = size-1
```

注意，由于扩展并不会真正给类添加了成员属性，因此也没有办法让扩展属性拥有一个备份字段.这也是为什么**初始化函数不允许有扩展属性**。扩展属性只能够通过明确提供 getter 和 setter方法来进行定义.

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
fun MyClass.Companion.foo(){

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
