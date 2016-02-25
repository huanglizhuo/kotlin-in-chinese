##基本类型

在 Kotlin 中，所有变量的成员方法和属性都是一个对象。一些类型是内建的，因为它们的实现是优化过的，但对用户来说它们就像普通的类一样。在这节中，我们将会讲到大多数的类型：数值，字符，布尔，以及数组。

###数值

Kotlin 处理数值的方法和 java 很相似，但不是完全一样。比如，不存在隐式转换数值的宽度，并且在字面上有一些小小的不同。

Kotlin 提供了如下内建数值类型(和 java 很相似)：


| **Type** | **Bitwidth** |
| -------------- | :------------------: |
| Double | 64 |
| Float | 32 |
| Long | 64 |
| Int | 32 |
| Short | 16 |
| Byte | 8 |

注意字符在 Kotlin 中不是数值类型

###字面值常量

主要是以下几种字面值常量和字符类型：

> --数型: `123`
> --长整型要加大写 `L` : `123L`
> --16进制：0x0f
> --二进制：0b00001011

注意不支持８进制

 Kotlin 也支持浮点数：

> -- 默认 Doubles : 123.5 , 123.5e10
> -- Floats 是通过加 `f` 或 `F` 来实现的：123.5f

###表示

在 java 平台上，数值被 JVM 虚拟机以字节码的方式物理存储的，除非我们需要做可空标识(比如说 Int?) 或者是泛型调用的。在后者中数值是装箱的。

注意装箱过的数值是不保留特征的：

```kotlin
val a: Int = 10000
print (a === a ) //打印 'true'
val boxedA: Int? =a
val anotherBoxedA: Int? = a
print (boxedA === anotherBoxedA ) //注意这里打印的是 'false'
```

另一方面，它们是值相等的：

```kotlin
val a: Int = 10000
print(a == a) // Prints 'true'
val boxedA: Int? = a
val anotherBoxedA: Int? = a
print(boxedA == anotherBoxedA) // Prints 'true'
```

###显式转换

由于不同的表示，短类型不是长类型的子类型。如果是的话我们就会碰到下面这样的麻烦了

```kotlin
//这是些伪代码，不能编译的
val a: Int? =1 //一个装箱过的 Int (java.lang.Integer)
val b: Long? = a // 一个隐式装箱的 Long (java.lang.Long)
pritn ( a == b )// 很惊讶吧　这次打印出的是 'false'
```
因此特性甚至值都会悄悄丢失掉

所以，短类型是不会隐式转换为长类型的。这意味着我们必须显式转换才能把 `Byte` 赋值给 `Int` 

```kotlin
val b: Byte = 1 // OK, literals are checked statically
val i: Int = b //ERROR
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

隐式转换一般情况下是不容易被发觉的，因为我们可以使用上下文推断出类型，并且算术运算会为合适的转换进行重载，比如

```kotlin
val l = 1.toLong + 1 //Long  + Int => Long
```

###运算符

Kotlin支持标准的算术运算表达式，这些运算符被声明为相应类的成员。参看[运算符重载](http://kotlinlang.org/docs/reference/operator-overloading.html)。

至于位运算，Kotlin 并没有提供特殊的操作符，只是提供了可以叫中缀形式的方法，比如：

val x = (1 shl 2) and 0x000FF000

下面是全部的位运算操作符(只可以用在 `Int` 和 `Long` 类型)：

> `shl(bits)` – 带符号左移 (相当于 Java’s `<<`)
> `shr(bits)` – 带符号右移 (相当于 Java’s `>>`)
> `ushr(bits)` – 无符号右移 (相当于 Java’s `>>>`)
> `and(bits)` – 按位与
> `or(bits)` – 按位或
> `xor(bits)` – 按位异或
> `inv(bits)` – 按位翻转

###字符

字符类型用 `Char` 表示。不能直接当做数值来使用

```Kotlin
fun check(c: Char) {
	if (c == 1) { //ERROR: 类型不匹配
		//...
	}
}
```

字符是单引号包起来的 `'1'`,`'\n'`,`'\uFF00'`。我们可以显示的把它转换为 `Int` 型

```kotlin
fun decimalDigitValue(c: Char): Int {
	if (c !in '0'..'9') 
		throw IllegalArgumentException("Out of range")
	return c.toInt() - '0'.toInt() //显示转换为数值类型
}
```
和数值类型一样，字符在空检查后会在需要的时候装箱。特性不会被装箱操作保留的。

###布尔值

布尔值只有 true 或者 false

布尔值的内建操作包括

> || – lazy disjunction
> && – lazy conjunction

###Array

Arrays在 Kotlin 中由 `Array` 类表示，有 `get` 和 `set` 方法(通过运算符重载可以由[]调用)，以及 `size` 方法，以及一些常用的函数：

```kotlin
class Array<T> private () {
	fun size(): Int
	fun get(index: Int): T
	fun set(Index: Int, value: T): Uint
	fun iterator(): Iterator<T>
	//...
}
```

我们可以给库函数 `arrayOf()` 传递每一项的值来创建Array，`arrayOf(1, 2, 3)` 创建了一个[1, 2, 3] 这样的数组。也可以使用库函数 `arrayOfNulls()` 创建一个指定大小的空Array。

或者通过指定Array大小并提供一个迭代器

(原文Another option is to use a factory function that takes the array size and the function that can return the initial value of each array element given its index)：

```kotlin
// 创建一个 Array<String>  内容为 ["0", "1", "4", "9", "16"]
val asc = Array(5, {i -> (i * i).toString() })
```

像我们上面提到的，`[]` 操作符表示调用　`get()` `set()` 函数

注意：和 java 不一样，arrays 在 kotlin 中是不可变的。这意味这 kotlin 不允许我们把 `Array<String>` 转为 `Array<Any>` ,这样就阻止了可能的运行时错误(但你可以使用 `Array<outAny>` , 参看 [Type Projections](http://kotlinlang.org/docs/reference/generics.html#type-projections))

Kotlin 有专门的类来表示原始类型从而避免过度装箱： ByteArray, ShortArray, IntArray 等等。这些类与 Array 没有继承关系，但它们有一样的方法与属性。每个都有对应的库函数：

```kotlin
val x: IntArray = intArray(1, 2, 3)
x[0] = x[1] + x[2]
```

###字符串

字符串是有 `String` 表示的。字符串是不变的。字符串的元素可以通过索引操作读取: `s[i]` 。字符串可以用 for 循环迭代：

```kotlin
for (c in str) {
	println(c)
}
```

Kotlin 有俩种类型的 string ：一种是可以带分割符的，一种是可以包含新行以及任意文本的。带分割符的 string 很像 java 的 string:
```kotlin
val s = "Hello World!\n"
```

行String 是由三个引号包(`"""`)裹的,不包含分割符并且可以包含其它字符：

```kotlin
val text = """
	for (c in "foo")
		print(c)
"""
```

###模板

字符串可以包含模板表达式。一个模板表达式由一个 $ 开始并包含另一个简单的名称：

```kotlin
val i = 10
val s = "i = $i" // 识别为 "i = 10"
```

或者是一个带大括号的表达式：

```kotlin
val s = "abc"
val str = "$s.length is ${s.length}" //识别为 "abc.length is 3"
```
