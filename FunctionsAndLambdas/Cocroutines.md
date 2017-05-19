## 协程(Coroutines)有些 APIs 是需要长时间运行，并且需要调用者阻塞直到这些调用完成（比如网络 IO ，文件 IO ，CPU 或者 GPU 比较集中的工作）。协程提供了一种避免线程阻塞并且用一种更轻量级，更易操控到操作：协程暂停。



协程把异步编程放入库中来简化这类操作。程序逻辑在协程中顺序表述，而底层的库会将其转换为异步操作。库会将相关的用户代码打包成回调，订阅相关事件，调度其执行到不同的线程（甚至不同的机器），而代码依然想顺序执行那么简单。

很多其它语言中的异步模型都可以用 Kotlin 协程实现为库。比如 C# ECMAScipt  中的　[async/wait](https://github.com/Kotlin/kotlinx.coroutines/blob/master/coroutines-guide.md#composing-suspending-functions) ，Go 语言中的  [channels](https://github.com/Kotlin/kotlinx.coroutines/blob/master/coroutines-guide.md#channels) 和 [`select`](https://github.com/Kotlin/kotlinx.coroutines/blob/master/coroutines-guide.md#select-expression) ，以及 C# 和 Python 的 [generators/`yield`](http://kotlinlang.org/docs/reference/coroutines.html##generators-api-in-kotlincoroutines)  模型。下面的描述会详细解释提供这些结构的库。



### 阻塞和挂起
一般来说，协程是一种可以不阻塞线程但却可以被挂起的计算过程。线程阻塞总是昂贵的，尤其是在高负载的情形下，因为只有小部分的线程是实际运行的，因此阻塞它们会导致一些重要的任务被延迟。

而协程的挂起基本没有什么开销。没有上下文切换或者任何的操作系统的介入。最重要的是，挂起是可以背用户库控制的，库的作者可以决定在挂起时根据需要进行一些优化／日志记录／拦截等操作。

另一个不同就是协程不能被任意的操作挂起，而仅仅可以在被标记为 *挂起点*  的地方进行挂起。

### 挂起函数
当一个函数被 `suspend` 修饰时表示可以背挂起。

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
还有一点，挂起函数可以是虚拟的，当覆盖它们时，必须指定 suspend 修饰符：

```kotlin
interface Base {
    suspend fun foo()
}

class Derived: Base {
    override suspend fun foo() { …… }
}
to be continue
```


