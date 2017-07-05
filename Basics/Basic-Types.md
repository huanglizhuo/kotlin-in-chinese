## 基本类型
在 Kotlin 中，所有的东西都是对象，这就意味着我们可以调用任何变量的成员函数和属性。一些类型是内建的，它们的实现是优化过的，但对用户来说它们就像普通的类一样。在这节中，我们将会讲到大多数的类型：数值，字符，布尔，以及数组。

### 数值
Kotlin 处理数值的方法和 java 很相似，但不是完全一样。比如，不存在隐式转换数值的精度，并且在字面上有一些小小的不同。

Kotlin 提供了如下内建数值类型(和 java 很相似)：


| **类型** | **位宽** |
| ------ | :----: |
| Double |   64   |
| Float  |   32   |
| Long   |   64   |
| Int    |   32   |
| Short  |   16   |
| Byte   |   8    |

注意字符在 Kotlin 中不是数值类型

### 字面值常量
主要是以下几种字面值常量：

> --十进制数值: `123`  
> 	--长整型要加大写 `L` : `123L`  
> --16进制：`0x0f`  
> --二进制：`0b00001011`

注意不支持８进制

Kotlin 也支持传统的浮点数表示：

> -- 默认双精度浮点数(Double) : `123.5` , `123.5e10`  
> -- 单精度浮点数(Float)要添加 `f` 或 `F` ：123.5f

### 数值常量中可以添加下划线分割(1.1版本新特性)
您可以使用下划线增加数值常量的可读性:

```kotlin
val oneMillion = 1_000_000
val creditCardNumber = 1234_5678_9012_3456L
val socialSecurityNumber = 999_99_9999L
val hexBytes = 0xFF_EC_DE_5E
val bytes = 0b11010010_01101001_10010100_10010010
```

### 表示
在 java 平台上，数值被 JVM 虚拟机以字节码的方式物理存储的，除非我们需要做可空标识(比如说 Int?) 或者涉及泛型。在后者中数值是装箱过的。

注意装箱过的数值是不保留特征的：

```kotlin
val a: Int = 10000
print (a === a ) // 打印 'true'
val boxedA: Int? =a
val anotherBoxedA: Int? = a
print (boxedA === anotherBoxedA ) // 注意这里打印的是 'false'
```

然而，它们是值相等的：

```kotlin
val a: Int = 10000
print(a == a) // 打印 'true'
val boxedA: Int? = a
val anotherBoxedA: Int? = a
print(boxedA == anotherBoxedA) // 打印 'true'
```

### 显式转换
由于不同的表示，低精度类型不是高精度类型的子类型。如果是的话我们就会碰到下面这样的麻烦了

```kotlin
// 这是些伪代码，不能编译的
val a: Int? =1 // 一个装箱过的 Int (java.lang.Integer)
val b: Long? = a // 一个隐式装箱的 Long (java.lang.Long)
print( a == b )// 很惊讶吧　这次打印出的是 'false' 这是由于 Long 类型的 equals() 只有和 Long 比较才会相同
```

因此不止是特征会丢失，有时候连值相等都会悄悄实效。

所以，低精度类型是不会隐式转换为高精度类型的。这意味着我们必须显式转换才能把 `Byte` 赋值给 `Int`

```kotlin
val b: Byte = 1 // OK, 字面值常量会被静态检查
val i: Int = b // ERROR
```

我们可以通过显式转换把数值类型提升

```kotlin
val i: Int = b.toInt() // 显式转换
```

每个数值类型都支持下面的转换：

> ` toByte(): Byte`

>  `toShort(): Short`

>  ` toInt(): Int`

> ` toLong(): Long`

> ` toFloat(): Float`

> ` toDouble(): Double`

> ` toChar(): Char`

隐式转换一般情况下是不容易被发觉的，因为我们使用了上下文推断出类型，并且算术运算会为合适的转换进行重载，比如

```kotlin
val l = 1.toLong + 1 // Long  + Int => Long
```

