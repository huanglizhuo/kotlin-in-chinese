## 协程(Coroutines)有些 APIs 是需要长时间运行，并且需要调用者阻塞直到这些调用完成（比如网络 IO ，文件 IO ，CPU 或者 GPU 比较集中的工作）。协程提供了一种避免线程阻塞并且用一种更轻量级，更易操控到操作：协程暂停。



协程把异步编程放入库中来简化这类操作。程序逻辑在协程中顺序表述，而底层的库会将其转换为异步操作。库会将相关的用户代码打包成回调，订阅相关事件，调度其执行到不同的线程（甚至不同的机器），而代码依然想顺序执行那么简单。

很多其它语言中的异步模型都可以用 Kotlin 协程实现为库。比如 C# ECMAScipt  中的　[async/wait](https://github.com/Kotlin/kotlinx.coroutines/blob/master/coroutines-guide.md#composing-suspending-functions) ，Go 语言中的  [channels](https://github.com/Kotlin/kotlinx.coroutines/blob/master/coroutines-guide.md#channels) 和 [`select`](https://github.com/Kotlin/kotlinx.coroutines/blob/master/coroutines-guide.md#select-expression) ，以及 C# 和 Python 的 [generators/`yield`](http://kotlinlang.org/docs/reference/coroutines.html##generators-api-in-kotlincoroutines)  模型。下面的描述会详细解释提供这些结构的库。



### 阻塞和挂起
一般来说，协程是一种可以不阻塞线程但却可以被挂起的计算过程。线程阻塞总是昂贵的，尤其是在高负载的情形下，因为只有小部分的线程是实际运行的，因此阻塞它们会导致一些重要的任务被延迟。

而协程的挂起基本没有什么开销。没有上下文切换或者任何的操作系统的介入。最重要的是，挂起是可以背用户库控制的，库的作者可以决定在挂起时根据需要进行一些优化／日志记录／拦截等操作。

另一个不同就是协程不能被任意的操作挂起，而仅仅可以在被标记为 *挂起点*  的地方进行挂起。

### 挂起函数
当一个函数被 `suspend` 修饰时表示可以被挂起。

