##函数

###函数声明

在 kotlin 中用关键字 fun 声明函数：

```kotlin
fun double(x: Int): Int {

}
```

###函数用法

通过传统的方法调用函数

```kotlin
val result = double(2)
```
通过 . 注解调用

```kotlin
Sample().foo()
```
###中缀标记
函数也可以通过中缀表达式调用，只要符和下面的规则

>　它们是成员函数或者是[扩展函数](http://kotlinlang.org/docs/reference/extensions.html)
>　只有一个参数

```kotlin
//给 Int 定义一个扩展方法
fun Int.shl(x: Int): Int {
...
}

1 shl 2 //用中缀注解调用扩展函数

1.shl(2)
```

###参数

函数参数是用 Pascal 符号定义的　name:type。参数之间用逗号隔开，每个参数必须指明类型。

```kotlin
fun powerOf(number: Int, exponent: Int) {
...
}
```

###默认参数

函数参数可以有默认参数。这样相比其他语言可以减少重载。

```kotlin
fun read(b: Array<Byte>, off: Int = 0, len: Int = b.size() ) {
...
}
```

###命名参数

在调用函数时可以参数可以命名。这对于有很多参数或只有一个的函数来说很方便。

下面是一个例子：

```kotlin
fun reformat(str: String, normalizeCase: Boolean = true,upperCaseFirstLetter: Boolean = true,
             divideByCamelHumps: Boolean = false,
             wordSeparator: Char = ' ') {
...
}
```

我们可以使用默认参数

>reformat(str)

然而当调用非默认参数是就会像下面这样：

```kotlin
reformat(str, true, true, false, '_')
```

使用命名参数我们可以让代码可读性更强：

```kotlin
reformat(str,
    normalizeCase = true,
    uppercaseFirstLetter = true,
    divideByCamelHumps = false,
    wordSeparator = '_'
  )
```

如果不需要全部参数的话可以这样：

```kotlin
reformat(str, wordSeparator = '_')
```

###不带返回值的参数

如果函数不返回一个有用的值，就可以返回一个 `Unit` .`Unit` 不必有明确的返回

```kotlin
fun printHello(name: String?): Unit {
    if (name != null)
        println("Hello ${name}")
    else
        println("Hi there!")
    // `return Unit` or `return` is optional
}
```

`Unit` 返回值也可以省略，比如下面这样：

```kotlin
fun printHello(name: String?) {
    ...
}
```
###单表达是函数

当函数只返回一个表达式时，大括号可以省略并且函数体可以在 ＝ 后面只直接指定

```kotlin
fun double(x: Int): Int = x*2
```

###明确返回类型

下面的例子中必须有明确返回值：

>带表达式的函数体必须是 public 或 protected 。这些都被认为是公共接口的 API 。没有明确的函数返回值会使得不小心就会改变类型值。这就是为什么 [属性](http://kotlinlang.org/docs/reference/properties.html#getters-and-setters) 必须要有明确的的类型

>函数体有大括号的话就必须有明确的返回值，除非它想要返回 `Uint` 。Kotlin 不会对带大括号的函数做返回类型推断，因为这样的函数可能会很复杂。

###变长参数

函数最后一个参数可以用 vararg 注解标记：

```kotlin
fun asList<T>(vararg ts: T): List<T> {
	val result = ArrayList<T>()
	for (t in ts)
		result.add(t)
	return result
}
```

允许给函数传递可变长度的参数：

```kotlin
val list = asList(1, 2, 3)
```

只有一个参数可以注解为 `vararg` 。可以是最后一个参数，或倒数第二个，因为最后一个参数有可能是一个函数(允许在外面传递 lambda 表达式)

当调用变长参数的函数时，我们可以一个一个的传递参数，比如 `asList(1, 2, 3)`，或者我们要传递一个 array 的内容给函数，我们就可以使用 * 前缀操作符：

```kotlin
val a = array(1, 2, 3)
val list = asList(-1, 0, *a, 4)
```

###函数范围

Kotlin 中可以在文件个根级声明函数，这就意味者你不用创建一个类来持有函数。除了顶级函数，Kotlin 函数可以声明为局部的，作为成员函数或扩展函数。

####局部函数

Kotlin 支持局部函数，比如在另一个函数使用另一函数。

```kotlin
fun dfs(graohL Graph) {
	fun dfs(current: Vertex, vistied: Set<Vertex>) {
		if (!visited.add(current)) return 
		for (v in currnt.neighbors)
			dfs(v, visited)
	}
	dfs(graoh,vertices[0], HashSet())
}
```

局部函数可以访问外部函数的局部变量(比如闭包)

```kotlin
fun dfs(graph: Graph) {
	val visited = HashSet<Vertex>()
	fun dfs(current: Vertex) {
		if (!visited.add(current)) return 
		for (v in current.neighbors)
			dfs(v)
	}
	dfs(graph.vertices[0])
}
```

局部函数甚至可以返回到外部函数 [qualified return expressions](http://kotlinlang.org/docs/reference/returns.html)

```kotlin
fun reachable(from: Vertex, to: Vertex): Boolean {
	val visited = HashSet<Vertex>()
	fun dfs(current: Vertex) {
		if (current == to) return@reachable true
		if (!visited.add(current)) return
		for (v  in current.neighbors)
			dfs(v)
	}
	dfs(from)
	return false
}
```

###成员函数

成员函数是定义在一个类或对象里边的

```kotlin
class Sample() {
	fun foo() { print("Foo") }
}
```

成员函数可以用 . 的方式调用

```kotlin
Sample.foo()
```

更多请参看[类](http://kotlinlang.org/docs/reference/classes.html)和[继承](http://kotlinlang.org/docs/reference/classes.html#inheritance)

###泛型函数

函数可以有泛型参数，样式是在函数后跟上尖括号。

```kotlin
fun sigletonArray<T>(item: T): Array<T> {
	return Array<T>(1, {item})
}
```

更多请参看[泛型](http://kotlinlang.org/docs/reference/generics.html)

###内联函数

参看[这里](http://kotlinlang.org/docs/reference/inline-functions.html)

###扩展函数

参看[这里](http://kotlinlang.org/docs/reference/extensions.html)

###高阶函数和 lambda 表达式

参看[这里](http://kotlinlang.org/docs/reference/lambdas.html)

###尾递归函数

Kotlin 支持函数式编程的尾递归。这个允许一些算法可以通过循环而不是递归解决问题，从而避免了栈溢出。当函数被标记为 `tailRecursive` 时，编译器会优化递归，并用高效迅速的循环代替它。

```kotlin
tailRecursive fun findFixPoint(x: Double = 1.0): Double 
	= if (x == Math.cos(x)) x else findFixPoint(Math.cos(x))
```

这段代码计算的是数学上的余弦不动点。Math.cos 从 1.0  开始不断重复，直到值不变为止，结果是 0.7390851332151607 
这段代码和下面的是等效的：

```kotlin
private fun findFixPoint(): Double {
	var x = 1.0
	while (true) {
		val y = Math.cos(x)
		if ( x == y ) return y
		x = y
	}
}
```

使用 `tailRecursive` 注解必须在最后一个操作中掉用自己。在递归调用代码后面是不允许有其它代码的，并且也不可以用 try/catch/finall 块。当前的尾递归只在 JVM 的后端中可以用