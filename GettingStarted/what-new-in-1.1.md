# Kotlin 1.1 新特性
新特性列表

-协程
-其它语言特性
-标准库
-JVM 后端
-JavaScript 后端

## JavaScript
从 Kotlin 1.1 开始，JavaScript 支持不再是实验性的了。所有特性均支持，并为前端开发环境提过了大量新的工具。参看下文了解详细改变

## 协程
Kotlin 1.1关键的新特性就是协程，带来了像 `async/wait`,`yield` 这样的编程模式。Kotlin 设计的关键特性是所有协程都是由库实现的，而不是语言。所以你不需要与任何特定的编程范例或者并行库进行绑定。

协程是一个高效轻量级的线程，可以挂起并稍后恢复执行。协程是又挂起函数支持的：调用这个函数可能会挂起一个协程，并开启一个新的协程，大多数情况下采用匿名挂起函数（也就是可挂起lambda表达式）

让我们看一下在[kotlinx.coroutines](https://github.com/kotlin/kotlinx.coroutines)库中实现的`async`/`await`:

```
// 在后台线程池中运行代码fun asyncOverlay() = async(CommonPool) {    // 开启两个异步操作    val original = asyncLoadImage("original")    val overlay = asyncLoadImage("overlay")    // and then apply overlay to both results    applyOverlay(original.await(), overlay.await())}// 在 UI 上下文中中开启新的协程launch(UI) {    // wait for async overlay to complete    val image = asyncOverlay().await()    // and then show it in UI    showImage(image)}```

`async{...}` 开启协程，当调用 `await()` 时协程被挂起，等待 original 和 overlay 的执行结果，当两者完成时恢复执行。

标准库用协程支持 *lazily generated sequences* 懒生成序列，主要使用 `yield` 和 `yieldAll` 函数。在这样的队列中，返回序列会在每次元素被取出后暂停，并在下次取元素是恢复，例子：

```kotlin
val seq = buildSequence {
    for (i in 1..5) {
        // yield a square of i
        yield(i * i)
    }
    // yield a range
    yieldAll(26..28)
}

// print the sequence
println(seq.toList())
```

更多信息请参看 [coroutine documentation ](https://kotlinlang.org/docs/reference/coroutines.html) 以及 [tutorials](https://kotlinlang.org/docs/tutorials/coroutines-basic-jvm.html)

## 其它语言特性
### 类型别名
类型别名允许给现有的类型定义别名。主要在集合，函数类型中很常用。例子：

```kotlin
typealias OscarWinners = Map<String, String>

fun countLaLaLand(oscarWinners: OscarWinners) =
oscarWinners.count { it.value.contains("La La Land") }

// Note that the type names (initial and the type alias) are interchangeable:
fun checkLaLaLandIsTheBestMovie(oscarWinners: Map<String, String>) =
oscarWinners["Best picture"] == "La La Land"
```
参看[documentation](https://kotlinlang.org/docs/reference/type-aliases.html)[KEEP](https://github.com/Kotlin/KEEP/blob/master/proposals/type-aliases.md) 了解更多详细信息。

### Bound callable references（绑定可执行引用  暂时没想好怎么翻译，如果你有好的建议请发issue）
使用 `::` 操作符可以获取一个指向特定对象实例的方法或者熟悉的成员引用。之前这只能用在 lambda 表达式上。例子：

```Kotlin
val numberRegex = "\\d+".toRegex()
val numbers = listOf("abc", "123", "456").filter(numberRegex::matches)
```

参看[documentation](https://kotlinlang.org/docs/reference/reflection.html#bound-function-and-property-references-since-11)[KEEP](https://github.com/Kotlin/KEEP/blob/master/proposals/type-aliases.md) 了解更多详细信息。

## 密封类和数据类
Kotlin 1.1 移除了一些密封类和数据类的限制。现在能在同一个文件中定义顶级密封类类的子类，而不必是内嵌类或者密封类类。数据类现在可以扩展其它类。这样可以更优雅简洁的定义表达式类的等级：

```Kotlin
sealed class Expr

data class Const(val number: Double) : Expr()
data class Sum(val e1: Expr, val e2: Expr) : Expr()
object NotANumber : Expr()

fun eval(expr: Expr): Double = when (expr) {
    is Const -> expr.number
    is Sum -> eval(expr.e1) + eval(expr.e2)
    NotANumber -> Double.NaN
}
val e = eval(Sum(Const(1.0), Const(2.0)))
```
参看 [documentation](https://kotlinlang.org/docs/reference/sealed-classes.html#relaxed-rules-for-sealed-classes-since-11) or [sealed class](https://github.com/Kotlin/KEEP/blob/master/proposals/sealed-class-inheritance.md) 和 [data class](https://github.com/Kotlin/KEEP/blob/master/proposals/data-class-inheritance.md) KEEPs 获取更详细的信息。

## lambdas 表达式的解构
现在可以使用[destucting declaration](https://kotlinlang.org/docs/reference/multi-declarations.html) 语法取出 lambda 中的参数：

```kotlin
val map = mapOf(1 to "one", 2 to "two")
// before
println(map.mapValues { entry ->
                       val (key, value) = entry
                       "$key -> $value!"
                      })
// now
println(map.mapValues { (key, value) -> "$key -> $value!" })
```

参看[documentation](https://kotlinlang.org/docs/reference/multi-declarations.html#destructuring-in-lambdas-since-11) 和 [KEEP](https://github.com/Kotlin/KEEP/blob/master/proposals/destructuring-in-parameters.md) 获取更详细的信息。

## 用下划线标注下未使用的参数
对于有多个参数的 lambda ，你可以用 `_` 字符代替你不使用的参数。

```kotlin
map.forEach { _, value -> println("$value!") }
```

在[destructuring declarations](https://kotlinlang.org/docs/reference/multi-declarations.html)　中也可以用

```kotlin
val (_, status) = getResult()
```

阅读[KEEP](https://github.com/Kotlin/KEEP/blob/master/proposals/underscore-for-unused-parameters.md)获取更详细信息

## 数字值中的下划线
像 java8 一样 Kotlin 现在支持在数字值中使用下划线划分组：

```Kotlin
val oneMillion = 1_000_000
val hexBytes = 0xFF_EC_DE_5E
val bytes = 0b11010010_01101001_10010100_10010010
```

阅读[[KEEP](https://github.com/Kotlin/KEEP/blob/master/proposals/underscores-in-numeric-literals.md)获取更详细信息

## 属性简写
用表达式定义个 get 属性，现在可以省略属性类型了：

```kotlin
data class Person(val name: String, val age: Int){
	val isAdult get() = age >= 20 // 属性类型被推断为 'Boolean'
}
```

## 内联属性访问器
现在没有 backing field 的属性访问器可以用 'inline' 修饰。

这些属性访问器和内联函数的编译方式是一样的。

```kotlin
public val <T> List<T>.lastIndex: Int
inline get() = this.size -1
```

[documentation](http://kotlinlang.org/docs/reference/inline-functions.html#inline-properties-since-11)[KEEP](https://github.com/Kotlin/KEEP/blob/master/proposals/inline-properties.md) 了解更多详细信息。

##  本地代理属性
现在可以在本地变量上使用 [代理属性](http://kotlinlang.org/docs/reference/delegated-properties.html) 了。一个应用场景就是定义一个懒求值的本地变量：

```kotlin
val answer by lazy {
    println("Calculating the answer...")
    42
}
if (needAnswer()) {                     // returns the random value
    println("The answer is $answer.")   // answer is calculated at this point
}
else {
    println("Sometimes no answer is the answer...")
}
```

阅读 [KEEP](https://github.com/Kotlin/KEEP/blob/master/proposals/local-delegated-properties.md) 了解更多详细信息。

## 拦截代理属性绑定
代理属性现在可以通过 'provideDelegate' 运算符拦截属性绑定。比如，在绑定属性名之前，想要检查属性名字，可以像下面这样做：

```kotlin 
class ResourceLoader<T>(id: ResourceID<T>) {
    operator fun provideDelegate(thisRef: MyUI, prop: KProperty<*>): ReadOnlyProperty<MyUI, T> {
        checkProperty(thisRef, prop.name)
        ... // property creation
    }

    private fun checkProperty(thisRef: MyUI, name: String) { ... }
}

fun <T> bindResource(id: ResourceID<T>): ResourceLoader<T> { ... }

class MyUI {
    val image by bindResource(ResourceID.image_id)
    val text by bindResource(ResourceID.text_id)
}
```

'provideDelegate' 方法可以可以在创建 MyUI 实例的每个属性时调用，并做相应的验证检查。

## 通用枚举值的访问
现在可以用通用方法罗列枚举类的所有值。

```kotlin
enum class RGB { RED, GREEN, BLUE }

inline fun <reified T : Enum<T>> printAllValues() {
    print(enumValues<T>().joinToString { it.name })
}

## 标准库
###  string to number 的转换
在 String 类新增了转为数字而不抛出异常的扩展： String.toIntOrNull(): Int?, String.toDoubleOrNull(): Double? etc.

val port = System.getenv("PORT")?.toIntOrNull() ?: 80

### onEach()
'onEach' 是一个很小的扩展，但却对集合和序列很有用，这样就可以对集合和队列采用链式调用执行一些操作。

to be continue
