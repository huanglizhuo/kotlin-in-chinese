[原文](http://kotlinlang.org/docs/reference/basic-syntax.html)

# 准备开始
## 基本语法
### 定义包名
在源文件的开头定义包名：

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
fun sum(a: Int , b: Int) : Int{
	return a + b
}

fun main(args: Array<String>) {
  print("sum of 3 and 5 is ")
  println(sum(3, 5))
}
```

该函数只有一个表达式函数体以及一个自推导型的返回值：

```kotlin
fun sum(a: Int, b: Int) = a + b

fun main(args: Array<String>) {
  println("sum of 19 and 23 is ${sum(19, 23)}")
}
```

返回一个没有意义的值：

```kotlin
fun printSum(a: Int, b: Int): Unit {
  println("sum of $a and $b is ${a + b}")
}

fun main(args: Array<String>) {
  printSum(-1, 8)
}
```

Unit 的返回类型可以省略：

```kotlin
fun printSum(a: Int, b: Int) {
  println("sum of $a and $b is ${a + b}")
}

fun main(args: Array<String>) {
  printSum(-1, 8)
}
```

更多请参看[函数](../FunctionsAndLambdas/Functions.md)

### 定义局部变量
声明常量：
```kotlin
fun main(args: Array<String>) {
  val a: Int = 1  // 立即初始化
  val b = 2   // 推导出Int型
  val c: Int  // 当没有初始化值时必须声明类型
  c = 3       // 赋值
  println("a = $a, b = $b, c = $c")
}
```

变量：

```kotlin
fun main(args: Array<String>) {
  var x = 5 // 推导出Int类型
  x += 1
  println("x = $x")
}
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

参看[文档化 Kotlin 代码](http://kotlinlang.org/docs/reference/kotlin-doc.html)学习更多关于文档化注释的语法。

### 使用字符串模板
```kotlin
fun main(args: Array<String>) {
  var a = 1
  // 使用变量名作为模板:
  val s1 = "a is $a"

  a = 2
  // 使用表达式作为模板:
  val s2 = "${s1.replace("is", "was")}, but now is $a"
  println(s2)
}
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

fun main(args: Array<String>) {
    println("max of 0 and 42 is ${maxOf(0, 42)}")
}
```

把if当表达式：

```kotlin
fun maxOf(a: Int, b: Int) = if (a > b) a else b

fun main(args: Array<String>) {
    println("max of 0 and 42 is ${maxOf(0, 42)}")
}
```

更多请参看[if表达式](../Basics/Control-Flow.md)

### 使用可空变量以及空值检查
当空值可能出现时应该明确指出该引用可空。

下面的函数是当 str 中不包含整数时返回空:

```kotlin
fun parseInt(str : String): Int?{
	//...
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


fun main(args: Array<String>) {
  printProduct("6", "7")
  printProduct("a", "7")
  printProduct("a", "b")
}
```

或者这样

```kotlin
fun parseInt(str: String): Int? {
  return str.toIntOrNull()
}

fun printProduct(arg1: String, arg2: String) {
  val x = parseInt(arg1)
  val y = parseInt(arg2)

  // ...
  if (x == null) {
    println("Wrong number format in arg1: '${arg1}'")
    return
  }
  if (y == null) {
    println("Wrong number format in arg2: '${arg2}'")
    return
  }

  // x 和 y 将会在空值检测后自动转换为非空值
  println(x * y)
}
```

更多请参看[空安全](../Other/Null-Safety.md)

### 使用值检查并自动转换
使用 is 操作符检查一个表达式是否是某个类型的实例。如果对不可变的局部变量或属性进行过了类型检查，就没有必要明确转换：

```kotlin
fun getStringLength(obj: Any): Int? {
  if (obj is String) {
    // obj 将会在这个分支中自动转换为 String 类型
    return obj.length
  }

  // obj 在种类检查外仍然是 Any 类型
  return null
}


fun main(args: Array<String>) {
  fun printLength(obj: Any) {
    println("'$obj' string length is ${getStringLength(obj) ?: "... err, not a string"} ")
  }
  printLength("Incomprehensibilities")
  printLength(1000)
  printLength(listOf(Any()))
}
```

或者这样

```kotlin
fun getStringLength(obj: Any): Int? {
  if (obj !is String) return null

  // obj 将会在这个分支中自动转换为 String 类型
  return obj.length
}


fun main(args: Array<String>) {
  fun printLength(obj: Any) {
    println("'$obj' string length is ${getStringLength(obj) ?: "... err, not a string"} ")
  }
  printLength("Incomprehensibilities")
  printLength(1000)
  printLength(listOf(Any()))
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


fun main(args: Array<String>) {
  fun printLength(obj: Any) {
    println("'$obj' string length is ${getStringLength(obj) ?: "... err, is empty or not a string at all"} ")
  }
  printLength("Incomprehensibilities")
  printLength("")
  printLength(1000)
}
```

更多请参看 [类](../ClassesAndObjects/Classes-and-Inheritance.md) 和 [类型转换](../Other/Type-Checks-and-Casts.md)

### 使用循环
```kotlin
fun main(args: Array<String>) {
  val items = listOf("apple", "banana", "kiwi")
  for (item in items) {
    println(item)
  }
}
```

或者

```kotlin
fun main(args: Array<String>) {
  val items = listOf("apple", "banana", "kiwi")
  for (index in items.indices) {
    println("item at $index is ${items[index]}")
  }
}
```

参看[for循环](http://kotlinlang.org/docs/reference/control-flow.html#for-loops)

### 使用 while 循环
```kotlin
fun main(args: Array<String>) {
  val items = listOf("apple", "banana", "kiwi")
  var index = 0
  while (index < items.size) {
    println("item at $index is ${items[index]}")
    index++
  }
}
```

参看[while循环](http://kotlinlang.org/docs/reference/control-flow.html#while-loops)

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

fun main(args: Array<String>) {
  println(describe(1))
  println(describe("Hello"))
  println(describe(1000L))
  println(describe(2))
  println(describe("other"))
}
```

参看[when表达式](http://kotlinlang.org/docs/reference/control-flow.html#when-expression)

###  使用ranges
使用 in 操作符检查数值是否在某个范围内：

```kotlin
fun main(args: Array<String>) {
  val x = 10
  val y = 9
  if (x in 1..y+1) {
      println("fits in range")
  }
}
```

检查数值是否在范围外：

```kotlin
fun main(args: Array<String>) {
  val list = listOf("a", "b", "c")

  if (-1 !in 0..list.lastIndex) {
    println("-1 is out of range")
  }
  if (list.size !in list.indices) {
    println("list size is out of valid list indices range too")
  }
}
```

在范围内迭代：

```kotlin
fun main(args: Array<String>) {
  for (x in 1..5) {
    print(x)
  }
}
```

或者使用步进：

```kotlin
fun main(args: Array<String>) {
  for (x in 1..10 step 2) {
    print(x)
  }
  for (x in 9 downTo 0 step 3) {
    print(x)
  }
}
```

参看[Ranges](http://kotlinlang.org/docs/reference/ranges.html)

### 使用集合
对一个集合进行迭代：

```kotlin
fun main(args: Array<String>) {
  val items = listOf("apple", "banana", "kiwi")
  for (item in items) {
    println(item)
  }
}
```

使用 in 操作符检查集合中是否包含某个对象

```kotlin
fun main(args: Array<String>) {
  val items = setOf("apple", "banana", "kiwi")
  when {
    "orange" in items -> println("juicy")
    "apple" in items -> println("apple is fine too")
  }
}
```

使用lambda表达式过滤和映射集合：

```kotlin
fun main(args: Array<String>) {
  val fruits = listOf("banana", "avocado", "apple", "kiwi")
  fruits
    .filter { it.startsWith("a") }
    .sortedBy { it }
    .map { it.toUpperCase() }
    .forEach { println(it) }
}
```

参看[高阶函数和lambda表达式](../FunctionsAndLambdas/Higher-OrderFunctionsAndLambdas.md)
