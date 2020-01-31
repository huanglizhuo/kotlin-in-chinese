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

val time = measureTimeMillis {
    foo（）
        .conflate（）//混合排放，不处理每个排放
        .collect {值->
            delay（300）//假设我们正在处理300毫秒
            println（值）
        }
}
println（“在$ time ms中收集”）
目标平台：在kotlin v.1.3.61上运行的JVM
您可以从此处获取完整的代码。

我们看到，虽然第一个数字仍在处理中，第二个已经被处理，而第三个已经产生，所以第二个被混合了，只有最新的（第三个）被交付给收集器：

1个
3
758毫秒内收集
处理最新价值
当发射器和收集器都很慢时，合并是加快处理速度的一种方法。它通过删除发射值来实现。另一种方法是取消缓慢的收集器，并在每次发出新值时重新启动它。有一组xxxLatest运算符，它们执行与xxx运算符相同的基本逻辑，但是会在其块上取消新值的代码。在上一个示例中，让我们尝试将comlate更改为collectLatest：

val time = measureTimeMillis {
    foo（）
        .collectLatest {value-> //取消并重新启动最新值
            println（“收集$ value”）
            delay（300）//假设我们正在处理300毫秒
            println（“完成$ value”）
        }
}
println（“在$ time ms中收集”）
目标平台：在kotlin v.1.3.61上运行的JVM
您可以从此处获取完整的代码。

由于collectLatest的主体需要300毫秒，但是每100毫秒会发出一个新值，因此我们可以看到该块在每个值上运行，但仅针对最后一个值才完成：

收集1
收集2
收集3
完成3
741毫秒内收集
组成多个流
有很多方法可以组成多个流。

压缩
就像Kotlin标准库中的Sequence.zip扩展功能一样，流具有zip运算符，该运算符结合了两个流的相应值：


val nums =（1..3）.asFlow（）//数字1..3
val strs = flowOf（“一个”，“两个”，“三个”）//字符串
nums.zip（strs）{a，b->“ $ a-> $ b”} //组成一个字符串
    .collect {println（it）} //收集并打印
目标平台：在kotlin v.1.3.61上运行的JVM
您可以从此处获取完整的代码。

该示例打印：

1->一
2->两个
3->三个
结合
当流表示变量或操作的最新值时（另请参阅有关合并的部分），可能需要执行依赖于相应流的最新值的计算，并在任何上游出现时重新进行计算。流发出一个值。相应的运算符家族称为合并。

例如，如果上一个示例中的数字每300毫秒更新一次，但是字符串每400毫秒更新一次，那么即使使用每400毫秒打印一次的结果，使用zip运算符对它们进行压缩仍会产生相同的结果：

在此示例中，我们使用onEach中间运算符来延迟每个元素，并使发出采样流的代码更具声明性且更短。