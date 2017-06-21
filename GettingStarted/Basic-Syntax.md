[原文](http://kotlinlang.org/docs/reference/basic-syntax.html)

## 基本语法

### 包定义
在源文件的开头定义包：

```kotlin
package my.demo
import java.util.*
//...
```

包名不必和文件夹路径一致：源文件可以放在任意位置。

更多请参看 [包(package)](../Basics/Packages.md)
### 定义函数
定义一个函数接受两个 int 型参数，返回值为 int ：

```kotlin
fun sum(a: Int, b: Int): Int {
    return a + b
}
```

只有一个表达式作为函数体，以及自推导型的返回值：

```kotlin
fun sum(a: Int, b: Int) = a + b
```

返回一个没有意义的值：

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

### 定义局部变量
一次赋值（只读）的局部变量：
```kotlin
val a: Int = 1  // 立刻赋值
val b = 2   // `Int` 类型是自推导的
val c: Int  // 没有初始化器时要指定类型
c = 3       // 推导型赋值
```

可修改的变量：

```kotlin
var x = 5 // `Int` type is inferred
x += 1
```

更多请参看[属性和字段](../ClassesAndObjects/Properties-and-Fields.md)

### 注释
与 java 和 javaScript 一样，Kotlin 支持单行注释和块注释。

```kotlin
// 单行注释

/*  哈哈哈哈
    这是块注释 */
```

与 java 不同的是 Kotlin 的 块注释可以级联。

参看[文档化 Kotlin 代码](../Tools/Documenting-Kotlin-Code.md)学习更多关于文档化注释的语法。

### 使用字符串模板
```kotlin
var a = 1
// simple name in template:
val s1 = "a is $a" 

a = 2
// arbitrary expression in template:
val s2 = "${s1.replace("is", "was")}, but now is $a"
```

更多请参看[字符串模板](../Basics/Basic-Types.md)

### 使用条件表达式
```kotlin
fun maxOf(a: Int, b: Int): Int {
    if (a > b) {
        return a
    } else {
        return b
    }
}
```

使用 if 作为表达式：

```kotlin
fun maxOf(a: Int, b: Int) = if (a > b) a else b
```

更多请参看 [if 表达式](../Basics/Control-Flow.md)

### 使用可空变量以及空值检查
当空值可能出现时应该明确指出该引用可空。

当 str 中不包含整数时返回空:

```kotlin
fun parseInt(str: String): Int? {
    // ...
}
```

使用一个返回可空值的函数：

```kotlin
fun parseInt(str: String): Int? {
  return str.toIntOrNull()
}

fun printProduct(arg1: String, arg2: String) {
  val x = parseInt(arg1)
  val y = parseInt(arg2)

  // 直接使用 x*y 会产生错误因为它们中有可能会有空值
  if (x != null && y != null) {
    // x 和 y 将会在空值检测后自动转换为非空值
    println(x * y)
  }
  else {
    println("either '$arg1' or '$arg2' is not a number")
  }    
}
```

或者这样

```kotlin
if (x == null) {
    println("Wrong number format in arg1: '${arg1}'")
    return
}
if (y == null) {
    println("Wrong number format in arg2: '${arg2}'")
    return
}

// x and y are automatically cast to non-nullable after null check
println(x * y)
```

更多请参看[空安全](../Other/Null-Safety.md)

### 使用值检查以及自动转换
使用 is 操作符检查一个表达式是否是某个类型的实例。如果对不可变的局部变量或属性进行过了类型检查，就没有必要明确转换：

```kotlin
fun getStringLength(obj: Any): Int? {
  if (obj is String) {
    // obj 将会在这个分支中自动转换为 String 类型
    return obj.length
  }

  // obj 在类型检查分支外仍然是 Any 类型
  return null
}
```

或者这样

```kotlin
fun getStringLength(obj: Any): Int? {
  if (obj !is String) return null
  
  // obj 将会在这个分支中自动转换为 String 类型
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

### 使用 for 循环
```kotlin
val items = listOf("apple", "banana", "kiwi")
for (item in items) {
    println(item)
}
```

或者

```kotlin
val items = listOf("apple", "banana", "kiwi")
for (index in items.indices) {
    println("item at $index is ${items[index]}")
}
```

参看[for循环](../Basics/Control-Flow.md)

### 使用 while 循环
```kotlin
val items = listOf("apple", "banana", "kiwi")
var index = 0
while (index < items.size) {
    println("item at $index is ${items[index]}")
    index++
}
```

参看[while循环](../Basics/Control-Flow.md)

### 使用 when 表达式
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

###  使用ranges
使用 in 操作符检查数值是否在某个范围内：

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
    println("list size is out of valid list indices range too")
}
```

使用范围内迭代：

```kotlin
for (x in 1..5) {
    print(x)
}
```

或者使用步进：

```kotlin
for (x in 1..10 step 2) {
    print(x)
}
for (x in 9 downTo 0 step 3) {
    print(x)
}
```

参看[Ranges](../Other/Ranges.md)

### 使用集合
对一个集合进行迭代：

```kotlin
for (item in items) {
    println(item)
}
```

使用 in 操作符检查集合中是否包含某个对象

```kotlin
when {
    "orange" in items -> println("juicy")
    "apple" in items -> println("apple is fine too")
}
```

使用lambda表达式过滤和映射集合：

```kotlin
fruits
.filter { it.startsWith("a") }
.sortedBy { it }
.map { it.toUpperCase() }
.forEach { println(it) }
```

参看[高阶函数和lambda表达式](../FunctionsAndLambdas/Higher-OrderFunctionsAndLambdas.md)

