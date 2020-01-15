Kotlin / Native 运行时不鼓励使用带有互斥代码块和条件变量的经典的面向线程的并发模型，因为该模型容易出错且不可靠。 相反，我们推荐使用一系列替代方案，来允许你使用硬件并发并实现阻塞IO。 这些方法如下，并将在后面的部分中详细说明：

- 可传递消息的 Worker
- 对象子图所有权转移
- 对象子图冻结
- 对象子图分离
- 使用C全局变量的原始共享内存
- 用于阻止操作的协程（本文档未涵盖）

## Works

Kotlin / Native运行提供 Workers 的概念替换线程：带有请求队列的并发可执行控制流。 Workers 与 Actor 模型中的 Actor 非常相似。 一个 Worker 可以与另一个 Worker 交换 Kotlin 对象，以便在任何时候每个可变对象都由一个 Worker 拥有，但是所有权可以转移。 请参阅[对象传输和冻结部分](https://kotlinlang.org/docs/reference/native/concurrency.html#transfer)。

一旦使用 `Worker.start` 函数调用启动 worker，便可以通过唯一整数 ID 进行寻址。 其他 worker 或非 worker 并发原语（例如OS线程）可以通过 `execute` 调用向 worker 发送消息。

```kotlin
val future = execute(TransferMode.SAFE, { SomeDataForWorker() }) {
   // data returned by the second function argument comes to the
   // worker routine as 'input' parameter.
   input ->
   // Here we create an instance to be returned when someone consumes result future.
   WorkerResult(input.stringParam + " result")
}

future.consume {
  // Here we see result returned from routine above. Note that future object or
  // id could be transferred to another worker, so we don't have to consume future
  // in same execution context it was obtained.
  result -> println("result is $result")
}
```