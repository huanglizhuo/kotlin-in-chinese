## 使用 Ant
### 获得 Ant 任务
Kotlin 提供了 Ant 三个任务:

> kotlinc : Kotlin 面向 JVM 的编译器

> kotlin2js: 面向 javaScript 的编译器

> withKotlin: 使用标准 javac Ant 任务时编译 Kotlin 文件的任务

这些任务定义学在 kotlin-ant.jar 标准库中，位与 [kotlin compiler](https://github.com/JetBrains/kotlin/releases/tag/build-0.12.613) 的 lib 文件夹下

### 面向 JVM 的只有 kotlin 源文件任务
当项目只有 kotlin 源文件时，最简单的方法就是使用 kotlinc 任务：

```ant
<project name="Ant Task Test" default="build">
    <typedef resource="org/jetbrains/kotlin/ant/antlib.xml" classpath="${kotlin.lib}/kotlin-ant.jar"/>

    <target name="build">
        <kotlinc src="hello.kt" output="hello.jar"/>
    </target>
</project>
```

${kotlin.lib} 指向 kotlin 单独编译器解压的文件夹

### 面向 JVM 的只有 kotlin 源文件但有多个根的任务
如果一个项目包含多个根源文件，使用 src 定义路径：

```ant
<project name="Ant Task Test" default="build">
    <typedef resource="org/jetbrains/kotlin/ant/antlib.xml" classpath="${kotlin.lib}/kotlin-ant.jar"/>

    <target name="build">
        <kotlinc output="hello.jar">
            <src path="root1"/>
            <src path="root2"/>
        </kotlinc>
    </target>
</project>
```

### 面向 JVM 的有 kotlin 和 java 源文件
如果项目包含 java kotlin 代码，使用 kotlinc 是可以的，但建议使用 withKotlin 任务

```ant
<project name="Ant Task Test" default="build">
    <typedef resource="org/jetbrains/kotlin/ant/antlib.xml" classpath="${kotlin.lib}/kotlin-ant.jar"/>

    <target name="build">
        <delete dir="classes" failonerror="false"/>
        <mkdir dir="classes"/>
        <javac destdir="classes" includeAntRuntime="false" srcdir="src">
            <withKotlin/>
        </javac>
        <jar destfile="hello.jar">
            <fileset dir="classes"/>
        </jar>
    </target>
</project>
```

### 面向 JavaScript 的只有一个源文件夹的
```ant
<project name="Ant Task Test" default="build">
    <typedef resource="org/jetbrains/kotlin/ant/antlib.xml" classpath="${kotlin.lib}/kotlin-ant.jar"/>

    <target name="build">
        <kotlin2js src="root1" output="out.js"/>
    </target>
</project>
```

### 面向 JavaScript 有前缀，后缀以及 sourcemap 选项
```ant
<project name="Ant Task Test" default="build">
    <taskdef resource="org/jetbrains/kotlin/ant/antlib.xml" classpath="${kotlin.lib}/kotlin-ant.jar"/>

    <target name="build">
        <kotlin2js src="root1" output="out.js" outputPrefix="prefix" outputPostfix="postfix" sourcemap="true"/>
    </target>
</project>
```

#### ##面向 JavaScript 只有一个源码文件夹并有元信息的选项
如果你想要描述 javaScript/Kotlin 库的转换结果，`mateInfo` 选项是很有用的。如果`mateInfo` 设置为 true 则编译附加 javaScript 文件时会创建二进制的元数据。这个文件会与转换结果一起发布

```ant
<project name="Ant Task Test" default="build">
    <typedef resource="org/jetbrains/kotlin/ant/antlib.xml" classpath="${kotlin.lib}/kotlin-ant.jar"/>

    <target name="build">
        <!-- out.meta.js will be created, which contains binary descriptors -->
        <kotlin2js src="root1" output="out.js" metaInfo="true"/>
    </target>
</project>
```

## 参考
下面是所有的元素和属性

#### #kotlinc 属性**名字**|**描述**|**必须性**|**默认值**
---|---|---|---|
src|要编译的Kotlin 文件或者文件夹|yes|
output|目标文件夹或 .jar 文件名 |yes|
classpath|类的完整路径|no|
classpathref|类的完整路径参考|no|
stdlib|"Kotlin-runtime.jar" 的完整路径|no|”“
includeRuntime|如果输出是 .jar 文件，是否 kotlin 运行时库是否包括在 jar 中|no|true

#### #withKotlin 属性
**名字**|**描述**|**必须性**|**默认值**
---|---|---|---|
src|要编译的Kotlin 文件或者文件夹|yes|
output|目标文件夹 |yes|
library|库文件(kt,dir,jar) |no|
outputPrefix|生成 javaScript 文件的前缀|no|
outputSufix|生成 javaScript 文件的后缀|no|
sourcemap|是否生成 sourcemap |no|
metaInfo |是否生成二进制元数据文件描述 |no|
main |是否生成调用主函数 |no|
