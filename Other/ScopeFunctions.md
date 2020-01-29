## 作用域函数

Kotlin标准库包含几个函数，其唯一目的是在指定对象的上下文中执行代码块。当在带有 lambda 表达式的对象上调用此类函数时，它将形成一个临时作用域。在此作用域内，可以访问不带名称的对象。这些功能称为作用域函数。其中有五个：`let`，`run`，`with`，`apply` 以及 `also`。

基本上，这些功能执行相同的操作：在 对应对象上执行代码块。不同之处在于此对象如何在块内可用以及整个表达式的返回值是什么。

下面是作用域函数的典型用法：

```Kotlin
Person("Alice", 20, "Amsterdam").let {
    println(it)
    it.moveTo("London")
    it.incrementAge()
    println(it)
}
```

如果不同 let 实现相同功能,则必须引入一个新变量，并在每次使用它时重复其名称。

```Kotlin
val alice = Person("Alice", 20, "Amsterdam")
println(alice)
alice.moveTo("London")
alice.incrementAge()
println(alice)
```

作用域函数没有引入任何新的技术能力，但是它们可以使您的代码更加简洁和易读。

由于作用域函数的相似性质，选择合适的作用域函数可能会有些棘手。主要取决于你的意图和项目中使用的一致性。下面，我们将详细描述作用域函数之间的区别及其用法约定。

## 区别

因为作用域函数本质上都非常相似，所以了解它们之间的差异很重要。每个作用域函数之间有两个主要区别：

- 引用上下文对象的方式
- 返回值

### 上下文对象：this或it

在作用域函数的 lambda 中，可以通过短引用而不是其实际名称来使用上下文对象。每个作用域函数使用两种访问上下文对象的方式之一：作为lambda接收器（this）或作为lambda自变量（it）。两者都提供相同的功能，因此我们将描述每种情况在每种情况下的优缺点，并提供有关其用法的建议。

```Kotlin
fun main() {
    val str = "Hello"
    // this
    str.run {
        println("The receiver string length: $length")
        //println("The receiver string length: ${this.length}") // does the same
    }

    // it
    str.let {
        println("The receiver string's length is ${it.length}")
    }
}
```

**this**

`run`，`with` 以及 `apply` 将上下文对象称为lambda接收器-即关键字 `this` 。因此，在其lambda中，该对象是可用的，就像在普通类函数中一样。在大多数情况下，访问接收器对象的成员时可以省略此代码，从而缩短代码。另一方面，如果省略，则很难区分接收器构件和外部对象或功能。因此，对于主要对对象成员进行操作的lambda，建议将上下文对象作为接收器（即`this`）：调用其函数或分配属性。

```Kotlin
val adam = Person("Adam").apply { 
    age = 20                       // same as this.age = 20 or adam.age = 20
    city = "London"
}
```

**it**

反过来，`let` 和 `also` 则将上下文对象作为lambda参数。如果未指定参数名称，则使用隐式默认名称 `it` 来访问该对象。`it` 比 `this` 更短，并且带有 `it` 的表达式通常更易于阅读。但是，在调用对象函数或属性时，没有像`this`那样隐式可用的对象。因此，将上下文对象作为`it`时,更适合当函数调用需要参数时, 如果在代码块中使用多个变量时,`it`也是更好的选择。

```Kotlin
fun getRandomInt(): Int {
    return Random.nextInt(100).also {
        writeToLog("getRandomInt() generated value $it")
    }
}

val i = getRandomInt()
```

此外，当您将上下文对象作为参数传递时，可以在作用域内为上下文对象提供自定义名称。

```Kotlin
fun getRandomInt(): Int {
    return Random.nextInt(100).also { value ->
        writeToLog("getRandomInt() generated value $value")
    }
}

val i = getRandomInt()
```

### 返回值

从返回结果上区分作用域函数如下:

- `apply` `also` 返回上下文对象
- `let` `run` `with` 返回 lambda 结果

这两个选项使您可以根据下一步在代码中的选择来选择适当的功能。

**上下文对象**

`apply` 和 `also` 返回值就是上下文对象本身。因此，它们可以作为副步骤包含在调用链中：可以在它们之后继续在同一对象上链接函数调用。

```Kotlin
val numberList = mutableListOf<Double>()
numberList.also { println("Populating the list") }
    .apply {
        add(2.71)
        add(3.14)
        add(1.0)
    }
    .also { println("Sorting the list") }
    .sort()
```

它们还可以用于返回上下文对象的函数的return语句中。

