##动态类型

作为静态类型的语言，kotlin任然拥有与无类型或弱类型语言的调用，比如 javaScript。为了方便使用，`dynamic`应而生：

```kotlin
val dyn: dynamic = ...
```

`dynamic` 类型关闭了 kotlin 的类型检查：

>这样的类型可以分配任意变量或者在任意的地方作为参数传递
>任何值都可以分配为`dynamic` 类型，或者作为参数传递给任何接受 `dynamic` 类型参数的函数
>这样的类型不做 null 检查

`dynamic` 最奇特的特性就是可以在 `dynamic` 变量上调用任何属性或任何方法：
(The most peculiar feature of dynamic is that we are allowed to call any property or function with any parameters on a dynamic variable:)
```kotlin
dyn.whatever(1, "foo", dyn) // 'whatever' is not defined anywhere
dyn.whatever(*array(1, 2, 3))
```

在 javaScript 平台上这样的代码会编译为 'as is': `dyn.whatever(1)` 在 kotlin 中变成了 javaScript 中的代码
(On the JavaScript platform this code will be compiled “as is”: dyn.whatever(1) in Kotlin becomes dyn.whatever(1) in the generated JavaScript code.)

动态调用可以返回 `dynamic` 作为结果，因此我们可以轻松实现链式调用：

```kotlin
dyn.foo().bar.bat(0
```

当给动态调用传递一个 lambda 表达式时，所有的参数默认都是 `dynamic`：

```kotlin
dyn.foo {
  x -> x.bar() // x is dynamic
}
```

更多细节参看[spec document](https://github.com/JetBrains/kotlin/blob/master/spec-docs/dynamic-types.md)
