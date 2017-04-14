## 类型检查和转换
## is !is 表达式
我们可以在运行是通过上面俩个操作符检查一个对象是否是某个特定类：

```kotlin
if (obj is String) {
	print(obj.length)
}

if (obj !is String) { // same as !(obj is String)
	print("Not a String")
}
else {
	print(obj.length)
}
```

### 智能转换
在很多情形中，需要使用非明确的类型，因为编译器会跟踪 `is` 检查静态变量，并在需要的时候自动插入安全转换：

```kotlin
fun demo(x: Any) {
	if (x is String) {
		print(x.length) // x is automatically cast to String
	}
}
```

编译器足够智能如何转换是安全的，如果不安全将会返回：

```kotlin
if (x !is String) return
print(x.length) //x 自动转换为 String
```

或者在 `||` `&&` 操作符的右边的值

```kotlin
 // x is automatically cast to string on the right-hand side of `||`
  if (x !is String || x.length == 0) return

  // x is automatically cast to string on the right-hand side of `&&`
  if (x is String && x.length > 0)
      print(x.length) // x is automatically cast to String
```

这样的转换在 when 表达式和 whie 循环中也会发生

```kotlin
when (x) {
	is Int -> print(x + 1)
	is String -> print(x.length + 1)
	is Array<Int> -> print(x.sum())
}
```

### “不安全”的转换符和
如果转换是不被允许的那么转换符就会抛出一个异常。因此我们称之为不安全的。在kotlin 中　我们用前缀 as 操作符

```kotlin
val x: String = y as String
```

注意 null 不能被转换为 `String` 因为它不是 [`nullable`](http://kotlinlang.org/docs/reference/null-safety.html)，也就是说如果 `y` 是空的，则上面的代码会抛出空异常。

为了 java 的转换语句匹配我们得像下面这样：

```kotlin
val x: String?= y as String?
```

###  "安全"转换符
为了避免抛出异常，可以用 as? 这个安全转换符，这样失败就会返回 null　：

```kotlin
val x: String ?= y as? String
```

不管 as? 右边的是不是一个非空 `String` 结果都会转换为可空的。
