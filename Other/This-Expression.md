## This 表达式
为了记录下当前接受者，我们使用 this 表达式：

> 在类的成员中，this 表示当前类的对象

> 在[扩展函数](http://kotlinlang.org/docs/reference/extensions.html)或[扩展字面函数](http://kotlinlang.org/docs/reference/lambdas.html#function-literals)中，this 表示 . 左边接收者参数

如果 this 没有应用者，则指向的是最内层的闭合范围。为了在其它范围中返回 this ，需要使用标签

### this使用范围
为了在范围外部(一个类，或者表达式函数，或者带标签的扩展字面函数)访问 this ，我们需要在使用 `this@lable` 作为 lable

```kotlin
class A { // implicit label @A
  inner class B { // implicit label @B
    fun Int.foo() { // implicit label @foo
      val a = this@A // A's this
      val b = this@B // B's this

      val c = this // foo()'s receiver, an Int
      val c1 = this@foo // foo()'s receiver, an Int

      val funLit = @lambda {String.() ->
        val d = this // funLit's receiver
        val d1 = this@lambda // funLit's receiver
      }


      val funLit2 = { (s: String) ->
        // foo()'s receiver, since enclosing function literal 
        // doesn't have any receiver
        val d1 = this 
      }
    }
  }
}
```
