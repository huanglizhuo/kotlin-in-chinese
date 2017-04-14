## 属性和字段
### 属性声明
在 Kotlin 中类可以有属性，我们可以使用 var 关键字声明可变属性，或者用 val 关键字声明只读属性。

```kotlin
public class Address { 	
	public var name: String = ...
  	public var street: String = ...
	public var city: String = ...
  	public var state: String? = ...
	public var zip: String = ...
}
```

我们可以像使用 java 中的字段那样,通过名字直接使用一个属性：

```kotlin
fun copyAddress(address: Address) : Address {
	val result = Address() //在 kotlin 中没有 new 关键字
	result.name = address.name //accessors are called
	result.street = address.street
}
```

### Getter 和 Setter 
声明一个属性的完整语法如下：

```kotlin
var <propertyName>: <PropertyType> [ = <property_initializer> ]
	<getter>
	<setter>
```

语法中的初始化语句，getter 和 setter 都是可选的。如果属性类型可以从初始化语句或者类的成员函数中推断出来,那么他的类型也是忽略的。

例子：

```kotlin
var allByDefault: Int? // 错误: 需要一个初始化语句, 默认实现了 getter 和 setter 方法
var initialized = 1 // 类型为 Int, 默认实现了 getter 和 setter
```

只读属性的声明语法和可变属性的声明语法相比有两点不同: 它以 val 而不是 var 开头，不允许 setter 函数：

```kotlin
val simple: Int? // 类型为 Int ，默认实现 getter ，但必须在构造函数中初始化

val inferredType = 1 // 类型为 Int 类型,默认实现 getter
```

我们可以像写普通函数那样在属性声明中自定义的访问器，下面是一个自定义的 getter 的例子:

```kotlin
var isEmpty: Boolean
	get() = this.size == 0
```

下面是一个自定义的setter:

```kotlin
var stringRepresentation: String
	get() = this.toString()
	set (value) {
		setDataFormString(value) // 格式化字符串,并且将值重新赋值给其他元素
}
```

为了方便起见,setter 方法的参数名是value,你也可以自己任选一个自己喜欢的名称.

如果你需要改变一个访问器的可见性或者给它添加注解，但又不想改变默认的实现，那么你可以定义一个不带函数体的访问器:

```kotlin
var settVisibilite: String = "abc"//非空类型必须初始化
	private set // setter 是私有的并且有默认的实现
var setterVithAnnotation: Any?
	@Inject set // 用 Inject 注解 setter
```

###  备用字段
在 kotlin 中类不可以有字段。然而当使用自定义的访问器时有时候需要备用字段。出于这些原因 kotlin 使用 `field` 关键词提供了自动备用字段，

```kotllin
var counter = 0 //初始化值会直接写入备用字段
	set(value) {
		if (value >= 0)
			field  = value
	}
```

`field` 关键词只能用于属性的访问器.

编译器会检查访问器的代码,如果使用了备用字段(或者访问器是默认的实现逻辑)，就会自动生成备用字段,否则就不会.

比如下面的例子中就不会有备用字段：

```kotlin
val isEmpty: Boolean
	get() = this.size == 0
```

### 备用属性
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

### 编译时常量
那些在编译时就能知道具体值的属性可以使用 `const` 修饰符标记为 *编译时常量*. 这种属性需要同时满足以下条件:

* Top-level or member of an object   //宝宝不会翻 :( :( :( 

* 以`String`或基本类型进行初始化

* 没有自定义getter
	
这种属性可以被当做注解使用:
```kotlin
const val SUBSYSTEM_DEPRECATED: String = "This subsystem is deprecated"
@Deprected(SUBSYSTEM_DEPRECATED) fun foo() { ... }
```

### 延迟初始化属性
通常,那些被定义为拥有非空类型的属性,都需要在构造器中初始化.但有时候这并没有那么方便.例如在单元测试中,属性应该通过依赖注入进行初始化,
或者通过一个 setup 方法进行初始化.在这种条件下,你不能在构造器中提供一个非空的初始化语句,但是你仍然希望在访问这个属性的时候,避免非空检查.

为了处理这种情况,你可以为这个属性加上 `lateinit` 修饰符

```kotlin
public class MyTest {
	lateinit var subject: TestSubject
	
	@SetUp fun setup() {
		subject = TestSubject()
	}
	
	@Test fun test() {
		subject.method() 
	}
}
```

这个修饰符只能够被用在类的 var 类型的可变属性定义中,不能用在构造方法中.并且属性不能有自定义的 getter 和 setter访问器.这个属性的类型必须是非空的,同样也不能为一个基本类型.

在一个延迟初始化的属性初始化前访问他,会导致一个特定异常,告诉你访问的时候值还没有初始化.

### 复写属性
参看[复写成员](http://kotlinlang.org/docs/reference/classes.html#overriding-members)

### 代理属性
最常见的属性就是从备用属性中读（或者写）。另一方面，自定义的 getter 和 setter 可以实现属性的任何操作。有些像懒值( lazy values )，根据给定的关键字从 map 中读出，读取数据库，通知一个监听者等等，像这些操作介于 getter setter 模式之间。

像这样常用操作可以通过代理属性作为库来实现。更多请参看[这里](http://kotlinlang.org/docs/reference/delegated-properties.html)。
