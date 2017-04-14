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
```

该函数只有一个表达式函数体以及一个自推导型的返回值：

```kotlin
fun sum(a: Int, b: Int) = a + b
```

返回一个没有意义的值：

```kotlin
fun printSum(a: Int, b: Int): Unit{
	print( a + b)
}
```

Uint 的返回类型可以省略：

```kotlin
fun printSum(a: Int, b: Int){
	print( a + b)
}
```

更多请参看[函数](../FunctionsAndLambdas/Functions.md)

### 定义局部变量
声明常量：
```kotlin
val a: Int = 1
val b = 1 //推导出Int型
val c: Int //当没有初始化值时必须声明 类型
c = 1 // 赋值
```

变量：

```kotlin
 var x = 5 //推导出Int型
x += 1
```

更多请参看[属性和字段](../ClassesAndObjects/Properties-and-Filds.md)

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
	if (args.size == 0) return

	print("First argument: ${args[0]}")
}
```

更多请参看[字符串模板](../Basics/Basic-Types.md)

### 使用条件表达式
```kotlin
fun max(a: Int, b: Int): Int {
	if (a > b)
		return a
	else
		return b
}
```

把if当表达式：

```kotlin
	fun max(a: Int,  b: Int) = if (a > b) a else b
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
fun main(args: Array<String>) {
	if (args.size <2 ){
		print("Two integers expected")
		return
	}
	
	val x = parseInt(args[0])
	val y = parseInt(args[1])

	//直接使用 x*y 会产生错误因为它们中有可能会有空值
	if (x != null && y !=null){
		//x 和 y 将会在空值检测后自动转换为非空值
		print(x * y)
	}
}
```

或者这样

```kotlin
	//...
	if (x == null) {
		print("Wrong number format in '${args[0]}' ")
		return
	}
	if (y == null) {
		print("Wrong number format in '${args[1]}' ")
		return
	}
	//x 和 y 将会在空值检测后自动转换为非空值
	print(x * y)
```

更多请参看[空安全](../Other/Null-Safety.md)

### 使用值检查并自动转换
使用 is 操作符检查一个表达式是否是某个类型的实例。如果对不可变的局部变量或属性进行过了类型检查，就没有必要明确转换：

```kotlin
fun getStringLength(obj: Any): Int? {
	if ( obj is string ){
		//obj 将会在这个分支中自动转换为 String 类型
		return obj.length
	}
	// obj 在种类检查外仍然是 Any 类型
	return null
}
```

或者这样

```kotlin
fun getStringLength(obj: Any): Int? {
	if ( obj is string )
		return obj.length

	//obj 将会在这个分支中自动转换为 String 类型
	return null
}
```

甚至可以这样

```kotlin
fun getStringLength(obj: Any): Int? {
	if (obj is String && obj.length > 0)
		return obj.Length
	return null
}
```

更多请参看 [类](../ClassesAndObjects/Classes-and-Inheritance.md) 和 [类型转换](../Other/Type-Checks-and-Casts.md)

### 使用循环
```kotlin
fun main(args: Array<String>){
	for (arg in args)
		print(arg)
}
```

或者

```kotlin
for (i in args.indices)
	print(args[i])
```

参看[for循环](http://kotlinlang.org/docs/reference/control-flow.html#for-loops)

### 使用 while 循环
```kotlin
fun main(args: Array<Atring>){
	var i = 0
	while (i < args.size){
		print(args[i++])
	}
}
```

参看[while循环](http://kotlinlang.org/docs/reference/control-flow.html#while-loops)

### 使用 when 表达式
```kotlin
fun cases(obj: Any) {
    when (obj) {
        1 -> print("one")
        "hello" -> print("Greeting")
        is Long -> print("Long")
        !is Long -> print("Not a string")
        else -> print("Unknown")
    }
}
```

参看[when表达式](http://kotlinlang.org/docs/reference/control-flow.html#when-expression)

###  使用ranges
检查 in 操作符检查数值是否在某个范围内：

```kotlin
if (x in 1..y-1)
	print("OK")
```

检查数值是否在范围外：

```kotlin
if (x !in 0..array.lastIndex)
	print("Out")
```

参看[Ranges](http://kotlinlang.org/docs/reference/ranges.html)

### 使用集合
对一个集合进行迭代：

```kotlin
for (name in names)
	println(name)
```

使用 in 操作符检查集合中是否包含某个对象

```kotlin
if (text in names) //将会调用nemes.contains(text)方法
	print("Yes)
```

使用字面函数过滤和映射集合：

```kotlin
names
     .filter { it.startsWith("A") }
     .sortedBy { it }
     .map { it.toUpperCase() }
     .forEach { print(it) }
```

参看[高级函数和lambda表达式](../FunctionsAndLambdas/Higher-OrderFunctionsAndLambdas.md)
