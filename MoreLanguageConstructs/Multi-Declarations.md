## 多重声明
有时候可以通过给对象插入多个成员函数做区别是很方便的，比如：

```kotlin
val (name, age) = person
```

这种语法叫多重声明。多重声明一次创建了多个变量。我们声明了俩个新变量：`name` `age` 并且可以独立使用：

```kotlin
println(name)
println(age)
```

多重声明被编译成下面的代码：

```kotlin
val name = persion.component1()
val age = persion.component2()
```

`component1()` `component2()`是另一个转换原则的例子。任何类型都可以在多重分配的右边。当然了，也可以有 `component3()` `component4()` 等等

多重声明也可以在 for 循环中用

```kotlin
for ((a, b) in collection) { ... }
```

参数 `a` 和 `b` 是 `component1()`　`component2()` 的返回值

### 例子：一个函数返回俩个值
要是一个函数想返回俩个值。比如，一个对象结果，一个是排序的状态。在 Kotlin 中的一个紧凑的方案是声明 [data](http://kotlinlang.org/docs/reference/data-classes.html) 类并返回实例：

```kotlin
data class Result(val result: Int, val status: Status)

fun function(...): Result {
	//...
	return Result(result, status)
}

val (result, status) = function(...)
```
数据类自动声明 `componentN()` 函数

注意：也可以使用标准类 `Pair` 并让函数返回 'Pair<Int , staus>'，但可读性不是很强

### 例子：多重声明和 Map
转换 map 的最好办法可能是下面这样：

```kotlin
for ((key, value) in map) {

}
```

为了让这个可以工作，我们需要

>通过提供 `iterator()` 函数序列化呈现 map
>通过 `component1()`和 `component1()` 函数是把元素成对呈现

事实上，标准库提供了这样的扩展：

```kotlin
fun <K, V> Map<K, V>.iterator(): Iterator<Map.Entry<K, V>> = entrySet().iterator()
fun <K, V> Map.Entry<K, V>.component1() = getKey()
fun <K, V> Map.Entry<K, V>.component2() = getValue()
```

因此你可以用 for 循环方便的读取 map (或者其它数据集合)
