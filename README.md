[![Analytics](https://ga-beacon.appspot.com/UA-80536214-1/readme)](https://github.com/huanglizhuo/kotlin-in-chinese)

2020.1.30

开始同步 Kotlin 1.3.61 文档

2017.5.17

Android Announces Support for Kotlin

By Mike Cleron, Director, Android Platform

Google 宣布官方支持 Kotlin [android-developers.googleblog](https://android-developers.googleblog.com/2017/05/android-announces-support-for-kotlin.html)

Android Studio 3.0 将默认集成 Kotlin plug-in ，博客中还说到 Kotlin 有着出色的设计，并相信 Kotlin 会帮助开发者更快更好的开发 Android 应用

Expedia, Flipboard, Pinterest, Square 等公司都有在自家的项目中使用 Kotlin 

2017.3.8

Kotlin 1.1 正式发布，这次最令人振奋的莫过于协程的发布，有了协程就可以更优雅的完成异步编程了

更多新特性请参看[what's new in kotlin 1.1](https://kotlinlang.org/docs/reference/whatsnew11.html)

[pdf下载](https://www.gitbook.com/download/pdf/book/huanglizhuo/kotlin-in-chinese)  [ePub下载](https://www.gitbook.com/download/epub/book/huanglizhuo/kotlin-in-chinese)

记得要点 star star star

发现有翻译的不好的或者错误欢迎到 github 提 [issue](https://github.com/huanglizhuo/kotlin-in-chinese/issues/new)

## 号外 号外 Kotlin 1.0 正式发布
Android 世界的 Swift 终于发布1.0版本 

Kotlin 是一个实用性很强的语言，专注于互通，安全，简洁，工具健全...

无缝支持 Java+Kotlin 项目，可以更少的使用样版代码，确保类型安全。

[Kotlin 1.0 更新日志](http://blog.jetbrains.com/kotlin/2016/02/kotlin-1-0-released-pragmatic-language-for-jvm-and-android/)

还换了logo :)

Kotlin LOC (软件规模代码行) 如下图

![kotlin](./kotlinLOC.png)

* [准备开始](GettingStarted/README.md) 
   * [基本语法](GettingStarted/Basic-Syntax.md) 
   * [习惯用语](GettingStarted/Idioms.md) 
   * [编码风格](GettingStarted/Coding-Conventions.md) 

* [基础](Basics/README.md) 
   * [基本类型](Basics/Basic-Types.md)
   * [包](Basics/Packages.md)
   * [控制流](Basics/Control-Flow.md)
   * [返回与跳转](Basics/Returns-and-Jumps.md)

* [类和对象](ClassesAndObjects/README.md)
   * [类和继承](ClassesAndObjects/Classes-and-Inheritance.md)　
   * [属性和字段](ClassesAndObjects/Properties-and-Fields.md)　
   * [接口](ClassesAndObjects/Interfaces.md) 
   * [可见性修饰词](ClassesAndObjects/Visibility-Modifiers.md) 
   * [扩展](ClassesAndObjects/Extensions.md) 
   * [数据对象](ClassesAndObjects/Data-Classes.md) 
   * [泛型](ClassesAndObjects/Generics.md)
   * [嵌套类](ClassesAndObjects/NestedClasses.md) 
   * [枚举类](ClassesAndObjects/EnumClasses.md) 
   * [对象表达式和声明](ClassesAndObjects/ObjectExpressicAndDeclarations.md) 
   * [代理模式](ClassesAndObjects/Delegation.md) 
   * [代理属性](ClassesAndObjects/DelegationProperties.md) 

* [函数和lambda表达式](FunctionsAndLambdas/README.md)
   * [函数](FunctionsAndLambdas/Functions.md) 
   * [高阶函数和lambda表达式](FunctionsAndLambdas/Higher-OrderFunctionsAndLambdas.md) 
   * [内联函数](FunctionsAndLambdas/InlineFunctions.md) 

* [集合](Collections/README.md)
   * [集合概览](Collections/CollectionsOverview.md)
   * [结构化集合](Collections/ConstructionCollections.md)
   * [迭代器](Collections/Iterators.md)
   * [范围和进度](Collections/RangesandProgressions.md)
   * [序列](Collections/Squences.md)
   * [操作概览](Collections/OperationsOverview.md)
   * [转化](Collections/Transformations.md)
   * [过滤](Collections/Filtering.md)
   * [加减操作符](Collections/PlusandMinusOperators.md)
   * [分组](Collections/Grouping.md)
   * [取得部分集合](Collections/RetrievingCollectionParts.md)
   * [取得单个元素](Collections/RetrivingSingleElements.md)
   * [排序](Collections/Ording.md)
   * [聚合操作](Collections/AggregateOperations.md)
   * [集合写曹锁](Collections/CollectionWriteOperations.md) 
   * [只针对于 list 的操作](Collections/ListSepcificOperations.md)
   * [只针对于 set 的操作](Collections/SetSepcificOperations.md)
   * [只针对于 map 的操作](Collections/MapSepcificOperations.md)

* [协程](Coroutines/README.md)
   * [协程指南](Coroutines/CoroutinesGuide.md)
   * [基础](Coroutines/Basics.md)
   * [取消和超时](Coroutines/CancellationAndTimeouts.md)
   * [频道](Coroutines/Channels.md)
   * [组合挂起函数](Coroutines/ComposingSuspendingFunctions.md)
   * [协程上下文和调度器](Coroutines/CoroutineContextAndDispatchers.md)
   * [异常处理](Coroutines/ExceptionHandling.md)
   * [Select 表达式](Coroutines/SelectExpression.md)
   * [共享可变状态与并发](Coroutines/SharedMutableStateAndConcurrency.md)


* [更多语言结构](MoreLanguageConstructs/README.md)
   * [解构声明](MoreLanguageConstructs/DestructuringDeclarations.md)
   * [类型检查和自动转换](MoreLanguageConstructs/Type-Checks-and-Casts.md)
   * [This表达式](MoreLanguageConstructs/This-Expression.md)
   * [等式](MoreLanguageConstructs/Equality.md)
   * [运算符重载](MoreLanguageConstructs/Opetator-overloading.md)
   * [空安全](MoreLanguageConstructs/Null-Safety.md)
   * [异常](MoreLanguageConstructs/Exceptions.md)
   * [注解](MoreLanguageConstructs/Annotations.md)
   * [反射](MoreLanguageConstructs/Reflection.md)
   * [作用域函数](MoreLanguageConstructs/ScopeFunctions.md)
   * [类型安全构造器](MoreLanguageConstructs/Type-SafeBuilders.md)
   * [试验性 API 标注](MoreLanguageConstructs/ExperimentalAPIMarkers.md)



* [参考](Reference/README.md)
    * [API](Reference/API-Reference.md) 
    * [语法](Reference/Grammar.md)
* [互用性](Interop/README.md)
   * [与 java 交互](Interop/Java-Interop.md)

* [工具](Tools/README.md) 
   * [Kotlin代码文档](Tools/Documenting-Kotlin-Code.md)
   * [使用Maven](Tools/Using-Maven.md) 
   * [使用Ant](Tools/Using-Ant.md) 
   * [使用Griffon](Tools/Using-Griffon.md) 
   * [使用Gradle](Tools/Using-Gradle.md)　

* [FAQ](FAQ/README.md)
   * [与java对比](FAQ/Comparison2java.md) 
   * [与Scala对比](FAQ/Comparison2Scala.md) 
