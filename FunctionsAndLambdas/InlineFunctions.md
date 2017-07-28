## 内联函数
使用[高阶函数](http://kotlinlang.org/docs/reference/lambdas.html)带来了相应的运行时麻烦：每个函数都是一个对象，它捕获闭包，即这些变量可以在函数体内被访问。内存的分配，虚拟调用的运行都会带来开销

但在大多数这种开销是可以通过内联文本函数避免。下面就是一个很好的例子。`lock()` 函数可以很容易的在内联点调用。思考一下下面的例子：

```kotlin
lock(i) { foo() }
```

(Instead of creating a function object for the parameter and generating a call)，编译器可以忽略下面的代码：

```kotlin
lock.lock()
try {
	foo()
}
finally {
	lock.lock()
}
```

这好像不是我们开始想要的

想要让编译器不这样做的话，我们需要用 `inline` 标记 `lock()` 函数：
```kotlin
inline fun lock<T>(lock: Lock,body: ()-> T): T {
	//...
}
```

`inline` 标记即影响函数本身也影响传递进来的 lambda 函数：所有的这些都将被关联到调用点。

内联可能会引起生成代码增长，但我们可以合理的解决它(不要内联太大的函数)

### @noinline
如果只需要在内联函数中内联部分Lambda表达式，可以使用`@noinline` 注解来标记不需要内联的参数：

```kotlin
inline fun foo(inlined: () -> Uint, @noinline notInlined: () -> Unit) {
	//...
}
```

内联的 lambda 只能在内联函数中调用，或者作为内联参数，但 `@noinline` 标记的可以通过任何我们喜欢的方式操控：存储在字段，( passed around etc)

注意如果内联函数没有内联的函数参数并且没有具体类型的参数，编译器会报警告，这样内联函数就没有什么优点的(如果你认为内联是必须的你可以忽略警告)

### 返回到非局部
在 kotlin 中，我们可以不加条件的使用 `return` 去退出一个命名函数或表达式函数。这意味这退出一个 lambda 函数，我们不得不使用[标签](http://kotlinlang.org/docs/reference/returns.html#return-at-labels)，而且空白的 `return` 在 lambda 函数中是禁止的，因为 lambda 函数不可以造一个闭合函数返回：

```kotlin
fun foo() {
	ordinaryFunction {
		return // 错误　不可以在这返回
	}
}
```

但如果 lambda 函数是内联传递的，则返回也是可以内联的，因此允许下面这样：

```kotlin
fun foo() {
	inlineFunction {
		return //
	]
}
```

注意有些内联函数可以调用传递进来的 lambda 函数，但不是在函数体，而是在另一个执行的上下文中，比如局部对象或者一个嵌套函数。在这样的情形中，非局部的控制流也不允许在lambda 函数中。为了表明，lambda 参数需要有 `InlineOptions.ONLY_LOCAL_RETURN` 注解：

```kotlin
inline fun f(inlineOptions(InlineOption.ONLY_LOCAL_RETURN) body: () -> Unit) {
    val f = object: Runnable {
        override fun run() = body()
    }
    // ...
}
```

内联 lambda 不允许用 break 或 continue ，但在以后的版本可能会支持。

### 实例化参数类型
有时候我们需要访问传递过来的类型作为参数：

```kotlin
fun <T> TreeNode.findParentOfType(clazz: Class<T>): T? {
	var p = parent
	while (p != null && !clazz.isInstance(p)) {
		p = p?.parent
	}
	@suppress("UNCHECKED_CAST")
	return p as T
}
```

现在，我们创立了一颗树，并用反射检查它是否是某个特定类型。一切看起来很好，但调用点就很繁琐了：

```kotlin
myTree.findParentOfType(javaClass<MyTreeNodeType>() )
```
我们想要的仅仅是给这个函数传递一个类型，即像下面这样：

```kotlin
myTree.findParentOfType<MyTreeNodeType>()

```

为了达到这个目的，内联函数支持具体化的类型参数，因此我们可以写成这样：


```kotlin
inline fun <reified T> TreeNode.findParentOfType(): T? {
	var p = parent
	while (p != null && p !is T) {
		p = p?.parent
	}
	return p as T
}
```

我们用 refied 修饰符检查类型参数，既然它可以在函数内部访问了，也就基本上接近普通函数了。因为函数是内联的，所以不许要反射，像 `!is` ｀as｀这样的操作都可以使用。同时，我们也可以像上面那样调用它了 `myTree.findParentOfType<MyTreeNodeType>()`

尽管在很多情况下会使用反射，我们仍然可以使用实例化的类型参数 `javaClass()` 来访问它：

```kotlin
inline fun methodsOf<reified T>() = javaClass<T>().getMethods()

fun main(s: Array<String>) {
	println(methodsOf<String>().joinToString('\n'))
}
```

普通的函数(没有标记为内联的)不能有实例化参数。

更底层的解释请看[spec document](https://github.com/JetBrains/kotlin/blob/master/spec-docs/reified-type-parameters.md)
