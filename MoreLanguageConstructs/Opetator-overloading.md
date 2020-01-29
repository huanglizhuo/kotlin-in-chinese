## 运算符重载
Kotlin 允许我们实现一些我们自定义类型的运算符实现。这些运算符有固定的表示，和固定的优先级。为实现这样的运算符，我们提供了固定名字的数字函数和扩展函数，比如二元运算符的左值和一元运算符的参数类型。

### 转换
这里我们描述了一些常用运算符的重载

#### 一元运算符
**表达式**|**转换**
---|---
+a|a.plus()
-a|a.minus()
!a|a.not()

这张表解释了当编译器运行时，比如，表达式 `+a` ，是这样运行的：

>决定 a 的类型，假设是 T
>寻找接收者是 T 的无参函数 `plus()` ，比如数字函数或者扩展函数
>如果这样的函数缺失或不明确，则返回错误。
>如果函数是当前函数或返回类型是 `R` 则表达式 `+a` 是 `R` 类型。

注意这些操作符和其它的一样，都被优化为基本类型并且不会产生多余的开销。

**表达式**|**转换**
---|---
a++| a.inc() + see below
a--| a.dec() + see below

这些操作符允许修改接收者和返回类型。

```kotlin
inc()/dec() shouldn’t mutate the receiver object.
By “changing the receiver” we mean the receiver-variable, not the receiver object.
```

编译器是这样解决有后缀的操作符的比如 `a++` :

>决定 a 的类型，假设是 T
>寻找无参函数 `inc()` ，作用在接收者T 
>如果返回类型是 R ，则必须是 T 的子类

计算表达式的效果是：

>把 a 的初始值存储在 a0 中
>把 a.inc() 的结果作用在 a 上
>把 a0 作为表达式的返回值

a-- 的步骤也是一样的

++a  --a 的解决方式也是一样的

#### 二元操作符
**表达式**|**转换**
---|---
a + b | a.plus(b)
a - b | a.minus(b)
a * b | a.times(b)
a / b | a.div(b)
a % b | a.mod(b)
a..b | a.rangeTo(b)

编译器只是解决了该表中翻译为列的表达式

**表达式**|**转换**
---|---
a in b | b.contains(a)
a !in b | !b.contains(a)

in 和 !in 的产生步骤是一样的，但参数顺序是相反的。

**标志** | **转换**
---|---
a[i] | a.get(i)
a[i, j] | a.get(i, j)
a[i_1, ..., i_n] | a.get(i_1, ... , i_n)
a[i] = b | a.set(i, b)
a[i,j] =b | a.set(i, j, b)
a[i_1, ... , i_n] = b | a.set(i_1,... ,o_n,b)

方括号被转换为 get set 函数

**标志** | **转换**
---|---
a(i) | a.invoke(i)
a(i, j) | a.invoke(i, j)
a(i_1, ... , i_n) | a.invoke(i_1, ..., i_n)

括号被转换为带有正确参数的 invoke 参数

**表达式** | **转换**
---|---
a += b	|a.plusAssign(b)
a -= b	|a.minusAssign(b)
a *= b	|a.timesAssign(b)
a /= b	|a.divAssign(b)
a %= b	|a.modAssign(b)

在分配 a+= b时编译器是下面这样实现的：

> 右边列的函数是否可用
>  对应的二元函数(比如 plus() )是否也可用,不可用在报告错误
> 确定它的返回值是 `Unit` 否则报告错误
> 生成 `a.plusAssign(b)` 
> 否则试着生成 a=a+b 代码

Note: assignments are NOT expressions in Kotlin.

**表达式** | **转换**
---|---
a == b	|a?.equals(b) ?: b.identityEquals(null)
a != b	|!(a?.equals(b) ?: b.identityEquals(null))

注意 ===   !== 是不允许重载的

== 操作符有俩点特别：

> 它被翻译成一个复杂的表达式，用于筛选空值，而且 null == null 是真

> 它需要带有特定签名的函数，而不仅仅是特定名称的函数，下面这样：

```kotlin
fun equals(other: Any?): Boolean
```

或者用相同的参数列表和返回类型的扩展功能

**标志** | **转换**
---|---

a > b	｜a.compareTo(b) > 0
a < b	｜a.compareTo(b) < 0
a >= b	｜a.compareTo(b) >= 0
a <= b	｜a.compareTo(b) <= 0

所有的比较都转换为 `compareTo` 的调用，这个函数需要返回 `Int` 值

### 命名函数的中缀调用
我们可以通过 [中缀函数的调用](http://kotlinlang.org/docs/reference/functions.html#infix-notation) 来模拟自定义中缀操作符
