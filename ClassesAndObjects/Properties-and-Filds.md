##属性和字段

###属性声明

在 Kotlin 中类可以有属性，可以 var 关键字声明为可变的，或者用 val 关键字声明为只读。

```kotlin
public class Address { 	
	public var name: String = ...
  	public var street: String = ...
	public var city: String = ...
  	public var state: String? = ...
	public var zip: String = ...
}
```

我们可以通过名字，就像使用 java 中的字段那样，一样直接使用它：

```kotlin
fun copyAddress(address: Address) : Address {
	val result = Address() //在 kotlin 中没有 new 关键字
	result.name = address.name //accessors are called
	result.street = address.street
}
```

###Getter 和 Setter 函数

声明变量的完整语法是这样的：

```kotlin
var <propertyName>: <PropertyType> [ = <property_initializer> ]
<getter>
<setter>
```

在初始化函数中，getter 和 setter 是可选的。如果属性在初始化可以被推断或者它是基类成员的复写则类型也是可选的。

例子：

```kotlin
var allByDefault: Int? // error: explicit initializer required, default getter and setter implied
var initialized = 1 // has type Int, default getter and setter
```

注意 public 或 protected 属性的类型是不能推断的,因为改变初始化时的值会对 public api 产生不想要的影响。比如：

```kotlin
public val example = 1 //错误 public 属性必须有明确的类型
```

只读属性的声明与可变属性的声明有俩点不通：以 val 开头，没有 setter 函数：

```kotlin
val simple: Int? // 有 Int 类型，默认有 getter ，必须在构造函数中初始化

val inferredType = 1//有 Int 类型以及默认的 getter
```

我们可以像普通函数那样在属性声明中写自定义的访问方式，下面是一个自定义的 getter :

```kotlin
var stringRepresentation: String
	get() = this.toString()
	set (value) {
		setDataFormString(value) // 
}
```

如果你需要改变一个访问者的可见性或者注解它，可以不用改变默认的实现，你定义一个不带函数体的访问者:

```kotlin
var settVisibilite: String = "abc"//非空类型必须初始化
	private set // setter 是私有的并且有默认的实现
var setterVithAnnotation: Any?
	@Inject set // 用 Inject 注解 setter
```

### 备用字段

在 kotlin 中不可以有字段。然而当使用自定义的访问者时需要备用字段。出于这些原因 kotlin 提供了自动备用字段，可以通过 $ 加名字使用：

```kotllin
var counter = 0 //在初始化中直接写入备用字段
	set(value) {
		if (value >= 0)
			$counter  = value
	}
```

`$counter` 字段只能在它定义的那个类中使用

编译器会检查访问体内部是否使用了备用字段，如果有就会生成

比如下面的例子中就不会有备用字段：

```kotlin
val isEmpty: Boolean
	get() = this.size == 0
```

###备用属性

如果你想要做一些事情但不适合这种 "隐含备用字段" 方案，你可以试着用备用属性的方式：

```kotlin
private var _table: Map<String, Int>? = null
public val table: Map<String, Int>
	get() {
	if (_table == null)
		_table = HashMap() //参数类型是推导出来的
		return _table ?: throw AssertionError("Set to null by another thread")
}
```

综合来讲，这些和 java 很相似，可以避免函数访问私有属性而破坏它的结构

###复写属性

参看[复写成员](http://kotlinlang.org/docs/reference/classes.html#overriding-members)

###代理属性

最常见的属性就是从备用属性中读（或者写）。另一方面，自定义的 getter 和 setter 可以实现属性的任何操作。有些像懒值( lazy values )，根据给定的关键字从 map 中读出，读取数据库，通知一个监听者等等，像这些操作介于 getter setter 模式之间。

像这样常用操作可以通过代理属性作为库来实现。更多请参看[这里](http://kotlinlang.org/docs/reference/delegated-properties.html)。