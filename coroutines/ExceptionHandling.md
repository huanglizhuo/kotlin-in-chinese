目录

- [异常处理](#异常处理)
    -[异常传播](#异常传播)
    -[协程异常处理器](#协程异常处理器)
    -[取消和异常](#取消和异常)
    -[异常聚合](#异常聚合)
    -[监管](#监管)
        -[监管] job(#监管)
        -[监管作用域](#监管作用域)
        -[监管协程中的异常](#监管协程中的异常)

## 异常处理

本节介绍异常处理和异常的取消。我们已经知道，取消的协程会在挂起点上引发 [CancellationException](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-cancellation-exception/index.html)，并且协程机制会忽略它。但是，如果在取消过程中引发异常或同一个协程的多个子协程发异常，会发生什么呢？

### 异常传播

协程构建器有两种形式：自动传播异常（[launch](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/launch.html)和[actor](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/launch.html)）或向用户暴露异常（[async](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/async.html)和[produce](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/produce.html)）。前者将异常视为未处理的异常，类似于Java的Thread.uncaughtExceptionHandler，而后者则依靠用户使用最终异常，例如通过 [await](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-deferred/await.html) 或 [receive](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/-receive-channel/receive.html)（produce和receive将在Channels部分中介绍）。

可以借助GlobalScope中创建协程的简单示例来演示：

```Kotlin
import kotlinx.coroutines.*

fun main() = runBlocking {
    val job = GlobalScope.launch {
        println("Throwing exception from launch")
        throw IndexOutOfBoundsException() // Will be printed to the console by Thread.defaultUncaughtExceptionHandler
    }
    job.join()
    println("Joined failed job")
    val deferred = GlobalScope.async {
        println("Throwing exception from async")
        throw ArithmeticException() // Nothing is printed, relying on user to call await
    }
    try {
        deferred.await()
        println("Unreached")
    } catch (e: ArithmeticException) {
        println("Caught ArithmeticException")
    }
}
```

这段代码的输出如下(附带调试信息)：

```shell
Throwing exception from launch
Exception in thread "DefaultDispatcher-worker-2 @coroutine#2" java.lang.IndexOutOfBoundsException
Joined failed job
Throwing exception from async
Caught ArithmeticException
```

### 协程异常处理器
但是，如果不想将所有异常打印到控制台怎么办？ [CoroutineExceptionHandler](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-coroutine-exception-handler/index.html) 上下文元素可以视为协程的通用 `catch` 块，在协程中可能发生自定义日志记录或异常处理。它类似于使用Thread.uncaughtExceptionHandler。

在JVM上，可以通过ServiceLoader注册CoroutineExceptionHandler来为所有协程重新定义全局异常处理程序。全局异常处理程序与Thread.defaultUncaughtExceptionHandler相似，在没有其他特定的处理程序注册时使用。在Android上，uncaughtExceptionPreHandler被安装为全局协程异常处理程序。

仅在用户未处理的异常上会调用CoroutineExceptionHandler，因此在[async](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/async.html)生成器等中注册它无效。

```Kotlin
val handler = CoroutineExceptionHandler { _, exception -> 
    println("Caught $exception") 
}
val job = GlobalScope.launch(handler) {
    throw AssertionError()
}
val deferred = GlobalScope.async(handler) {
    throw ArithmeticException() // Nothing will be printed, relying on user to call deferred.await()
}
joinAll(job, deferred)
```

此代码的输出是：

```shell
Caught java.lang.AssertionError

```

### 取消和异常

取消与异常紧密相关。协程在内部使用CancellationException进行取消，所有处理程序都将忽略这些异常，因此它们仅应用作其他调试信息的源，可以通过catch块获取这些信息。使用 [Job.cancel](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-job/cancel.html) 取消协程时，协程终止，但不取消其父级。

```Kotlin
val job = launch {
    val child = launch {
        try {
            delay(Long.MAX_VALUE)
        } finally {
            println("Child is cancelled")
        }
    }
    yield()
    println("Cancelling child")
    child.cancel()
    child.join()
    yield()
    println("Parent is not cancelled")
}
job.join()
```

此代码的输出是：

```shell
Cancelling child
Child is cancelled
Parent is not cancelled
```

如果协程遇到除CancellationException以外的其他异常，它将取消具有该异常的父对象。此行为不能被覆盖，并且用于为不依赖于CoroutineExceptionHandler实现的结构化并发提供稳定的协程层次结构。当父级的所有子级终止时，父级将处理原始异常。

这也是为什么在这些示例中，CoroutineExceptionHandler始终注册到在GlobalScope中创建的协程的原因。将异常处理程序注册到在主runBlocking范围内启动的协程中没有意义，因为尽管安装了该处理程序，但当其子级异常完成时，主协程将始终被取消。

```Kotlin
val handler = CoroutineExceptionHandler { _, exception -> 
    println("Caught $exception") 
}
val job = GlobalScope.launch(handler) {
    launch { // the first child
        try {
            delay(Long.MAX_VALUE)
        } finally {
            withContext(NonCancellable) {
                println("Children are cancelled, but exception is not handled until all children terminate")
                delay(100)
                println("The first child finished its non cancellable block")
            }
        }
    }
    launch { // the second child
        delay(10)
        println("Second child throws an exception")
        throw ArithmeticException()
    }
}
job.join()
```

此代码的输出是：

```shell
Second child throws an exception
Children are cancelled, but exception is not handled until all children terminate
The first child finished its non cancellable block
Caught java.lang.ArithmeticException
```

### 异常聚合

如果协程的多个子级抛出异常会怎样？一般规则是“第一个异常获胜”，因此第一个引发的异常向处理器暴露。但这可能导致丢失的异常，例如，协程在其 `finally` 块中抛出异常。因此，抑制了其他异常。

解决方案之一是分别报告每个异常，但是Deferred.await应该具有相同的机制来避免行为不一致，这将导致协程的实现细节（无论它是否将工作的一部分委派给了孩子或不）泄漏到其异常处理程序。

```Kotlin
import kotlinx.coroutines.*
import java.io.*

fun main() = runBlocking {
    val handler = CoroutineExceptionHandler { _, exception ->
        println("Caught $exception with suppressed ${exception.suppressed.contentToString()}")
    }
    val job = GlobalScope.launch(handler) {
        launch {
            try {
                delay(Long.MAX_VALUE)
            } finally {
                throw ArithmeticException()
            }
        }
        launch {
            delay(100)
            throw IOException()
        }
        delay(Long.MAX_VALUE)
    }
    job.join()  
}
```

> 注意：以上代码仅在支持抑制异常的JDK7 +上才能正常工作


此代码的输出是：

```shell
Caught java.io.IOException with suppressed [java.lang.ArithmeticException]
```

> 请注意，该机制当前仅在Java版本1.7+上有效。 JS和Native的限制是暂时的，将来会修复。

取消异常是透明的，默认情况下是未包装的：

```Kotlin
val handler = CoroutineExceptionHandler { _, exception ->
    println("Caught original $exception")
}
val job = GlobalScope.launch(handler) {
    val inner = launch {
        launch {
            launch {
                throw IOException()
            }
        }
    }
    try {
        inner.join()
    } catch (e: CancellationException) {
        println("Rethrowing CancellationException with original cause")
        throw e
    }
}
job.join()
```

此代码的输出是：

```shell
Rethrowing CancellationException with original cause
Caught original java.io.IOException
```

### 监管

正如我们之前研究的那样，取消是在整个协程层次中传播的双向关系。但是，如果需要单向取消怎么办？

此类需求的一个很好的例子是在其范围内定义了工作的UI组件。如果UI的任何子任务失败，则不必总是取消（有效地杀死）整个UI组件，但是如果UI组件被破坏（并且其工作被取消），则必须使所有子任务失败，因为他们的结果不再需要。

另一个示例是一个服务器进程，该进程产生多个子作业，并且需要监督它们的执行，跟踪其失败并仅重新启动那些失败的子作业。

### 监管 job

为此，可以使用[SupervisorJob](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-supervisor-job.html)。 它与常规Job相似，唯一的不同是取消仅向下传播。 举一个例子：

```Kotlin
import kotlinx.coroutines.*

fun main() = runBlocking {
    val supervisor = SupervisorJob()
    with(CoroutineScope(coroutineContext + supervisor)) {
        // launch the first child -- its exception is ignored for this example (don't do this in practice!)
        val firstChild = launch(CoroutineExceptionHandler { _, _ ->  }) {
            println("First child is failing")
            throw AssertionError("First child is cancelled")
        }
        // launch the second child
        val secondChild = launch {
            firstChild.join()
            // Cancellation of the first child is not propagated to the second child
            println("First child is cancelled: ${firstChild.isCancelled}, but second one is still active")
            try {
                delay(Long.MAX_VALUE)
            } finally {
                // But cancellation of the supervisor is propagated
                println("Second child is cancelled because supervisor is cancelled")
            }
        }
        // wait until the first child fails & completes
        firstChild.join()
        println("Cancelling supervisor")
        supervisor.cancel()
        secondChild.join()
    }
}
```

此代码的输出是：

```shell
First child is failing
First child is cancelled: true, but second one is still active
Cancelling supervisor
Second child is cancelled because supervisor is cancelled
```

### 监管作用域

对于有作用域的并发，出于相同的目的，可以使用 [supervisorScope](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/supervisor-scope.html) 代替 [coroutineScope](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/coroutine-scope.html)。它仅在一个方向上传播取消，并且仅在失败后才取消所有子级。它也像coroutineScope一样等待所有孩子完成任务。

```Kotlin
import kotlin.coroutines.*
import kotlinx.coroutines.*

fun main() = runBlocking {
    try {
        supervisorScope {
            val child = launch {
                try {
                    println("Child is sleeping")
                    delay(Long.MAX_VALUE)
                } finally {
                    println("Child is cancelled")
                }
            }
            // Give our child a chance to execute and print using yield 
            yield()
            println("Throwing exception from scope")
            throw AssertionError()
        }
    } catch(e: AssertionError) {
        println("Caught assertion error")
    }
}
```

此代码的输出是：

```shell
Child is sleeping
Throwing exception from scope
Child is cancelled
Caught assertion error
```

### 监管协程中的异常

常规job和监管job之间的另一个关键区别是异常处理。每个子job都应通过异常处理机制自行处理其异常。这种差异来自于子job的失败不会传播给父job.

```Kotlin
import kotlin.coroutines.*
import kotlinx.coroutines.*

fun main() = runBlocking {
    val handler = CoroutineExceptionHandler { _, exception -> 
        println("Caught $exception") 
    }
    supervisorScope {
        val child = launch(handler) {
            println("Child throws an exception")
            throw AssertionError()
        }
        println("Scope is completing")
    }
    println("Scope is completed")
}
```

此代码的输出是：

```shell
Scope is completing
Child throws an exception
Caught java.lang.AssertionError
Scope is completed
```