## 流程控制
###  if 表达式
在 Kotlin 中，if 是表达式，比如它可以返回一个值。是除了condition ? then : else)之外的唯一一个三元表达式

```kotlin
//传统用法
var max = a
if (a < b)
	max = b

//带 else 
var max: Int
if (a > b)
	max = a
else
	max = b

//作为表达式
val max = if (a > b) a else b
```

if  分支可以作为块，最后一个表达是是该块的值：

```kotlin
val max = if (a > b){
	print("Choose a")
	a
}
else{
	print("Choose b")
	b
}
```


如果 if 表达式只有一个分支，或者分支的结果是 `Unit` , 它的值就是 `Unit` 。

参看[if语法](http://kotlinlang.org/docs/reference/grammar.html#if)

### When 表达式
when 取代了 C 风格语言的 switch 。最简单的用法像下面这样

```kotlin
when (x) {
	1 -> print("x == 1")
	2 -> print("x == 2")
	else -> { //Note the block
		print("x is neither 1 nor 2")
	}
}
```

when会对所有的分支进行检查直到有一个条件满足。when 可以用做表达式或声明。如果用作表达式的话，那么满足条件的分支就是总表达式。如果用做声明，那么分支的的的值会被忽略。(像 if 表达式一样，每个分支是一个语句块，而且它的值就是最后一个表达式的值)

在其它分支都不匹配的时候默认匹配 else 分支。如果把 when 做为表达式的话 else 分支是强制的，除非编译器可以提供所有覆盖所有可能的分支条件。

如果有分支可以用同样的方式处理的话，分支条件可以连在一起：

```kotlin
when (x) {
	0,1 -> print("x == 0 or x == 1")
	else -> print("otherwise")
}
```

可以用任意表达式作为分支的条件

```kotlin
when (x) {
	parseInt(s) -> print("s encode x")
	else -> print("s does not encode x")
}
```

甚至可以用 in 或者 !in 检查值是否值在一个集合中：

```kotlin
when (x) {
	in 1..10 -> print("x is in the range")
	in validNumbers -> print("x is valid")
	!in 10..20 -> print("x is outside the range")
	else -> print("none of the above")
}
```

也可以用 is 或者 !is 来判断值是否是某个类型。注意，由于 [smart casts](http://kotlinlang.org/docs/reference/typecasts.html#smart-casts) ，你可以不用另外的检查就可以使用相应的属性或方法。

```kotlin
val hasPrefix = when (x) {
	is String -> x.startsWith("prefix")
	else -> false
}
```

when 也可以用来代替 if-else if 。如果没有任何参数提供，那么分支的条件就是简单的布尔表达式，当条件为真时执行相应的分支：

```kotlin
when {
	x.isOdd() -> print("x is odd")
	x.isEven() -> print("x is even")
	else -> print("x is funny")
}
```

参看[when语法](http://kotlinlang.org/docs/reference/grammar.html#when)

### for 循环
for 循环通过任何提供的迭代器进行迭代。语法是下面这样的：

```kotlin
for (item in collection)
	print(item)
```

内容可以是一个语句块

```kotlin
for (item: Int in ints){
	//...
}
```

像之前提到的， for 可以对任何提供的迭代器进行迭代，比如：

> has an instance- or extension-function iterator(), whose return type

> has an instance- or extension-function next(), and

> has an instance- or extension-function hasNext() that returns Boolean.

如果你想通过 list 或者 array 的索引进行迭代，你可以这样做：

```kotlin
for (i in array.indices)
	print(array[i])
```

在没有其它对象创建的时候 "iteration through a range " 会被自动编译成最优的实现。

### while 循环
while 和 do...while 像往常那样

```kotlin
while (x > 0) {
	x--
}

do {
	val y = retrieveData()
} while (y != null) // y 在这是可见的
```

参看[while 语法](has an instance- or extension-function hasNext() that returns Boolean.)

### 在循环中使用 break 和 continue
kotlin 支持传统的 break 和 continue 操作符。参看[返回和跳转](http://kotlinlang.org/docs/reference/returns.html)
