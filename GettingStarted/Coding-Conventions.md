[原文](http://kotlinlang.org/docs/reference/coding-conventions.html)

## 编码规范

本页包含了当前 kotlin 语言的代码风格。

### 命名风格
如有疑惑，默认为Java编码约定，比如：

> --使用骆驼命名法(在命名中避免下划线)

> --类型名称首字母大写

> --方法和属性首字母小写

> --缩进用四个空格

> --public 方法要写说明文档，这样它就可以出现在 Kotllin Doc 中

### 冒号
在冒号区分类型和父类型中要有空格，在实例和类型之间是没有空格的：

```kotlin
interface Foo<out T : Any> : Bar {
    fun foo(a: Int): T
}
```

### Lambdas
在 Lambdas 表达式中，大括号与表达式间要有空格，箭头与参数和函数体间要有空格。lambda表达应尽可能不要写在圆括号中

```Kotlin
list.filter { it > 10 }.map { element -> element * 2 }
```

在使用简短而非嵌套的lambda中，建议使用`it`而不是显式地声明参数。在使用参数的嵌套lambda中，参数应该总是显式声明

### 类声明格式
参数比较少的类可以用一行表示：

```Kotlin
class Person(id: Int, name: String)
```

具有较多的参数的类应该格式化成每个构造函数的参数都位于与缩进的单独行中。此外，结束括号应该在新行上。如果我们使用继承，那么超类构造函数调用或实现的接口列表应该位于与括号相同的行中

```Kotlin
class Person(
    id: Int,
    name: String,
    surname: String
) : Human(id, name) {
    // ...
}
```

对于多个接口，应该首先定位超类构造函数调用，然后每个接口应该位于不同的行中

```Kotlin
class Person(
    id: Int,
    name: String,
    surname: String
) : Human(id, name),
    KotlinMaker {
    // ...
}
```

构造函数参数可以使用常规缩进或连续缩进(双倍正常缩进)。

### Unit
如果函数返回 Unit ，返回类型应该省略：

```kotlin
fun foo() { // ": Unit"被省略了
}
```

### 函数 vs 属性
在某些情况下，没有参数的函数可以与只读属性互换。尽管语义是相似的，但是有一些风格上的约定在什么时候更偏向于另一个。

在下面的情况下，更偏向于属性而不是一个函数:
> -- 不需要抛出异常
> -- 复杂度为O(1)
> -- 低消耗的计算(或首次运行结果会被缓存)
> -- 返回与调用相同的结果