```Kotlin
fun getRandomInt(): Int {
    return Random.nextInt(100).also {
        writeToLog("getRandomInt() generated value $it")
    }
}

val i = getRandomInt()
```

**Lambda结果**

`let` `run` `with` 返回lambda结果。因此，在将结果分配给变量，对结果进行链接操作等时，可以使用它们。

```Kotlin
val numbers = mutableListOf("one", "two", "three")
val countEndsWithE = numbers.run { 
    add("four")
    add("five")
    count { it.endsWith("e") }
}
println("There are $countEndsWithE elements that end with e.")
```

此外，您可以忽略返回值，并使用范围函数为变量创建临时范围。

```Kotlin
val numbers = mutableListOf("one", "two", "three")
with(numbers) {
    val firstItem = first()
    val lastItem = last()        
    println("First item: $firstItem, last item: $lastItem")
}
```

### 函数

为了帮助你选择合适的作用域函数，我们将详细描述它们并提供使用建议。从技术上讲，函数在许多情况下是可以互换的，因此示例提供了定义常用用法样式的约定。

**let**

上下文对象可通过作参数 `it` 访问。返回值是lambda结果。

`let` 可以用于在调用链的结果上调用一个或多个函数。例如，以下代码在集合上打印两个操作的结果：

```Kotlin
val numbers = mutableListOf("one", "two", "three", "four", "five")
val resultList = numbers.map { it.length }.filter { it > 3 }
println(resultList)   
```

使用let，可以将其重写为：

```Kotlin
val numbers = mutableListOf("one", "two", "three", "four", "five")
numbers.map { it.length }.filter { it > 3 }.let { 
    println(it)
    // and more function calls if needed
} 
```

如果代码块包含一个函数作为参数，则可以使用方法引用（::)代替lambda：

```Kotlin
val numbers = mutableListOf("one", "two", "three", "four", "five")
numbers.map { it.length }.filter { it > 3 }.let(::println)
```

`let` 通常用于仅使用非空值执行代码块。要对非空对象执行操作，请使用安全调用运算符`?`对其进行调用，并在lambda 中用`it`操作。

```Kotlin
val str: String? = "Hello"   
//processNonNullString(str)       // compilation error: str can be null
val length = str?.let { 
    println("let() called on $it")        
    processNonNullString(it)      // OK: 'it' is not null inside '?.let { }'
    it.length
}
```

使用let的另一种情况是引入范围有限的局部变量以提高代码的可读性。要为上下文对象定义一个新变量，请提供其名称作为lambda参数，以便使用而不是默认值`it`。

```Kotlin
val numbers = listOf("one", "two", "three", "four")
val modifiedFirstItem = numbers.first().let { firstItem ->
    println("The first item of the list is '$firstItem'")
    if (firstItem.length >= 5) firstItem else "!" + firstItem + "!"
}.toUpperCase()
println("First item after modifications: '$modifiedFirstItem'")
```

**with**

一个非扩展函数：上下文对象作为参数传递，但是在lambda内部，它可用作接收器（this）使用。返回值是lambda结果。

我们建议在不需要 lambda 结果的情况下在上下文对象上调用函数。在代码中，with可以理解为“使用此对象，请执行以下操作”。

```Kotlin
val numbers = mutableListOf("one", "two", "three")
with(numbers) {
    println("'with' is called with argument $this")
    println("It contains $size elements")
}
```

`with` 的另一个用例是引入一个辅助对象，该对象的属性或函数将用于计算值。

```Kotlin
val numbers = mutableListOf("one", "two", "three")
val firstAndLast = with(numbers) {
    "The first element is ${first()}," +
    " the last element is ${last()}"
}
println(firstAndLast)
```

**run**

上下文对象可用作接收者（this）使用。返回值是lambda结果。

`run` 与 `with` 相同，作为 `let` 调用 - 上下文对象的扩展函数 (run does the same as with but invokes as let - as an extension function of the context object.)

当lambda同时包含对象初始化和返回值的计算时，`run` 很有用。

```Kotlin
val service = MultiportService("https://example.kotlinlang.org", 80)

val result = service.run {
    port = 8080
    query(prepareRequest() + " to port $port")
}

// the same code written with let() function:
val letResult = service.let {
    it.port = 8080
    it.query(it.prepareRequest() + " to port ${it.port}")
}
```

除了在接收者对象上调用`run`之外，还可以将其用作非扩展函数。非扩展`run`可以在需要表达式的地方执行包含多个语句的块。

