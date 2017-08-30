## 数据类
我们经常创建一个只保存数据的类。在这样的类中一些函数只是机械的对它们持有的数据进行一些推导。在 kotlin 中这样的类称之为 data 类，用 `data` 标注:

```kotlin
data class User(val name: String, val age: Int)
```

编译器会自动根据主构造函数中声明的所有属性添加如下方法：

>`equals()`/`hashCode` 函数

> `toString` 格式是 "User(name=john, age=42)"

> [compontN()functions] (http://kotlinlang.org/docs/reference/multi-declarations.html) 对应按声明顺序出现的所有属性

> `copy()` 函数

如果在类中明确声明或从基类继承了这些方法，编译器不会自动生成。

为确保这些生成代码的一致性，并实现有意义的行为，数据类要满足下面的要求：

注意如果构造函数参数中没有 `val` 或者 `var` ，就不会在这些函数中出现；

> 主构造函数应该至少有一个参数；

> 主构造函数的所有参数必须标注为 `val` 或者 `var` ；

> 数据类不能是 abstract，open，sealed，或者 inner ；

> 数据类不能继承其它的类（但可以实现接口）。
>
> (在1.1之前)数据类只能实现接口。

从1.1开始数据类可以继承其它类（参考 [Sealed classes](http://kotlinlang.org/docs/reference/sealed-classes.html) ）



在 JVM 中如果构造函数是无参的，则所有的属性必须有默认的值，(参看[Constructors](http://kotlinlang.org/docs/reference/classes.html#constructors));

> data class User(val name: String = "", val age: Int = 0)

### 复制
我们经常会对一些属性做修改但想要其他部分不变。这就是 `copy()` 函数的由来。在上面的 User 类中，实现起来应该是这样：

```kotlin
fun copy(name: String = this.name, age: Int = this.age) = User(name, age)
```

有了 copy 我们就可以像下面这样写了：

```kotlin
val jack = User(name = "jack", age = 1)
val olderJack = jack.copy(age = 2)
```

### 数据类和多重声明
组件函数允许数据类在[多重声明](http://kotlinlang.org/docs/reference/multi-declarations.html)中使用：

```kotlin
val jane = User("jane", 35)
val (name, age) = jane
println("$name, $age years of age") //打印出 "Jane, 35 years of age"
```

### 标准数据类
标准库提供了 `Pair` 和 `Triple`。在大多数情形中，命名数据类是更好的设计选择，因为这样代码可读性更强而且提供了有意义的名字和属性。
