## Ranges
range 表达式拥有 `rangeTo` 函数操作符是 `..`  。 Range 可以对任何可比较的类型做操作，但对很多原语是优化过的。下面是些例子：

```kotlin
if (i in 1..10) {
	println(i)
}

if (x !in 1.0..3.0) println(x)

if (str in "island".."isle") println(str)
```

数字的范围有个附加的特性：它们可以迭代。编译器会把它转成类似于 java 的 for 循环的形式，且不用担心越界：

```kotlin
for (i in 1..4) print(i) // prints "1234"

for (i in 4..1) print(i) // prints nothing

for (x in 1.0..2.0) print("$x ") // prints "1.0 2.0 "
```

如果你想迭代数字并想反过来，这个相当简单，你可以使用 `downTo()` 函数

```kotlin
for (i in 4 downTo 1) print(i)
```

也可以使用指定步数的迭代，这个用到 `step()` 

```kotlin
for (i in 1..4 step 2) print(i) // prints "13"

for (i in 4 downTo 1 step 2) print(i) // prints "42"

for (i in 1.0..2.0 step 0.3) print("$i ") // prints "1.0 1.3 1.6 1.9 "
```

### 工作原理
在标准库中有俩种接口：Range<T> 和 Progression<N>

Range<T> 表示数学范围上的一个间隔。它有俩个端点：start 和 end 。主要的操作符是 contains 通常在 in/!in 操作符内：

Progression<N> 表示一个算数级数。它有一个 start 和 end 以及一个非零 increment 。Progression<N>　是Iterable<N> 的一个子类，因此可以使用在 for 循环中，或者 map filter 等等。第一个元素是 start 下一个元素都是前一个元素的 increment 。`Progression` 的迭代与 java/javaScript 的 for 循环相同：

```kotlin
// if increment > 0
for (int i = start; i <= end; i += increment) {
  // ...
}
// if increment < 0
for (int i = start; i >= end; i += increment) {
  // ...
}
```

### 范围指标
使用例子：

```kotlin
// Checking if value of comparable is in range. Optimized for number primitives.
if (i in 1..10) println(i)

if (x in 1.0..3.0) println(x)

if (str in "island".."isle") println(str)

// Iterating over arithmetical progression of numbers. Optimized for number primitives (as indexed for-loop in Java).
for (i in 1..4) print(i) // prints "1234"

for (i in 4..1) print(i) // prints nothing

for (i in 4 downTo 1) print(i) // prints "4321"

for (i in 1..4 step 2) print(i) // prints "13"

for (i in (1..4).reversed()) print(i) // prints "4321"

for (i in (1..4).reversed() step 2) print(i) // prints "42"

for (i in 4 downTo 1 step 2) print(i) // prints "42"

for (x in 1.0..2.0) print("$x ") // prints "1.0 2.0 "

for (x in 1.0..2.0 step 0.3) print("$x ") // prints "1.0 1.3 1.6 1.9 "

for (x in 2.0 downTo 1.0 step 0.3) print("$x ") // prints "2.0 1.7 1.4 1.1 "

for (str in "island".."isle") println(str) // error: string range cannot be iterated over
```

### 常见的接口的定义
有俩种基本接口：`Range` `Progression`

`Range` 接口定义了一个范围，或者是数学意义上的一个间隔。

```kotlin
interface Range<T : Comparable<T>> {
	val start: T
	val end: T
	fun contains(Element : T): Boolean
}
```

`Progression` 定义了数学上的级数。包括 start end increment 端点。最大的特点就是它可以迭代，因此它是 Iterable 的子类。end 不是必须的。

```kotlin
interface Progression<N : Number> : Iterable<N> {
	val start : N
	val end : N
	val increment : Number
}
```
与 java 的 for 循环类似：

```kotlin
// if increment > 0
for (int i = start; i <= end; i += increment) {
  // ...
}

// if increment < 0
for (int i = start; i >= end; i += increment) {
  // ...
}
```

### 类的实现
为避免不需要的重复，让我们先考虑一个数字类型　`Int` 。其它的数字类型也一样。注意这些类的实例需要用相应的构造函数来创建，使用 rangeTo() downTo() reversed() stop() 实用函数。

IntProgression 类很直接也很简单：

```kotlin
class IntProgression(override val start: Int, override val end: Int, override val increment: Int ): Progression<Int> {
	override fun iterator(): Iterator<Int> = IntProgressionIteratorImpl(start, end, increment)
}
```

`IntRange` 有些狡猾：它实现了 `Progression<Int>` `Range<Int>` 接口，因为它天生以通过 range 迭代(默认增加值是 1 )：

```kotlin
class IntRange(override val start: Int, override val end: Int): Range<Int>, Progression<Int> {
  override val increment: Int
    get() = 1
  override fun contains(element: Int): Boolean = start <= element && element <= end
  override fun iterator(): Iterator<Int> = IntProgressionIteratorImpl(start, end, increment)
}
```

`ComparableRange` 也很简单：

```kotlin
class ComparableRange<T : Comparable<T>>(override val start: T, override val end: T): Range<T> {
  override fun contains(element: T): Boolean = start <= element && element <= end
}
```

### 一些实用的函数
**rangeTo()**

`rangeTo()` 函数仅仅是调用 *Range 的构造函数，比如：

```kotlin
class Int {
	fun rangeTo(other: Byte): IntRange = IntRange(this, Other)
	fun rangeTo(other: Int): IntRange = IntRange(this, other)
}
```

**downTo()**

`downTo()` 扩展函数可以为任何数字类型定义，这里有俩个例子：

```kotlin
fun Long.downTo(other: Double): DoubleProgression {
	return DoubleProgression(this, other, -1.0)
}

fun Byte.downTo(other: Int): IntProgression {
	return IntProgression(this, other, -1)
}
```

**reversed()**

`reversed()` 扩展函数是给所有的 `*Range`和`*Progression` 类定义的，并且它们都返回反向的级数。

```kotlin
fun IntProgression.reversed(): IntProgression {
	return IntProgression(end, start, -increment)
}

fun IntRange.reversed(): IntProgression {
	return IntProgression(end, start, -1)
}
```
**step()**

`step()` 扩展函数是给所有的 `*Range`和`*Progression` 类定义的，所有的返回级数都修改了 setp 值。注意 step 值总是正的，否则函数不会改变迭代的方向。

```kotlin
fun IntProgression.step(step: Int): IntProgression {
  if (step <= 0) throw IllegalArgumentException("Step must be positive, was: $step")
  return IntProgression(start, end, if (increment > 0) step else -step)
}

fun IntRange.step(step: Int): IntProgression {
  if (step <= 0) throw IllegalArgumentException("Step must be positive, was: $step")
  return IntProgression(start, end, step)
}
```