```Kotlin
val hexNumberRegex = run {
    val digits = "0-9"
    val hexDigits = "A-Fa-f"
    val sign = "+-"

    Regex("[$sign]?[$digits$hexDigits]+")
}

for (match in hexNumberRegex.findAll("+1234 -FFFF not-a-number")) {
    println(match.value)
}
```

`apply`

上下文对象可用作接收者（this）。返回值是对象本身。

`apply` 适合用于不返回值且主要在接收者对象的成员上运行的代码块。适用的常见情况是对象配置。此类调用可以理解为“将以下赋值操作应用于对象”。

```Kotlin
val adam = Person("Adam").apply {
    age = 32
    city = "London"        
}
```

将接收者作为返回值，可以轻松地将`apply`应用于调用链以进行更复杂的处理。

**also**

上下文对象可用作参数（`it`）。返回值是对象本身。

`also` 有助于执行一些将上下文对象作为参数的操作。`also` 用于不改变对象的其他操作，例如日志或打印调试信息。通常可以在不打破程序逻辑的情况下从调用链中将 `also` 调用删除。

当在代码中看到时 `also` 时，可以将其理解为“并且还可以执行以下操作”。

```Kotlin
val numbers = mutableListOf("one", "two", "three")
numbers
    .also { println("The list elements before adding new one: $it") }
    .add("four")
```

## 选择函数

为了帮助选择合适的作用域函数，我们提供了它们之间的主要区别表。

| 函数 | 对象应用 | 返回值 | 是否是扩展函数 |
| ---- | ---- | ---- | ---- |
| let | it | lambda 结果 | 是 |
| run | this | lambda 结果 | 是 |
| run | _ | lambda 结果 | 否(不需要上下文对象进行调用) |
| with | this | lambda 结果 | 否(需要上下文对象作为参数) |
| apply | this | 上下文对象 | 是 |
| also | it | 上下文对象 | 是 |

以下是根据预期目的选择作用域函数的简短指南：

- 在非空对象上执行lambda：let
- 将表达式引入为局部作用域中的变量：let
- 对象配置：apply
- 对象配置和计算结果：run
- 需要表达式的运行语句：非扩展 run
- 附加效果：also
- 对对象进行成组的函数调用：with

不同功能的用例重叠，因此您可以根据项目或团队中使用的特定约定选择功能。

尽管作用域函数是使代码更简洁的一种方法，但请避免过度使用它们：这会降低代码的可读性并导致错误。 避免嵌套作用域函数，并在链接它们时要小心：很容易混淆当前上下文对象及其值。

## takeIf 和 takeUnless

除了范围函数外，标准库还包含函数 `takeIf` 和 `takeUnless` 。这些功能可以将对对象状态的检查嵌入到调用链中。

在提供谓词的对象上调用时，takeIf 返回与谓词匹配的对象。否则，它返回null。因此，takeIf 是单个对象的过滤功能。反过来，如果takeUnless与谓词不匹配，则返回该对象；如果与谓词不匹配，则返回null。该对象可用作lambda参数（it）。

```Kotlin
val number = Random.nextInt(100)

val evenOrNull = number.takeIf { it % 2 == 0 }
val oddOrNull = number.takeUnless { it % 2 == 0 }
println("even: $evenOrNull, odd: $oddOrNull")
```

在takeIf和takeUnless之后链接其他函数时，不要忘记执行空检查或安全调用（?.），因为它们的返回值是可为空的。

```Kotlin
val str = "Hello"
val caps = str.takeIf { it.isNotEmpty() }?.toUpperCase()
//val caps = str.takeIf { it.isNotEmpty() }.toUpperCase() //compilation error
println(caps)
```

takeIf和takeUnless与作用域函数结合很有用。一个很好的例子是让它们链接在一起，以便在与给定谓词匹配的对象上运行代码块。为此，请在对象上调用takeIf，然后使用安全调用（？）调用let。对于与谓词不匹配的对象，takeIf返回null且不调用let。

```Kotlin
fun displaySubstringPosition(input: String, sub: String) {
    input.indexOf(sub).takeIf { it >= 0 }?.let {
        println("The substring $sub is found in $input.")
        println("Its start position is $it.")
    }
}

displaySubstringPosition("010000011", "11")
displaySubstringPosition("010000011", "12")
```

在不适用标准库函数时相同功能的实现如下:

```Kotlin
fun displaySubstringPosition(input: String, sub: String) {
    val index = input.indexOf(sub)
    if (index >= 0) {
        println("The substring $sub is found in $input.")
        println("Its start position is $index.")
    }
}

displaySubstringPosition("010000011", "11")
displaySubstringPosition("010000011", "12")
```