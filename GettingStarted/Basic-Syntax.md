[原文](http://kotlinlang.org/docs/reference/basic-syntax.html)

## 基本语法

### 包定义和引入

在源文件的开头定义包：

```kotlin
package my.demo

import kotlin.text.*

// ...
```

包名不必和文件夹路径一致：源文件可以放在任意位置。

更多请参看 [包(package)](../Basics/Packages.md)

### 程序入口

Kotlin 应用的入口是 `main` 函数.

```kotlin
fun main() {
    println("Hello world!")
}
```

### 函数

下面的函数接受两个 `Int` 型参数，返回值为 `Int` ：

```kotlin
fun sum(a: Int, b: Int): Int {
    return a + b
}
```

具有表达式主体和推断的返回类型的函数：

```kotlin
fun sum(a: Int, b: Int) = a + b
```

返回无意义的值：

```kotlin
fun printSum(a: Int, b: Int): Unit {
    println("sum of $a and $b is ${a + b}")
}
```

`Unit` 的返回类型可以省略：

```kotlin
fun printSum(a: Int, b: Int) {
  println("sum of $a and $b is ${a + b}")
}
```

更多请参看[函数](../FunctionsAndLambdas/Functions.md)

### 变量

只读的本地变量通过`val`关键字定义.该类变量只能赋值一次：

```kotlin
val a: Int = 1  // 立刻赋值
val b = 2   // `Int` 类型是自推导的
val c: Int  // 没有初始化器时要指定类型
c = 3       // 推断型赋值
```

被关键字`var`修饰的变量可以重新赋值：

```kotlin
var x = 5 // `Int` type is inferred
x += 1
```

顶级变量:

```kotlin
val PI = 3.14
var x = 0

fun incrementX() { 
    x += 1 
}
```

更多请参看[属性和字段](../ClassesAndObjects/Properties-and-Fields.md)

### 注释
与多数现代语言一样，Kotlin 支持单行注释(行尾注释)和多行注释(块注释)。

```kotlin
// 行尾注释

/*  这是块注释
    可以在多行注释 */
```

Kotlin 块注释可以嵌套.

```kotlin 
/* The comment starts here
/* contains a nested comment */     
and ends here. */
```

参看[文档化 Kotlin 代码](../Tools/Documenting-Kotlin-Code.md)更多关于文档化注释的语法。

### 字符串模板

```kotlin
var a = 1
// simple name in template:
val s1 = "a is $a" 

a = 2
// arbitrary expression in template:
val s2 = "${s1.replace("is", "was")}, but now is $a"
```

更多请参看[字符串模板](../Basics/Basic-Types.md)

### 条件表达式

```kotlin
fun maxOf(a: Int, b: Int): Int {
    if (a > b) {
        return a
    } else {
        return b
    }
}
```

kotin 中可以使用 if 作为表达式：

```kotlin
fun maxOf(a: Int, b: Int) = if (a > b) a else b
```

更多请参看 [if 表达式](../Basics/Control-Flow.md)

### 可空变量以及空值检查

当空值可能出现时必须明确标注该引用可空。

当 str 中不包含整数时返回空:

```kotlin
fun parseInt(str: String): Int? {
    // ...
}
```

使用函数返回空值：

```kotlin
fun printProduct(arg1: String, arg2: String) {
    val x = parseInt(arg1)
    val y = parseInt(arg2)

    // Using `x * y` yields error because they may hold nulls.
    if (x != null && y != null) {
        // x and y are automatically cast to non-nullable after null check
        println(x * y)
    }
    else {
        println("'$arg1' or '$arg2' is not a number")
    }    
}
```

或者

```kotlin
if (x == null) {
    println("Wrong number format in arg1: '$arg1'")
    return
}
if (y == null) {
    println("Wrong number format in arg2: '$arg2'")
    return
}

// x and y are automatically cast to non-nullable after null check
println(x * y)
```

更多请参看[空安全](../Other/Null-Safety.md)

### 类型检查以及自动转换

`is` 操作符可以检查表达式是否是是某个类型的实例。如果不可变的局部变量或属性进行过了类型检查，就没有必要显示转换：

```kotlin
fun getStringLength(obj: Any): Int? {
  if (obj is String) {
    // obj 将会在这个分支中自动转换为 `String` 类型
    return obj.length
  }

  // obj 在类型检查分支外仍然是 Any 类型
  return null
}
```

或者

```kotlin
fun getStringLength(obj: Any): Int? {
  if (obj !is String) return null
  
  // obj 将会在这个分支中自动转换为 `String` 类型
  return obj.length
}
```

甚至可以这样

```kotlin
fun getStringLength(obj: Any): Int? {
	// obj 将会在&&右边自动转换为 String 类型
  if (obj is String && obj.length > 0) {
    return obj.length
  }

  return null
}
```

更多请参看 [类](../ClassesAndObjects/Classes-and-Inheritance.md#3) 和 [类型转换](../Other/Type-Checks-and-Casts.md)

### for 循环
```kotlin
val items = listOf("apple", "banana", "kiwifruit")
for (item in items) {
    println(item)
}
```

或者

```kotlin
val items = listOf("apple", "banana", "kiwifruit")
for (index in items.indices) {
    println("item at $index is ${items[index]}")
}
```

参看[for循环](../Basics/Control-Flow.md)

### while 循环
```kotlin
val items = listOf("apple", "banana", "kiwifruit")
var index = 0
while (index < items.size) {
    println("item at $index is ${items[index]}")
    index++
}
```

参看[while循环](../Basics/Control-Flow.md)

### when 表达式
```kotlin
fun describe(obj: Any): String =
    when (obj) {
        1          -> "One"
        "Hello"    -> "Greeting"
        is Long    -> "Long"
        !is String -> "Not a string"
        else       -> "Unknown"
    }
```

参看[when表达式](../Basics/Control-Flow.md)

###  ranges

使用 `in` 操作符判断数值是否在某个范围内：

```kotlin
val x = 10
val y = 9
if (x in 1..y+1) {
    println("fits in range")
}
```

检查数值是否在范围外：

```kotlin
val list = listOf("a", "b", "c")

if (-1 !in 0..list.lastIndex) {
    println("-1 is out of range")
}
if (list.size !in list.indices) {
    println("list size is out of valid list indices range, too")
}
```

参看[Ranges](../Other/Ranges.md)

### 集合

遍历集合：

```kotlin
for (item in items) {
    println(item)
}
```

使用 `in` 操作符检查集合中是否包含某个对象:

```kotlin
when {
    "orange" in items -> println("juicy")
    "apple" in items -> println("apple is fine too")
}
```

使用lambda表达式过滤和映射集合：

```kotlin
val fruits = listOf("banana", "avocado", "apple", "kiwifruit")
fruits
  .filter { it.startsWith("a") }
  .sortedBy { it }
  .map { it.toUpperCase() }
  .forEach { println(it) }
```

参看[集合概述](../Collections/CollectionsOverview.md)

### 创建基本类以及实例: 

```kotlin
val rectangle = Rectangle(5.0, 2.0)
val triangle = Triangle(3.0, 4.0, 5.0)
```

参看[类](../ClassesAndObjects/Classes-and-Inheritance.md)和[对象及实例](../ClassesAndObjects/ObjectExpressicAndDeclarations.md)