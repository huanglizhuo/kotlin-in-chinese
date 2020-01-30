当值可以具有受限集中的一种类型但不能具有任何其他类型时，密封类用于表示受限类层次结构。从某种意义上讲，它们是枚举类的扩展：枚举类型的值集也受到限制，但是每个枚举常量仅作为单个实例存在，而密封类的子类可以具有多个实例，这些实例可以包含状态。

要声明一个密封类，可以将 `sealed` 修饰符放在该类的名称之前。密封类可以具有子类，但是所有子类必须与密封类本身在同一文件中声明。 （在Kotlin 1.1之前，规则更加严格：类必须嵌套在密封类的声明中）。

```Kotlin
sealed class Expr
data class Const(val number: Double) : Expr()
data class Sum(val e1: Expr, val e2: Expr) : Expr()
object NotANumber : Expr()
```

（上面的示例使用了Kotlin 1.1的另一项新特性：数据类可以扩展其他类，包括密封类）。

密封类本身是抽象的，不能直接实例化，并且可以具有抽象成员。

密封的类不允许具有非私有的构造函数（默认情况下，它们的构造函数是私有的）。

请注意，继承密封类的子类的类（间接继承程序）可以放置在任何位置，而不必放在同一文件中。

当在when表达式中使用密封类时，使用密封类的主要好处就发挥了作用。如果可以验证该语句是否涵盖所有情况，则无需在该语句中添加 else 子句。但仅当将 when 用作表达式（使用结果）而不用作语句时，此方法才有效。

```Kotlin
fun eval(expr: Expr): Double = when(expr) {
    is Const -> expr.number
    is Sum -> eval(expr.e1) + eval(expr.e2)
    NotANumber -> Double.NaN
    // the `else` clause is not required because we've covered all the cases
}
```