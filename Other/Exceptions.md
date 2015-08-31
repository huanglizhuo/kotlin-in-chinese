##异常

###异常类

所有的异常类都是 `Exception` 的子类。每个异常都有一个消息，栈踪迹和可选的原因。

使用 throw 表达式，抛出异常

```kotlin
throw MyException("Hi There!")
```
使用 try 捕获异常

```kotlin
try {
  // some code
}
catch (e: SomeException) {
  // handler
}
finally {
  // optional finally block
}
```

有可能有不止一个的 catch 块。finally 块可以省略。

###try 是一个表达式

try 可以有返回值：

```kotlin
val a: Int? = try { parseInt(input) } catch (e: NumberFormatException) { null }
```

try 返回值要么是 try 块的最后一个表达式，要么是 catch 块的最后一个表达式。finally 块的内容不会对表达式有任何影响。

###检查异常

Kotlin 中没有异常检查。这是由多种原因造成的，我们这里举个简单的例子

下面是 JDK `StringBuilder` 类实现的一个接口

```java
Appendable append(CharSequence csq) throws IOException;
```

这个签名说了什么？　它说每次我把 string 添加到什么东西(StringBuilder 或者 log console 等等)上时都会捕获 `IOExceptions` 为什么呢？因为可能涉及到 IO 操作(Writer 也实现了 Appendable)...　所以导致所有实现 Appendable 的接口都得捕获异常

```java
try {
  log.append(message)
}
catch (IOException e) {
  // Must be safe
}
```

这样是不利的，参看[Effective java ](http://www.oracle.com/technetwork/java/effectivejava-136174.html)

Bruce Eckel 在[java 需要异常检查吗?](http://www.mindview.net/Etc/Discussions/CheckedExceptions)说到：

> Examination of small programs leads to the conclusion that requiring exception specifications could both enhance developer productivity and enhance code quality, but experience with large software projects suggests a different result – decreased productivity and little or no increase in code quality.

###java 互动

参看 [Java Interoperability section](http://kotlinlang.org/docs/reference/java-interop.html)