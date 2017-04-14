## 与 Scala 对比
Kotlin 设计时的俩个主要目标是：

> 至少和 java 运行速度一样快

> 在保证语言尽量简单的情况下在易用性上提高

考虑到这俩点，如果你喜欢 Scala ，你可能不需要 Kotlin

### Scala 有而 Kotlin 没有的
> 隐式转换，隐式参数
	--在 Scala 中，在不适用 debugger 的时候很难知道代码发生了什么，因为太多的东西是隐式的
	--通过函数增加类型在 kotlin 中需要使用[扩展函数](http://kotlinlang.org/docs/reference/extensions.html)

> 可重载和类型成员

> 路径依赖的类型

> 宏

> Existential types
	--类型推断是很特殊的情形

> 特征的初始化逻辑很复杂
	--参看[类和继承](http://kotlinlang.org/docs/reference/classes.html)

>自定义象征操作
	--参看[操作符重载](http://kotlinlang.org/docs/reference/operator-overloading.html)

> 内建 xml
	--参看[Type-safe Groovy-style builders](http://kotlinlang.org/docs/reference/type-safe-builders.html)

以后 kotlin可能会添加的特性：

> 结构类型

> 值类型

> Yield 操作符

> Actors

> 并行集合(Parallel collections)

### Kotlin 有而 Scala 没有的
>零开销的null安全
- Scala的是Option，是在句法和运行时的包装

>[ Smart casts](http://kotlinlang.org/docs/reference/typecasts.html)

>[Kotlin 的内联函数非局部的跳转](http://kotlinlang.org/docs/reference/inline-functions.html#inline-functions)

> [First-class delegation](http://kotlinlang.org/docs/reference/delegation.html)。也通过第三方插件：Autoproxy实现
