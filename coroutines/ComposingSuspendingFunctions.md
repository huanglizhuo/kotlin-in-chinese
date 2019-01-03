
**目录**

<!--- TOC -->

* [组合挂起函数](#组合挂起函数)
  * [默认顺序调用](#默认顺序调用)
  * [使用 async 并发](#使用 async 并发)
  * [惰性启动的 async](#惰性启动的-async)
  * [async 风格的函数](#async-风格的函数)
  * [使用 async 的结构化并发](#使用-async-的结构化并发)

<!--- END_TOC -->

## 组合挂起函数

本节介绍了挂起函数的组合方法。

### 默认顺序调用

假设我们在其他地方定义了两个挂起函数，它们可以像某种远程服务调用或计算一样有用。 我们只是假装它们很有用，但实际上每个只是为了这个例子的目的而延迟一秒：

```kotlin
suspend fun doSomethingUsefulOne(): Int {
    delay(1000L) // pretend we are doing something useful here
    return 13
}

suspend fun doSomethingUsefulTwo(): Int {
    delay(1000L) // pretend we are doing something useful here, too
    return 29
}
```

如果按顺序调用它们,首先调用 doSomethingUsefulOne , 接下来调用 doSomethingUsefulTwo 并且计算它们结果的和？ 实际上，如果我们要根据第一个函数的结果来决定是否我们需要调用第二个函数或者决定如何调用它时，我们就会这样做。

我们使用普通的顺序来进行调用，因为这些代码是运行在协程中的，只要像常规的代码一样顺序都是默认的。下面的示例展示了测量执行两个挂起函数所需要的总时间：

```kotlin
import kotlinx.coroutines.*
import kotlin.system.*

fun main() = runBlocking<Unit> {
    val time = measureTimeMillis {
        val one = doSomethingUsefulOne()
        val two = doSomethingUsefulTwo()
        println("The answer is ${one + two}")
    }
    println("Completed in $time ms")    
}

suspend fun doSomethingUsefulOne(): Int {
    delay(1000L) // pretend we are doing something useful here
    return 13
}

suspend fun doSomethingUsefulTwo(): Int {
    delay(1000L) // pretend we are doing something useful here, too
    return 29
}
```

结果会像下面这样:

> The answer is 42 
> Completed in 2017 ms

### 使用 async 并发

如果 doSomethingUsefulOne 与 doSomethingUsefulTwo 之间没有依赖，如何能更快的得到结果，让它们进行并发执行呢? [async](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/async.html) 可以很好的完成这点。

概念上讲 async 与 launch 是一致的. 它启动了一个单独的协程,该协程是一个轻量级的线程并可以和其它所有的协程并发工作. 不同的是 launch 会返回一个 job 并且不带有任何结果值, 而 async 会返回一个 延期[Deffered](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-deferred/index.html) 一个轻量级的非阻塞的 future 表示一个会稍后提供结果的 promise. 你可以在延期值上使用 .await() 以获得结果,同时延期也是一个 job ,可以执行取消操作.

```kotlin
import kotlinx.coroutines.*
import kotlin.system.*

fun main() = runBlocking<Unit> {
    val time = measureTimeMillis {
        val one = async { doSomethingUsefulOne() }
        val two = async { doSomethingUsefulTwo() }
        println("The answer is ${one.await() + two.await()}")
    }
    println("Completed in $time ms")    
}

suspend fun doSomethingUsefulOne(): Int {
    delay(1000L) // pretend we are doing something useful here
    return 13
}

suspend fun doSomethingUsefulTwo(): Int {
    delay(1000L) // pretend we are doing something useful here, too
    return 29
}
```

结果如下:

```
The answer is 42
Completed in 1017 ms
```

速度提升了两倍,这是因为我们并发执行了两个协程.记住协程的并发永远是显示的.

### 惰性启动的async

使用值为 CoroutineStart.LAZY 的可选启动参数进行异步时设置惰性选项。 它仅在出现 await 调用或者调用 start 函数时才启动协同程序。 运行以下示例：

```kotlin
import kotlinx.coroutines.*
import kotlin.system.*

fun main() = runBlocking<Unit> {
    val time = measureTimeMillis {
        val one = async(start = CoroutineStart.LAZY) { doSomethingUsefulOne() }
        val two = async(start = CoroutineStart.LAZY) { doSomethingUsefulTwo() }
        // some computation
        one.start() // start the first one
        two.start() // start the second one
        println("The answer is ${one.await() + two.await()}")
    }
    println("Completed in $time ms")    
}

suspend fun doSomethingUsefulOne(): Int {
    delay(1000L) // pretend we are doing something useful here
    return 13
}

suspend fun doSomethingUsefulTwo(): Int {
    delay(1000L) // pretend we are doing something useful here, too
    return 29
}
```

运行结果如下:

The answer is 42
Completed in 1017 ms

所以，这里定义了两个协同程序，但是没有像前面的例子那样执行，但是程序员在完全通过调用start开始执行时会给出控制权。 我们首先启动一个，然后启动两个，然后等待各个协同程序完成。

注意，如果我们在println中调用了await并且在各个协程上省略了start，那么我们就会得到顺序行为，因为await启动协程执行并等待执行完成，这不是懒惰的预期用例。 在计算值涉及挂起函数的情况下，async（start = CoroutineStart.LAZY）的用例是标准惰性函数的替代。

### async风格的函数

我们可以使用具有显式 GlobalScope 引用的异步协程生成器来定义异步样式函数，这些函数可以异步调用doSomethingUsefulOne和doSomethingUsefulTwo。 我们为这些函数名添加“Async”后缀，以突出显示它们只启动异步计算的事实，并且需要使用生成的延迟值来获取结果。

```kotlin
// The result type of somethingUsefulOneAsync is Deferred<Int>
fun somethingUsefulOneAsync() = GlobalScope.async {
    doSomethingUsefulOne()
}

// The result type of somethingUsefulTwoAsync is Deferred<Int>
fun somethingUsefulTwoAsync() = GlobalScope.async {
    doSomethingUsefulTwo()
}
```

请注意，这些xxxAsync函数不是挂起函数。 它们可以在任何地方使用。 但是，它们的使用总是意味着它们的动作与调用代码的异步（这里意味着并发）。

以下示例显示了它们在协同程序之外的用法：

```kotlin
import kotlinx.coroutines.*
import kotlin.system.*

    // note, that we don't have `runBlocking` to the right of `main` in this example
    fun main() {
        val time = measureTimeMillis {
            // we can initiate async actions outside of a coroutine
            val one = somethingUsefulOneAsync()
            val two = somethingUsefulTwoAsync()
            // but waiting for a result must involve either suspending or blocking.
            // here we use `runBlocking { ... }` to block the main thread while waiting for the result
            runBlocking {
                println("The answer is ${one.await() + two.await()}")
            }
        }
        println("Completed in $time ms")
    }

fun somethingUsefulOneAsync() = GlobalScope.async {
    doSomethingUsefulOne()
}

fun somethingUsefulTwoAsync() = GlobalScope.async {
    doSomethingUsefulTwo()
}

suspend fun doSomethingUsefulOne(): Int {
    delay(1000L) // pretend we are doing something useful here
    return 13
}

suspend fun doSomethingUsefulTwo(): Int {
    delay(1000L) // pretend we are doing something useful here, too
    return 29
}
```

> 这里提供了具有异步功能的编程风格，仅用于说明，因为在其他编程语言中很流行。 由于下面解释的原因，强烈建议不要将这种风格与Kotlin协同程序一起使用。

考虑一下如果 val one = somethingUsefulOneAsync() 这一行和 one.await() 表达式这里在代码中有逻辑错误， 并且程序抛出了异常以及程序在操作的过程中被中止，将会发生什么。 通常情况下，一个全局的异常处理者会捕获这个异常，将异常打印成日记并报告给开发者，但是反之该程序将会继续执行其它操作。但是这里我们的 somethingUsefulOneAsync 仍然在后台执行， 尽管如此，启动它的那次操作也会被终止。这个程序将不会进行结构化并发，如下一小节所示。

### 使用async的结构化并发

让我们使用使用 async 的并发这一小节的例子并且提取出一个函数并发的调用 doSomethingUsefulOne 与 doSomethingUsefulTwo 并且返回它们两个的结果之和。 由于 async 被定义为了 CoroutineScope 上的扩展，我们需要将它写在作用域内，并且这是 coroutineScope 函数所提供的：

```kotlin
suspend fun concurrentSum(): Int = coroutineScope {
    val one = async { doSomethingUsefulOne() }
    val two = async { doSomethingUsefulTwo() }
     one.await() + two.await()
}
```

这种情况下，如果在 concurrentSum 函数内部发生了错误，并且它抛出了一个异常， 所有在作用域中启动的协程都将会被取消。

```kotlin
import kotlinx.coroutines.*
import kotlin.system.*

fun main() = runBlocking<Unit> {
    val time = measureTimeMillis {
        println("The answer is ${concurrentSum()}")
    }
    println("Completed in $time ms")    
}

suspend fun concurrentSum(): Int = coroutineScope {
    val one = async { doSomethingUsefulOne() }
    val two = async { doSomethingUsefulTwo() }
     one.await() + two.await()
}

suspend fun doSomethingUsefulOne(): Int {
    delay(1000L) // 假设我们在这里做了些有用的事
    return 13
}

suspend fun doSomethingUsefulTwo(): Int {
    delay(1000L) // 假设我们在这里也做了些有用的事
    return 29
}
```

从上面的 main 函数的输出可以看出，我们仍然可以同时执行这两个操作：

```kotlin
The answer is 42
Completed in 1017 ms
```

取消始终通过协程的层次结构来进行传递：

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking<Unit> {
    try {
        failedConcurrentSum()
    } catch(e: ArithmeticException) {
        println("Computation failed with ArithmeticException")
    }
}

suspend fun failedConcurrentSum(): Int = coroutineScope {
    val one = async<Int> { 
        try {
            delay(Long.MAX_VALUE) // 模拟一个长时间的运算
            42
        } finally {
            println("First child was cancelled")
        }
    }
    val two = async<Int> { 
        println("Second child throws an exception")
        throw ArithmeticException()
    }
        one.await() + two.await()
}
```

注意，当第一个子协程失败的时候第一个 async 是如何等待父线程被取消的：

```
Second child throws an exception
First child was cancelled
Computation failed with ArithmeticException
```