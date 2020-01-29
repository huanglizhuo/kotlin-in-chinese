## 解构声明

有时将对象解构为多个变量会很方便，例如：
```Kotlin
val (name, age) = person
```

此语法称为解构声明。解构声明一次创建多个变量。我们已经声明了两个新变量：name和age，并且可以独立使用它们：

```Kotlin
println(name)
println(age)
```
解构声明被编译为以下代码：

```Kotlin
val name = person.component1()
val age = person.component2()
```

component1（）和component2（）函数是Kotlin中广泛使用的约定原理的另一个示例（请参阅+和*等运算符，for循环等）。只要可以在其上调用所需数量的 component 函数，任何内容都可以在解构声明的右侧。当然了可以有component3（）和component4（）等等。

请注意，componentN（）函数需要使用`operator`关键字进行标记，以允许在解构声明中使用它们。

解构声明也可以在for循环中使用,比如：

```Kotlin
for（（（a，b）in collection）{...}
```

变量a和b获得在集合的元素上调用的component1（）和component2（）返回的值。

## 示例：函数返回两个值

假设我们需要从一个函数返回两个值。例如，结果对象和某种状态。在Kotlin中执行此操作的一种简便方法是声明一个[数据类](https://kotlinlang.org/docs/reference/data-classes.html)并返回其实例：

```Kotlin
data class Result(val result: Int, val status: Status)
fun function(...): Result {
    // computations
    
    return Result(result, status)
}

// Now, to use this function:
val (result, status) = function(...)
```

由于数据类会自动声明componentN（）函数，因此在这里可以执行解构声明。

注意：我们也可以使用标准类`Pair`并让 function() 返回Pair <Int，Status>，但是通常最好对数据进行正确命名。

## 示例：解构声明和 map

遍历 map 最优的方法可能是这样:

```Kotlin
for ((key, value) in map) {
   // do something with the key and the value
}
```

为了使其工作，我们应该

- 提供 iterator（）函数，将 map 呈现为值序列；
- 提供函数 component1（）和component2（）将每个元素成对呈现。

实际上，标准库提供了以下扩展：

```Kotlin
operator fun <K, V> Map<K, V>.iterator(): Iterator<Map.Entry<K, V>> = entrySet().iterator()
operator fun <K, V> Map.Entry<K, V>.component1() = getKey()
operator fun <K, V> Map.Entry<K, V>.component2() = getValue()
```

因此，你可以在for循环中对 map 使用解构声明（当然也包括其它集合数据类）。

## 下划线表示未使用的变量（从1.1开始）

如果在解构声明中不需要某个变量，则可以使用下划线代替其名称：

```Kotlin
val (_, status) = getResult()
```

这种方式跳过的组件将不会调用 `componentN()` 运算符。

## Lambda 中的解构（自1.1开始）

可以对lambda参数使用解构声明语法。如果lambda具有 `Pair` 类型的参数（或Map.Entry或具有componentN函数的任何其他类型），则可以通过在括号中引入几个新参数来代替它：

```Kotlin
map.mapValues { entry -> "${entry.value}!" }
map.mapValues { (key, value) -> "$value!" }
```

注意声明两个参数和声明一个解构对而不是一个参数之间的区别：

```Kotlin
{ a -> ... } // one parameter
{ a, b -> ... } // two parameters
{ (a, b) -> ... } // a destructured pair
{ (a, b), c -> ... } // a destructured pair and another parameter
```

如果未使用已分解结构参数的组件，则可以将其替换为下划线，以避免重新为其取名：

```Kotlin
map.mapValues { (_, value) -> "$value!" }
```

可以为整个已解构参数指定类型或为特定组件指定类型：

```Kotlin
map.mapValues { (_, value): Map.Entry<Int, String> -> "$value!" }

map.mapValues { (_, value: String) -> "$value!" }
```
