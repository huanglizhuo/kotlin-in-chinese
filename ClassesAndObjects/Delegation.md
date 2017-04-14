## 代理
### 类代理
[代理模式](https://en.wikipedia.org/wiki/Delegation_pattern) 给实现继承提供了很好的代替方式， Kotlin 在语法上支持这一点，所以并不需要什么样板代码。`Derived` 类可以继承 `Base` 接口并且指定一个对象代理它全部的公共方法：

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

在 `Derived` 的父类列表中的 by 从句会将 `b` 存储在 `Derived` 内部对象，并且编译器会生成 `Base` 的所有方法并转给 `b`。
