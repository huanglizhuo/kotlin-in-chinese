- [共享的可变状态和并发](#共享的可变状态和并发)
    - [问题](#问题)
    - [Volatiles不起作用](#Volatiles不起作用)
    - [线程安全的数据结构](#线程安全的数据结构)
    - [细粒度线程限制](#细粒度线程限制)
    - [粗粒度线程限制](#粗粒度线程限制)
    - [互斥](#互斥)
    - [Actors](#Actors)

## 共享的可变状态和并发

协程可以使用诸如Dispatchers.Default之类的多线程调度程序并发执行。 它代表了所有常见的并发问题。 主要问题是同步访问共享可变状态。 协程领域中针对此问题的某些解决方案与多线程世界中的解决方案相似，但其他解决方案却是独一无二的。

### 问题

让我们启动一百个协程，它们都执行相同的动作数千次。 我们还将测量它们的完成时间以进行进一步的比较：

```Kotlin
suspend fun massiveRun(action: suspend () -> Unit) {
    val n = 100  // number of coroutines to launch
    val k = 1000 // times an action is repeated by each coroutine
    val time = measureTimeMillis {
        coroutineScope { // scope for coroutines 
            repeat(n) {
                launch {
                    repeat(k) { action() }
                }
            }
        }
    }
    println("Completed ${n * k} actions in $time ms")    
}
```

我们从一个非常简单的操作开始，该操作使用多线程Dispatchers.Default增加共享的可变变量。


```Kotlin
var counter = 0

fun main() = runBlocking {
    withContext(Dispatchers.Default) {
        massiveRun {
            counter++
        }
    }
    println("Counter = $counter")
}
```

最后他会输出什么呢? 很大可能上不是 "Counter - 100000", 因为上百个协程会在多线中非同步的对 `counter` 进行并发操作.

### Volatiles不起作用

常见的误解是，标记变量为 `volatile` 可解决并发问题。 让我们尝试一下：

```Kotlin
@Volatile // in Kotlin `volatile` is an annotation 
var counter = 0

fun main() = runBlocking {
    withContext(Dispatchers.Default) {
        massiveRun {
            counter++
        }
    }
    println("Counter = $counter")
}
```

该代码的运行速度较慢，但是最后我们仍然无法获得“ Counter = 100000”，因为 volatile 变量可确保线性化（ linearizable 这是“原子”的技术术语）可读写相应的变量，但不提供原子性较大的动作（在我们的情况下为增加）。

### 线程安全的数据结构

适用于线程和协程的通用解决方案是使用线程安全（aka同步，线性化或原子）数据结构，该结构为需要在共享状态下执行的相应操作提供所有必需的同步。 对于简单的计数器，我们可以使用 `AtomicInteger` 类，该类具有原子级 `incrementAndGet` 操作：

```Kotlin
var counter = AtomicInteger()

fun main() = runBlocking {
    withContext(Dispatchers.Default) {
        massiveRun {
            counter.incrementAndGet()
        }
    }
    println("Counter = $counter")
}
```

这是针对此特定问题的最快解决方案。 它适用于简单计数器，集合，队列和其他标准数据结构以及对其的基本操作。 但是，它不容易扩展到复杂状态或者复杂操作,它们并没有一个可以立即使用的线程安全的实现.

### 细粒度线程限制

线程限制是一种解决共享可变状态的方法，其中对特定共享状态的所有访问都限于一个线程。它通常用于UI应用程序中，其中所有UI状态都限制在单个事件调度/应用程序线程中。通过使用
单线程上下文。

```Kotlin
val counterContext = newSingleThreadContext("CounterContext")
var counter = 0

fun main() = runBlocking {
    withContext(Dispatchers.Default) {
        massiveRun {
            // confine each increment to a single-threaded context
            withContext(counterContext) {
                counter++
            }
        }
    }
    println("Counter = $counter")
}
```

这段代码非常慢，因为它可以进行细粒度的线程约束。每个单独的增量都使用withContext（counterContext）块从多线程 Dispatchers.Default 上下文切换到单线程上下文。

### 粗粒度线程限制

实际上，线程限制是大块执行的，例如状态更新业务逻辑的大部分都局限于单个线程中。下面的示例就是这样做的，首先在单线程上下文中运行每个协程。

```Kotlin
val counterContext = newSingleThreadContext("CounterContext")
var counter = 0

fun main() = runBlocking {
    // confine everything to a single-threaded context
    withContext(counterContext) {
        massiveRun {
            counter++
        }
    }
    println("Counter = $counter")
}
```

现在这可以更快地工作并产生正确的结果。

### 互斥

解决该问题的互斥解决方案是使用永远不会并行的临界区来保护共享状态的所有修改。在阻塞的世界中，通常会为此使用 `synchronized` 或 `ReentrantLock` 。协程的替代品称为 [Mutex](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.sync/-mutex/index.html) 。它具有 [lock](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.sync/-mutex/lock.html) 和 [unlock](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.sync/-mutex/unlock.html) 函数来界定临界区。关键区别在于Mutex.lock（）是一个可挂起函数。它不会阻塞线程。

还有一个 [withLock](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.sync/with-lock.html) 扩展函数，可以方便地表示 `mutex.lock(); try { ... } finally { mutex.unlock() }` 模式：

```Kotlin
val mutex = Mutex()
var counter = 0

fun main() = runBlocking {
    withContext(Dispatchers.Default) {
        massiveRun {
            // protect each increment with lock
            mutex.withLock {
                counter++
            }
        }
    }
    println("Counter = $counter")
}
```

此示例中的锁定是细粒度的，因此要付出代价。 但是，在某些情况下，须定期修改某些共享状态，但是没有限制该状态的自然线程，这是一个不错的选择。

### Actors

[actor](https://en.wikipedia.org/wiki/Actor_model) 是一个结合协程创建的实体，该协程限制并封装的状态以及与其他协程通信的通道组成的实体。可以将简单的actor编写为函数，但是状态复杂的actor更适合用类表示。

[actor](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/actor.html) 协程生成器，可以方便地将actor的邮箱通道合并到其作用域中，以从中接收消息，并将send通道合并到结果 job 对象中，actor的单个引用作为其句柄。

使用actor的第一步是定义actor将要处理的消息类。 Kotlin的 [sealed 类](https://kotlinlang.org/docs/reference/sealed-classes.html) 非常适合该目的。我们定义 CounterMsg 密封类,IncCounter 消息以增加计数器，GetCounter 消息获取其值。后者需要发送响应。为此，此处使用了 [CompletableDeferred](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-completable-deferred/index.html) 通信原语，该原语表示将来将要知道（传递）的单个值。

```Kotlin
// Message types for counterActor
sealed class CounterMsg
object IncCounter : CounterMsg() // one-way message to increment counter
class GetCounter(val response: CompletableDeferred<Int>) : CounterMsg() // a request with reply
```

然后，我们定义一个使用actor协程生成器启动actor的函数：

```Kotln
// This function launches a new counter actor
fun CoroutineScope.counterActor() = actor<CounterMsg> {
    var counter = 0 // actor state
    for (msg in channel) { // iterate over incoming messages
        when (msg) {
            is IncCounter -> counter++
            is GetCounter -> msg.response.complete(counter)
        }
    }
}
```

主要代码很简单：

```Kotlin
fun main() = runBlocking<Unit> {
    val counter = counterActor() // create the actor
    withContext(Dispatchers.Default) {
        massiveRun {
            counter.send(IncCounter)
        }
    }
    // send a message to get a counter value from an actor
    val response = CompletableDeferred<Int>()
    counter.send(GetCounter(response))
    println("Counter = ${response.await()}")
    counter.close() // shutdown the actor
}
```

正确执行 actor 本身在什么上下文中都没有关系（正确性）。actor 是一个协程,而协程是顺序执行的，因此将状态限制为特定协程可以解决共享可变状态的问题。实际上，参与者可以修改自己的私有状态，但只能通过消息相互影响（避免使用任何锁）。

Actor比在负载锁定更有效，因为在这种情况下，Actor总是有工作要做，并且根本不必切换到其他上下文。

请注意， [actor](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/actor.html) 协程构建器是 [produce](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/produce.html) 协程构建器的对偶。actor 有一个接收消息的通道相关联，而生产者与其发送元素的通道相关联。