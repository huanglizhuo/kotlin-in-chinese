## 包
一个源文件以包声明开始：

```kotlin
package foo.bar

fun bza() {}

class Goo {}

//...
```

源文件的所有内容(比如类和函数)都被包声明包括。因此在上面的例子中， `bza() ` 的全名应该是 `foo.bar.bza` ，`Goo` 的全名是 `foo.bar.Goo`。

如果没有指定包名，那这个文件的内容就从属于没有名字的 "default" 包。

### 默认导入
许多包被默认导入到每个Kotlin文件中：

> -- kotlin.*  
> -- kotlin.annotation.*  
> -- kotlin.collections.*  
> -- kotlin.comparisons.* (since 1.1)  
> -- kotlin.io.*  
> -- kotlin.ranges.*  
> -- kotlin.sequences.*  
> -- kotlin.text.*

一些增强包会根据平台来决定是否默认导入：

> -- JVM:  
> ---- java.lang.*  
> ---- kotlin.jvm.*  

> -- JS:  
> ---- kotlin.js.*

### Imports
除了模块中默认导入的包，每个文件都可以有它自己的导入指令。导入语法的声明在[grammar](http://kotlinlang.org/docs/reference/grammar.html#imports)中描述。

我们可以导入一个单独的名字，比如下面这样：

```kotlin
import foo.Bar // Bar 现在可以不用条件就可以使用
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
>-- 在[对象声明](http://kotlinlang.org/docs/reference/object-declarations.html#object-declarations)中声明的函数和属性  
>-- [枚举常量](http://kotlinlang.org/docs/reference/enum-classes.html)

### 可见性和包嵌套
如果最顶的声明标注为 private , 那么它是自己对应包私有 (参看[ Visibility Modifiers](http://kotlinlang.org/docs/reference/visibility-modifiers.html))。如果包内有私有的属性或方法，那它对所有的子包是可见的。

注意包外的的成员是默认不导入的，比如在导入 `foo.bar` 后我们不能获得 `foo` 的成员
