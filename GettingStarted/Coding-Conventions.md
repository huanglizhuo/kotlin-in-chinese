##代码风格

本页包含了当前 kotlin 语言的代码风格。

###命名风格

如果对 java 默认的代码风格有疑惑，比如下面这些 ：

> --使用骆驼命名法(在命名中避免下划线)

> --类型名称首字母大写

> --方法和属性首字母小写

> --缩进用四个空格

> --public 方法要写说明文档，这样它就可以出现在 Kotllin Doc 中

###冒号

在冒号区分类型和父类型中要有空格，在实例和类型之间是没有空格的：

```kotlin
interface Foo<out T : Any> : Bar {
	fun foo(a: Int): T
}
```

###Lambdas

在 Lambdas 表达式中，大括号与表达式间要有空格，箭头与参数和函数体间要有空格。尽可能的把 lambda 放在括号外面传入

```Kotlin
list.filter { it > 10 }.map { element -> element * 2 }
```

在 lambdas 中建议使用 `it` 而不是申明参数。有嵌套多个参数的 lambdas 中参数是必须明确申明的

###Unit

如果函数返回 Unit ，返回类型应该省略：

```kotlin
fun foo() {
}
```
