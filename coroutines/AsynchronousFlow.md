### **目录**

- [异步Flow](#异步Flow)
    - [表达多个值](#表达多个值)
        - [序列](#序列)
        - [挂起函数](#挂起函数)
        - [Flow](#Flow)
    - [Flow是冷的](#Flow是冷的)
    - [Flow取消](#Flow取消)
    - [Flow构建器](#Flow构建器)
    - [中间Flow操作符](#中间Flow操作符)
        - [转换操作符](#转换操作符)
        - [大小限制操作符](#大小限制操作符)
    - [尾端Flow操作符](#尾端Flow操作符)
    - [Flow是顺序的](#Flow是顺序的)
    - [Flow上下文](#Flow上下文)
        - [withContex 错误发射](#withContex错误发射)
        - [flowOn 操作符](#flowOn操作符)
    - [缓冲](#缓冲)
        - [合并](#合并)
        - [处理最后一个值](#处理最后一个值)
    - [组合多个Flow](#组合多个Flow)
        - [zip](#zip) 
        - [Combine](#Combine)
    - [扁平化Flow](#扁平化Flow)
        - [flatMapConcat](#flatMapConcat)
        - [flatMapMerge](#flatMapMerge)
        - [flatMapLatest](#flatMapLatest)
    - [Flow异常](#Flow异常)
        - [透明捕获](#透明捕获)
        - [延迟捕获](#延迟捕获)
    - [Flow完成](#Flow完成)
        - [必定执行的 finaly 块](#必定执行的finaly块)
        - [声明式处理](#声明式处理)
        - [仅上游处理Flow异常](#仅上游处理Flow异常)
    - [命令式与声明式](#命令式与声明式)
    - [启动Flow](#启动Flow)
    - [Flow与响应式Stream](#Flow与响应式Stream)


## 异步Flow

挂起函数异步返回一个值，但是我们如何返回多个异步计算的值呢？这就是为什么要引入 Kotlin Flows

### 代表多个值

可以使用集合在Kotlin中表示多个值。例如，我们可以有一个 `foo()` 函数，该函数返回三个数字的列表，然后使用forEach将它们全部打印出来：

```Kotlin
package kotlinx.coroutines.guide.flow01

fun foo(): List<Int> = listOf(1, 2, 3)
 
fun main() {
    foo().forEach { value -> println(value) } 
}
```

输出如下：

```shell
1
2
3
```

#### 序列

如果我们使用一些占用CPU的阻塞代码来进行计算（每次计算需要100毫秒），那么我们可以使用Sequence来表示数字：

```Kotlin
package kotlinx.coroutines.guide.flow02

fun foo(): Sequence<Int> = sequence { // sequence builder
    for (i in 1..3) {
        Thread.sleep(100) // pretend we are computing it
        yield(i) // yield next value
    }
}

fun main() {
    foo().forEach { value -> println(value) } 
}
```

该代码输出与上面相同，但在打印每个数字之前要等待100毫秒。

#### 挂起函数

但是，此计算将阻止正在运行代码的主线程。当这些值由异步代码计算时，我们可以将函数foo标记为suspend，这样它就可以在不阻塞的情况下执行其工作，并将结果作为列表返回：

```Kotlin
suspend fun foo(): List<Int> {
    delay(1000) // pretend we are doing something asynchronous here
    return listOf(1, 2, 3)
}

fun main() = runBlocking<Unit> {
    foo().forEach { value -> println(value) } 
}
```

此代码在等待一秒钟后打印数字。

#### Flow

使用 `List<Int>` 作为结果类型，意味着我们只能一次返回所有值。为了表示异步计算的值Flow，我们可以使用 [Flow<Int>](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flow.html) 类型，就像对同步计算的值使用 `Sequence<Int>` 类型一样：

```Kotlin
package kotlinx.coroutines.guide.flow04

import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*

fun foo(): Flow<Int> = flow { // flow builder
    for (i in 1..3) {
        delay(100) // pretend we are doing something useful here
        emit(i) // emit next value
    }
}

fun main() = runBlocking<Unit> {
    // Launch a concurrent coroutine to check if the main thread is blocked
    launch {
        for (k in 1..3) {
            println("I'm not blocked $k")
            delay(100)
        }
    }
    // Collect the flow
    foo().collect { value -> println(value) } 
}
```

该代码在打印每个数字之前等待100毫秒，而不会阻塞主线程。这是通过每100毫秒从主线程中运行的单独协程打印“我未被阻止”来验证的：

```shell
I'm not blocked 1
1
I'm not blocked 2
2
I'm not blocked 3
3
```

请注意，代码与先前示例中的Flow有以下区别：

- [Flow](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flow.html) 类型的构建器函数称为[flow](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flow.html)。
- `flow{...}` 构建器块中的代码是可挂起的。
- 函数 `foo（）` 不再标记为suspend。
- 使用 [emit](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/-flow-collector/emit.html) 函数从flow中发射值。
- 使用 [collect](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/collect.html) 函数从flow中收集值。

我们可以在foo的flow{...}的主体中将Thread.sleep替换为delay，在这种情况下主线程将被阻塞。

### Flow是冷的

Flow 是类似于序列的冷流 - 流构建器中的代码在开始收集前不会运行。 在以下示例中将体现这一特性：

```Kotlin
package kotlinx.coroutines.guide.flow05

import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*

fun foo(): Flow<Int> = flow { 
    println("Flow started")
    for (i in 1..3) {
        delay(100)
        emit(i)
    }
}

fun main() = runBlocking<Unit> {
    println("Calling foo...")
    val flow = foo()
    println("Calling collect...")
    flow.collect { value -> println(value) } 
    println("Calling collect again...")
    flow.collect { value -> println(value) } 
}
```

输出入下:

```shell
Calling foo...
Calling collect...
Flow started
1
2
3
Calling collect again...
Flow started
1
2
3
```

这是 `foo()` 函数（返回流）未使用 suspend 修饰符标记的主要原因。 就其本身而言，foo（）快速返回并且不等待任何东西。 该流在每次收集时启动，这就是为什么当我们再次调用collect时看到“流已开始”的原因。

### Flow取消

Flow 必须与协程合作取消。 但流程基础结构不会引入其他取消点。 取消是完全透明。 与往常一样，当将流挂起在可取消的挂起函数（如delay）中时，可以取消流收集，否则不能取消。

以下示例显示了在 [withTimeoutOrNull](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/with-timeout-or-null.html) 块中运行代码时，如何在超时时取消该流：

```Kotlin
fun foo(): Flow<Int> = flow { 
    for (i in 1..3) {
        delay(100)          
        println("Emitting $i")
        emit(i)
    }
}

fun main() = runBlocking<Unit> {
    withTimeoutOrNull(250) { // Timeout after 250ms 
        foo().collect { value -> println(value) } 
    }
    println("Done")
}
```

注意这里 `foo()` 仅发出两个数字,输出如下:

```shell
Emitting 1
1
Emitting 2
2
Done
```

### Flow构建器

上一个例子中的 `flow{...}` 是最基本的一个, 还有其他构建器可以实现更简单的声明Flow:

- [flowOF](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flow-of.html) 构建器,定义一个生产固定集合值的流
- 可以使用.asFlow（）扩展函数将各种集合和序列转换为流。

因此，例子中打印从1到3的数字可以这样写：

```Kotlin
// Convert an integer range to a flow
(1..3).asFlow().collect { value -> println(value) }
```

### 中间Flow运算符

可以使用运算符来转换流，就像使用集合和序列一样。中间运算符应用于上游流，并返回下游流。这些运算符是冷的，就像流一样。调用此类运算符本身并非暂停函数。它会迅速返回新的转换流的定义。

基本运算符具有熟悉的名称，例如 [map](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/map.html) 和 [filter](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/filter.html) 。序列的重要区别是这些运算符中的代码块可以调用挂起函数。

例如，即使执行请求是由挂起函数实现的长时间运行的操作，也可以使用map运算符将传入请求的流映射到结果。

```Kotlin
suspend fun performRequest(request: Int): String {
    delay(1000) // imitate long-running asynchronous work
    return "response $request"
}

fun main() = runBlocking<Unit> {
    (1..3).asFlow() // a flow of requests
        .map { request -> performRequest(request) }
        .collect { response -> println(response) }
}
```
输出如下,每行直接相隔 1 秒打印:

```
response 1
response 2
response 3
```

#### 转换操作符

在流转换操作符中,最常见的是 [transform](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/transform.html). 它可以用来模仿简单的转换，例如map和filter，也可以实现更复杂的转换。使用 `transform` 运算符，我们可以发出任意次数的任意值。

例如，使用 transform 我们可以在执行长时间运行的异步请求之前发出一个字符串，并在其后添加一个响应：

```Kotlin
(1..3).asFlow() // a flow of requests
    .transform { request ->
        emit("Making request $request") 
        emit(performRequest(request)) 
    }
    .collect { response -> println(response) }
```

此代码的输出是：

```Kotlin
Making request 1
response 1
Making request 2
response 2
Making request 3
response 3
```

#### 大小限制操作符

限制大小的中间运算符（如[take](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/take.html)）会在达到相应的限制时取消流程的执行。协程中的取消总是通过抛出异常来执行的，因此在取消的情况下，所有资源管理功能（例如try {...}finally{...}块）都可以正常运行：

```Kotlin
fun numbers(): Flow<Int> = flow {
    try {                          
        emit(1)
        emit(2) 
        println("This line will not execute")
        emit(3)    
    } finally {
        println("Finally in numbers")
    }
}

fun main() = runBlocking<Unit> {
    numbers() 
        .take(2) // take only the first two
        .collect { value -> println(value) }
}            
```

输出结果显示 `flow{...}` 块中发射第二个数字后就停止了:

```
1
2
Finally in numbers
```

### 尾端Flow操作符

流上的尾端操作符是一个挂起函数，该函数启动流收集工作。 collect 操作符是最基本的操作符，但还有其他尾端操作符，可以使操作变得更简单：

- 转换为各种集合，如[toList](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/to-list.html) 和 [toSet](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/to-set.html)。
- 运算符获取[first](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/first.html)值并确保流只发射[single](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/single.html)值。
- 压缩合并流为某个值 [reduce](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/reduce.html) [fold](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/fold.html)

例如：

```Kotlin
val sum = (1..5).asFlow()
    .map { it * it } // squares of numbers from 1 to 5                           
    .reduce { a, b -> a + b } // sum them (terminal operator)
println(sum)
```

输出

```
55
```

### Flow是顺序的

除非使用对多个流进行操作的特殊运算符，否则将依次执行流的每个单独集合。集合直接在协程中工作，该协程调用终尾端操作符。默认情况下，不启动新的协程。每个发出的值都由所有中间操作符从上游到下游进行处理，然后再传递给尾端操作符。

请参见以下示例，该示例过滤偶数整数并将其映射到字符串：

```Kotlin
(1..5).asFlow()
    .filter {
        println("Filter $it")
        it % 2 == 0              
    }              
    .map { 
        println("Map $it")
        "string $it"
    }.collect { 
        println("Collect $it")
    }    
```

输出：

```Kotlin
Filter 1
Filter 2
Map 2
Collect string 2
Filter 3
Filter 4
Map 4
Collect string 4
Filter 5
```

### Flow上下文

流的收集总是发生在调用协程的上下文中。例如，如果有一个foo流，那么以下代码将在该代码的作者指定的上下文中运行，而不管foo流的实现细节如何：

```Kotlin
withContext（context）{
    foo.collect {value->
        println（value）//在指定的上下文中运行
    }
}
```

流的此属性称为上下文保留。

因此，默认情况下，flow{...} 构建器中的代码在相应流的收集器提供的上下文中运行。例如，考虑foo的实现，该实现打印被调用的线程并发出三个数字：

```kotlin
fun foo(): Flow<Int> = flow {
    log("Started foo flow")
    for (i in 1..3) {
        emit(i)
    }
}  

fun main() = runBlocking<Unit> {
    foo().collect { value -> log("Collected $value") } 
}            
```

代码输出如下:

```shell
[main @coroutine#1] Started foo flow
[main @coroutine#1] Collected 1
[main @coroutine#1] Collected 2
[main @coroutine#1] Collected 3
```

因为 `foo().collect` 是在主线程调用, foo 流主体也将在主线程中调用.这是快速运行或异步代码的理想默认值，这些代码不关心执行上下文并且不会阻塞调用者。

### withContex错误发射

然而，长时间运行的CPU消耗代码可能需要在 [Dispatchers.Default](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-dispatchers/-default.html) 上下文中执行。UI更新代码可能需要在 [Dispatchers.Main](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/-dispatchers/-main.html) 的上下文中执行。 通常，[withContext](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines/with-context.html) 将用于Kotlin协程更改代码的上下文，但是 flow{...} 构建器中的代码必须遵守上下文保留属性，并且不允许从其他上下文中发出。

尝试运行以下代码：

```Kotlin
fun foo(): Flow<Int> = flow {
    // The WRONG way to change context for CPU-consuming code in flow builder
    kotlinx.coroutines.withContext(Dispatchers.Default) {
        for (i in 1..3) {
            Thread.sleep(100) // pretend we are computing it in CPU-consuming way
            emit(i) // emit next value
        }
    }
}

fun main() = runBlocking<Unit> {
    foo().collect { value -> println(value) } 
}            
```

代码将产生如下异常:


```shell
Exception in thread "main" java.lang.IllegalStateException: Flow invariant is violated:
        Flow was collected in [CoroutineId(1), "coroutine#1":BlockingCoroutine{Active}@5511c7f8, BlockingEventLoop@2eac3323],
        but emission happened in [CoroutineId(1), "coroutine#1":DispatchedCoroutine{Active}@2dae0000, DefaultDispatcher].
        Please refer to 'flow' documentation or use 'flowOn' instead
    at ...
```

### flowOn操作符

这里的异常可以借助 [flowOn](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flow-on.html) 函数来改变 flow 发射上下文.正确修改 flow 上下文方法如下, 例子中还会打印出对应线程:

```Kotlin
fun foo(): Flow<Int> = flow {
    for (i in 1..3) {
        Thread.sleep(100) // pretend we are computing it in CPU-consuming way
        log("Emitting $i")
        emit(i) // emit next value
    }
}.flowOn(Dispatchers.Default) // RIGHT way to change context for CPU-consuming code in flow builder

fun main() = runBlocking<Unit> {
    foo().collect { value ->
        log("Collected $value") 
    } 
}            
```

这里 `flow{...}` 在后台线程工作,而收集发生在主线程:

另一个需要注意的是,这里的 [flowOn](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flow-on.html) 操作符改变了流天然的顺序性. 现在，收集发生在一个协程（“协程1”）中，发射发生在另一个协程（“协程2”）中，该协程与收集协程同时在另一个线程中运行。 当FlowOn运算符必须在其上下文中更改CoroutineDispatcher时，它会为上游流创建另一个协程。

### 缓冲

从收集流程所花费的总时间来看，尤其是在涉及长时间运行的异步操作时，以不同的协程运行流程的不同部分可能会有所帮助。 例如，考虑以下情况：foo（）流的发射速度很慢，花费100毫秒来生成一个元素； 收集器也很慢，需要300毫秒来处理一个元素。 让我们看看用流收集三个数字需要多长时间：

```Kotlin
fun foo(): Flow<Int> = flow {
    for (i in 1..3) {
        delay(100) // pretend we are asynchronously waiting 100 ms
        emit(i) // emit next value
    }
}

fun main() = runBlocking<Unit> { 
    val time = measureTimeMillis {
        foo().collect { value -> 
            delay(300) // pretend we are processing it for 300 ms
            println(value) 
        } 
    }   
    println("Collected in $time ms")
}
```

产生结果如下，整个集合大约需要1200毫秒（三个数字，每个数字400毫秒）：

```shell
1
2
3
Collected in 1220 ms
```

我们可以在流上使用 [buffer](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/buffer.html) 运算符，可以在收集开始时并行发送数据, 而不是顺序进行:

```Kotlin
val time = measureTimeMillis {
    foo()
        .buffer() // buffer emissions, don't wait
        .collect { value -> 
            delay(300) // pretend we are processing it for 300 ms
            println(value) 
        } 
}   
println("Collected in $time ms")
```

由于我们已经有效地创建了处理管道，因此它只需要等待100毫秒即可处理第一个数字，然后只需花费300毫秒来处理每个数字，因此它会更快地产生相同的数字。这样大约需要1000毫秒才能运行：

```Kotlin
1
2
3
Collected in 1071 ms
```

请注意，flowOn运算符在必须更改CoroutineDispatcher时使用相同的缓冲机制，但是在这里，我们显式地请求缓冲而不更改执行上下文。

**合并**

当流表示操作的部分结果或操作状态更新时，可能不必处理每个值，而只需要处理最近的值即可。在这种情况下，当收集器太慢而无法处理中间值时，可以使用合并运算符跳过中间值。以前面的示例为基础：

```Kotlin
package kotlinx.coroutines.guide.flow18

import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*
import kotlin.system.*

fun foo(): Flow<Int> = flow {
    for (i in 1..3) {
        delay(100) // pretend we are asynchronously waiting 100 ms
        emit(i) // emit next value
    }
}

fun main() = runBlocking<Unit> { 
    val time = measureTimeMillis {
        foo()
            .conflate() // conflate emissions, don't process each one
            .collect { value -> 
                delay(300) // pretend we are processing it for 300 ms
                println(value) 
            } 
    }   
    println("Collected in $time ms")
}
```

我们看到，虽然第一个数字仍在处理中，第二个和第三个已经发出，所以第二个被合并了，只有最新的（第三个）被交付给收集器：

```shell
1
3
Collected in 758 ms
```

### 处理最新值

当发射器和收集器都很慢时，合并是加快处理速度的一种方法。它通过删除发射值来实现。另一种方法是取消缓慢的收集器，并在每次发出新值时重新启动它。有一组 `xxxLatest` 运算符，它们执行与 `xxx` 运算符相同的基本逻辑，但是会在其块上取消新值的代码。在上一个示例中，让我们尝试将 [conflate](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/conflate.html)  改为 [collectLatest](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/collect-latest.html):

```Kotlin
val time = measureTimeMillis {
    foo()
        .collectLatest { value -> // cancel & restart on the latest value
            println("Collecting $value") 
            delay(300) // pretend we are processing it for 300 ms
            println("Done $value") 
        } 
}   
println("Collected in $time ms")
```

由于 collectLatest 的主体需要300毫秒，但是每100毫秒会发出一个新值，因此我们可以看到该块在每个值上运行，但仅针对最后一个值才完成：

```shell
Collecting 1
Collecting 2
Collecting 3
Done 3
Collected in 741 ms
```

### 组和多个流

有很多方法可以组和多个流。

**zip (压缩)**

就像Kotlin标准库中的 [Sequence.zip](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.sequences/zip.html) 扩展功能一样，流具有zip运算符，该运算符结合了两个流的相应值：

```Kotlin
val nums = (1..3).asFlow() // numbers 1..3
val strs = flowOf("one", "two", "three") // strings 
nums.zip(strs) { a, b -> "$a -> $b" } // compose a single string
    .collect { println(it) } // collect and print
```

输出如下：

```shell
1 -> one
2 -> two
3 -> three
```

**Combine (组合)**

当流表示变量或操作的最新值时（另请参阅有关[合并](https://kotlinlang.org/docs/reference/coroutines/flow.html#conflation)的部分），可能需要执行依赖于相应流的最新值的计算，并在任何上游有变动时重新进行计算。流发出一个值。相应的运算符族称为[Combine](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/combine.html)。

例如，如果上一个示例中的数字每300毫秒更新一次，但是字符串每400毫秒更新一次，使用zip运算符对它们进行压缩仍会产生相同的结果,尽管结果任然是没 400
毫秒打印一次：

在此示例中，我们使用onEach中间运算符来延迟每个元素，并使发出采样流的代码更具声明性且更短。

```Kotlin
val nums = (1..3).asFlow().onEach { delay(300) } // numbers 1..3 every 300 ms
val strs = flowOf("one", "two", "three").onEach { delay(400) } // strings every 400 ms
val startTime = System.currentTimeMillis() // remember the start time 
nums.zip(strs) { a, b -> "$a -> $b" } // compose a single string with "zip"
    .collect { value -> // collect and print 
        println("$value at ${System.currentTimeMillis() - startTime} ms from start") 
    } 
```

然而当使用 [combine](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/combine.html) 运算符取代 zip:

```Kotlin
val nums = (1..3).asFlow().onEach { delay(300) } // numbers 1..3 every 300 ms
val strs = flowOf("one", "two", "three").onEach { delay(400) } // strings every 400 ms          
val startTime = System.currentTimeMillis() // remember the start time 
nums.combine(strs) { a, b -> "$a -> $b" } // compose a single string with "combine"
    .collect { value -> // collect and print 
        println("$value at ${System.currentTimeMillis() - startTime} ms from start") 
    } 
```

我们会得到完全不同的结果,每次 `nums` 或者 `strs` 流都会打印一行:


```shell
1 -> one at 452 ms from start
2 -> one at 651 ms from start
2 -> two at 854 ms from start
3 -> two at 952 ms from start
3 -> three at 1256 ms from start
```

### 扁平化流

流表示异步接收的值序列，因此会有每产生一个值,就会触发对另一个值序列的请求的情形。例如，我们可以具有以下函数，该函数返回两个字符串，每个字符串的间隔为500 ms：

```Kotlin
fun requestFlow(i: Int): Flow<String> = flow {
    emit("$i: First") 
    delay(500) // wait 500 ms
    emit("$i: Second")    
}
```

现在，如果我们有三个整数流，并为每个整数调用 `requestFlow` ，如下所示：

```Kotlin
(1..3).asFlow().map { requestFlow(it) }
```

然后，我们得到一个流的流（Flow<Flow<String>>），该流需要扁平化为单个流以进行进一步处理。集合和序列为此具有 [flatten](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.sequences/flatten.html) 和 [flatMap](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.sequences/flat-map.html) 运算符。但是，由于流的异步性质，它们要求使用不同的扁平化模式，因此，在流上有一系列扁平化运算符。

**flatMapConcat**

串联模式由 [flatMapConcat](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flat-map-concat.html) 和 [flattenConcat](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flatten-concat.html) 运算符实现。它们是相应序列运算符的最直接类似物。他们等待内部流完成，然后开始收集下一个，如以下示例所示：

```Kotlin
package kotlinx.coroutines.guide.flow23

import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*

fun requestFlow(i: Int): Flow<String> = flow {
    emit("$i: First") 
    delay(500) // wait 500 ms
    emit("$i: Second")    
}

fun main() = runBlocking<Unit> { 
    val startTime = currentTimeMillis() // remember the start time 
    (1..3).asFlow().onEach { delay(100) } // a number every 100 ms 
        .flatMapConcat { requestFlow(it) }                                                                           
        .collect { value -> // collect and print 
            println("$value at ${currentTimeMillis() - startTime} ms from start") 
        } 
}
```

从输出中可以清楚地看到f latMapConcat 的顺序性质：

```shell
1: First at 121 ms from start
1: Second at 622 ms from start
2: First at 727 ms from start
2: Second at 1227 ms from start
3: First at 1328 ms from start
3: Second at 1829 ms from start
```

**flatMapMerge**

另一种扁平模式是同时收集所有传入流并将其值合并为单个流，以便尽快发出值。它由 [flatMapMerge](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flat-map-merge.html) 和 [flattenMerge](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flatten-merge.html) 运算符实现。它们都接受一个可选的 `concurrency` 参数，该参数限制了同时收集的并发流的数量（默认情况下它等于 [DEFAULT_CONCURRENCY](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/-d-e-f-a-u-l-t_-c-o-n-c-u-r-r-e-n-c-y.html) ）。

```Kotlin
val startTime = System.currentTimeMillis() // remember the start time 
(1..3).asFlow().onEach { delay(100) } // a number every 100 ms 
    .flatMapMerge { requestFlow(it) }                                                                           
    .collect { value -> // collect and print 
        println("$value at ${System.currentTimeMillis() - startTime} ms from start") 
    } 
```

flatMapMerge 的并发本质是显而易见的：

```Kotlin
1: First at 136 ms from start
2: First at 231 ms from start
3: First at 333 ms from start
1: Second at 639 ms from start
2: Second at 732 ms from start
3: Second at 833 ms from start
```

请注意，flatMapMerge 顺序调用其代码块（在此示例中为{requestFlow（it）}），但同时收集结果流，这等效于先执行顺序映射{requestFlow（it）}，然后在对结果调用[flattenMerge](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flatten-merge.html)

**flatMapLatest**

用与“处理最新值”一节中所示的 collectLatest 运算符类似的方式，存在对应的“最新”扁平模式，在该模式下，一旦发出新流，就会取消先前流的集合。它由 [flatMapLatest](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/flat-map-latest.html) 运算符实现。

```Kotlin
package kotlinx.coroutines.guide.flow25

import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*

fun requestFlow(i: Int): Flow<String> = flow {
    emit("$i: First") 
    delay(500) // wait 500 ms
    emit("$i: Second")    
}

fun main() = runBlocking<Unit> { 
    val startTime = currentTimeMillis() // remember the start time 
    (1..3).asFlow().onEach { delay(100) } // a number every 100 ms 
        .flatMapLatest { requestFlow(it) }                               
        .collect { value -> // collect and print 
            println("$value at ${currentTimeMillis() - startTime} ms from start") 
        } 
}
```

此示例中的输出很好地演示了flatMapLatest的工作方式：

```Kotlin
1: First at 142 ms from start
2: First at 322 ms from start
3: First at 425 ms from start
3: Second at 931 ms from start
```

请注意，flatMapLatest取消了其块（在此示例中为{requestFlow（it）}）上的所有代码的新值。在此特定示例中，这没有什么区别，因为对requestFlow本身的调用是快速，非挂起的并且无法取消。但是，如果要在其中使用诸如delay之类的暂停功能，它的优点才显示出来。

### Flow异常

当运算符中的发射器或代码引发异常时，流收集将提前完成并带有返回异常。有几种处理这些异常的方法。

**收集器 try catch**

收集器可以使用 Kotlin 的 try/catch 块来处理异常：

```Kotlin
fun foo(): Flow<Int> = flow {
    for (i in 1..3) {
        println("Emitting $i")
        emit(i) // emit next value
    }
}

fun main() = runBlocking<Unit> {
    try {
        foo().collect { value ->         
            println(value)
            check(value <= 1) { "Collected $value" }
        }
    } catch (e: Throwable) {
        println("Caught $e")
    } 
}        
```

这段代码成功地在collect终端操作符中捕获了一个异常，并且正如我们所看到的，此后不再发出任何值：

```Kotlin
Emitting 1
1
Emitting 2
2
Caught java.lang.IllegalStateException: Collected 2
```

**所有的东西都可以捕获**

前面的示例实际上捕获了在发射器或任何中间或终端运算符中发生的任何异常。例如，让我们更改代码，以便将发出的值[映射](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/map.html)到字符串，但是相应的代码会产生异常：

```Kotlin
fun foo(): Flow<String> = 
    flow {
        for (i in 1..3) {
            println("Emitting $i")
            emit(i) // emit next value
        }
    }
    .map { value ->
        check(value <= 1) { "Crashed on $value" }                 
        "string $value"
    }

fun main() = runBlocking<Unit> {
    try {
        foo().collect { value -> println(value) }
    } catch (e: Throwable) {
        println("Caught $e")
    } 
}        
```

仍然会捕获此异常，并且停止收集：

```shell
Emitting 1
string 1
Emitting 2
Caught java.lang.IllegalStateException: Crashed on 2
```

### 异常透明化

但是，发射器的代码如何封装其异常处理行为？

流必须对异常透明，`try/catch` 块内部的 `flow{...}` 构建器中发出值,违反了异常透明性。这样可以保证引发异常的收集器始终可以使用上一个示例中的 try/catch 捕获异常。

发射器可以使用 [catch](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/catch.html) 运算符，该运算符保留此异常透明性并允许对其异常处理进行封装。 catch运算符的主体可以分析异常并根据捕获到的异常以不同的方式对异常作出反应：

- 可以使用 `throw` 再次抛出异常。
- 异常可以在 [catch](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/catch.html) 块内借助 [emit](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/-flow-collector/emit.html) 转换为异常为发射值。
- 异常可以被记录,处理,或者被其他代码处理.

比如,我们可以在捕获异常时发出一段文本:

```Kotlin
package kotlinx.coroutines.guide.flow28

import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*

fun foo(): Flow<String> = 
    flow {
        for (i in 1..3) {
            println("Emitting $i")
            emit(i) // emit next value
        }
    }
    .map { value ->
        check(value <= 1) { "Crashed on $value" }                 
        "string $value"
    }

fun main() = runBlocking<Unit> {
    foo()
        .catch { e -> emit("Caught $e") } // emit on exception
        .collect { value -> println(value) }
}        
```

即使我们不使用 try/catch 代码，示例的输出也相同。

**透明 catch**

[catch](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/catch.html) 中间操作符遵循异常透明性，仅捕获上游异常（这是catch之上而非之下所有运算符的异常）。如果collect {...}中的块（放置在catch下方）抛出异常，则它不会被捕获：

```Kotlin
fun foo(): Flow<Int> = flow {
    for (i in 1..3) {
        println("Emitting $i")
        emit(i)
    }
}

fun main() = runBlocking<Unit> {
    foo()
        .catch { e -> println("Caught $e") } // does not catch downstream exceptions
        .collect { value ->
            check(value <= 1) { "Collected $value" }                 
            println(value) 
        }
}            
```

尽管有捕获操作符，但不会打印“Caught…”消息：

**声明式捕捉**

我们可以将catch操作符的声明性与处理所有异常的目的结合起来,将collect操作符的主体移到[onEach](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/on-each.html) 内并将其放在catch操作符之前。这时必须通过不带参数的调用collect（）来触发此流的收集：

```Kotlin
foo()
    .onEach { value ->
        check(value <= 1) { "Collected $value" }                 
        println(value) 
    }
    .catch { e -> println("Caught $e") }
    .collect()
```

现在我们看到打印了一条“ Caught…”消息，因此我们可以捕获所有异常，而无需显式使用try / catch块：

### 流完成

流收集完成时（正常或异常），可能需要执行一个操作。你可能已经注意到，它可以通过两种方式完成：命令式或声明式。

**必定执行的finaly块**

除了try/catch之外，收集器还可以使用 `finally` 块在收集完成后执行操作。

```Kotlin
fun foo(): Flow<Int> = (1..3).asFlow()

fun main() = runBlocking<Unit> {
    try {
        foo().collect { value -> println(value) }
    } finally {
        println("Done")
    }
}        
```

这段代码打印出foo()流产生的三个数字，后跟一个“Done”字符串：

```shell
1
2
3
Done
```

**声明式处理**

对于声明性方法，流具有 [onCompletion](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/on-completion.html) 中间操作符，该操作符在流已全部收集完成时被调用。

可以使用onCompletion运算符重写前面的示例，并产生相同的输出：

```Kotlin
foo()
    .onCompletion { println("Done") }
    .collect { value -> println(value) }
```

onCompletion 的主要优点是 lambda 的可空的Throwable参数，可用于确定流收集是正常完成还是异常完成。在下面的示例中，foo() 流在发出数字1之后引发异常：

```Kotlin
fun foo(): Flow<Int> = flow {
    emit(1)
    throw RuntimeException()
}

fun main() = runBlocking<Unit> {
    foo()
        .onCompletion { cause -> if (cause != null) println("Flow completed exceptionally") }
        .catch { cause -> println("Caught exception") }
        .collect { value -> println(value) }
}           
```

如你所料，它将打印：

```Kotlin
1
Flow completed exceptionally
Caught exception
```

与catch不同，onCompletion运算符不处理异常。从上面的示例代码可以看出，异常仍然向下游流动。它会交付给其他onCompletion运算符，并且可以由catch运算符处理。


**仅上游处理Flow异常**

就像catch运算符一样，仅来自上游的异常对 onCompletion 可见，而下游异常对其不可见。 例如，运行以下代码：

```Kotlin
fun foo(): Flow<Int> = (1..3).asFlow()

fun main() = runBlocking<Unit> {
    foo()
        .onCompletion { cause -> println("Flow completed with $cause") }
        .collect { value ->
            check(value <= 1) { "Collected $value" }                 
            println(value) 
        }
}
```

我们可以看到完成原因为null，但收集失败，出现以下异常：

```Kotlin
1
Flow completed with null
Exception in thread "main" java.lang.IllegalStateException: Collected 2
```

### 命令式与声明式

现在我们知道了如何收集流，并以命令式和声明式方式处理流的完成和异常。那么问题来了，首选哪种方法，为什么？作为一个库，我们不主张采用任何特定的方法，并且认为这两个选项都是有效的，应根您自己的喜好和代码风格进行选择。

### 启动流

使用流来表示来自某个源的异步事件很容易。在这种情况下，我们需要一个 `addEventListener` 函数的类似物，该函数通过对传入事件的反应来注册一段代码，并继续进行进一步的工作。 onEach 运算符可以担任此角色。但是，onEach是中间运算符。我们还需要尾端操作符来收集流。否则，仅调用onEach无效。

如果我们在onEach之后使用collect尾端操作符，那么它后面的代码将直到流收集完成后触发：

```Kotlin
// Imitate a flow of events
fun events(): Flow<Int> = (1..3).asFlow().onEach { delay(100) }

fun main() = runBlocking<Unit> {
    events()
        .onEach { event -> println("Event: $event") }
        .collect() // <--- Collecting the flow waits
    println("Done")
}            
```

输出如下:

```Kotlin
Event: 1
Event: 2
Event: 3
Done
```

[launchIn](https://kotlin.github.io/kotlinx.coroutines/kotlinx-coroutines-core/kotlinx.coroutines.flow/launch-in.html) 尾端操作符出现了。通过用 launchIn 代替collect，我们可以在单独的协程中启动流的集合，以便立即继续执行其他代码：

```Kotlin
package kotlinx.coroutines.guide.flow36

import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*

// Imitate a flow of events
fun events(): Flow<Int> = (1..3).asFlow().onEach { delay(100) }

fun main() = runBlocking<Unit> {
    events()
        .onEach { event -> println("Event: $event") }
        .launchIn(this) // <--- Launching the flow in a separate coroutine
    println("Done")
```

它打印：

```Kotlin
Done
Event: 1
Event: 2
Event: 3
```

launchIn的必需参数必须指定一个 CoroutineScope ，在其中启动用于收集流的协程。在上面的示例中，此作用域来自 runBlocking 协程构建器，因此，在运行流程时，此runBlocking范围等待其子协程完成，并防止main函数返回并终止此示例。

在实际应用中，范围将来自生命周期有限的实体。一旦此实体的生命周期终止，则将取消相应的作用域，从而取消相应流的收集。这样，一对 onEach { ... }.launchIn(scope) 就像addEventListener 一样工作。但是，由于取消和结构化并发达到了此目的，因此不需要相应的removeEventListener函数。

请注意，launchIn还返回一个Job，该Job仅可在不取消整个作用域或不加入整个作用域的情况下用于取消相应的流程集合协程。

### Flow 和反应式 Streams

对于那些熟悉[reactive stream](https://www.reactive-streams.org/)或反应式框架（例如RxJava和 project Reactor）的人来说，Flow的设计可能看起来非常熟悉。

确实，它的设计受到了Reactive Streams及其各种实现的启发。但是Flow的主要目标是拥有尽可能简单的设计，是Kotlin和挂起(suspesion)友好且遵循结构化并发。没有 其它框架及其出色大量的工作，就不不会有 Kotlin 中 flowd 的实现。你可以在 [Reactive Streams和Kotlin Flows](https://medium.com/@elizarov/reactive-streams-and-kotlin-flows-bfd12772cda4) 文章中阅读完整故事。

从概念上讲，Flow虽然有所不同，但它是反应性流，而且也可以将其转换为反应性（符合规范和TCK规范）的发布者，反之亦然。这样的转换器是由kotlinx.coroutines开箱即用地提供的，可以在相应的反应模块中找到（针对 Reactive Streams 的kotlinx-coroutines-active，用于P roject Reactor 的kotlinx-coroutines-reactor和针对RxJava2的kotlinx-coroutines-rx2） 。集成模块包括与 Flow 的相互转换，与Reactor的Context集成以及与各种反应式实体一起使用的易于挂起的方式。