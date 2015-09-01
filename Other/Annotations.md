##注解

###注解声明

注解是一种将元数据附加到代码中的方法。声明注解需要在类前面使用 annotation 关键字：

```kotlin
annotation class fancy
```

###用法

```kotlin
@fancy class Foo {
	@fancy fun baz(@fancy foo: Int): Int {
		return (@fancy 1)
	}
}
```

在多数情形中 @ 标识是可选的。只有在注解表达式或本地声明中才必须：

```kotlin
fancy class Foo {
	fancy fun baz(fancy foo: Int): Int {
		@fancy fun bar() { ... }
		return (@fancy 1)
	}
}
```

如果要给构造函数注解，就需要在构造函数声明时添加 constructor 关键字，并且需要在前面添加注解：

```kotlin
class Foo @inject constructor (dependency: MyDependency)
	//...
```

也可以注解属性访问者：

```kotlin
class Foo {
	var x: MyDependency?=null
		@inject set
}
```

###构造函数

注解可以有带参数的构造函数。

```kotlin
annotation class special(val why: String)
special("example") class Foo {}
```

###Lambdas

注解也可以用在 Lambda 中。这将会应用到 lambda 生成的 invoke() 方法。这对 [Quasar](http://www.paralleluniverse.co/quasar/)框架很有用，在这个框架中注解被用来并发控制

```kotlin
annotation class Suspendable
val f = @Suspendable { Fiber.sleep(10) }
```

###java 注解

java 注解在 kotlin 中是完全兼容的：

```kotlin
import org.junit.Test
import org.junit.Assert.*

class Tests {
  Test fun simple() {
    assertEquals(42, getTheAnswer())
  }
}
```

java 注解也可以通过在导入是重命名实现像修改者那样：

```kotlin
import org.junit.Test as test

class Tests {
  test fun simple() {
    ...
  }
}
```

因为 java 中注解参数顺序是没定义的，你不能通过传入参数的方法调用普通函数。相反，你需要使用命名参数语法：

```java
//Java
public @interface Ann {
	int intValue();
	String stringValue(0;
}

//kotlin
Ann(intValue = 1, stringValue = "abc") class C
```

像 java 中那样，值参数是特殊的情形；它的值可以不用明确的名字。

```java
public @interface AnnWithValue {
	String value();
}

//kotlin
AnnWithValue("abc") class C
```

如果java 中的 value 参数有数组类型，则在 kotlin 中变成 vararg 参数：

```kotlin
// Java
public @interface AnnWithArrayValue {
    String[] value();
}
// Kotlin
AnnWithArrayValue("abc", "foo", "bar") class C

```

如果你需要明确一个类作为一个注解参数，使用 Kotlin 类[KClass](http://kotlinlang.org/api/latest/jvm/stdlib/kotlin.reflect/-k-class/index.html)。Kotlin 编译器会自动把它转为 java 类，因此 java 代码就可以正常看到注解和参数了。

```kotlin
import kotlin.reflect.KClass

annotation class Ann(val arg1: KClass<*>, val arg2: KClass<out Any?>)

Ann(String::class, Int::class) class MyClass
```

注解实例的值在 kotlin 代码中是暴露属性。

```kotlin
// Java
public @interface Ann {
    int value();
}
// Kotlin
fun foo(ann: Ann) {
    val i = ann.value
}
```