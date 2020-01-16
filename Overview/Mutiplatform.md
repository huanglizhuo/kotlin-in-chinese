## 多平台编程

**多平台项目是 Kotlin 1.2 和 1.3 中的实验性特性。本文档中描述的所有语言和工具功能都可能在将来的Kotlin版本中发生变更**

使 Kotlin 可以在所有平台上使用是我们的目标，但我们认为这是一个更为重要的目标的前提：在平台之间共享代码。 借助对JVM，Android，JavaScript，iOS，Linux，Windows，Mac甚至STM32等嵌入式系统的支持，Kotlin可以处理现代应用程序的任何和所有组件。这为代码和知识经验重用带来了极大便利， 节省了大量用于多平台而需要两次或多次实现的时间，从而可以把时间用在更具挑战性的任务上。

### 它是怎么工作的

首先，多平台不是编译代码到所有平台。这种模型有着显而易见的缺点，我们都知道现代应用往往需要调用它们运行所在平台的独有特性。Kotlin不会限制你使用世界上所有API的共同子集。每个组件都可以尽可能的共享代码，同时可以借助语言提供的 [预期/实际（expect/actual） 机制]{https://kotlinlang.org/docs/reference/platform-specific-declarations.html}调用平台相关代码。

下面是一段示例代码，一个简单的日志记录框架，可以共享和互相调用公共逻辑和平台逻辑。公共代码如下：

```Kotlin

enum class LogLevel {
    DEBUG, WARN, ERROR
}

// 预期用于平台特有的API
internal expect fun writeLogMessage(message: String, logLevel: LogLevel)

// 预期通用代码API
fun logDebug(message: String) = writeLogMessage(message, LogLevel.DEBUG)
fun logWarn(message: String) = writeLogMessage(message, LogLevel.WARN)
fun logError(message: String) = writeLogMessage(message, LogLevel.ERROR)

```


它预期目标为 writeLogMessage 提供特定于平台的实现，并且通用代码现在可以使用此声明，而无需考虑如何实现。

在JVM上，可以提供一种将日志写入标准输出的实现：

```Kotlin
internal actual fun writeLogMessage(message: String, logLevel: LogLevel) {
    println("[$logLevel]: $message")
}
```

在 JavaScript 中，有着完全不同的API集合，所以也可以这样实现：

```Kotlin
internal actual fun writeLogMessage(message: String, logLevel: LogLevel) {
    when (logLevel) {
        LogLevel.DEBUG -> console.log(message)
        LogLevel.WARN -> console.warn(message)
        LogLevel.ERROR -> console.error(message)
    }
}
```

在 1.3 中，我们重新设计了整个多平台模型。我们用于描述多平台Gradle项目的新DSL更加灵活，并将持续致力于提供简单明了的配置。

### 多平台库

通用代码可以依靠一组涵盖日常任务（例如HTTP，序列化和管理协程）的库。 此外，所有平台上都提供了广泛的标准库。

你也可以编写自己的库，以提供通用的API，并在每个平台上以不同的方式实现它。

### 使用场景

#### Android - iOS

在移动平台中共享代码是 Kotlin 多平台的主要使用场景，如今可以在 Android 和 iOS 中共享业务逻辑，网络链接等部分的代码。

参看：

[移动多平台的特性，场景学习以及示例代码](https://www.jetbrains.com/lp/mobilecrossplatform/?_ga=2.87157787.2065084095.1578977891-351818486.1577339287)

[新建一个移动多平台项目](https://play.kotlinlang.org/hands-on/Targeting%20iOS%20and%20Android%20with%20Kotlin%20Multiplatform/01_Introduction)

#### Client - Server

另一个场景就是在浏览器中共享客户端和服务端共享代码逻辑，Kotlin 多平台也支持这种场景。

[Ktor框架](https://ktor.io/?_ga=2.20579131.2065084095.1578977891-351818486.1577339287) 适合在连接的系统中构建异步服务器和客户端。

### 怎么开始

刚开始学 Kotlin ？ 可以去[开始](https://kotlinlang.org/docs/reference/basic-syntax.html)页面。

建议的文档页面：

- [设置多平台项目](https://kotlinlang.org/docs/reference/building-mpp-with-gradle.html#setting-up-a-multiplatform-project)
- [平台相关声明](https://kotlinlang.org/docs/reference/platform-specific-declarations.html)

推荐教程：

- [多平台 Kotlin 库](https://kotlinlang.org/docs/tutorials/mpp/multiplatform-library.html)
- [多平台项目：iOS 和 Android](https://kotlinlang.org/docs/tutorials/native/mpp-ios-android.html)

示例项目：

- [KotlinConf app](https://github.com/JetBrains/kotlinconf-app)
- [KotlinConf Spinner app](https://github.com/jetbrains/kotlinconf-spinner)

所有例子都可以在 [Github](https://github.com/Kotlin/kotlin-examples) 上找到