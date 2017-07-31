## 函数
### 函数声明
在 kotlin 中用关键字 `fun` 声明函数：

```kotlin
fun double(x: Int): Int {

}
```

### 函数用法
通过传统的方法调用函数

```kotlin
val result = double(2) 
```

通过`.`调用成员函数 

```kotlin
Sample().foo() // 创建Sample类的实例,调用foo方法
```

### 中缀符号
在满足以下条件时,函数也可以通过中缀符号进行调用:

>　它们是成员函数或者是[扩展函数](http://kotlinlang.org/docs/reference/extensions.html)
>　只有一个参数
>  使用`infix`关键词进行标记

```kotlin
//给 Int 定义一个扩展方法
infix fun Int.shl(x: Int): Int {
...
}

1 shl 2 //用中缀注解调用扩展函数

1.shl(2)
```

### 参数
函数参数是用 Pascal 符号定义的　name:type。参数之间用逗号隔开，每个参数必须指明类型。

```kotlin
fun powerOf(number: Int, exponent: Int) {
...
}
```

### 默认参数
函数参数可以设置默认值,当参数被忽略时会使用默认值。这样相比其他语言可以减少重载。

```kotlin
fun read(b: Array<Byte>, off: Int = 0, len: Int = b.size ) {
...
}
```

默认值可以通过在type类型后使用`=`号进行赋值

### 命名参数
在调用函数时可以参数可以命名。这对于那种有大量参数的函数是很方便的.

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

然而当调用非默认参数是就需要像下面这样：

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

注意,命名参数语法不能够被用于调用Java函数中,因为Java的字节码不能确保方法参数命名的不变性

### 不带返回值的参数
如果函数不会返回任何有用值，那么他的返回类型就是 `Unit` .`Unit` 是一个只有唯一值`Unit`的类型.这个值并不需要被直接返回:

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
### 单表达式函数
当函数只返回单个表达式时，大括号可以省略并在 = 后面定义函数体

```kotlin
fun double(x: Int): Int = x*2
```
在编译器可以推断出返回值类型的时候,返回值的类型可以省略:

```kotlin
fun double(x: Int) = x * 2

```

### 明确返回类型
下面的例子中必须有明确返回类型,除非他是返回 `Unit`类型的值,Kotlin 并不会对函数体重的返回类型进行推断,因为函数体中可能有复杂的控制流,他的返回类型未必对读者可见(甚至对编译器而言也有可能是不可见的)：

### 变长参数
函数的参数(通常是最后一个参数)可以用 vararg 修饰符进行标记：

```kotlin
fun <T> asList(vararg ts: T): List<T> {
	val result = ArrayList<T>()
	for (t in ts)
		result.add(t)
	return result
}
```

标记后,允许给函数传递可变长度的参数：

```kotlin
val list = asList(1, 2, 3)
```

只有一个参数可以被标注为 `vararg` 。加入`vararg`并不是列表中的最后一个参数,那么后面的参数需要通过命名参数语法进行传值,再或者如果这个参数是函数类型,就需要通过lambda法则.

当调用变长参数的函数时，我们可以一个一个的传递参数，比如 `asList(1, 2, 3)`，或者我们要传递一个 array 的内容给函数，我们就可以使用 * 前缀操作符：

```kotlin
val a = array(1, 2, 3)
val list = asList(-1, 0, *a, 4)
```

### 函数范围
Kotlin 中可以在文件顶级声明函数，这就意味者你不用像在Java,C#或是Scala一样创建一个类来持有函数。除了顶级函数，Kotlin 函数可以声明为局部的，作为成员函数或扩展函数。

#### 局部函数
Kotlin 支持局部函数，比如在一个函数包含另一函数。

```kotlin
fun dfs(graph: Graph) {
  fun dfs(current: Vertex, visited: Set<Vertex>) {
    if (!visited.add(current)) return
    for (v in current.neighbors)
      dfs(v, visited)
  }

  dfs(graph.vertices[0], HashSet())
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

### 成员函数
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

### 泛型函数
函数可以有泛型参数，样式是在函数名前加上尖括号。

```kotlin
fun <T> sigletonArray(item: T): Array<T> {
	return Array<T>(1, {item})
}
```

更多请参看[泛型](http://kotlinlang.org/docs/reference/generics.html)

### 内联函数
参看[这里](http://kotlinlang.org/docs/reference/inline-functions.html)

### 扩展函数
参看[这里](http://kotlinlang.org/docs/reference/extensions.html)

### 高阶函数和 lambda 表达式
参看[这里](http://kotlinlang.org/docs/reference/lambdas.html)

### 尾递归函数
Kotlin 支持函数式编程的尾递归。这个允许一些算法可以通过循环而不是递归解决问题，从而避免了栈溢出。当函数被标记为 `tailrec` 时，编译器会优化递归，并用高效迅速的循环代替它。

```kotlin
tailrec fun findFixPoint(x: Double = 1.0): Double 
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

使用 `tailrec` 修饰符必须在最后一个操作中调用自己。在递归调用代码后面是不允许有其它代码的，并且也不可以在 try/catch/finall 块中进行使用。当前的尾递归只在 JVM 的后端中可以用
