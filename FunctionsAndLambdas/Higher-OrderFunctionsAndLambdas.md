## 高阶函数与 lambda 表达式
### 高阶函数
高阶函数就是可以接受函数作为参数或者返回一个函数的函数。比如 `lock()` 就是一个很好的例子，它接收一个 lock 对象和一个函数，运行函数并释放 lock;

```kotlin
fun <T> lock(lock: Lock, body: () -> T ) : T {
	lock.lock()
	try {
		return body()
	}
	finally {
		lock.unlock()
	}
}
```

现在解释一下上面的代码吧：`body` 有一个函数类型 `() -> T`,把它设想为没有参数并返回 T 类型的函数。它引发了内部的 try 函数块，并被 `lock` 保护，结果是通过 `lock()` 函数返回的。

如果我们想调用 `lock()` ，函数，我们可以传给它另一个函数做参数，参看[函数参考](http://kotlinlang.org/docs/reference/reflection.html#function-references)：

```kotlin
fun toBeSynchroized() = sharedResource.operation()

val result = lock(lock, ::toBeSynchroized)
```

其实最方便的办法是传递一个字面函数(通常是 lambda 表达式)：

```kotlin
val result = lock(lock, {
sharedResource.operation() })
```

字面函数经常描述有更多[细节](http://kotlinlang.org/docs/reference/lambdas.html#function-literals-and-function-expressions)，但为了继续本节，我们看一下更简单的预览吧：

> 字面函数被包在大括号里

> 参数在 `->` 前面声明(参数类型可以省略)

> 函数体在 `->` 之后

在 kotlin 中有一个约定，如果最后一个参数是函数，可以省略括号：

```kotlin
lock (lock) {
	sharedResource.operation()
}
```

最后一个高阶函数的例子是 `map()` (of MapReduce):

```kotlin
fun <T, R> List<T>.map(transform: (T) -> R):
List<R> {
	val result = arrayListOf<R>()
	for (item in this)
		result.add(transform(item))
	return result
}
```

函数可以通过下面的方式调用

```kotlin
val doubled = ints.map {it -> it * 2}
```

如果字面函数只有一个参数，则声明可以省略，名字就是 `it` :

```kotlin
ints map {it * 2}
```

这样就可以写[LINQ-风格](http://msdn.microsoft.com/en-us/library/bb308959.aspx)的代码了：

```kotlin
strings filter {it.length == 5} sortBy {it} map {it.toUpperCase()}
```

### 内联函数
有些时候可以用 [内联函数](http://kotlinlang.org/docs/reference/inline-functions.html) 提高高阶函数的性能。

### 字面函数和函数表达式
字面函数或函数表达式就是一个 "匿名函数"，也就是没有声明的函数，但立即作为表达式传递下去。想想下面的例子：

```kotlin
max(strings, {a, b -> a.length < b.length })
```
`max` 函数就是一个高阶函数,它接受函数作为第二个参数。第二个参数是一个表达式所以本生就是一个函数，即字面函数。作为一个函数，相当于：

```kotlin
fun compare(a: String, b: String) : Boolean = a.length < b.length
```

### 函数类型
一个函数要接受另一个函数作为参数，我们得给它指定一个类型。比如上面的 `max` 定义是这样的：

```kotlin
fun max<T>(collection: Collection<out T>, less: (T, T) -> Boolean): T? {
	var max: T? = null
	for (it in collection)
		if (max == null || less(max!!, it))
			max = it
	return max
}
```

参数 `less` 是 `(T, T) -> Boolean`类型，也就是接受俩个 `T` 类型参数返回一个 `Boolean`:如果第一个参数小于第二个则返回真。

在函数体第四行， `less` 是用作函数

一个函数类型可以像上面那样写，也可有命名参数，更多参看[命名参数](http://kotlinlang.org/docs/reference/functions.html#named-arguments)

```kotlin
val compare: (x: T,y: T) -> Int = ...
```

### 函数文本语法
函数文本的完全写法是下面这样的：

```kotlin
val sum = {x: Int,y: Int -> x + y}
```

函数文本总是在大括号里包裹着，在完全语法中参数声明是在括号内，类型注解是可选的，函数体是在　`->` 之后，像下面这样：

```kotlin
val sum: (Int, Int) -> Int = {x, y -> x+y }
```

函数文本有时只有一个参数。如果 kotlin 可以从它本生计算出签名，那么可以省略这个唯一的参数，并会通过 `it` 隐式的声明它：

```kotlin
ints.filter {it > 0}//这是 (it: Int) -> Boolean  的字面意思
```

注意如果一个函数接受另一个函数做为最后一个参数，该函数文本参数可以在括号内的参数列表外的传递。参看 [callSuffix](http://kotlinlang.org/docs/reference/grammar.html#call-suffix)

### 函数表达式
上面没有讲到可以指定返回值的函数。在大多数情形中，这是不必要的，因为返回值是可以自动推断的。然而，如果你需要自己指定，可以用函数表达式来做：

```kotlin
fun(x: Int, y: Int ): Int = x + y
```

函数表达式很像普通的函数声明，除了省略了函数名。它的函数体可以是一个表达式(像上面那样)或者是一个块：

```kotlin
fun(x: Int, y: Int): Int {
	return x + y
}
```

参数以及返回值和普通函数是一样的，如果它们可以从上下文推断出参数类型，则参数可以省略：

```kotlin
ints.filter(fun(item) = item > 0)
```

返回值类型的推导和普通函数一样：函数返回值是通过表达式自动推断并被明确声明

注意函数表达式的参数总是在括号里传递的。 The shorthand syntax allowing to leave the function outside the parentheses works only for function literals.

字面函数和表达式函数的另一个区别是没有本地返回。没有 lable 的返回总是返回到 fun 关键字所声明的地方。这意味着字面函数内的返回会返回到一个闭合函数，而表达式函数会返回到函数表达式自身。

### 闭包
一个字面函数或者表达式函数可以访问闭包，即访问自身范围外的声明的变量。不像 java 那样在闭包中的变量可以被捕获修改：

```kotlin
var sum = 0

ins filter {it > 0} forEach {
	sum += it
}
print(sum)
```

### 函数表达式扩展
除了普通的功能，kotlin 支持扩展函数。这种方式对于字面函数和表达式函数都是适用的。它们最重要的使用是在 [Type-safe Groovy-style builders](http://kotlinlang.org/docs/reference/type-safe-builders.html)。

表达式函数的扩展和普通的区别是它有接收类型的规范。

```kotlin
val sum = fun Int.(other: Int): Int = this + other
```

接收类型必须在表达式函数中明确指定，但字面函数不用。字面函数可以作为扩展函数表达式，但只有接收类型可以通过上下文推断出来。

表达式函数的扩展类型是一个带接收者的函数：

```kotlin
sum : Int.(other: Int) -> Int
```
可以用 . 或前缀来使用这样的函数：

```kotlin
1.sum(2)
1 sum 2
```
