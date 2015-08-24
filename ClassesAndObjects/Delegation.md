##代理

###类代理

[代理模式](https://en.wikipedia.org/wiki/Delegation_pattern) 给实现继承提供了很好的代替方式， Kotlin 原生支持它，所以并不需要什么样板代码。`Derived` 类可以继承 `Base` 接口并且代理了它全部的公共方法：

```kotlin
interface Base {
	fun print()
}

class BaseImpl(val x: Int) : Base {
	override fun print() { printz(x) }
}

class Derived(b: Base) : Base by b

fun main() {
	val b = BaseImpl(10)
	Derived(b).print()
}
```

在 `Derived` 的父类列表中的条款意味这 `b` 将会存储在 `Derived` 对象中并且编译器会生成 `Base` 的所有方法并转给 `b`。