##反射

反射是一系列语言和库的特性，允许在运行是获取你代码结构。 Kotlin 把函数和属性作为语言的头等类，而且反射它们和使用函数式编程或反应是编程风格很像。

>On the Java platform, the runtime component required for using the reflection features is distributed as a separate JAR file (kotlin-reflect.jar). This is done to reduce the required size of the runtime library for applications that do not use reflection features. If you do use reflection, please make sure that the .jar file is added to the classpath of your project.

###类引用

最基本的反射特性就是得到运行时的类引用。要获取引用并使之成为静态类可以使用字面类语法：

```kotlin
val c = MyClass::class
```

引用是一种 [KClass](http://kotlinlang.org/api/latest/jvm/stdlib/kotlin.reflect/-k-class/index.html)类型的值。你可以使用 `KClass.properties` 和 `KClass.extensionProperties` 获取类和父类的所有属性引用的列表。

注意这与 java 类的引用是不一样的。参看 [java interop section](http://kotlinlang.org/docs/reference/java-interop.html#object-methods)


###函数引用

当有一个像下面这样的函数声明时：

```kotlin
fun isOdd(x: Int) =x % 2 !=0
```

我们可以通过 `isOdd(5)` 轻松调用，同样我们也可以把它作为一个值传递给其它函数。我们可以使用 `::` 操作符

```kotlin
val numbers = listOf(1, 2, 3)
println(numbers.filter( ::isOdd) ) //prints [1, 3]
```

这里 `::isOdd` 是是一个函数类型的值 `(Int) -> Boolean`

注意现在 `::` 操作符右边不能用语重载函数。将来，我们计划提供一个语法明确参数类型这样就可以使用明确的重载函数了。

如果需要使用一系列类，或者扩展函数，必须是需合格的，并且结果是扩展函数类型，比如。`String::toCharArray` 就带来一个 `String: String.() -> CharArray` 类型的扩展函数。

###例子：函数组合

考虑一下下面的函数：

```kotlin
fun compose<A, B, C>(f: (B) -> C, g: (A) -> B): (A) -> C {
    return {x -> f(g(x))}
}
```

它返回一个由俩个传递进去的函数的组合。现在你可以把它用在可调用的引用上了：

```kotlin
fun length(s: String) = s.size
val oddLength = compose(::isOdd, ::length)
val strings = listOf("a", "ab", "abc")

println(strings.filter(oddLength)) // Prints "[a, abc]"
```

###属性引用

在 kotlin 中访问顶级类的属性，我们也可以使用 `::` 操作符：

```kotlin
var x = 1
fun main(args: Array<String>) {
	println(::x.get())
	::x.set(2)
	println(x)
}
```

`::x` 表达式评估为 `KProperty<Int>` 类型的属性，它允许我们使用 `get()` 读它的值或者使用名字取回它的属性。更多请参看[docs on the KProperty class](http://kotlinlang.org/api/latest/jvm/stdlib/kotlin.reflect/-k-property.html)

对于可变的属性比如 `var y =1`,`::y`返回类型为 `[KMutableProperty<Int>](http://kotlinlang.org/api/latest/jvm/stdlib/kotlin.reflect/-k-mutable-property.html)`，它有 `set()` 方法

访问一个类的属性成员，我们这样修饰：

```kotlin
class A(val p: Int)

fun main(args: Array<String>) {
    val prop = A::p
    println(prop.get(A(1))) // prints "1"
}
```

对于扩展属性：

```kotlin
val String.lastChar: Char
  get() = this[size - 1]

fun main(args: Array<String>) {
  println(String::lastChar.get("abc")) // prints "c"
}
```

###与 java 反射调用

在 java 平台上，标准库包括反射类的扩展，提供了到 java 反射对象的映射(参看 kotlin.reflect.jvm 包)。比如，想找到一个备用字段或者 java getter 方法，你可以这样写：

```kotlin
import kotlin.reflect.jvm.*

class A(val p: Int)

fun main(args: Array<String>) {
    println(A::p.javaGetter) // prints "public final int A.getP()"
    println(A::p.javaField)  // prints "private final int A.p"
}
```

###构造函数引用

构造函数可以像方法或属性那样引用。只需要使用 `::` 操作符并加上类名。下面的函数是一个没有参数并且返回类型是 `Foo`:

```kotlin
calss Foo
fun function(factory : () -> Foo) {
	val x: Foo = factory()
}
```

我们可以像下面这样使用：

```kotlin
function(:: Foo)
```
