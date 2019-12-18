Kotlin 作为一门语言，只在标准库中提供了最基本的底层 APIs 以供其它库使用协程。与许多其他具有类似功能的语言不同，`async` 和 `await` 并不是 Kotlin 的关键字，甚至都没有包含在标准库中。此外，Kotlin的挂起函数的概念为异步操作提供了比异步操作 futures 和 promises 更安全且更不容易出错的抽象。

`kotlinx.coroutines` 是由 Jetbrains 为协程开发的一个功能丰富的库。包含了许多支持协程高阶原函数，包括 `launch` `async`以及其它。这些均将在本册指南中描述。

这部 `kotlinx.coroutines` 的指南包含了许多关于关键特性的例子，并被分为不同的主题。

为了使用协程以及按照本指南中的例子练习，需要添加 kotlinx-coroutines-core 模块依赖。参考 [README](https://github.com/kotlin/kotlinx.coroutines/blob/master/README.md#using-in-your-projects)

内容目录

* [协程](README.md)
   * [协程指南](CoroutinesGuide.md)
   * [基础](Basics.md)
   * [取消和超时](CancellationAndTimeouts.md)
   * [频道](Channels.md)
   * [组合挂起函数](ComposingSuspendingFunctions.md)
   * [协程上下文和调度器](CoroutineContextAndDispatchers.md)
   * [异常处理](ExceptionHandling.md)
   * [Select 表达式](SelectExpression.md)
   * [共享可变状态与并发](SharedMutableStateAndConcurrency.md)

附加参考

* [Guide to UI programming with coroutines](https://github.com/kotlin/kotlinx.coroutines/blob/master/ui/coroutines-guide-ui.md)
* [Guide to reactive streams with coroutines](https://github.com/kotlin/kotlinx.coroutines/blob/master/reactive/coroutines-guide-reactive.md)
* [Coroutines design document (KEEP)](https://github.com/Kotlin/kotlin-coroutines-examples/blob/master/kotlin-coroutines-informal.md)
* [Full kotlinx.coroutines API reference](https://kotlin.github.io/kotlinx.coroutines/)