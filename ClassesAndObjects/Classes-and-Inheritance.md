##类和继承

###类

在 Kotlin 中类用 class 声明：

```kotlin
class Invoice{
}
```

类的声明包含类名，类头(指定类型参数，主构造函数等等)，以及类主体，用大括号包裹。类头和类体是可选的；如果没有类体可以省略大括号。

```kotlin
class Empty
```

###构造函数

在 Kotlin 中类可以有一个主构造函数以及多个二级构造函数。主构造函数是类头的一部分：跟在类名后面(可以有可选的参数)。

```kotlin
class Person constructor(firstName: String){
}
```

如果主构造函数没有注解或可见性说明，则 constructor 关键字是可以省略：

```korlin
class Person(firstName: String){
}
```

主构造函数不能包含任意代码。初始化代码可以放在以 init 做前缀的初始化块内

```kotlin
class Customer(name: String){
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

事实上，声明属性并在主构造函数中初始化它们可以更简单：

```kotlin
class Person(val firstName: String, val lastName: String, var age: Int){
}
```

像平常的属性，在主构造函数中的属性可以是可变或只读。

如果构造函数有注解或可见性声明，则 constructor 关键字是不可少的，并且注解应该在前：

```kotlin
class Customer public inject constructor (name: String) {...}
```

参看[可见性](http://kotlinlang.org/docs/reference/visibility-modifiers.html#constructors)

###二级构造函数

类也可以有二级构造函数，该函数前缀是 constructor:

```kotlin
class Person {
	constructor(parent: Person) {
		parent.children.add(this)
	}
}
```

如果类有主构造函数，每个二级构造函数都要，或直接或间接通过另一个二级构造函数代理主构造函数。在同一个类中代理另一个构造函数使用 this 关键字：

```kotlin
class Person(val name: String) {
	constructor (name: String, paret: Person) : this(name) {
		parent.children.add(this)
	}
}
```

如果一个非抽象类没有声明构造函数(主构造函数或二级构造函数)，它会产生一个没有参数的构造函数。构造函数是 public 。如果你不想你的类有公共的构造函数，你就得声明一个空的主构造函数：

```kotlin
class DontCreateMe private constructor () {
}
```

>注意：在 JVM 虚拟机中，如果主构造函数的所有参数都有默认值，编译器会生成一个附加的无参的构造函数。

```kotlin
class Customer(val customerName: String = "")
```

###创建类的实例

我们可以像使用普通函数那样使用构造函数创建类实例：

```kotlin
val invoice = Invoice()
val customer = Customer("Joe Smith")
```

注意 Kotlin 没有 new 关键字。

###类成员

类可以包含：
>构造函数和初始化代码块
[函数](http://kotlinlang.org/docs/reference/functions.html)
[属性](http://kotlinlang.org/docs/reference/properties.html)
[包含内部类](http://kotlinlang.org/docs/reference/nested-classes.html)
[对象声明](http://kotlinlang.org/docs/reference/object-declarations.html)

###继承

Kotin 中所有的类都有共同的父类 Any ，下面是一个没有父类声明的类：

```kotlin
class Example //　隐式继承于 Any
```

`Any` 不是 `java.lang.Object`；事实上它除了 `equals()`,`hashCode()`以及`toString()`外没有任何成员了。参看[ Java interoperability]( Java interoperability)了解更多详情。

声明一个明确的父类，需要在类头后加冒号再加父类：

```kotlin
open class Base(p: Ont)

class Derived(p: Int) : Base(p)
```

如果类有主构造函数，则基类可以而且是必须在主构造函数中立即初始化。

如果类没有主构造函数，则必须在每一个构造函数中用 super 关键字初始化基类，或者在代理另一个构造函数做这件事。注意在这种情形中不同的二级构造函数可以调用基类不同的构造方法：

```kotlin
class MyView : View {
	constructor(ctx: Context) : super(ctx) {
	}
	constructor(ctx: Context, attrs: AttributeSet) : super(ctx,attrs) {
	}
}
```

open 注解与java 中的 final相反:它允许别的类继承这个类。默认情形下，kotlin 中所有的类都是 final 对应 [Effective Java](http://www.oracle.com/technetwork/java/effectivejava-136174.html) ：Design and document for inheritance or else prohibit it.

###复写成员

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

对于 `Derived.v()` 来说 override 注解是必须的。如果没有加的话，编译器会提示／如果没有 open 注解，像 `Base.nv()` ,在子类中声明一个同样的函数是不合法的，要么加 override 要么不要复写。在 final 类(就是没有open注解的类)中，open 类型的成员是不允许的。

标记为 override 的成员是 open的，它可以在子类中被复写。如果你不想被重写就要加 final:

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

###复写规则

在 kotlin 中，实现继承通常遵循如下规则：如果一个类从它的直接父类继承了同一个成员的多个实现，那么它必须复写这个成员并且提供自己的实现(或许只是直接用了继承来的实现)。为表示使用父类中提供的方法我们用 `super<Base>`表示:

```kotlin
open class A {
	open fun f () { print("A") }
	fun a() { print("a") }
}

interface B {
	fun f() { print("B") } //接口的成员变量默认是 open 的
	fun b() { print("b") }
}

class C() : A() , B{
	override fun f() {
		super<A>.f()//调用 A.f()
		super<B>.f()//调用 B.f()
	}
}
```

可以同时从 A B 中继承方法，而且 C 继承 a() 或 b() 的实现没有任何问题，因为它们都只有一个实现。但是 f() 有俩个实现，因此我们在 C 中必须复写 f() 并且提供自己的实现来消除歧义。

###抽象类

一个类或一些成员可能被声明成 abstract 。一个抽象方法在它的类中没有实现方法。因此当子类继承抽象成员时，它并不算一个实现：

```kotlin
abstract class A {
	abstract fun f()
}

interface B {
	open fun f() { print("B") }
}

class C() : A() , B {
	//我们是必须复写 f() 方法
}
```
记住我们不用给一个抽象类或函数添加 open 注解，它默认是带着的。

我们可以用一个抽象成员去复写一个带 open 注解的非抽象方法。

```kotlin
open class Base {
	open fun f() {}
}

abstract class Derived : Base() {
	override abstract fun f()
}
```

###伴随对象

在 kotlin 中不像 java 或者 C# 它没有静态方法。在大多数情形下，我们建议只用包级别的函数。

如果你要写一个没有实例类就可以调用的方法，但需要访问到类内部(比如说一个工厂方法)，你可以把它写成它所在类的一个成员(you can write it as a member of an object declaration inside that class)

更高效的方法是，你可以在你的类中声明一个[伴随对象](http://kotlinlang.org/docs/reference/object-declarations.html#companion-objects)，这样你就可以像 java/c# 那样把它当做静态方法调用，只需要它的类名做一个识别就好了
