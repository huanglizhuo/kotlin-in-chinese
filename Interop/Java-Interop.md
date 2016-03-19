##交互

### java 交互

Kotlin 在设计时就是以与 java 交互为中心的。现存的 java 代码可以在 kotlin 中使用，并且 Kotlin 代码也可以在 java 中流畅运行。这节我们会讨论在 kotlin 中调用 java 代码。

####在 kotlin 中调用 java 代码

基本所有的 Java 代码都可以运行

```kotlin
import java.util.*
fun demo(source: List<Int>) {
	val list = ArrayList<Int>()
	for (item in source )
		list.add(item)
	for (i in 0..source.size() - 1)
		list[i] = source[i]
}
```

**空的返回**

如果 Java 方法返回空，则在 Kotlin 调用中返回 `Unit`。如果
