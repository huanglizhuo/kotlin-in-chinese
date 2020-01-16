## 使用 Gradle 构建

**多平台项目是 Kotlin 1.2 和 1.3 中的实验性特性。本文档中描述的所有语言和工具功能都可能在将来的Kotlin版本中发生变更**

这篇文档解释了 [Kotlin 多平台项目](https://kotlinlang.org/docs/reference/multiplatform.html) 并且描述了gradle 中是如何配置的。

内容目录

- 项目结构
- 设置多平台项目
- gradle 插件
- 设置 Targets（目标）
    - 支持的平台
    - 配置编译
- 配置源集
    - 链接源集
    - 添加依赖
    - 语言设置
- 默认项目结构
- 运行测试
- 发布多平台库
- JVM 目标的Java 支持
- Android支持
    - 发布 Android 库
- 使用 Kotlin/Native 目标
    - 目标捷径
    - 构建最终native二进制

## 项目结构

Kotlin 多平台项目的结构由下面的构建块组成：

- [targets(目标)](https://kotlinlang.org/docs/reference/building-mpp-with-gradle.html#setting-up-targets) 是构建的一部分，负责构建，测试和打包用于其中一个平台的完整软件。因此，一个多平台项目通常包含多个目标。

- 建立每个目标涉及一次或多次编译 Kotlin 资源。 换句话说，目标可能具有一个或多个[编译](https://kotlinlang.org/docs/reference/building-mpp-with-gradle.html#configuring-compilations)。 例如，一个编译用于生产资源，另一个用于测试。

- Kotlin 资源被安排在[源集(源码集合)](https://kotlinlang.org/docs/reference/building-mpp-with-gradle.html#configuring-source-sets)中。 除了Kotlin源文件和资源外，每个源集可能都有自己的依赖性。 源集根据“依赖”关系构建层次结构。 源集本身与平台无关，但如果仅针对单个平台进行编译，则它可能包含特定于平台的代码和依赖项。

每个编译都有默认的源集，这部分用于设置特定于该编译的源和依赖项。 默认源集还用于通过“依赖”关系将其他源集定向到编译。

这是一个针对JVM和JS的项目的图示：

![mpp-structure-default-jvm-js](./mpp-structure-default-jvm-js.png)

这里 `jvm` 和 `js`, 两个目标分别编译了生产和测试源，并且其中的一些源是共享的。 通过仅创建两个目标即可实现此布局，而无需为编译和源集进行其他配置：这些是为该目标构建的默认设置。

在上面的示例中，JVM目标的生产源通过其main编译进行编译，因此包括来自源集`jvmMain`和 `commonMain` 的源和依赖（由于取决于关系）：

![mpp-one-compilation](./mpp-one-compilation.png)

这里，`jvmMain` 源集为共享 `commonMain` 源中的预期(expected)API提供了平台特定的实现。 这是在平台之间以灵活的方式在需要的地方使用特定于平台的实现方式共享代码的方式。

在后面章节中，将更详细地描述这些概念以及用于在项目中对其进行配置的DSL。

## 设置多平台项目

可以在IDE中创建一个新的多平台项目,在 "Kotlin" 下的“新建项目”对话框中选择一个多平台项目模板.

例如，如果选择 "Kotlin (Multiplatform Library)"，则会创建一个库项目，该项目具有三个目标，一个目标用于JVM，一个目标用于JS，一个目标用于正在使用的Native平台。 这些是通过以下方式在build.gradle脚本中配置的：

```Kotlin

plugins {
    kotlin("multiplatform") version "1.3.61"
}

repositories {
    mavenCentral()
}

kotlin {
    jvm() // Creates a JVM target with the default name 'jvm'
    js()  // JS target named 'js'
    mingwX64("mingw") // Windows (MinGW X64) target named 'mingw'
    
    sourceSets { /* ... */ }
}

```

这三个目标是使用预设函数jvm（），js（）和mingwX64（）创建的,它们各自有一些默认配置。 每个[受支持的平台](https://kotlinlang.org/docs/reference/building-mpp-with-gradle.html#supported-platforms)都有预设。

源集及其依赖通过如下配置：

```Kotlin
plugins { /* ... */ }

kotlin {
    /* Targets declarations omitted */

    sourceSets {
        val commonMain by getting {
            dependencies {
                implementation(kotlin("stdlib-common"))
            }
        }
        val commonTest by getting {
            dependencies {
                implementation(kotlin("test-common"))
                implementation(kotlin("test-annotations-common"))
            }
        }
        
        // Default source set for JVM-specific sources and dependencies:
        jvm().compilations["main"].defaultSourceSet {
            dependencies {
                implementation(kotlin("stdlib-jdk8"))
            }
        }
        // JVM-specific tests and their dependencies:
        jvm().compilations["test"].defaultSourceSet {
            dependencies {
                implementation(kotlin("test-junit"))
            }
        }
        
        js().compilations["main"].defaultSourceSet  { /* ... */ }
        js().compilations["test"].defaultSourceSet { /* ... */ }
        
        mingwX64("mingw").compilations["main"].defaultSourceSet { /* ... */ }
        mingwX64("mingw").compilations["test"].defaultSourceSet { /* ... */ }
    }
}
```

这些是上面配置目标生产和测试源的默认源集合。 源集commonMain和commonTest分别包含在所有目标的生产和测试编译中。 请注意，公共源集commonMain和commonTest的依赖项是公共的，平台相关库应该添加到特定目标的源集。

## Gradle 插件

Kotlin Multiplatform项目需要Gradle 4.7及更高版本，不支持较旧的Gradle版本。

要从零开始在Gradle项目中设置多平台项目，请首先在build.gradle文件的开头添加以下内容，将kotlin-multiplatform插件应用于该项目：

```Kotlin
plugins {
    kotlin("multiplatform") version "1.3.61"
}
```

这将在顶层创建 `kotlin` 扩展。接下来你就可以在构建脚本中使用它：

- 为多平台设置目标(默认是不会创建目标的)
- 配置源集和各自依赖

## 设置目标

目标是构建的一部分，负责编译，测试和打包针对某个受支持平台的软件。

所有目标都可以共享某些资源，也可以具有特定于平台的资源。

由于平台不同，目标也以不同的方式构建，并具有各种特定于平台的设置。 Gradle插件为支持的平台捆绑了许多预设。

要创建目标，请使用预设函数，这些函数根据目标进行命名，并可以选择接受目标名称和配置代码块：

```Kotlin
kotlin {
    jvm() // Create a JVM target with the default name 'jvm'
    js("nodeJs") // Create a JS target with a custom name 'nodeJs'
        
    linuxX64("linux") {
        /* Specify additional settings for the 'linux' target here */
    }
}
```

预设函数会根据目标函数是否存在返回目标,该返回可用于配置已存在目标:

```Kotlin
kotlin {
    /* ... */
    
    // Configure the attributes of the 'jvm6' target:
    jvm("jvm6").attributes { /* ... */ }
}
```

注意这里的目标平台和名字都很重要,如果目标是用 `jvm('jvm6')` 创建,再次使用 `jvm()` 会创建另一个mub(名字默认是 `jvm` ). 如果用于创建该名称下的目标的预设功能不同，则会报错。

从预设创建的目标将添加到kotlin.targets域对象集合中，该集合可用于按其名称访问它们或配置所有目标：


```kotlin
kotlin {
    jvm()
    js("nodeJs")
    
    println(targets.names) // Prints: [jvm, metadata, nodeJs]
    
    // Configure all targets, including those which will be added later:
    targets.all {
        compilations["main"].defaultSourceSet { /* ... */ }
    }
}
```

在多个预设中动态创建或者访问不同目标,可以使用 `targetFromPreset` 函数,该函数接收一个预设(包含在 `kotlin.presets` 域对象集合中),以及可选的,目标名字和一段配置代码块.

比如下面的代码,可以为所有 Kotlin/Native 支持的平台都创建一个目标:

```Kotlin
mport org.jetbrains.kotlin.gradle.plugin.mpp.KotlinNativeTargetPreset

/* ... */

kotlin {
    presets.withType<KotlinNativeTargetPreset>().forEach {
        targetFromPreset(it) { 
            /* Configure each of the created targets */
        }
    }
}
```

### 支持的平台

下面是一些预设目标,这些可以通过预设函数设置:

- `jvm` Kotlin/JVM 
- `js` Kotlin/JS
- `android` 安卓应用和库,注意在该目标创建前要先应用安卓 Gradle 插件.
- Kotlin/Native 预设目标(参看下面的[笔记](https://kotlinlang.org/docs/reference/building-mpp-with-gradle.html#using-kotlinnative-targets))
    - `androidNativeArm32` `androidNativeArm64` 对应安卓 NDK
    - `iosArm32` `iosArm64` `iosX64` 对应 iOS
    - `watchosArm32` `watchosArm64` `watchosX86` 对应 watchOS
    - `tvosArm64` `tvosX64` 对应 tvOS
    - `linuxArm32Hfp`, `linuxMips32`, `linuxMipsel32`, `linuxX64` 对应 linux
    - `macosX64` 对应 macOS
    - `mingwX64` `mingwX86` 对应 Windows
    - `wasm32` 对应 WebAssemnly

    注意这里有些 Kotlin/Native 目标需要[适当的硬件主机](https://kotlinlang.org/docs/reference/building-mpp-with-gradle.html#using-kotlinnative-targets)构建

一些目标需要附加的配置. 关于 Android 和 iOS 例子,可以参看[多平台项目:Android 和 iOS](https://kotlinlang.org/docs/tutorials/native/mpp-ios-android.html)教程

### 配置编译

构建目标需要一次或多次编译Kotlin。 目标的每个Kotlin编译都可以用于不同的目的（例如，生产代码，测试），并包含不同的源集。 可以在DSL中可以访问目标的编译，例如，获取任务，配置Kotlin编译器选项或获取依赖文件和编译输出：

```Kotlin
kotlin {
    jvm {
        val main by compilations.getting {
            kotlinOptions { 
                // Setup the Kotlin compiler options for the 'main' compilation:
                jvmTarget = "1.8"
            }
        
            compileKotlinTask // get the Kotlin task 'compileKotlinJvm' 
            output // get the main compilation output
        }
        
        compilations["test"].runtimeDependencyFiles // get the test runtime classpath
    }
    
    // Configure all compilations of all targets:
    targets.all {
        compilations.all {
            kotlinOptions {
                allWarningsAsErrors = true
            }
        }
    }
}
```

每个编译都附带一个默认源集，该默认源集存储特定于该编译的源和依赖项。 目标 `bar` 的编译 `foo` 的默认源集名字是 `barFoo` . 可以使用 `defaultSourceSet` 从编译中访问它：

```Kotlin
kotlin {
    jvm() // Create a JVM target with the default name 'jvm'
    
    sourceSets {
        // The default source set for the 'main` compilation of the 'jvm' target:
        val jvmMain by getting {
            /* ... */
        }
    }
    
    // Alternatively, access it from the target's compilation:
    jvm().compilations["main"].defaultSourceSet { 
        /* ... */
    }
}
```

要收集所有参与编译的源集，包括通过依赖关系添加的源集，可以使用属性 `allKotlinSourceSets` 。

对于某些特定用例，可能需要创建自定义编译。 这可以在目标的 `compilations` 域对象集合中完成。 请注意，需要为所有定制编译手动设置依赖项，并且定制编译输出的使用取决于构建作者。 例如，下面可以针对jvm（）目标的集成测试的定制编译：

```Kotlin
kotlin {
    jvm() {
        compilations {
            val main by getting

            val integrationTest by compilations.creating {
                defaultSourceSet {
                    dependencies {
                        // Compile against the main compilation's compile classpath and outputs:
                        implementation(main.compileDependencyFiles + main.output.classesDirs)
                        implementation(kotlin("test-junit"))
                        /* ... */
                    }
                }

                // Create a test task to run the tests produced by this compilation:
                tasks.create<Test>("integrationTest") {
                    // Run the tests with the classpath containing the compile dependencies (including 'main'),
                    // runtime dependencies, and the outputs of this compilation:
                    classpath = compileDependencyFiles + runtimeDependencyFiles + output.allOutputs

                    // Run only the tests from this compilation's outputs:
                    testClassesDirs = output.classesDirs
                }
            }
        }
    }
}
```

还要注意，默认情况下，自定义编译的默认源集既不依赖于commonMain也不依赖于commonTest。

## 配置源集

Kotlin 源集是 Kotlin源及其资源，依赖关系和语言设置的集合，这些源可能会参与一个或多个目标的Kotlin编译。

源集不限于平台特定的或“共享的”； 允许包含的内容取决于其用法：添加到多个编译中的源集仅限于通用语言功能和依赖项，而仅由单个目标使用的源集可以具有特定于平台的依赖项，并且其代码可以使用特定于目标平台的语言功能。

默认情况下会创建和配置一些源集：commonMain，commonTest和编译的默认源集。 请参阅默认项目结构。

源集在kotlin {...}扩展的sourceSets {...}块内配置：


```Kotlin
kotlin { 
    sourceSets { 
        val foo by creating { /* ... */ } // create a new source set by the name 'foo'
        val bar by getting { /* ... */ } // configure an existing source set by the name 'bar' 
    }
}
```

**注意：创建源集不会将其链接到任何目标。 一些源集是预定义的，因此会默认进行编译。 但是，自定义的源集必须明确指向编译,请参阅：[链接源集](https://kotlinlang.org/docs/reference/building-mpp-with-gradle.html#connecting-source-sets)。**

源集名称是大小写敏感的. 当通过名称引用默认源集时,确定前缀和目标名称匹配,比如源集 `iosX64Main` 对应 `iosX64`目标.

源集本身是平台无关的,但如果它只编译到某个平台则可以看做平台相关的. 源集可以既包含平台共享公共代码,也可以包含平台相关代码.

每个源集对于 Kotlin 源都有默认的源码目录: `src/<source set name>/kotlin` . 给源集添加 Kotlin 源码和资源,可以通过 `kotlin` `resource` 的 `SourceDirectorySet`:

默认源集的文件存储在如下目录:

- 源代码文件: `src/<source set name>/kotlin`
- 资源文件: `src/<source set name>/resources`

你需要手动创建这些目录

添加自定义 Kotlin 源码目录和资源目录可以通过以下方式:

```Kotlin
kotlin { 
    sourceSets { 
        val commonMain by getting {
            kotlin.srcDir("src")
            resources.srcDir("res")
        } 
    }
}
```

### 链接源集

Kotlin 源集可以通过依赖关系链接,如果源集 `foo` 依赖 源集`bar` 则会如下:

- 不论 `foo` 编译到任何平台, `bar` 都会参与到编译中,并会编译为目标二进制形式,比如 JVM class 文件或者 JS 代码
- 源集 `foo` 可以'看到' `bar` 的声明,也包括 `internal`声明,以及 `bar` 的依赖,尽管这些是通过 `implementation` 指定的依赖
- `foo` 可能包含针对 `bar` 的预期(expected)声明在特定于平台上的实现
- 源集 `bar` 总是与 `foo` 的资源一起处理和复制；
- `foo` `bar` 的语言设置应该是一致的

不允许源集循环依赖

源集的 DSL 可以定义源集间的链接:

```Kotlin
kotlin { 
    sourceSets { 
        val commonMain by getting { /* ... */ }
        val allJvm by creating {
            dependsOn(commonMain)
            /* ... */
        } 
    }
}
```

除了默认源集之外，还应将创建的自定义源集显式包含到依赖关系层次结构中，以便能够使用来自其他源集的声明，并且最重要的是可以参与编译。 大多数情况下，它们需要一个dependsOn（commonMain）或dependsOn（commonTest）语句，并且某些特定于平台的默认源集应直接或间接依赖于自定义源集：

```Kotlin
kotlin { 
    mingwX64()
    linuxX64()
    
    sourceSets {
        // custom source set with tests for the two targets
        val desktopTest by creating { 
            dependsOn(getByName("commonTest"))
            /* ... */
        }
        // Make the 'windows' default test source set for depend on 'desktopTest'
        mingwX64().compilations["test"].defaultSourceSet { 
            dependsOn(desktopTest)
            /* ... */
        }
        // And do the same for the other target:
        linuxX64().compilations["test"].defaultSourceSet {
            dependsOn(desktopTest)
            /* ... */
        }
    }
}
```

### 添加依赖

要将依赖项添加到源集，请使用源集DSL的 `dependencies { ... } ` 块。支持四种依赖项：

- `api` 依赖项在编译期间和运行时均会使用，并会导出到库使用者。如果当前模块的公共API中使用了依赖关系中的任何类型，则它应该是api依赖关系；

- `implementation` 依赖在当前模块的编译期间和运行时使用，但对其他模块通过 `implementation` 依赖当前模块时，不会导出该依赖。`implementation` 依赖关系类型应用于模块内部逻辑所需的依赖关系。如果模块是未发布的终结点应用程序，则它应该使用 `implementation` 而不是 `api` 依赖。

- `compileOnly` 依赖项仅用于编译当前模块，并且在运行时或其他模块的编译期间均不可用。这些依赖项应用于在运行时具有第三方实现的API。

- `runtimeOnlyOnly` 仅运行时可用，但在任何模块的编译过程中都不可见。

每个源集都指定了依赖种类，如下所示：

```Kotlin
kotlin {
    sourceSets {
        val commonMain by getting {
            dependencies {
                api("com.example:foo-metadata:1.0")
            }
        }
        val jvm6Main by getting {
            dependencies {
                api("com.example:foo-jvm6:1.0")
            }
        }
    }
}
```

注意，为了使IDE能够正确分析公共源的依赖关系，除了平台特定源集需要声明与平台特定组件依赖之外，公共源集还必须具有与 Kotlin 元数据包相对应的依赖关系。 通常，在使用已发布的库时（除非它与Gradle元数据一起发布，如下所述），需要后缀为-common（如kotlin-stdlib-common）或-metadata的组件。

然而 project（'...'）依赖于另一个多平台项目时会自动解析为适当的目标。 在源集的依赖项中指定单个project（'...'）依赖项就足够了，并且包含源集的编译将收到该项目的对应平台特定产物，只要它具有兼容的目标即可：

```Kotlin
kotlin {
    sourceSets {
        val commonMain by getting {
            dependencies {
                // All of the compilations that include source set 'commonMain'
                // will get this dependency resolved to a compatible target, if any:
                api(project(":foo-lib"))
            }
        }
    }
}
```

同样，如果以实验性 Gradle 元数据发布模式发布了多平台库，并且该项目也设置为使用元数据，那么只需为公共源集指定一次依赖项就足够了。 否则，除了公共模块之外，每个平台特定的源集还应提供库的相应平台模块，如上所示。

指定依赖关系的另一种方法是在顶层使用Gradle内置DSL，其配置名称遵循模式<sourceSetName> <DependencyKind>：

```Kotlin
dependencies {
    "commonMainApi"("com.example:foo-common:1.0")
    "jvm6MainApi"("com.example:foo-jvm6:1.0")
}
```

源集合依赖项DSL中不提供某些Gradle内置依赖项，例如 `gradleApi()` ,`localGroovy()`或`gradleTestKit()`。 但你可以将它们添加到顶级依赖项块中，如上所示。

可以使用 kotlin("stdlib") 添加像 `kotlin-stdlib` 或 `kotlin-reflect` 之类的Kotlin 依赖模块，这是 `org.jetbrains.kotlin：kotlin-stdlib` 的缩写。

### 语言设定

可以如下指定源集的语言设置：

```Kotlin
kotlin {
    sourceSets {
        val commonMain by getting {
            languageSettings.apply {
                languageVersion = "1.3" // possible values: '1.0', '1.1', '1.2', '1.3'
                apiVersion = "1.3" // possible values: '1.0', '1.1', '1.2', '1.3'
                enableLanguageFeature("InlineClasses") // language feature name
                useExperimentalAnnotation("kotlin.ExperimentalUnsignedTypes") // annotation FQ-name
                progressiveMode = true // false by default
            }
        }
    }
}
```

也可以为所有源集配置语言:

```Kotlin
kotlin.sourceSets.all {
    languageSettings.progressiveMode = true
}
```

源集的语言设置会影响在IDE中分析来源的方式。 由于当前的限制，在Gradle构建中，仅使用编译的默认源集的语言设置并将其应用于参与编译的所有源。

检查语言设置是否相互依赖，以确保源集之间的一致性。 即，如果 `foo` 依赖 `bar`：

- foo应该将languageVersion设置为大于或等于bar的语言；
- foo应该启用bar启用的所有不稳定的语言功能（错误修正功能没有这种要求）；
- foo应该使用bar使用的所有实验性注释；
- 可以任意设置apiVersion，错误修正语言功能和ProgressiveMode。

## 默认项目结构

默认情况下，每个项目都包含两个源集，`commonMain` 和 `commonTest`，在其中可以放置应在所有目标平台之间共享的所有代码。这些源集分别添加到各自生产和测试编译中。

添加目标后，将为其创建默认编译：

- 为 JAM,JS,以及原生目标创建 `main` 和 `test` 编译
- 针对Android目标的每个[Android构建变体](https://developer.android.com/studio/build/build-variants)的编译；

对于每个编译，在由 `<targetName> <CompilationName>` 组成的名称下都有一个默认源集。此默认源集参与了编译，因此应将其用于特定于平台的代码和依赖项，并通过“依赖于”的方式将其他源集添加到编译中。例如，目标为jvm6（JVM）和nodeJs（JS）的项目将具有源集：commonMain，commonTest，jvm6Main，jvm6Test，nodeJsMain，nodeJsTest。

默认用例集涵盖了绝大多数用例，不需要自定义用例集。

默认情况下，每个源集在 `src/<sourceSetName>/kotlin` 目录是Kotlin源码，在 `src/<sourceSetName>/resources下` 有资源文件。

在Android项目中，会为不同 [Android 源集](https://developer.android.com/studio/build/#sourcesets)创建对应的 Kotlin 源集。如果Android目标的名称为foo，则Android源集 bar 将会有 与Kotlin源集合对应的fooBar。但是，Kotlin编译能够从所有目录 `src/bar/java`，`src/bar/kotlin` 以及 ``src/foobar/kotlin`中使用Kotlin源。 Java源仅从这些目录中的第一个读取。

## 运行测试

JVM，Android，Linux，Windows和macOS当前默认支持在Gradle构建中运行测试。 JS和其他Kotlin/Native 目标需要手动配置以在适当的环境，模拟器或测试框架下运行测试。

每个适合测试的目标以名称 `<targetName>Test` 创建一个测试任务。`Check` 会运行所有目标的测试。

在将 `commonTest` 默认源集添加到所有测试编译后，所有目标平台上所需的测试和测试工具都可以添加在此处。

[kotlin.test API](https://kotlinlang.org/api/latest/kotlin.test/index.html) 可用于多平台测试。将 `kotlin-test-common` 和 `kotlin-test-annotations-common` 依赖项添加到 `commonTest` 以使用诸如 `kotlin.test.assertTrue（...）` 断言函数,以及 `@Test` / `@Ignore` / `@BeforeTest` / `@AfterTest` 等注解。

对于JVM目标，可以使用 `kotlin-test-junit` 或 `kotlin-test-testng` 用于相应的断言器实现和注解映射。

对于 Kotlin/JS 目标，添加 `kotlin-test-js` 作为测试依赖项。至此，创建了 Kotlin/JS 的测试任务，但默认情况下不运行测试。应该手动配置它们以使用JavaScript测试框架运行测试。

Kotlin/Native 目标不需要其他测试依赖项，内置了 `kotlin.test` API实现。

## 发布多平台库

