## 包
代码文件以包声明开始：

```kotlin
package foo.bar

fun bza() {}

class Goo {}

//...
```

代码文件的所有内容(比如类和函数)都被包含在包声明中。因此在上面的例子中， `bza() ` 的全名应该是 `foo.bar.bza` ，`Goo` 的全名是 `foo.bar.Goo`。

如果没有指定包名，那这个文件的内容就从属于没有名字的 "default" 包。

### 默认导入
许多包被默认导入到每个Kotlin文件中：

> -- kotlin.*  
>
> -- kotlin.annotation.*  
>
> -- kotlin.collections.*  
>
> -- kotlin.comparisons.* (since 1.1)  
>
> -- kotlin.io.*  
>
> -- kotlin.ranges.*  
>
> -- kotlin.sequences.*  
>
> -- kotlin.text.*

一些附加包会根据平台来决定是否默认导入：

> -- JVM:  
>
> ---- java.lang.*  
>
> ---- kotlin.jvm.*  

> -- JS:  
>
> ---- kotlin.js.*

### Imports
除了模块中默认导入的包，每个文件都可以有导入自己需要的包。导入语法可以在 [grammar](Reference/Grammar.md) 查看。

我们可以导入一个单独的名字，比如下面这样：

```kotlin
import foo.Bar // Bar 现在可以直接使用了
```

或者范围内的所有可用的内容 (包，类，对象，等等):

```kotlin
import foo.*/ /foo 中的所有都可以使用
```

如果命名有冲突，我们可以使用 `as` 关键字局部重命名解决冲突

```kotlin
import foo.Bar // Bar 可以使用
import bar.Bar as bBar // bBar 代表 'bar.Bar'
```

import关键字不局限于导入类;您也可以使用它来导入其他声明:

>-- 顶级函数与属性  
>
>-- 在[对象声明](http://kotlinlang.org/docs/reference/object-declarations.html#object-declarations)中声明的函数和属性  
>
>-- [枚举常量](http://kotlinlang.org/docs/reference/enum-classes.html)

与 Java 不同的是，Koting 没有[静态导入](https://docs.oracle.com/javase/8/docs/technotes/guides/language/static-import.html)的语法，所有的导入都是通过`import`关键字声明的。

### 顶级声明的可见性
如果最顶的声明标注为 private , 那么它是声明文件私有的 (参看[ Visibility Modifiers](http://kotlinlang.org/docs/reference/visibility-modifiers.html))。
