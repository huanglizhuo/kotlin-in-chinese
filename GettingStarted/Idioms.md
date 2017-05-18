## 习惯用语
这里是一些在 Kotlin 中经常使用的习语。如果你有特别喜欢的习语想要贡献出来，赶快发起 pull request 吧。

### 创建DTOs(POJOs/POCOs)  数据类
```kotlin
data class Customer(val name: String, val email: String)
```

给 Customer 类提供如下方法：

>  --为所有属性添加 getters ，如果为 var 类型同时添加 setters
>  --`equals()`
>  --`haseCode()`
>  --`toString()`
>  --`copy()`
>  --`component1()` , `component1()` , ... 参看[数据类](../ClassesAndObjects/Data-Classes.md)

### 函数默认值
```kotlin
fun foo(a: Int = 0, b: String = "") {...}
```

### 过滤 list
```kotlin
val positives = list.filter { x -> x >0 }
```
或者更短：

```kotlin
val positives = list.filter { it > 0 }
```

### 字符串插值
```kotlin
println("Name $name")
```

### 实例检查
```kotlin
when (x) {
	is Foo ->  ...
	is Bar -> ...
	else -> ...
}
```

### 遍历 map/list
```kotlin
for ((k, v) in map) {
	print("$k -> $v")
}
```
k,v 可以随便命名

### 使用 ranges
```kotlin
for (i in 1..100) { ... }  // 闭区间: 包括100
for (i in 1 until 100) { ... } // 半开区间: 不包括100
for (x in 2..10 step 2) { ... }
for (x in 10 downTo 1) { ... }
if (x in 1..10) { ... }
for (i in 1..100) { ... }
for (i in 2..10) { ... }
```

### 只读 list
```kotllin
val list = listOf("a", "b", "c")
```

### 只读map
```kotllin
val map = mapOf("a" to 1, "b" to 2, "c" to 3)
```

### 访问 map
```kotllin
println(map["key"])
map["key"] = value
```

### 懒属性(延迟加载)
```kotlin
val p: String by lazy {
	// 生成string的值
}
```

### 扩展函数
```kotlin
fun String.spcaceToCamelCase() { ... }
"Convert this to camelcase".spcaceToCamelCase()
```

### 创建单例模式
```kotlin
object Resource {
	val name = "Name"
}
```

### 如果不为空则... 的简写
```kotlin
val files = File("Test").listFiles()
println(files?.size)
```

### 如果不为空...否则... 的简写
```kotlin
val files = File("test").listFiles()
println(files?.size ?: "empty")
```

### 如果声明为空执行某操作
```kotlin
val data = ...
val email = data["email"] ?: throw IllegalStateException("Email is missing!")
```

### 如果不为空执行某操作
```kotlin
val date = ...
data?.let{
	...//如果不为空执行该语句块
}
```

### 返回 when 判断
```kotlin
fun transform(color: String): Int {
	return when(color) {
		"Red" -> 0
		"Green" -> 1
		"Blue" -> 2
		else -> throw IllegalArgumentException("Invalid color param value")
	}
}
```

### try-catch 表达式
```kotlin
fun test() {
	val result = try {
		count()
	}catch (e: ArithmeticException) {
		throw IllegaStateException(e)
	}
	//处理 result
}
```
###  if 表达式
```kotlin
fun foo(param: Int){
	val result = if (param == 1) {
		"one"
	} else if (param == 2) {
		"two"
	} else {
		"three"
	}
}
```

### 方法使用生成器模式返回 Unit
```kotlin
fun arrOfMinusOnes(size: Int): IntArray{
	return IntArray(size).apply{ fill(-1) }
}
```

### 只有一个表达式的函数
```kotlin
fun theAnswer() = 42
```
与下面的语句是等效的

```kotlin
fun theAnswer(): Int {
	return 42
}
```
这个可以和其它习惯用语组合成高效简洁的代码。譬如说 when 表达式：

```kotlin
fun transform(color: String): Int = when (color) {
	"Red" -> 0
	"Green" -> 1
	"Blue" -> 2
	else -> throw IllegalArgumentException("Invalid color param value")
}
```
### 利用 with 调用一个对象实例的多个方法
```kotlin
class Turtle {
	fun penDown()
	fun penUp()
	fun turn(degrees: Double)
	fun forward(pixels: Double)
}
val myTurtle = Turtle()
with(myTurtle) { // 画一个100像素的正方形
	penDown()
	for(i in 1..4) {
    forward(100.0)
		turn(90.0)
	}
	penUp()
}
```

### Java 7’s try with resources
```kotlin
val stream = Files.newInputStream(Paths.get("/some/file.txt"))
stream.buffered().reader().use { reader ->
	println(reader.readText())
}
```

### 需要泛型信息的泛型函数的方便形式
```kotlin
//  public final class Gson {
//     ...
//     public <T> T fromJson(JsonElement json, Class<T> classOfT) throws JsonSyntaxException {
//     ...

inline fun <reified T: Any> Gson.fromJson(json): T = this.fromJson(json, T::class.java)
```

### 产生一个可能为空的布尔值
```kotlin
val b: Boolean? = ...
if (b == true) {
  ...
} else {
  // `b` 是false或者null
}
```
