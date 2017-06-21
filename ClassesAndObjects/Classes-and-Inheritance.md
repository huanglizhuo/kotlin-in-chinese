## 类和继承
### 类

在 Kotlin 中类用 `class` 声：

```kotlin
class Invoice {
}
```

类的声明包含类名，类头(指定类型参数，主构造函数等等)，以及类主体，用大括号包裹。类头和类体是可选的；如果没有类体可以省略大括号。

```kotlin
class Empty
```

### 构造函数
在 Kotlin 中类可以有一个主构造函数以及多个二级构造函数。主构造函数是类头的一部分：跟在类名后面(可以有可选的类型参数)。

```kotlin
class Person constructor(firstName: String) {
}
```

如果主构造函数没有注解或可见性说明，则 `constructor` 关键字是可以省略：

```korlin
class Person(firstName: String){
}
```

主构造函数不能包含任意代码。初始化代码可以放在以 `init` 做前缀的初始化块内

```kotlin
class Customer(name: String) {
	init {
		logger,info("Customer initialized with value ${name}")
	}
}
```

注意主构造函数的参数可以用在初始化块内，也可以用在类的属性初始化声明处：

```kotlin
class Customer(name: String) {
	val customerKry = name.toUpperCase()
}
```

事实上，声明属性并在主构造函数中初始化,在 Kotlin 中有更简单的语法：

```kotlin
class Person(val firstName: String, val lastName: String, var age: Int) {
}
```

就像普通的属性，在主构造函数中的属性可以是可变的(`var`)或只读的(`val`)。

如果构造函数有注解或可见性声明，则 `constructor` 关键字是不可少的，并且可见性应该在前：

```kotlin
class Customer public @inject constructor (name: String) {...}
```

