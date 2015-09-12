##使用 Gradle

###插件和版本

kotlin-gradle-plugin 可以编译 Kotlin 文件和模块

> X.Y.SNAPSHOT:对应版本 X.Y 的快照，在 CI 服务器上的每次成功构建的版本。这些版本不是真正的稳定版，只是推荐用来测试新编辑器的功能的。现在所有的构建都是作为 0.1-SNAPSHOT 发表的。你可以参看[configure a snapshot repository in the pom file]()

>X.Y.X: 对应版本 X.Y.Z 的 release 或 milestone ，自动升级。它们是文件构建。Release 版本发布在 Maven Central 仓库。在 pom 文件里不需要多余的配置。

milestone 和 版本的对应关系如下：

**Milestone**|**Version**
---|---
M12.1|0.12.613
M12|0.12.200
M11.1|0.11.91.1
M11|0.11.91
M10.1|0.10.195
M10|0.10.4
M9|0.9.66
M8|0.8.11
M7|0.7.270
M6.2|0.6.1673
M6.1|0.6.602
M6|0.6.69
M5.3|0.5.998

###面向 Jvm

对于 jvm，需要应用 kotlin 插件

> apply plugin: "kotlin"

至于 M11 ，kotlin 文件可以与 java 混用。默认使用不同文件夹：

```
project
    - src
        - main (root)
            - kotlin
            - java
```

如果不使用默认的设置则对应的文件属性要修改：

```gradle
sourceSets {
    main.kotlin.srcDirs += 'src/main/myKotlin'
    main.java.srcDirs += 'src/main/myJava'
}
```

###面向JavaScript

但目标是 JavaScript 时：

> apply plugin: "kotln2js"

这个插件只对 kotlin 文件起作用，因此建议把 kotlin 和 java 文件分开。对于 jvm 如果不用默认的值则需要修改源文件夹：

```
sourceSets {
    main.kotlin.srcDirs += 'src/main/myKotlin'
}
```

如果你想建立一个复用的库，使用 `kotlinOptions.metaInfo` 生成附加的带附加二进制描述的 js 文件

```
compileKotlin2Js {
	kotlinOptions.metaInfo = true
}
```

###目标是 android

Android Gradle 模块与普通的 Gradle 模块有些不同，所以如果你想建立 kotlin 写的android 项目，则需要下面这样：

```
buildscript {
    ...
}
apply plugin: 'com.android.application'
apply plugin: 'kotlin-android'
```

####  Android Studio

如果使用 Android Studio,需要添加下面的代码：

```
android {
  ...

  sourceSets {
    main.java.srcDirs += 'src/main/kotlin'
  }
}
```

这是告诉 android studio kotlin 文件的目录位置方便 IDE 识别

###配置依赖

我们需要添加 kotlin-gradle-plugin 和 kotlin 标准库依赖

```
buildscript {
  repositories {
    mavenCentral()
  }
  dependencies {
    classpath 'org.jetbrains.kotlin:kotlin-gradle-plugin:<version>'
  }
}

apply plugin: "kotlin" // or apply plugin: "kotlin2js" if targeting JavaScript

repositories {
  mavenCentral()
}

dependencies {
  compile 'org.jetbrains.kotlin:kotlin-stdlib:<version>'
}
```

###使用快照版本

如果使用快照版本则如下所示：

```
buildscript {
  repositories {
    mavenCentral()
    maven {
      url 'http://oss.sonatype.org/content/repositories/snapshots'
    }
  }
  dependencies {
    classpath 'org.jetbrains.kotlin:kotlin-gradle-plugin:0.1-SNAPSHOT'
  }
}

apply plugin: "kotlin" // or apply plugin: "kotlin2js" if targeting JavaScript

repositories {
  mavenCentral()
  maven {
    url 'http://oss.sonatype.org/content/repositories/snapshots'
  }
}

dependencies {
  compile 'org.jetbrains.kotlin:kotlin-stdlib:0.1-SNAPSHOT'
}
```


###例子

[Kotlin](https://github.com/jetbrains/kotlin)仓库有如下例子：

>[Kotlin](https://github.com/JetBrains/kotlin-examples/tree/master/gradle/hello-world)
>[Mixed java and Kotlin](https://github.com/JetBrains/kotlin-examples/tree/master/gradle/mixed-java-kotlin-hello-world)
>[Android](https://github.com/JetBrains/kotlin-examples/tree/master/gradle/android-mixed-java-kotlin-project)
>[javaScript](https://github.com/JetBrains/kotlin/tree/master/libraries/tools/kotlin-gradle-plugin/src/test/resources/testProject/kotlin2JsProject)