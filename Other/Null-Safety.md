##空安全

###可空类型和非空类型

Kotlin 类型系统致力于消灭空引用。

在许多语言中都存在的一个大陷阱包括 java ，就是访问一个空引用的成员，结果会有空引用异常。在 java 中这就是 `NullPointerException` 或者叫 NPE

Kotlin 类型系统致力与消灭 `NullPointerException` 异常。唯一可能引起 NPE 异常的可能是：

>明确调用 `throw NullPointerException()`
>外部 java 代码引起
>一些前后矛盾的初始化(在构造函数中没初始化的成员在其它地方使用)

在 Kotlin 类型系统中可以为空和不可为空的引用是不同的。比如，普通的 `String` 类型的变量不能为空：

```kotlin
var a: String ="abc"
a = null //编译错误　
```

允许为空，我们必须把它声明为可空的变量：

```kotlin
var b: String? = "abc"
b = null
```

现在你可以调用 a 的方法，而不用担心 NPE 异常了：

```kotlin
val l = a.length()
```

但如果你想使用 b 调用同样的方法就有可能报错了：

```kotlin
val l = b.length() //错误：b 不可为空
```

但我们任然想要调用方法，有些办法可以解决。

###在条件中检查 null

首先，你可以检查 `b` 是否为空，并且分开处理下面选项：

```kotlin
val l = if (b != null) b.length() else -1
```

编译器会跟踪你检查的信息并允许在 if 中调用 length()。更复杂的条件也是可以的：

```kotlin
if (b != null && b.length() >0)
  print("Stirng of length ${b.length}")
else
  print("Empty string")
```

注意只有在 b 是不可变时才可以

###安全调用

第二个选择就是使用安全操作符，`?.`:

```kotlin
b?.length()
```

如果 b 不为空则返回长度，否则返回空。这个表达式的的类型是 Int?

安全调用在链式调用是是很有用的。比如，如果 Bob 是一个雇员可能分配部门(也可能不分配)，如果我们想获取 Bob 的部门名作为名字的前缀，就可以这样做：

```kotlin
bob?.department?.head?.name
```
这样的调用链在任何一个属性为空都会返回空。

### Elvis 操作符

当我们有一个 r 的可空引用时，我们可以说如果 `r` 不空则使用它，否则使用使用非空的 x :

```kotlin
val l: Int = if (b != null) b.length() else -1
```
尽管使用 if 表达式我们也可以使用　Elvis 操作符，`?:`

```kotlin
val l = b.length()?: -1
```
如果 ?: 左边表达式不为空则返回，否则返回右边的表达式。注意右边的表带式只有在左边表达式为空是才会执行

注意在 Kotlin 中 throw return 是表达式，所以它们也可以在 Elvis 操作符右边。这是非常有用的，比如检查函数参数是否为空；

```kotlin
fun foo(node: Node): String? {
  val parent = node.getParent() ?: return null
  val name = node.getName() ?: throw IllegalArgumentException("name expected")

  //...
}
```
### !! 操作符

第三个选择是 NPE-lovers。我们可以用 b!! ，这会返回一个非空的 b 或者抛出一个 b 为空的 NPE

```kotlin
val l = b !!.length()
```
###安全转换

普通的转换可能产生 `ClassCastException` 异常。另一个选择就是使用安全转换，如果不成功就返回空：

```kotlin
val aInt: Int? = a as? Int
```
