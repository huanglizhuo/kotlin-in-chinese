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

### 上下文中的 Job

协程的 [Job](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-job/index.html) 是其上下文的一部分，可以使用 `coroutineContext[Job]` 表达式从其检索：

```Kotlin
package kotlinx.coroutines.guide.context05

import kotlinx.coroutines.*

fun main() = runBlocking<Unit> {
    println("My job is ${coroutineContext[Job]}")
}
```

在[调试模式](https://kotlinlang.org/docs/reference/coroutines/coroutine-context-and-dispatchers.html#debugging-coroutines-and-threads)下，它输出如下内容：

```Kotlin
My job is "coroutine#1":BlockingCoroutine{Active}@6d311334
```

请注意，CoroutineScope中的 [isActive](My job is "coroutine#1":BlockingCoroutine{Active}@6d311334
) 只是 `coroutineContext[Job]?.isActive == true` 的快捷方式。

### 子协程

当协程在另一个协程的 CoroutineScope 中启动时，它会通过 [CoroutineScope.coroutineContext](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-coroutine-scope/coroutine-context.html)继承其上下文，新协程的Job成为父协程工作的子级。当父协程被取消时，其所有子进程也将被递归取消。

但是，当使用 GlobalScope 启动协程时，新协程的 Job 没有父项。因此，它不依赖于它的发布范围和独立运行。

```Kotlin
package kotlinx.coroutines.guide.context06

import kotlinx.coroutines.*

fun main() = runBlocking<Unit> {
    // launch a coroutine to process some kind of incoming request
    val request = launch {
        // it spawns two other jobs, one with GlobalScope
        GlobalScope.launch {
            println("job1: I run in GlobalScope and execute independently!")
            delay(1000)
            println("job1: I am not affected by cancellation of the request")
        }
        // and the other inherits the parent context
        launch {
            delay(100)
            println("job2: I am a child of the request coroutine")
            delay(1000)
            println("job2: I will not execute this line if my parent request is cancelled")
        }
    }
    delay(500)
    request.cancel() // cancel processing of the request
    delay(1000) // delay a second to see what happens
    println("main: Who has survived request cancellation?")
}
```

此代码的输出是：
```Kotlin
job1: I run in GlobalScope and execute independently!
job2: I am a child of the request coroutine
job1: I am not affected by cancellation of the request
main: Who has survived request cancellation?
```

### 父协程的责任

父协程总是等待所有子进程完成。父级不必显式跟踪其启动的所有子级，也不必使用Job.join在末尾等待它们：

```Kotlin
package kotlinx.coroutines.guide.context07

import kotlinx.coroutines.*

fun main() = runBlocking<Unit> {
    // launch a coroutine to process some kind of incoming request
    val request = launch {
        repeat(3) { i -> // launch a few children jobs
            launch  {
                delay((i + 1) * 200L) // variable delay 200ms, 400ms, 600ms
                println("Coroutine $i is done")
            }
        }
        println("request: I'm done and I don't explicitly join my children that are still active")
    }
    request.join() // wait for completion of the request, including all its children
    println("Now processing of the request is complete")
}
```

结果将是：

```Kotlin
request: I'm done and I don't explicitly join my children that are still active
Coroutine 0 is done
Coroutine 1 is done
Coroutine 2 is done
Now processing of the request is complete
```

### 命名协程以进行调试

当协程经常记录日志时，自动分配的ID是很好的方式，你只需要关联来自同一协程的日志记录。但是，当协程与特定请求的处理或执行某些特定的后台任务相关时，最好为调试目的明确命名它。 [CoroutineName](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-coroutine-name/index.html) 上下文元素的作用与线程名称相同。当打开[调试模式](https://kotlinlang.org/docs/reference/coroutines/coroutine-context-and-dispatchers.html#debugging-coroutines-and-threads)时，它包含在执行此协程的线程名称中。

下面的示例演示了此概念：

```Kotlin
package kotlinx.coroutines.guide.context08

import kotlinx.coroutines.*

fun log(msg: String) = println("[${Thread.currentThread().name}] $msg")

fun main() = runBlocking(CoroutineName("main")) {
    log("Started main coroutine")
    // run two background value computations
    val v1 = async(CoroutineName("v1coroutine")) {
        delay(500)
        log("Computing v1")
        252
    }
    val v2 = async(CoroutineName("v2coroutine")) {
        delay(1000)
        log("Computing v2")
        6
    }
    log("The answer for v1 / v2 = ${v1.await() / v2.await()}")
}
```

使用-Dkotlinx.coroutines.debug JVM选项生成的输出如下：

```shell
[main @main#1] Started main coroutine
[main @v1coroutine#2] Computing v1
[main @v2coroutine#3] Computing v2
[main @main#1] The answer for v1 / v2 = 42
```

### 结合上下文元素

有时我们需要为协程环境定义多个元素。我们可以使用`+`运算符。例如，我们可以同时使用指定的调度器和指定的名称启动协程：

```Kotlin
package kotlinx.coroutines.guide.context09

import kotlinx.coroutines.*

fun main() = runBlocking<Unit> {
    launch(Dispatchers.Default + CoroutineName("test")) {
        println("I'm working in thread ${Thread.currentThread().name}")
    }
}
```

使用-Dkotlinx.coroutines.debug JVM选项生成的输出如下：

```shell
I'm working in thread DefaultDispatcher-worker-1 @test#2
```

### 协程作用域

让我们将有关上下文，子协程和 Job 的知识放在一起。假设我们的应用程序有一个具有生命周期的对象，但是该对象不是协程。例如，我们正在编写一个Android应用程序，并在一个 Activity 的上下文中启动各种协程，以执行异步操作来获取和更新数据，制作动画等。在销毁该 Activity 时，必须取消所有这些协程以避免内存泄漏。当然，我们可以手动操作上下文和作业以绑定Activity及其协程的生命周期，但是kotlinx.coroutines提供了一个封装以下内容的抽象: [CoroutineScope](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-coroutine-scope/index.html) 。你应该已经熟悉了协程作用域，因为所有协程构建器都被声明为它的扩展。

我们通过创建与Activity的生命周期相关联的CoroutineScope实例来管理协程的生命周期。可以通过 [CoroutineScope()](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-coroutine-scope.html) 或 [MainScope()](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-main-scope.html)工厂函数创建CoroutineScope实例。前者创建通用作用域，而后者创建UI应用程序的作用域，并使用 Dispatchers.Main 作为默认调度器：

```Kotlin
class Activity {
    private val mainScope = MainScope()
    
    fun destroy() {
        mainScope.cancel()
    }
    // to be continued ...
```

另外，我们可以在Activity类中实现CoroutineScope接口。最好的方法是将委托与默认的工厂功能一起使用。我们还可以将所需的调度器（在此示例中使用Dispatchers.Default）与作用域相结合:

```Kotlin
class Activity : CoroutineScope by CoroutineScope(Dispatchers.Default) {
    // to be continued ...
```
现在，我们可以在此Activity的范围内启动协程，而不必显式指定它们的上下文。为了演示，我们启动了十个协程，它们会在不同的时间延迟：

```Kotlin
// class Activity continues
    fun doSomething() {
        // launch ten coroutines for a demo, each working for a different time
        repeat(10) { i ->
            launch {
                delay((i + 1) * 200L) // variable delay 200ms, 400ms, ... etc
                println("Coroutine $i is done")
            }
        }
    }
} // class Activity ends
```

在我们的主函数中，我们创建 Activity ，调用测试doSomething函数，并在500毫秒后销毁。这将取消从doSomething启动的所有协程。我们可以看到，因为在Activity被破坏后，即使等待了更长的时间，也不会再打印任何消息。

```Kotlin
package kotlinx.coroutines.guide.context10

import kotlin.coroutines.*
import kotlinx.coroutines.*

class Activity : CoroutineScope by CoroutineScope(Dispatchers.Default) {

    fun destroy() {
        cancel() // Extension on CoroutineScope
    }
    // to be continued ...

    // class Activity continues
    fun doSomething() {
        // launch ten coroutines for a demo, each working for a different time
        repeat(10) { i ->
            launch {
                delay((i + 1) * 200L) // variable delay 200ms, 400ms, ... etc
                println("Coroutine $i is done")
            }
        }
    }
} // class Activity ends

fun main() = runBlocking<Unit> {
    val activity = Activity()
    activity.doSomething() // run test function
    println("Launched coroutines")
    delay(500L) // delay for half a second
    println("Destroying activity!")
    activity.destroy() // cancels all coroutines
    delay(1000) // visually confir
```

该示例的输出为：

```Kotlin
Launched coroutines
Coroutine 0 is done
Coroutine 1 is done
Destroying activity!
```

如你所见，只有前两个协程打印一条消息，而其他两个在Activity.destroy（）中的一次job.cancel（）调用中被取消。

### Thradd-local数据

有时，能够将一些线程局部数据传递到协程或在协程之间很方便。但是，由于它们未绑定到任何特定线程，因此如果手动完成，则可能会导致模板式代码。

对于 [ThreadLocal](https://docs.oracle.com/javase/8/docs/api/java/lang/ThreadLocal.html) ，此处提供了 [asContextElement](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/java.lang.-thread-local/as-context-element.html) 扩展函数以进行恢复。它创建一个额外的context元素，该元素保留给定ThreadLocal，并在协程每次切换其上下文时将其恢复。

请看下面演示：

```Kotlin
package kotlinx.coroutines.guide.context11

import kotlinx.coroutines.*

val threadLocal = ThreadLocal<String?>() // declare thread-local variable

fun main() = runBlocking<Unit> {
    threadLocal.set("main")
    println("Pre-main, current thread: ${Thread.currentThread()}, thread local value: '${threadLocal.get()}'")
    val job = launch(Dispatchers.Default + threadLocal.asContextElement(value = "launch")) {
        println("Launch start, current thread: ${Thread.currentThread()}, thread local value: '${threadLocal.get()}'")
        yield()
        println("After yield, current thread: ${Thread.currentThread()}, thread local value: '${threadLocal.get()}'")
    }
    job.join()
    println("Post-main, current thread: ${Thread.currentThread()}, thread local value: '${threadLocal.get()}'")
}
```

在此示例中，我们使用Dispatchers.Default在后台线程池中启动了一个新协程，因此它在与线程池不同的线程上工作，但是它仍然具有我们使用 `threadLocal.asContextElement(value = "launch")`，无论协程在哪个线程上执行。 因此，输出（带有调试）为：

```shell
Pre-main, current thread: Thread[main @coroutine#1,5,main], thread local value: 'main'
Launch start, current thread: Thread[DefaultDispatcher-worker-1 @coroutine#2,5,main], thread local value: 'launch'
After yield, current thread: Thread[DefaultDispatcher-worker-2 @coroutine#2,5,main], thread local value: 'launch'
Post-main, current thread: Thread[main @coroutine#1,5,main], thread local value: 'main'
```

很容易忘记设置相应的上下文元素。如果运行协程的线程不同，则从协程访问的线程局部变量可能会具有非期待值。为避免此类情况，建议使用 [ensurePresent](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/java.lang.-thread-local/ensure-present.html) 方法，并在使用不当时进行快速失败。

ThreadLocal 具有顶级的支持，可与任何原始 kotlinx.coroutines 一起使用。但是，它有一个关键限制：当对线程局部变量进行更改时，不会将新值传播到协程调用者（因为上下文元素无法跟踪所有ThreadLocal对象访问），并且在下一次挂起时更新的值会丢失。使用 [withContext](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/with-context.html) 更新协程中 Thread-local 的值，有关更多详细信息，请参见 [asContextElement](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/java.lang.-thread-local/as-context-element.html)

同样，可以将值存储在可变类如 `class Counter(var i: Int)` 中，然后将其存储在线程局部变量中。但是，在这种情况下，你有责任将可能并发的修改同步到此可变的变量中。

对于高级用法，例如与日志记录MDC，事务上下文或内部使用线程本地传递数据的任何其他库的集成，请参阅 [ThreadContextElement](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-thread-context-element/index.html) 接口应实现的文档。