**目录**

* [取消与超时](#取消与超时)
  * [取消协程的执行](#取消协程的执行)
  * [取消是协作的](#取消是协作的)
  * [使计算代码可取消](#使计算代码可取消)
  * [在 finally 中释放资源](#在-finally-中释放资源)
  * [运行不可取消的代码块](#运行不可取消的代码块)
  * [超时](#超时)

## 取消与超时

这一部分包含了协程的取消与超时。

### 取消协程的执行

在长时间运行的应用程序中，你可能需要对后台协程进行细粒度控制。 例如，用户可能已关闭启动协程的页面，现在不再需要其结果，并且该操作是可取消的。`launch` 函数会返回一个可用于取消运行协程的 `job`：

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking {
    val job = launch {
        repeat(1000) { i ->
                println("I'm sleeping $i ...")
            delay(500L)
        }
    }
    delay(1300L) // delay a bit
    println("main: I'm tired of waiting!")
    job.cancel() // cancels the job
    job.join() // waits for job's completion 
    println("main: Now I can quit.")    
}
```

运行结果如下:

```kotlin
I'm sleeping 0 ...
I'm sleeping 1 ...
I'm sleeping 2 ...
main: I'm tired of waiting!
main: Now I can quit.
```

job.cancel 调用后,其它协程将不能从它获取任何结果,因为该协程已经取消.还有一个Job扩展函数[cancelAndJoin](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/cancel-and-join.html)，它结合cancel和join调用。

### 取消是协作的

协程的取消是协作的,协程代码必须合作才能取消. `kotlinx.coroutines`中所有的挂起函数都是可取消的.它们会检查协程是否可取消,若不可取消则抛出 `[CancellationException](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-cancellation-exception/index.html)` 异常. 然而如果协程正在进行计算并且没有检查可取消性, 那么它是不可取消的,比如下面的例子:

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking {
    val startTime = System.currentTimeMillis()
    val job = launch(Dispatchers.Default) {
        var nextPrintTime = startTime
        var i = 0
        while (i < 5) { // computation loop, just wastes CPU
            // print a message twice a second
            if (System.currentTimeMillis() >= nextPrintTime) {
                println("I'm sleeping ${i++} ...")
                nextPrintTime += 500L
            }
        }
    }
    delay(1300L) // delay a bit
    println("main: I'm tired of waiting!")
    job.cancelAndJoin() // cancels the job and waits for its completion
    println("main: Now I can quit.")    
}
```

可以试试它是否会在取消后继续打印“I'm sleeping”，直到作业在五次迭代后自行完成。

### 使计算代码可取消 

有两种方法可以使计算代码可以取消。 一种是定期调用检查取消的挂起功能。 [yield](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/yield.html)函数是一个很好的选择。 另一种是明确检查取消状态。 本例使用后一种方法:

把 `while (i < 5)` 改为 `while (isActive)` 并运行

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking {
    val startTime = System.currentTimeMillis()
    val job = launch(Dispatchers.Default) {
        var nextPrintTime = startTime
        var i = 0
        while (isActive) { // cancellable computation loop
            // print a message twice a second
            if (System.currentTimeMillis() >= nextPrintTime) {
                println("I'm sleeping ${i++} ...")
                nextPrintTime += 500L
            }
        }
    }
    delay(1300L) // delay a bit
    println("main: I'm tired of waiting!")
    job.cancelAndJoin() // cancels the job and waits for its completion
    println("main: Now I can quit.")    
}
```

正如你所看到的,现在循环可以取消. `[isActive](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/is-active.html)` 是一个扩展属性,可通过CoroutineScope对象在coroutine代码中使用。

### 在 finally 中释放资源

可取消的挂起函数会在取消时抛出CancellationException，这可以通过常规方式处理。 例如，try{...} finally {...} 表达式或者Kotlin `use` 函数在取消协程时正常执行其终止操作：

```kotlin
val job = launch {
    try {
        repeat(1000) { i ->
                println("I'm sleeping $i ...")
            delay(500L)
        }
    } finally {
        println("I'm running finally")
    }
}
delay(1300L) // delay a bit
println("main: I'm tired of waiting!")
job.cancelAndJoin() // cancels the job and waits for its completion
println("main: Now I can quit.")
```

join和cancelAndJoin都等待所有终结操作完成，因此上面的示例生成以下输出：

```
I'm sleeping 0 ...
I'm sleeping 1 ...
I'm sleeping 2 ...
main: I'm tired of waiting!
I'm running finally
main: Now I can quit.
```

### 运行不可取消的代码块

在前一个示例的finally块中使用挂起函数的任何尝试都会导致CancellationException，因为取消了运行此代码的协程。 通常，这没关系，因为所有表现良好的关闭操作（关闭文件，取消作业或关闭任何类型的通信通道）通常都是非阻塞的，并且不涉及任何挂起函数。 但是，在极少数情况下，当您需要在取消的协同程序中挂起时，可以使用[withContext](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/with-context.html)函数和[NonCancellable](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-non-cancellable.html)上下文将相应的代码包装在withContext（NonCancellable）{...}中，如下例所示：

```kotlin
val job = launch {
    try {
        repeat(1000) { i ->
                println("I'm sleeping $i ...")
            delay(500L)
        }
    } finally {
        withContext(NonCancellable) {
            println("I'm running finally")
            delay(1000L)
            println("And I've just delayed for 1 sec because I'm non-cancellable")
        }
    }
}
delay(1300L) // delay a bit
println("main: I'm tired of waiting!")
job.cancelAndJoin() // cancels the job and waits for its completion
println("main: Now I can quit.")
```

### 超时

在实践中取消协程执行的最主要的原因是因为它的执行时间超过了限制。 虽然可以手动跟踪对相应作业的引用并启动单独的协程以在延迟后取消跟踪的协程，但[withTimeout](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/with-timeout.html)函数是一个开箱即用的操作。 请看以下示例：

```kotlin
withTimeout(1300L) {
    repeat(1000) { i ->
            println("I'm sleeping $i ...")
        delay(500L)
    }
}
```

输出结果如下:

```
I'm sleeping 0 ...
I'm sleeping 1 ...
I'm sleeping 2 ...
Exception in thread "main" kotlinx.coroutines.TimeoutCancellationException: Timed out waiting for 1300 ms
```

withTimeout抛出的TimeoutCancellationException是CancellationException的子类。 我们之前没有看到它的堆栈跟踪打印在控制台上。 这是因为在取消的协程中，CancellationException被认为是协程完成的正常原因。 但是，在这个例子中，我们在main函数中使用了withTimeout。

因为取消只是一个异常，所有资源都以通常的方式关闭。 如果你可以在任何类型的超时上做一些额外的操作或者使用类似于withTimeout的withTimeoutOrNull函数，你可以在try {...} catch（e：TimeoutCancellationException）{...}块中用超时包装代码， 这样在超时时将返回null而不是抛出异常：


```kotlin
val result = withTimeoutOrNull(1300L) {
    repeat(1000) { i ->
            println("I'm sleeping $i ...")
        delay(500L)
    }
    "Done" // will get cancelled before it produces this result
}
println("Result is $result")
```

这样运行以上代码就不会抛出异常了:

```
I'm sleeping 0 ...
I'm sleeping 1 ...
I'm sleeping 2 ...
Result is null
```