参看[可见性](http://kotlinlang.org/docs/reference/visibility-modifiers.html#constructors)

### 二级构造函数
类也可以有二级构造函数，需要加前缀 `constructor`:

```kotlin
class Person {
	constructor(parent: Person) {
		parent.children.add(this)
	}
}
```

如果类有主构造函数，每个二级构造函数都要，或直接或间接通过另一个二级构造函数代理主构造函数。在同一个类中代理另一个构造函数使用 `this` 关键字：

```kotlin
class Person(val name: String) {
	constructor (name: String, paret: Person) : this(name) {
		parent.children.add(this)
	}
}
```

如果一个非抽象类没有声明构造函数(主构造函数或二级构造函数)，它会产生一个没有参数的构造函数。该构造函数的可见性是 public 。如果你不想你的类有公共的构造函数，你就得声明一个拥有非默认可见性的空主构造函数：

```kotlin
class DontCreateMe private constructor () {
}
```

>注意：在 JVM 虚拟机中，如果主构造函数的所有参数都有默认值，编译器会生成一个附加的无参的构造函数，这个构造函数会直接使用默认值。这使得 Kotlin 可以更简单的使用像 Jackson 或者 JPA 这样使用无参构造函数来创建类实例的库。

```kotlin
class Customer(val customerName: String = "")
```

### 创建类的实例
我们可以像使用普通函数那样使用构造函数创建类实例：

```kotlin
val invoice = Invoice()
val customer = Customer("Joe Smith")
```

注意 Kotlin 没有 `new` 关键字。

创建嵌套类、内部类或匿名类的实例参见[嵌套类](http://kotlinlang.org/docs/reference/nested-classes.html)

### 类成员
类可以包含：
>-- 构造函数和初始化代码块

>-- [函数](FunctionsAndLambdas/Functions.md)

>-- [属性](ClassesAndObjects/Properties-and-Fields.md)　

>-- [内部类](ClassesAndObjects/NestedClasses.md)

>-- [对象声明](ClassesAndObjects/ObjectExpressicAndDeclarations.md)

### 继承
Kotlin 中所有的类都有共同的父类 `Any` ，它是一个没有父类声明的类的默认父类：

```kotlin
class Example //　隐式继承于 Any
```

`Any` 不是 `java.lang.Object`；事实上它除了 `equals()`,`hashCode()`以及`toString()`外没有任何成员了。参看[Java interoperability]( Java interoperability)了解更多详情。

声明一个明确的父类，需要在类头后加冒号再加父类：

```kotlin
open class Base(p: Int)

class Derived(p: Int) : Base(p)
```

如果类有主构造函数，则基类可以而且是必须在主构造函数中使用参数立即初始化。

如果类没有主构造函数，则必须在每一个构造函数中用 `super` 关键字初始化基类，或者在代理另一个构造函数做这件事。注意在这种情形中不同的二级构造函数可以调用基类不同的构造方法：

```kotlin
class MyView : View {
	constructor(ctx: Context) : super(ctx) {
	}
	constructor(ctx: Context, attrs: AttributeSet) : super(ctx,attrs) {
	}
}
```

`open`注解与java中的`final`相反:它允许别的类继承这个类。默认情形下，kotlin 中所有的类都是 final ,对应 [Effective Java](http://www.oracle.com/technetwork/java/effectivejava-136174.html) ：Design and document for inheritance or else prohibit it.

### 复写方法
像之前提到的，我们在 kotlin 中坚持做明确的事。不像 java ，kotlin 需要把可以复写的成员都明确注解出来，并且重写它们：

```kotlin
open class Base {
	open fun v() {}
	fun nv() {}
}

class Derived() : Base() {
	override fun v() {}
}
```

对于 `Derived.v()` 来说`override`注解是必须的。如果没有加的话，编译器会提示。如果没有`open`注解，像 `Base.nv()` ,在子类中声明一个同样的函数是不合法的，要么加`override`要么不要复写。在 final 类(就是没有open注解的类)中，`open` 类型的成员是不允许的。

标记为`override`的成员是open的，它可以在子类中被复写。如果你不想被重写就要加 final:

```kotlin
open class AnotherDerived() : Base() {
	final override fun v() {}
}
```

**等等！我现在怎么hack我的库？！**

有个问题就是如何复写子类中那些作者不想被重写的类，下面介绍一些令人讨厌的方案。

我们认为这是不好的，原因如下：

> 最好的实践建议你不应给做这些 hack

>人们可以用其他的语言成功做到类似的事情

>如果你真的想 hack 那么你可以在 java 中写好 hack 方案，然后在 kotlin 中调用 (参看[java调用](http://kotlinlang.org/docs/reference/java-interop.html))，专业的构架可以很好的做到这一点

### 复写属性
复写属性与复写方法类似，在一个父类上声明的属性在子类上被重新声明，必须添加`override`，并且它们必须具有兼容的类型。每个被声明的属性都可以被一个带有初始化器的属性或带有getter方法的属性覆盖

```kotlin
open class Foo {
  open val x: Int get { ... }
}

class Bar1 : Foo() {
  override val x: Int = ...
}
```

您还可以使用`var`属性覆盖一个`val`属性，但反之则不允许。这是允许的，因为`val`属性本质上声明了一个getter方法，并将其重写为`var`，另外在派生类中声明了setter方法。

注意，可以在主构造函数中使用`override`关键字作为属性声明的一部分。

```kotlin
interface Foo {
    val count: Int
}

class Bar1(override val count: Int) : Foo

class Bar2 : Foo {
    override var count: Int = 0
}
```

### 复写规则
在 kotlin 中，实现继承通常遵循如下规则：如果一个类从它的直接父类继承了同一个成员的多个实现，那么它必须复写这个成员并且提供自己的实现(或许只是直接用了继承来的实现)。为表示使用父类中提供的方法我们用 `super<Base>`表示:

```kotlin
open class A {
	open fun f () { print("A") }
	fun a() { print("a") }
}

interface B {
	fun f() { print("B") } // 接口的成员变量默认是 open 的
	fun b() { print("b") }
}

class C() : A() , B {
	// 编译器会要求复写f()
	override fun f() {
		super<A>.f() // 调用 A.f()
		super<B>.f() // 调用 B.f()
	}
}
```

可以同时从 A 和 B 中继承方法，而且 C 继承 a() 或 b() 的实现没有任何问题，因为它们都只有一个实现。但是 f() 有俩个实现，因此我们在 C 中必须复写 f() 并且提供自己的实现来消除歧义。

### 抽象类
一个类或一些成员可能被声明成 abstract 。一个抽象方法在它的类中没有实现方法。记住我们不用给一个抽象类或函数添加 open 注解，它默认是带着的。

我们可以用一个抽象成员去复写一个带 open 注解的非抽象方法。

```kotlin
open class Base {
	open fun f() {}
}

abstract class Derived : Base() {
	override abstract fun f()
}
```

### 伴随对象
在 kotlin 中不像 java 或者 C# 它没有静态方法。在大多数情形下，我们建议只用包级别的函数。

如果你要写一个没有实例类就可以调用的方法，但需要访问到类内部(比如说一个工厂方法)，你可以把它写成它所在类的一个[成员](http://kotlinlang.org/docs/reference/object-declarations.html)(you can write it as a member of an object declaration inside that class)

更高效的方法是，你可以在你的类中声明一个[伴随对象](http://kotlinlang.org/docs/reference/object-declarations.html#companion-objects)，这样你就可以像 java/c# 那样把它当做静态方法调用，只需要它的类名做一个识别就好了

### 密封类
密封类用于代表严格的类结构，值只能是有限集合中的某中类型，不可以是任何其它类型。这就相当于一个枚举类的扩展：枚举值集合的类型是严格限制的，但每个枚举常量只有一个实例，而密封类的子类可以有包含不同状态的多个实例。

声明密封类需要在 class 前加一个 sealed 修饰符。密封类可以有子类但必须全部嵌套在密封类声明内部、

```Kotlin
sealed class Expr {
	class Const(val number: Double) : Expr()
	class Sum(val e1: Expr, val e2: Expr) : Expr()
	object NotANumber : Expr()
}
```

注意密封类子类的扩展可以在任何地方，不必在密封类声明内部进行。

使用密封类的最主要的的好处体现在你使用 [when 表达式]()。可以确保声明可以覆盖到所有的情形，不需要再使用 else 情形。

```Kotlin
fun eval(expr: Expr): Double = when(expr) {
	is Const -> expr.number
	is Sum -> eval(expr.e1) + eval(expr.e2)
	NotANumber -> Double.NaN
    // the `else` clause is not required because we've covered all the cases
}
```
