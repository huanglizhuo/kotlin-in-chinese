##枚举类

枚举类最基本的用法就是实现类型自举

```kotlin
enum class Direction {
	NORTH,SOUTH,WEST
}
```

每个自举常量都是一个对象。枚举常量通过逗号分开。

###初始化

因为每个枚举都是枚举类的一个实例，它们是可以初始化的。

```kotlin
enum class Color(val rgb: Int) {
	RED(0xFF0000),
	GREEN(0x00FF00),
	BLUE(0x0000FF)
}
```

###匿名类

枚举实例也可以声明它们自己的匿名类

```kotlin
enum class ProtocolState {
	WAITING {
		override fun signal() = Taking
	},
	Taking{
		override fun signal() = WAITING
	};
	abstract fun signal(): ProtocolState
}
```

通过对应的方法，以及复写基本方法。注意如果枚举定义了任何成员，你需要像在 java 中那样把枚举常量定义和成员定义分开。

###使用枚举常量

像 java 一样，Kotlin 中的枚举类有合成方法允许列出枚举常量的定义并且通过名字获得枚举常量。这些方法的签名就在下面列了出来(假设枚举类名字是 EnumClass)：

```kotlin
EnumClass.valueOf(value: String): EnumClass
EnumClass.values(): Array<EnumClass>
```

如果指定的名字不匹配枚举类中任何定义个枚举常量那么`valueOf()`方法将会抛出参数异常。

每个枚举常量都有或取在枚举类中声明的名字和位置的方法：

```kotlin
name(): Sting
ordinal(): Int
```

枚举类也实现了 [Comparable](http://kotlinlang.org/api/latest/jvm/stdlib/kotlin/-comparable/index.html) 接口，比较时使用的是它们在枚举类定义的自然顺序。