##数据类

我们经常创建一个只保存数据的类。在这样的类中一些函数只是机械的对它们持有的数据进行一些推导。在 kotlin 中这样的类可以标注为 `data` :

```kotlin
data class User(val name: String, val age: Int)
```

这叫做数据对象。编译器会根据主构造函数自动给所有属性添加如下方法：

>`equals()`/`hashCode`　

> `toString` 格式是 "User(name=john, age=42)"

> [compontN()functions] (http://kotlinlang.org/docs/reference/multi-declarations.html) 对应按声明顺序出现的所有属性

> `copy()` 

如果在类中明确声明或从基类继承了这些方法，编译器就不会自动生成了／

注意如果构造函数参数中没有 `val` 或者 `var` ，就不会在这些函数中出现；

>在 JVM 中如果构造函数是无参的，则所有的属性必须有默认的值
data class User(val name: String = "", val age: Int = 0)

###复制

我们经常会对一些属性做修改但想要其他部分不变。这就是 `copy()` 函数的由来。在上面的 User 类中，实现起来应该是这样：

```kotlin
fun copy(name: String = this.name, age: Int = this.age) = User(name, age)
```

这样就允许改写了

```kotlin
val jack = User(name = "jack", age = 1)
val olderJack = jack.copy(age = 2)
```

###数据类和多重声明

组件函数允许数据类在[多重声明](http://kotlinlang.org/docs/reference/multi-declarations.html)中使用：

```kotlin
val jane = User("jane", 35)
val (name, age) = jane
println("$name, $age years of age")
```

###标准数据类

标准库提供了 `Pair` 和 `Triple`。在大多数情形中，命名数据类是更好的设计选择，因为这样代码可读性更强而且提供了有意义的名字和属性。