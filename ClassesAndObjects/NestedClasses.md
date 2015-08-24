##嵌套类

类可以嵌套在其他类中

```kotlin
class Outer {
	private val bar: Int = 1
	class Nested {
		fun foo() = 2
	}
}

val demo = Outer.Nested().foo() //==2
```

###内部类

类可以标记为 inner 这样就可以访问外部类的成员。内部类拥有外部类的一个对象引用：

```kotlin
class Outer {
	private val bar: Int = 1
	inner class Inner {
		fun foo() = bar
	}
}

val demo = Outer().Inner().foo() //==1
```

参看[这里](http://kotlinlang.org/docs/reference/this-expressions.html)了解更过 this 在内部类的用法