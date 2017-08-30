## 返回与跳转
Kotlin 有三种结构跳转表达式：

> -- return   
> -- break 结束最近的闭合循环  
> -- continue 跳到最近的闭合循环的下一次循环  

上述表达式都可以作为更大的表达式的一部分：

```kotlin
val s = person.name ?: return
```

这些表达式的类型是 [Nothing type](http://kotlinlang.org/docs/reference/exceptions.html#the-nothing-type)

### break 和 continue 标签
在 Kotlin 中表达式可以添加标签。标签通过 @ 结尾来表示，比如：`abc@`，`fooBar@` 都是有效的(参看[语法](http://kotlinlang.org/docs/reference/grammar.html#label))。使用标签语法只需像这样：

```kotlin
loop@ for (i in 1..100){
	// ...
}
```

我们可以用标签实现 break 或者 continue 的快速跳转：

```kotlin
loop@ for (i in 1..100) {
	for (j in i..100) {
		if (...)
			break@loop
	}
}
```

break 是跳转标签后面的表达式，continue 是跳转到循环的下一次迭代。

###  返回到标签


在字面函数，局部函数，以及对象表达式中，函数在Kotlin中是可以嵌套的。合法的return 允许我们返回到外层函数。最重要的使用场景就是从lambda表达式中返回，还记得我们之前的写法吗：

```kotlin
fun foo() {
	ints.forEach {
		if (it  == 0) return
		print(it)
	}
}
```

return 表达式返回到最近的闭合函数，比如 `foo` (注意这样非局部返回仅仅可以在[内联函数](http://kotlinlang.org/docs/reference/inline-functions.html)中使用)。如果我们需要从一个字面函数返回可以使用标签修饰 return :

```kotlin
fun foo() {
	ints.forEach lit@ {
		if (it ==0) return＠lit
		print(it)
	}
}
```

现在它仅仅从字面函数中返回。经常用一种更方便的含蓄的标签：比如用和传入的 lambda 表达式名字相同的标签。

```kotlin
fun foo() {
	ints.forEach {
		if (it == 0) return@forEach
		print(it)
	}
}
```

另外，我们可以用函数表达式替代匿名函数。在函数表达式中使用 return 语句可以从函数表达式中返回。

```kotlin
fun foo() {
	ints.forEach(fun(value:  Int){
		if (value == 0) return
		print(value)
	})
}
```


当返回一个值时，解析器给了一个参考，比如(原文When returning a value, the parser gives preference to the qualified return, i.e.)：

```kotlin
return@a 1
```

表示 “在标签 `@a` 返回 `1` ” 而不是返回一个标签表达式 `(@a 1)`

命名函数自动定义标签：

```kotlin
foo outer() {
	foo inner() {
		return@outer
	}
}
```
