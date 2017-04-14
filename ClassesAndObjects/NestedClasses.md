## 嵌套类
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

### 内部类
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

参看[这里](http://kotlinlang.org/docs/reference/this-expressions.html)了解更多 this 在内部类的用法

###  匿名内部类
匿名内部类的实例是通过 [对象表达式](ClassesAndObjects/ObjectExpressicAndDeclarations.md)  创建的：

```kotlin
window.addMouseListener(object: MouseAdapter() {
    override fun mouseClicked(e: MouseEvent) {
        // ...
    }
                                                                                                            
    override fun mouseEntered(e: MouseEvent) {
        // ...
    }
})
```

如果对象是函数式的 java 接口的实例（比如只有一个抽象方法的 java 接口），你可以用一个带接口类型的 lambda 表达式创建它。



```kot
val listener = ActionListener { println("clicked") }
```