```kotlin
suspend fun doSomething(foo: Foo): Bar{
 ... 
}
```
这样的函数被称为 *挂起函数*，因为调用它可能导致挂起协程（库可以在调用结果已经存在的情形下决定取消挂起）。挂起函数可以想正常函数那样接受参数返回结果，但只能在协程中调用或着被其他挂起函数调用。事实上启动一个协程至少需要一个挂起函数，而且常常时匿名的（比如lambda）。下面这个例子是一个简单的`async()` 函数（来自[`kotlinx.coroutines`](http://kotlinlang.org/docs/reference/coroutines.html#generators-api-in-kotlincoroutines)  库）：

```kotlin
fun <T> async(block: suspend() -> T)
```

这里的 `async()`只是一个普通的函数（不是挂起函数），但 `block` 参数是一个带有 `suspend` 修饰的函数类型，所以当传递一个 lambda 给`async()`时，这会是一个挂起 lambda ，这样我们就可以在这里调用一个挂起函数了。

```kotlin
async{
  doSomething(foo)
}
```

继续类比，`await()` 函数可以是一个挂起函数(因此在 `await(){}` 语句块内仍然可以调用)，该函数会挂起协程直至指定操作完成并返回结果：

```kotlin
async {
    ...
    val result = computation.await()
    ...
}
```

更多关于 `async/await` 原理的内容请看[这里](https://github.com/Kotlin/kotlinx.coroutines/blob/master/coroutines-guide.md#composing-suspending-functions)

注意 `await()`和 `doSomething()` 不能在像 main() 这样的普通函数中调用： 

```kotlin
fun main(args: Array<String>) {
    doSomething() // 错误：挂起函数从非协程上下文调用
}
```
还有一点，挂起函数可以是虚函数，当覆写它们时，必须指定 suspend 修饰符：

```kotlin
interface Base {
    suspend fun foo()
}

class Derived: Base {
    override suspend fun foo() { …… }
}
```

### `@RestrictsSuspension` 注解

扩展函数（以及lambda）可以被标记为`suspend`。这样方便了用户创建其他[DSLs](http://kotlinlang.org/docs/reference/type-safe-builders.html)以及扩展其它API。有些情况下，库的作者需要阻止用户添加新的挂起线程的方案。

这时就需要`@RestrictsSuspension`注解了。当一个接收者类或者接口`R`被标注时，所有可挂起扩展都需要代理`R`的成员或者其它扩展。由于扩展时不能互相无限代理（会导致程序终止），这就保障了所有挂起都是通过调用`R`的成员，这样库作者就能完全掌控挂起方式了。

不过这样的场景不常见，它需要所有的挂起都通过库的特殊方式实现。比如，用下面的 [`buildSequence()`](http://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental/build-sequence.html) 函数实现生成器时，必须保证协程中所有的挂起都是通过调用`yield()`或者`yieldAll()`来实现。这就是为什么[`SequenceBuilder`](http://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental/-sequence-builder/index.html) 被标注为 `@RestrictsSuspension`:

```kotlin
@RestrictsSuspension
public abstract class SequenceBuilder<in T> {
    ...
}
```

可以参看[Github](https://github.com/JetBrains/kotlin/blob/master/libraries/stdlib/src/kotlin/coroutines/experimental/SequenceBuilder.kt) 源码

### 协程内部机制

这里并不打算全盘解释协程内部的工作原理，而是给大家一个整体上的概念。

协程完全时通过编译技术（并不需要 VM 或者 OS 方面的支持）实现，挂起时借由代码转换实现。基本上所有的挂起函数（当然是有些优化措施，但这里我们不会深入说明）都被转换为状态机。在挂起前，下一个状态会存储在编译器生成的与本地变量关联的类中。到恢复协程时，本地变量会被恢复为挂起之前的状态。

挂起的协程可以存储以及作为一个对象进行传递，该协程会继续持有其状态和本地变量。这样的对象的类型时`Continuation`，代码转换的整体实现思路是基于经典的 [Continuation-passing style](https://en.wikipedia.org/wiki/Continuation-passing_style) 。所有挂起函数要有一个额外的参数类型`Continuation`。

更多的细节可以参看[设计文档](https://github.com/Kotlin/kotlin-coroutines/blob/master/kotlin-coroutines-informal.md) 。其它语言（比如C# ECMASript2016）中类似的 async/await 模型在这里都有描述，当然了其它语言的实现机制和 Kotlin 有所不同

### 协程的实验状态

协程的设计是[实验性](http://kotlinlang.org/docs/reference/compatibility.html#experimental-features)的，也就是说在后面的 releasees 版本中可能会有所变更。当在 Kotlin1.1 中编译协程时，默认会有警告：*The feature "coroutines" is experimental* 。可以通过 [opt-in flag](http://kotlinlang.org/docs/diagnostics/experimental-coroutines.html) 来移除警告。

由于处于实验状态，协程相关的标准库都在`kotlin.coroutines.experimental`包下。当设计确定时实验状态将会取消，最后的API将会移到 `kotlin.coroutines`,实验性的包将会保留（或许是作为一个单独的构建中）以保持兼容。

**千万注意**: 建议库作者可以采用同样的转换：为基于协程的 API 采用 "experimental" 前缀作包名（比如`com.example.experimental`）。当最终 API 发布时，遵循下面的步骤：

- 复制所有 API 到 `com.example`包下
- 保留实验性大包做兼容。

这样可以减少用户的迁移问题。

### 标准 API

协程主要在三种层级中支持：

- 语言层面的支持（比如支持函数挂起）
- Kotlin 标准库中核心底层 API
- 可以直接在代码中使用的高级 API

#### 底层 API：`kotlin.coroutines`

底层 API 比较少，强烈建议不要使用，除非要创建高级库。这部分 API 主要在两个包中：

- [`kotlin.coroutines.experimental`](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental/index.html) 带有主要类型与下述原语
  - [`createCoroutine()`](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental/create-coroutine.html)
  - [`startCoroutine()`](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental/start-coroutine.html)
  - [`suspendCoroutine()`](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental/suspend-coroutine.html)
- [`kotlin.coroutines.experimental.intrinsics`](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental.intrinsics/index.html) 带有更底层的内联函数如 [`suspendCoroutineOrReturn`](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental.intrinsics/suspend-coroutine-or-return.html)

关于这些 API 用法的更多细节可以在[这里](https://github.com/Kotlin/kotlin-coroutines/blob/master/kotlin-coroutines-informal.md)找到。

#### `kotlin.coroutines`中的生成器API：

 `kotlin.coroutines.experimental` 中唯一的“应用层面”的函数是：

- [`buildSequence()`](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental/build-sequence.html)
- [`buildIterator()`](https://kotlinlang.org/api/latest/jvm/stdlib/kotlin.coroutines.experimental/build-iterator.html)

这些和 `kotlin-stdlib` 打包在一起，因为和序列相关。事实上，这些函数（这里单独以 `buildSequence()` 作为事例）实现生成器提供了一种更加简单的构造延迟序列的方法：



```kotlin
val fibonacciSeq = buildSequence {
    var a = 0
    var b = 1

    yield(1)

    while (true) {
        yield(a + b)

        val tmp = a + b
        a = b
        b = tmp
    }
}
```
这里通过调用 `yield()`函数生成新的斐波那契数，就可以生成一个无限的斐波那契数列。当遍历这样的数列时，每遍历一步就生成一个斐波那契数，这样就可以从中取出无限的斐波那契数。比如 `fibonacciSeq.take(8).toList()`会返回`[1, 1, 2, 3, 5, 8, 13, 21]`。协程让这一实现开销更低。

为了演示正真的延迟序列，在`buildSequence()`中打印一些调试信息：

```kotlin 
val lazySeq = buildSequence {
    print("START ")
    for (i in 1..5) {
        yield(i)
        print("STEP ")
    }
    print("END")
}

// Print the first three elements of the sequence
lazySeq.take(3).forEach { print("$it ") }
```

运行上面的代码运，如果我们输出前三个元素的数字与生成循环的 `STEP` 有交叉。这意味着计算确实是惰性的。要输出 `1`，我们只执行到第一个 `yield(i)`，并且过程中会输出 `START`。然后，输出 `2`，我们需要继续下一个 `yield(i)`，并会输出 `STEP`。`3` 也一样。永远不会输出再下一个 `STEP`（以及`END`），因为我们没有请求序列的后续元素。

使用 `yieldAll()` 函数可以一次性生成序列所有值：

```kotlin
val lazySeq = buildSequence {
    yield(0)
    yieldAll(1..10) 
}

lazySeq.forEach { print("$it ") }
```

 `buildIterator()` 与 `buildSequence()`作用相似，只不过返回值时延迟迭代器。

通过给`SequenceBuilder`类写挂起扩展，可以给 `buildSequence()`添加自定义生成逻辑：

```kotlin
suspend fun SequenceBuilder<Int>.yieldIfOdd(x: Int) {
    if (x % 2 != 0) yield(x)
}

val lazySeq = buildSequence {
    for (i in 1..10) yieldIfOdd(i)
}
```



#### 其它高级API：`kotlinx.coroutines`

Kotlin 标准库只提供与协程相关的核心 API 。主要有基于协程的库核心原语和接口可以使用。

大多数基于协程的应用程序级API都作为单独的库发布：[`kotlinx.coroutines`](https://github.com/Kotlin/kotlinx.coroutines)。这个库覆盖了

-  平台无关的异步编程此模块`kotlinx-coroutines-core` 
  - 包括类似 Go 语言的`select` 和其他便利原语
  - 这个库的综合指南[在这里](https://github.com/Kotlin/kotlinx.coroutines/blob/master/coroutines-guide.md)查看。
- 基于 JDK 8 中的 `CompletableFuture` 的 API：`kotlinx-coroutines-jdk8`
- 基于 JDK 7 及更高版本 API 的非阻塞 IO（NIO）：`kotlinx-coroutines-nio`
- 支持 Swing (`kotlinx-coroutines-swing`) 和 JavaFx (`kotlinx-coroutines-javafx`)
- 支持 RxJava：`kotlinx-coroutines-rx`

这些库既提供了方便的 API ，也可以作为构建其它基于协程库的样板参考。