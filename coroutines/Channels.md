**目录**

* [通道（实验性的）](#通道实验性的)
  * [通道基础](#通道基础)
  * [通道的关闭与迭代](#通道的关闭与迭代)
  * [构建通道生产者](#构建通道生产者)
  * [管道](#管道)
  * [使用管道生产素数](#使用管道生产素数)
  * [扇出](#扇出)
  * [扇入](#扇入)
  * [带缓冲的通道](#带缓冲的通道)
  * [通道是公平的](#通道是公平的)
  * [计时器通道](#计时器通道)


### 通道（实验性的)

延迟值提供了在协程之间传输单个值的便捷方法。通道提供了一种传输值流的方法。

> 通道是kotlinx.coroutines的实验性功能。他们的API预计将在即将发布的kotlinx.coroutines库更新中发展，并且可能会发生重大变化。

### 通道基础

[Channel](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/-channel/index.html)在概念上与BlockingQueue非常相似。一个关键的区别是队列的 put 操作是非阻塞的，而channel 的 [send](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/-send-channel/send.html) 操作是可挂起的，队列的 take 操作是阻塞的,channel 的 [reveive](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/-receive-channel/receive.html) 是可挂起的。

```kotlin
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*

fun main() = runBlocking {
    val channel = Channel<Int>()
    launch {
        // this might be heavy CPU-consuming computation or async logic, we'll just send five squares
        for (x in 1..5) channel.send(x * x)
    }
    // here we print five received integers:
    repeat(5) { println(channel.receive()) }
    println("Done!")
}
```

代码的输出如下:

```
1
4
9
16
25
Done!
```

### 通道的关闭与迭代

与队列不同,通道可以关闭以表示没有即将到来的元素. 在接收者一边,可以使用普通的 `for` 循环从通道中读取元素.

概念上讲, [close](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/-send-channel/close.html) 类似于给通道发送一个特殊的关闭令牌. 一旦收到结束令牌循环既停止,因此可以保证前面发送的元素都被接收到了:

```kotlin
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*

fun main() = runBlocking {
    val channel = Channel<Int>()
    launch {
        for (x in 1..5) channel.send(x * x)
        channel.close() // we're done sending
    }
    // here we print received values using `for` loop (until the channel is closed)
    for (y in channel) println(y)
    println("Done!")
```

### 构建通道生产者

协程生产元素队列的模式很常见. 这是生产者-消费者模式的一部分,在并发代码中十分常见. 你把生产者抽象为一个接收通道作为参数的一个函数,但这与编程常识有些不同,因为通常结果必须从函数中返回.

[produce](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/produce.html) 是一个很方便的协程生产者,在生产者端可以很好的工作， 并且扩展函数 [consumeEach](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/consume-each.html) 在消费者端可以很方便的替代 for 循环：

```kotlin
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*

fun CoroutineScope.produceSquares(): ReceiveChannel<Int> = produce {
    for (x in 1..5) send(x * x)
}

fun main() = runBlocking {
    val squares = produceSquares()
    squares.consumeEach { println(it) }
    println("Done!")
}
```

### 管道

管道是一种协程生成的模式，可以生产无穷多个元素：

```kotlin
fun CoroutineScope.produceNumbers() = produce<Int> {
    var x = 1
    while (true) send(x++) // infinite stream of integers starting from 1
}
```

其它的协程可以消费该流,做些操作,或者生产其它结果.下面的例子中对流中的数字进行了求平方操作:

```kotlin
fun CoroutineScope.square(numbers: ReceiveChannel<Int>): ReceiveChannel<Int> = produce {
    for (x in numbers) send(x * x)
}
```

主代码开始并连接该通道:

```kotlin
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*

fun main() = runBlocking {
    val numbers = produceNumbers() // produces integers from 1 and on
    val squares = square(numbers) // squares integers
    for (i in 1..5) println(squares.receive()) // print first five
    println("Done!") // we are done
    coroutineContext.cancelChildren() // cancel children coroutines
}

fun CoroutineScope.produceNumbers() = produce<Int> {
    var x = 1
    while (true) send(x++) // infinite stream of integers starting from 1
}

fun CoroutineScope.square(numbers: ReceiveChannel<Int>): ReceiveChannel<Int> = produce {
    for (x in numbers) send(x * x)
}
```

> 所有创建协程的函数都被定义为CoroutineScope的扩展函数，因此我们可以依赖结构化并发来确保我们的应用程序中没有延迟的全局协同程序。

### 使用管道生产素数

让我们来展示一个极端的例子,在协程中使用一个管道来生成素数。首先创建一个无限数字序列。

```kotlin
fun CoroutineScope.numbersFrom(start: Int) = produce<Int> {
    var x = start
    while (true) send(x++) // infinite stream of integers from start
}
```

接下来的管道中对接收到数字队列进行过滤,移除所有可以被给定素数整除的数字:

```kotlin
fun CoroutineScope.filter(numbers: ReceiveChannel<Int>, prime: Int) = produce<Int> {
    for (x in numbers) if (x % prime != 0) send(x)
}
```

接下来构建一个从数字2开始的管道,从当前通道中获得素数,并不断的用当前发现的素数用于新的通道进行过滤:

```
numbersFrom(2) -> filter(2) -> filter(3) -> filter(5) -> filter(7) ... 
```

下面的例子打印了前十个素数， 在主线程的上下文中运行整个管道。直到所有的协程在该主协程 runBlocking 的作用域中被启动完成。 我们不必使用一个显式的列表来保存所有被我们已经启动的协程。 我们使用 cancelChildren 扩展函数在我们打印了前十个素数以后来取消所有的子协程。

```kotlin
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*

fun main() = runBlocking {
    var cur = numbersFrom(2)
    for (i in 1..10) {
        val prime = cur.receive()
        println(prime)
        cur = filter(cur, prime)
    }
    coroutineContext.cancelChildren() // cancel all children to let main finish    
}

fun CoroutineScope.numbersFrom(start: Int) = produce<Int> {
    var x = start
    while (true) send(x++) // infinite stream of integers from start
}

fun CoroutineScope.filter(numbers: ReceiveChannel<Int>, prime: Int) = produce<Int> {
    for (x in numbers) if (x % prime != 0) send(x)
}
```

代码输出如下:

```
2
3
5
7
11
13
17
19
23
29
```

注意，你可以使用标准库中的 [buildIterator](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines/build-iterator.html) 协程构建器来构建一个相似的管道。 使用 buildIterator 替换 produce、用 yield 替换 send、用next 替换 receive、 以及Iterator 替换 ReceiveChannel 来摆脱协程作用域，你将不再需要 runBlocking。 然而，如上所示，如果你在 Dispatchers.Default 上下文中运行它，使用通道的管道的好处在于它可以充分利用多核心 CPU。

无论如何，这是找到素数的极不切实际的方法。 实际上，管道确实涉及一些其他的挂起调用（比如对远程服务的异步调用），并且这些管道不能使用buildSequence / buildIterator构建，因为它们不允许任意挂起，这与完全异步的 `produce` 操作不同。

### 扇出

多个协程可能从同一个通道获取数据,并在它们之间进行分布式工作. 我们先建立一个生产者协程,该协程定期产生一个整数(每秒钟10个):

```kotlin
fun CoroutineScope.produceNumbers() = produce<Int> {
    var x = 1 // start from 1
    while (true) {
        send(x++) // produce next
        delay(100) // wait 0.1s
    }
}
```

接下来我们可以有不同的消费者协程. 在这个例子中,只是打印自己 id 以及收到的数字:

```kotlin
fun CoroutineScope.launchProcessor(id: Int, channel: ReceiveChannel<Int>) = launch {
    for (msg in channel) {
        println("Processor #$id received $msg")
    }    
}
```

接下来启动五个消费者并让它们不间断工作.看看会发生什么:

```kotlin
val producer = produceNumbers()
repeat(5){
    launchProcessor(it, producer)
}

delay(950)
producer.cancel()
```

输出结果可能会和下面的相似,处理器收到的 id 或许会不同:

```
Processor #2 received 1
Processor #4 received 2
Processor #0 received 3
Processor #1 received 4
Processor #3 received 5
Processor #2 received 6
Processor #4 received 7
Processor #0 received 8
Processor #1 received 9
Processor #3 received 10
```

主义取消生产者协程并关闭他的通道,会导致对该通道正在进行的迭代终止.

还有，注意我们在 `launchProcessor` 中是怎样使用 for 循环显式迭代通道并执行扇出的。 与 consumeEach 不同，for 循环是可以在多个协程中十分安全地使用.如果其中一个处理者协程执行失败，其它的处理器协程仍然会继续处理通道，而用 consumeEach 编写的处理器始终在正常或非正常完成时消耗（取消）底层通道。

### 扇入

多个协程可以发送消息到同一个通道。 比如说，创建一个字符串的通道，和一个往该通道中以指定的延迟发送指定字符串的挂起函数：

```kotlin
suspend fun sendString(channel: SendChannel<String>, s: String, time: Long) {
    while (true) {
        delay(time)
        channel.send(s)
    }
}
```

接下来让我们看看如果我们同时开启多个协程发送字符串会发生什么(本例中,我们在主线程的上下文中作为主协程的子协程来启动它们):

```kotlin
val channel = Channel<String>()
launch { sendString(channel, "foo", 200L) }
launch { sendString(channel, "BAR!", 500L) }
repeat(6) { // receive first six
    println(channel.receive())
}
coroutineContext.cancelChildren() // cancel all children to let main finish
```

代码输出如下:

```
foo
foo
BAR!
foo
foo
BAR!
```

### 带缓冲的通道

上面介绍的通道都是没有缓冲的.无缓冲的通道需要发送者和接受者相互匹配(既互相约定).如果先调起发送者,则会被挂起,直到接收者被调用, 如果接收者先被调用,则在发送者调起前它一直被挂起.

[Channel()](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/-channel.html) 工厂函数和 [produce](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/produce.html) 构建器都可以接收一个可选参数 `capacity` 来指定缓冲大小. 缓冲使得发送者在挂起前可以发送多个元素, 这个 `BlockingQueue` 指定容量是一致的, 它们会在缓冲慢之前一直阻塞.

看看下面代码的执行:

```kotlin
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*

fun main() = runBlocking<Unit> {
    val channel = Channel<Int>(4) // create buffered channel
    val sender = launch { // launch sender coroutine
        repeat(10) {
            println("Sending $it") // print before sending each element
            channel.send(it) // will suspend when buffer is full
        }
    }
    // don't receive anything... just wait....
    delay(1000)
    sender.cancel() // cancel sender coroutine    
}
```

使用缓冲通道并给 capacity 参数传入 四 它将打印 “sending” 五次：

```
Sending 0
Sending 1
Sending 2
Sending 3
Sending 4
```

前四个元素被加入到了缓冲区并且发送者在试图发送第五个元素的时候被挂起。

### 通道是公平的

发送和接收操作是公平的 并且严格按照调用它们的协程顺序进行。它们遵守先进先出原则，可以看到第一个协程调用 receive 并得到了元素。在下面的例子中两个协程 “乒” 和 "乓" 都从共享的“桌子”通道接收到这个“球”元素。


```kotlin
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*

data class Ball(var hits: Int)

fun main() = runBlocking {
    val table = Channel<Ball>() // a shared table
    launch { player("ping", table) }
    launch { player("pong", table) }
    table.send(Ball(0)) // serve the ball
    delay(1000) // delay 1 second
    coroutineContext.cancelChildren() // game over, cancel them
}

suspend fun player(name: String, table: Channel<Ball>) {
    for (ball in table) { // receive the ball in a loop
        ball.hits++
        println("$name $ball")
        delay(300) // wait a bit
        table.send(ball) // send the ball back
    }
}
```

“ping”协程首先启动，因此它是第一个接收球的人。 即使“ping”coroutine在将球送回桌面后立即再次接球，球也会被“pong”协程接收，因为它已经在等待了：

```
ping Ball(hits=1)
pong Ball(hits=2)
ping Ball(hits=3)
pong Ball(hits=4)
```

请注意，由于正在使用的执行程序的性质，有时通道可能会产生看起来不公平的执行。 有关详细信息，请参阅此[issue](https://github.com/Kotlin/kotlinx.coroutines/issues/111)

### 计时器通道

计时器通道是一种特别的会合通道，每次经过特定的延迟都会从该通道进行消费并产生 Unit。 虽然它看起来似乎没用，它被用来构建分段来创建复杂的基于时间的 produce 管道和进行窗口化操作以及其它时间相关的处理。 [select](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.selects/select.html) 中可以用计时器通道来进行“on tick”操作.

用 [ticker](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/ticker.html) 工厂方法可以创建一个这样的通道. 可以用 [ReceiveChannel.cancel](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.channels/-receive-channel/cancel.html) 方法表示接下来没有元素输出.

下面我们实践一下:

```kotlin
import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*

fun main() = runBlocking<Unit> {
    val tickerChannel = ticker(delayMillis = 100, initialDelayMillis = 0) // create ticker channel
    var nextElement = withTimeoutOrNull(1) { tickerChannel.receive() }
    println("Initial element is available immediately: $nextElement") // initial delay hasn't passed yet

    nextElement = withTimeoutOrNull(50) { tickerChannel.receive() } // all subsequent elements has 100ms delay
    println("Next element is not ready in 50 ms: $nextElement")

    nextElement = withTimeoutOrNull(60) { tickerChannel.receive() }
    println("Next element is ready in 100 ms: $nextElement")

    // Emulate large consumption delays
    println("Consumer pauses for 150ms")
    delay(150)
    // Next element is available immediately
    nextElement = withTimeoutOrNull(1) { tickerChannel.receive() }
    println("Next element is available immediately after large consumer delay: $nextElement")
    // Note that the pause between `receive` calls is taken into account and next element arrives faster
    nextElement = withTimeoutOrNull(60) { tickerChannel.receive() } 
    println("Next element is ready in 50ms after consumer pause in 150ms: $nextElement")

    tickerChannel.cancel() // indicate that no more elements are needed
}
```

结果如下:

```
Initial element is available immediately: kotlin.Unit
Next element is not ready in 50 ms: null
Next element is ready in 100 ms: kotlin.Unit
Consumer pauses for 150ms
Next element is available immediately after large consumer delay: kotlin.Unit
Next element is ready in 50ms after consumer pause in 150ms: kotlin.Unit
```

请注意，ticker 知道可能的消费者暂停，并且默认情况下会调整下一个生成的元素如果发生暂停则延迟，试图保持固定的生成元素率。

给可选的 mode 参数传入 TickerMode.FIXED_DELAY 可以保持固定元素之间的延迟。