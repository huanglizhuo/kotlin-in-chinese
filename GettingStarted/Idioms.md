##习惯用语

这里是一些在kotlin中随机但经常使用的习语。如果你有特别喜欢的习语想要贡献出来，赶快发起 pull request 吧。

###创建DTO's(POJO's/POCO's)

```kotlin
data class Customer(val name: String,val email: String)
```

给 Customer 类提供如下方法：

>  --为所有属性添加 getters 并为所有变量添加 setters
>  --equals()
>  --haseCode()
>  --toString()
>  --copy()
>  --component1() , component1() , ... 参看[数据类](http://kotlinlang.org/docs/reference/data-classes.html)

###声明局部 final 变量

```kotlin
val a = foo()
```
###函数默认值

```kotlin
fun foo(a: Int = 0, b: String = "") {...}
```

###过滤 list

```kotlin
val positives = list.filter { x -> x >0 }

```
或者更短：

```kotlin
val positives = list.filter { it > 0 }
```

###字符串插值

```kotlin
println( "Name $name" )
```

###实例检查

```kotlin
when (x) {
	is Foo ->  ...
	is Bar -> ...
	else -> ...
}
```

###遍历 map/list
```kotlin
for ((k, v) in map) {
	print("$k -> $v")
}
```
k,v 可以随便命名

###使用 ranges
```kotlin
for (i in 1..100) { ... }
for (i in 2..10) { ... }
```

###只读 list
```kotllin
val list = listOf("a", "b", "c")
```

###只读map

```kotllin
val map = maoOf("a" to 1, "b" to 2, "c" to 3)
```

###获取map中的值

```kotllin
println(map["key"])
map["key"] = value
```

###Lazy property(不知道怎么翻译 :(  )

```kotlin
val p: String by Delegates.lazy {

}
```

###扩展函数(给现有类增添新函数)
```kotlin
fun String.spcaceToCamelCase() { ... }
"Convert this to camelcase".spcaceToCamelCase()
```

###创建单例模式
```kotlin
object Resource {
	val name = "Name"
}
```

###If not null shorthand(没想到怎么翻译)
```kotlin
 val files = File("Test").listFiles()
println(files?.size)
```

###If not null and else shorthand(没想到怎么翻译)
```kotlin
 val files = File("test").listFiles()
println(files?.size ?: "empty")
```

###如果为空执行某操作
```kotlin
val data = ...
val email = data["email"] ?: throw
IllegalStateException("Email is missing!")
```

###如果不为空执行某操作
```kotlin
val date = ...
data?.let{
	...//如果不为空执行该语句块
}
```

###返回 when 判断
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

###返回 try-catch 语句块

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
###返回 if 判断
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

###只有一个表达式的函数
```kotlin
fun theAnswer() = 42
```
与下面的语句是等效的

```kotlin
fun theAnswer(): Int {
	return 42
}
```
这个可以和其它习语组合成高效简洁的代码。譬如说 when 表达式：

```kotlin
fun transform(color: String): Int = when (color) {
	"Red" -> 0
	"Green" -> 1
	"Blue" -> 2
	else -> throw IllegalArgumentException("Invalid color param value")
}
```