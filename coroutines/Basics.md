**内容目录**

* [协程基础](#协程基础)
  * [你的第一个协程](#你的第一个协程)
  * [桥接阻塞与非阻塞的世界](#桥接阻塞与非阻塞的世界)
  * [等待任务](#等待任务)
  * [结构化并发](#结构化并发)
  * [作用域构建器](#作用域构建器)
  * [提取函数重构](#提取函数重构)
  * [协程是轻量级的](#协程是轻量级的)
  * [全局协程就像一个守护线程](#全局协程就像一个守护线程)


## 协程基础

这部分将包含协程的基础概念。

### 你的第一个协程

运行下面的代码：

```kotlin
import kotlinx.coroutines.*

fun main() {
    GlobalScope.launch { // launch a new coroutine in background and continue
        delay(1000L) // non-blocking delay for 1 second (default time unit is ms)
        println("World!") // print after delay
    }
    println("Hello,") // main thread continues while coroutine is delayed
    Thread.sleep(2000L) // block main thread for 2 seconds to keep JVM alive
}
```

你会得到如下结果：

```text
Hello,
World!
```
本质上讲，协程是轻量级的线程。它们由 `launch` 协程构建器在对应的 `[CoroutineScope](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/coroutine-scope.html)` 上下文中启动。在本例中是 `[GlobalScope](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-global-scope/index.html)` , 意味着新协程的生命周期受限于整个应用生命周期。

你可以用
`GlobalScope.launch { …… }` 替换为 `thread { …… }`，将 `delay(……)` 替换为 `Thread.sleep(……)` 达到同样目的。

如果你用 `GlobalScope.launch` 替换为 `thread`，编译器会报以下错误：

```
Error: Kotlin: Suspend functions are only allowed to be called from a coroutine or another suspend function
```

这是因为 [delay] 是一个特殊的 _挂起函数_ ，它不会造成线程阻塞，但是会 _挂起_ 协程，并且只能在协程中使用。

### 桥接阻塞与非阻塞的世界

第一个例子混合了非阻塞 `delay(...)` 和阻塞 `Thread.sleep(...)` 。这会让人搞混哪个是阻塞哪个是非阻塞。下面用 `[runBlocking](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/run-blocking.html)` 协程构建器来说明什么是阻塞：

```kotlin
import kotlinx.coroutines.*

fun main() { 
    GlobalScope.launch { // launch a new coroutine in background and continue
        delay(1000L)
        println("World!")
    }
    println("Hello,") // main thread continues here immediately
    runBlocking {     // but this expression blocks the main thread
        delay(2000L)  // ... while we delay for 2 seconds to keep JVM alive
    } 
}
```

结果是一样的，但是代码只用了非阻塞的函数 delay。调用了 `runBlocking` 的主线程会一直阻塞直到 `runBlocking` 内部的协程执行完毕。

这个例子还可以改写为更加惯用的方式。使用 `runBlocking` 包装主函数的执行：

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking<Unit> { // start main coroutine
    GlobalScope.launch { // launch a new coroutine in background and continue
        delay(1000L)
        println("World!")
    }
    println("Hello,") // main coroutine continues here immediately
    delay(2000L)      // delaying for 2 seconds to keep JVM alive
}
```

这里的 `runBlocking<Unit> {...}` 作为用来启动顶级主协程的适配器。显示声明了返回值是 `Unit`，因为 Kotlin 的 main 函数返回值是 Unit

也可以这样给挂起函数写单元测试：

```kotlin
class MyTest {
    @Test
    fun testMySuspendingFunction() = runBlocking<Unit> {
        // here we can use suspending functions using any assertion style that we like
    }
}
```

### 等待一个任务

指定延迟时间去等待另一个协程结束并不是一个好办法。我们可以显式的等待一个`[Job](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-job/index.html)` 是否完成：

```kotlin
val job = GlobalScope.launch { // launch a new coroutine and keep a reference to its Job
    delay(1000L)
    println("World!")
}
```

结果仍然是一致的，但主协程与后台任务的执行时间解绑了，代码更加整洁。

### 结构化并发

在使用协程时还需要些东西.当用 `GlobalScope.launch` 会创建一个顶级协程.尽管协程是轻量级的,运行时仍然会带来内存资源的消耗.如果我们忘记对新启动协程的引用,他将一直运行下去. 如果协程中的代码挂起了(比如我们delay 时间过长),万一启动过多协程而内存不够了怎么办? 手动保持所有已启动的协程的引用,以及调用`join`来避免出错.

这里有更好的解决方案.可以借助结构化并发来解决。我们可以在执行操作所在的指定作用域内启动协程，而不是像使用线程（线程总是全局的）那样在 [GlobalScope] 中启动.

在下面的例子中,主函数被`runBlocking`协程构建器转换为协程.每个协程构建器,包括`runBlocking`都会添加一个 `CoroutineScope` 到自己的代码块中.我们可以在这个范围内启动协程而不用显示调用`join`,因为外部协程(本例中的`runBlocking`) 直到它启动的所有协程均执行结束才会结束.因此,我们可以把代码写的更加简洁:

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking { // this: CoroutineScope
    launch { // launch a new coroutine in the scope of runBlocking
        delay(1000L)
        println("World!")
    }
    println("Hello,")
}
```

### 作用域构建器

除了由不同构建器提供的协程作用域之外，还可以使用coroutineScope构建器声明自己的作用域。 它会创建新的协程范围，并且在所有已启动的子项完成之前不会完成。 

runBlocking和coroutineScope 看起来很相似，他们都会等自身作用域内代码和子协程完成后结束。他们之间的主要区别在 runBlocking 会这阻塞当前线程并等待，而 coroutineScope 只会挂起，并释放当前线程给其它人用。由于这个不同，runBlocking 是一个普通函数，而 coroutineScope 是一个挂起函数。

下面等例子演示了区别。

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking { // this: CoroutineScope
    launch { 
        delay(200L)
        println("Task from runBlocking")
    }
    
    coroutineScope { // Creates a coroutine scope
        launch {
            delay(500L) 
            println("Task from nested launch")
        }
    
        delay(100L)
        println("Task from coroutine scope") // launch 块内结束前就打印
    }
    
    println("Coroutine scope is over") // 直到 launch 块内完成才会打印
}
```

注意这里 "Task from coroutine scope" 消息打印后，在等待内嵌 launch 时 “Task from runBlocking” 将会执行并且打印，尽管 coroutineScope 还没有结束。


### 提取函数重构

让我们将launch {...}中的代码块提取到一个单独的函数中。 当对此代码执行“Extract function”重构时，需要创建一个带有suspend修饰符的新函数。 这是你的第一个挂起函数。 挂起函数可以在协同程序内部使用，就像普通函数一样，但它们的附加功能是它们可以使用其他挂起函数（例如本例中的延迟）来挂起协程的执行。

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking {
    launch { doWorld() }
    println("Hello,")
}

// this is your first suspending function
suspend fun doWorld() {
    delay(1000L)
    println("World!")
}
```

但如果提取的函数包含在当前作用域上调用的协程构建器，该怎么办？ 这种情况下，只有suspend修饰符是不够的。 在CoroutineScope上添加 doWorld 扩展方法是其中一个解决方案，但这并不是一个好的方式，因为它不会使API更清晰。 惯用解决方案是将显式CoroutineScope作为包含目标函数的类中的字段，或者在外部类实现CoroutineScope达到隐式实现。 作为最后的手段，可以使用 [CoroutineScope（coroutineContext）](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-coroutine-scope.html)，但是这种方法在结构上是不安全的，这样你将不再能够控制此方法的执行范围。 只有私有API才能使用此构建器。

### 协程是轻量级的

运行下面的代码:

```kotlin
import kotlinx.coroutines.*

fun main() = runBlocking {
    repeat(100_000) { // launch a lot of coroutines
        launch {
            delay(1000L)
            print(".")
        }
    }
}
```

它会启动100K协程,每个协程在一秒后打印 ".". 如果换成线程去实现,会发生什么呢? (最可能得是出现内存不足的错误)

### 全局协程就像一个守护线程

下面的代码在GlobalScope中启动一个长时间运行的协程，它会每秒打印“I'm sleeping”两次，然后在一段时延后返回：

```kotlin
GlobalScope.launch {
    repeat(1000) { i ->
            println("I'm sleeping $i ...")
        delay(500L)
    }
}
delay(1300L) // just quit after delay
```

试着运行下代码他会打印三行然后结束:

```
I'm sleeping 0 ...
I'm sleeping 1 ...
I'm sleeping 2 ...
```

在GlobalScope中启动的协程不会使进程保持活动状态。它们就像守护线程（当非守护线程结束后守护线程就会结束）。