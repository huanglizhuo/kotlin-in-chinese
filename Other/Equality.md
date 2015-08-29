##相等

在 kotlin 中有俩中相等：

>参照相等(指向相同的对象)
>结构相等

###参照相等

参照相等是通过 `===` 操作符判断的(不等是`!==` ) a===b 只有 a b 指向同一个对象是判别才成立。


另外，你可以使用内联函数 `identityEquals()` 判断参照相等：

```kotlin
a.identityEquals(b)
a identityEquals b
```

###结构相等

结构相等是通过 `==` 判断的。像 `a == b` 将会翻译成：

```kotlin
a?.equals(b) ?: b === null
```

如果 a 不是 null 则调用 `equals(Any?)` 函数，否则检查 b 是否参照等于 null

注意完全没有必要为优化你的代码而将 `a == null` 写成 `a === null` 编译器会自动帮你做的。