### 运算符
Kotlin支持标准的算术运算表达式，这些运算符被声明为相应类的成员(但是编译器将调用优化到相应的指令)。参看[运算符重载](http://kotlinlang.org/docs/reference/operator-overloading.html)。

至于位运算，Kotlin 并没有提供特殊的操作符，只是提供了命名函数，可以采用中缀形式调用，比如：

```kotlin
val x = (1 shl 2) and 0x000FF000
```

下面是全部的位运算操作符(只可以用在 `Int` 和 `Long` 类型)：

> `shl(bits)` – 有符号左移 (相当于 Java’s `<<`)
> `shr(bits)` – 有符号右移 (相当于 Java’s `>>`)
> `ushr(bits)` – 无符号右移 (相当于 Java’s `>>>`)
> `and(bits)` – 按位与
> `or(bits)` – 按位或
> `xor(bits)` – 按位异或
> `inv(bits)` – 按位翻转

### 字符
字符类型用 `Char` 表示。不能直接当做数值来使用

```Kotlin
fun check(c: Char) {
	if (c == 1) { // ERROR: 类型不匹配
		// ...
	}
}
```
字符是由单引号包裹的：'1'，特殊的字符通过反斜杠\\转义，下面的字符序列支持转义：`\t`,`\b`,`\n`,`\r`,`\'`,`\"`,`\\`和`\$`。编码任何其他字符，使用 Unicode 转义语法：`\uFF00`。

我们可以将字符显示的转义为Int数字：

```kotlin
fun decimalDigitValue(c: Char): Int {
	if (c !in '0'..'9')
		throw IllegalArgumentException("Out of range")
	return c.toInt() - '0'.toInt() //显示转换为数值类型
}
```
和数值类型一样，需要一个可空引用时，字符会被装箱。特性不会被装箱保留。

### 布尔值
布尔值只有 true 或者 false

如果需要一个可空引用，将会对布尔值装箱

布尔值的内建操作包括

> `||` – 短路或  
>
> `&&` – 短路与  
>
> `!` - 取反

### 数组
数组在 Kotlin 中由 `Array` 类表示，有 `get` 和 `set` （通过运算符重载为`[]` ）方法，和 `size` 属性，以及一些常用的函数：

```kotlin
class Array<T> private constructor() {
	val size: Int
	operator fun get(index: Int): T
  operator fun set(index: Int, value: T): Unit

  operator fun iterator(): Iterator<T>
  // ...
}
```

给库函数 `arrayOf()` 传递每一项的值来创建Array，`arrayOf(1, 2, 3)` 创建了一个[1, 2, 3] 这样的数组。也可以使用库函数 `arrayOfNulls()` 创建一个指定大小的空Array。

另一种方式就是使用工厂函数，接受一个数组大小参数以及一个可以根据给定索引创建初始值的函数：

```kotlin
// 创建一个 Array<String>  内容为 ["0", "1", "4", "9", "16"]
val asc = Array(5, {i -> (i * i).toString() })
```

像我们上面提到的，`[]` 操作符表示调用　`get()` `set()` 函数

注意：和 java 不一样，arrays 在 kotlin 中是不可变的。这意味这 kotlin 不允许我们把 `Array<String>` 转为 `Array<Any>` ,这样就阻止了可能的运行时错误(但你可以使用 `Array<outAny>` , 参看 [Type Projections](http://kotlinlang.org/docs/reference/generics.html#type-projections))

Kotlin 有专门的类来表示原始类型从而避免过度装箱： ByteArray, ShortArray, IntArray 等等。这些类与 Array 没有继承关系，但它们有一样的方法与属性。每个都有对应的库函数：

```kotlin
val x: IntArray = intArrayOf(1, 2, 3)
x[0] = x[1] + x[2]
```

### 字符串
字符串是由 `String` 表示的。字符串是不变的。字符串的元素可以通过索引操作读取: `s[i]` 。字符串可以用 for 循环迭代：

```kotlin
for (c in str) {
	println(c)
}
```

#### 字符串字面值
Kotlin 有两种类型的字符串字面值：一种是可以带转义符的，一种是可以包含新行以及任意文本的。带转义符的 string 很像 java 的 string:

```kotlin
val s = "Hello World!\n"
```

转义是使用传统的反斜线的方式。参见[Characters](#characters)，查看支持的转义序列。

整行String 是由三个引号包裹的(`"""`),不可以包含转义符但可 以包含其它字符：

```kotlin
val text = """
	for (c in "foo")
		print(c)
"""
```

你可以通过 [trim-margin()](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.text/trim-margin.html) 函数移除空格：

```kotlin
val text = """
    |Tell me and I forget.
    |Teach me and I remember.
    |Involve me and I learn.
    |(Benjamin Franklin)
    """.trimMargin()
```

默认采用`|`标注起始前缀，也可以传递其它的字符做起始前缀，比如`trimMargin(">")`

#### 字符串模板

字符串可以包含模板表达式，即可求值的代码片段，并将其结果连接到字符串中。模板表达式由 $ 开始并包含另一个简单的名称：

```kotlin
val i = 10
val s = "i = $i" // 求值为 "i = 10"
```

或者是一个带大括号的表达式：

```kotlin
val s = "abc"
val str = "$s.length is ${s.length}" // 结果为 "abc.length is 3"
```

模板既可以原始字符串中使用，也可以在转义字符串中使用。如果需要在原始字符串(不支持反斜杠转义)中表示一个文字$字符，那么可以使用以下语法：

```kotlin
val price = """
${'$'}9.99
"""
```
