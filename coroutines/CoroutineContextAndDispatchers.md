**目录**

* [协程上下文和调度器](#协程上下文和调度器)
    * [调度器和线程](#调度器和线程)
    * [不确定与确定调度器](#不确定与确定调度器)
    * [调试协程和线程](#调试协程和线程)
    * [在线程之间跳转](#在线程之间跳转)
    * [在上下文中的Job](#在上下文中的Job)
    * [子协程](#子协程)
    * [父协程的责任](#父协程的责任)
    * [命名协程以进行调试](#命名协程以进行调试)
    * [结合上下文元素](#结合上下文元素)
    * [协程作用域](#协程作用域)
    * [Thradd-local数据](#Thradd-local数据)

## 协程上下文和调度器

协程始终在由Kotlin标准库中定义的CoroutineContext类型的值表示的上下文中执行。

协程上下文是一组各种元素。主要元素是协程的 Job（我们之前已经看到过）及其调度器，本节将对此进行介绍。

### 调度器和线程

协程上下文包括一个协程调度器（请参阅 [CoroutineDispatcher](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-coroutine-dispatcher/index.html) ），该调度器确定相应协程用于执行哪个线程或多个线程。协程调度器可以将协程执行限制在特定线程中，将其分派到线程池，或者让它不确定地运行。

所有协程构建器（如 launch 和 async ）都接受可选的CoroutineContext参数，该参数可用于为新协程和其他上下文元素显式指定调度器。

尝试以下示例：

```Kotlin
launch { // context of the parent, main runBlocking coroutine
    println("main runBlocking      : I'm working in thread ${Thread.currentThread().name}")
}
launch(Dispatchers.Unconfined) { // not confined -- will work with main thread
    println("Unconfined            : I'm working in thread ${Thread.currentThread().name}")
}
launch(Dispatchers.Default) { // will get dispatched to DefaultDispatcher 
    println("Default               : I'm working in thread ${Thread.currentThread().name}")
}
launch(newSingleThreadContext("MyOwnThread")) { // will get its own new thread
    println("newSingleThreadContext: I'm working in thread ${Thread.currentThread().name}")
}
```

输出如下（可能以不同的顺序）：

```shell

Unconfined            : I'm working in thread main
Default               : I'm working in thread DefaultDispatcher-worker-1
newSingleThreadContext: I'm working in thread MyOwnThread
main runBlocking      : I'm working in thread main
```


当不带参数使用 launch{...}时，它将从要启动的 CoroutineScope 继承上下文（并因此继承调度器）。在这种情况下，它继承了在主线程中运行的主 runBlocking 协程的上下文。

[Dispatchers.Unconfined](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-dispatchers/-unconfined.html) 是一种特殊的调度器，它似乎也运行在主线程中，但实际上，这是一种不同的机制，稍后将进行说明。

在 [GlobalScope](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-global-scope/index.html) 中启动协程时使用的默认调度器由[Dispatchers.Default](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-dispatchers/-default.html)表示，并使用共享的后台线程池，因此 launch（Dispatchers.Default）{...}与GlobalScope.launch {..使用相同的调度器。 }。

[newSingleThreadContext](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/new-single-thread-context.html) 为协程创建一个线程来运行。专用线程是非常昂贵的资源。在实际的应用程序中，必须在不需要时使用 [close](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-executor-coroutine-dispatcher/close.html) 函数将其释放，或者将其存储在顶级变量中，然后在整个应用程序中重复使用。

### 不确定与确定调度器
[Dispatchers.Unconfined](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-dispatchers/-unconfined.html) 协程调度器在调用者线程中启动协程，但仅直到第一个挂起点为止。挂起后，它将在线程中恢复协程，该协程完全由调用的挂起函数确定。不确定调度器适用于既不占用CPU时间也不更新确定于特定线程的任何共享数据（如UI）的协程。

另一方面，默认情况下，调度器是从外部CoroutineScope继承的。特别是，runBlocking协程的默认调度器仅限于调用程序线程，因此继承它的作用是通过可预测的FIFO调度将执行限制在此线程中。

```Kotlin
launch(Dispatchers.Unconfined) { // not confined -- will work with main thread
    println("Unconfined      : I'm working in thread ${Thread.currentThread().name}")
    delay(500)
    println("Unconfined      : After delay in thread ${Thread.currentThread().name}")
}
launch { // context of the parent, main runBlocking coroutine
    println("main runBlocking: I'm working in thread ${Thread.currentThread().name}")
    delay(1000)
    println("main runBlocking: After delay in thread ${Thread.currentThread().name}")
}
```

输出：

```shell
Unconfined      : I'm working in thread main
main runBlocking: I'm working in thread main
Unconfined      : After delay in thread kotlinx.coroutines.DefaultExecutor
main runBlocking: After delay in thread main
```

因此，具有从 `runBlocking {...}` 继承的上下文的协程将继续在主线程中执行，而不受约束的协程将在延迟函数使用的默认执行程序线程中继续执行。

不确定调度程序是一种高级机制，在某些情况下很有用，因为在这种情况下不需要协程进行分派以便以后执行，否则会产生不良的副作用，因为协程中的某些操作必须立即执行。不确定的调度程序不应在通用代码中使用。

### 调试协程和线程

协程可以在一个线程上挂起，并在另一个线程上恢复。即使使用单线程调度程序，也可能很难弄清楚协程在做什么，在何时何地进行。使用线程调试应用程序的常用方法是在每个log语句的日志文件中打印线程名称。日志记录框架普遍支持此功能。当使用协程时，仅线程名不会提供太多上下文，因此 `kotlinx.coroutines` 包含调试工具以使其更容易。

使用 `-Dkotlinx.coroutines.debug` JVM选项运行以下代码：

```Kotlin
package kotlinx.coroutines.guide.context03

import kotlinx.coroutines.*

fun log(msg: String) = println("[${Thread.currentThread().name}] $msg")

fun main() = runBlocking<Unit> {
    val a = async {
        log("I'm computing a piece of the answer")
        6
    }
    val b = async {
        log("I'm computing another piece of the answer")
        7
    }
    log("The answer is ${a.await() * b.await()}")
}
```

有三个协程。 runBlocking 内部的主要协程（＃1）和两个协程计算延迟值a（＃2）和b（＃3）。它们都在runBlocking上下文中执行，并且仅限于主线程。此代码的输出是：

```shell
[main @coroutine#2] I'm computing a piece of the answer
[main @coroutine#3] I'm computing another piece of the answer
[main @coroutine#1] The answer is 42
```

log函数将线程的名称打印在方括号中，可以看到它是主线程，并附加了当前正在执行的协程的标识符。调试模式打开时，此标识符将连续分配给所有创建的协程。

> 当使用 `-ea` 选项运行JVM时，调试模式也会打开。您可以在 [DEBUG_PROPERTY_NAME](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-d-e-b-u-g_-p-r-o-p-e-r-t-y_-n-a-m-e.html) 属性的文档中阅读有关调试功能的更多信息。

### 在线程之间跳转

使用-Dkotlinx.coroutines.debug JVM选项运行以下代码：

```Kotin
package kotlinx.coroutines.guide.context04

import kotlinx.coroutines.*

fun log(msg: String) = println("[${Thread.currentThread().name}] $msg")

fun main() {
    newSingleThreadContext("Ctx1").use { ctx1 ->
        newSingleThreadContext("Ctx2").use { ctx2 ->
            runBlocking(ctx1) {
                log("Started in ctx1")
                withContext(ctx2) {
                    log("Working in ctx2")
                }
                log("Back to ctx1")
            }
        }
    }
}

```

它演示了几种新技术。一种是在具有明确指定的上下文的情况下使用 [runBlocking](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/run-blocking.html)，另一种是使用 withContext 函数更改协程的上下文，同时仍保留在同一协程中，如下面的输出所示：

```shell
[Ctx1 @coroutine#1] Started in ctx1
[Ctx2 @coroutine#1] Working in ctx2
[Ctx1 @coroutine#1] Back to ctx1
```

请注意，此示例还使用Kotlin标准库中的 `use` 函数，以在不再需要使用newSingleThreadContext创建的线程时释放它们。

在上下文中的工作
协程的Job是其上下文的一部分，可以使用coroutineContext [Job]表达式从其检索：

println（“我的工作是$ {coroutineContext [Job]}”）
目标平台：在kotlin v.1.3.61上运行的JVM
您可以在此处获取完整的代码。

在调试模式下，它输出如下内容：

我的工作是“协程＃1”：BlockingCoroutine {Active} @ 6d311334
请注意，CoroutineScope中的isActive只是coroutineContext [Job] ?. isActive == true的便捷快捷方式。

协程的孩子
当协程在另一个协程的CoroutineScope中启动时，它会通过CoroutineScope.coroutineContext继承其上下文，新协程的Job成为父协程工作的子级。当父协程被取消时，其所有子进程也将被递归取消。

但是，当使用GlobalScope启动协程时，新协程的工作没有父项。因此，它不依赖于它的发布范围和独立运行。

//启动协程以处理某种传入请求
val请求=启动{
    //它产生了另外两个作业，其中一个是GlobalScope
    GlobalScope.launch {
        println（“ job1：我在GlobalScope中运行并独立执行！”）
        延迟（1000）
        println（“ job1：我不受取消请求的影响”）
    }
    //并且另一个继承父上下文
    发射{
        延迟（100）
        println（“ job2：我是协程的孩子”）
        延迟（1000）
        println（“ job2：如果我的父请求被取消，我将不执行此行”）
    }
}
延迟（500）
request.cancel（）//取消请求处理
delay（1000）//延迟一秒钟看看会发生什么
println（“ main：谁在取消请求后还幸存？”）
目标平台：在kotlin v.1.3.61上运行的JVM
您可以在此处获取完整的代码。

此代码的输出是：

job1：我在GlobalScope中运行并独立执行！
job2：我是协程的孩子
job1：我不受取消请求的影响
main：谁在请求取消中幸存下来？
父母的责任
父协程总是等待所有子进程完成。父级不必显式跟踪其启动的所有子级，也不必使用Job.join在末尾等待它们：

//启动协程以处理某种传入请求
val请求=启动{
    repeat（3）{i-> //启动几个孩子的工作
        发射{
            delay（（i + 1）* 200L）//可变延迟200ms，400ms，600ms
            println（“协程$ i已完成”）
        }
    }
    println（“请求：我完成了，并且我没有明确加入仍然活跃的孩子”）
}
request.join（）//等待请求完成，包括所有子请求
println（“请求的处理已完成”）
目标平台：在kotlin v.1.3.61上运行的JVM
您可以在此处获取完整的代码。

结果将是：

要求：我已经完成，但没有明确加入仍然活跃的孩子
协程0完成
协程1完成
协程2完成
至此，请求处理完成
命名协程以进行调试
当协程经常记录日志时，自动分配的ID很好，您只需要关联来自同一协程的日志记录。但是，当协程与特定请求的处理或执行某些特定的后台任务相关时，最好为调试目的明确命名它。 CoroutineName上下文元素的作用与线程名称相同。当打开调试模式时，它包含在执行此协程的线程名称中。

下面的示例演示了此概念：