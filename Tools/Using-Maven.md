## 使用Maven
### 插件和版本
Kotlin-maven-plugin 可以编译 Kotlin 资源和模块。现在只有 Maven V3 支持

通过 Kotlin.version 定义你想要的 Kotlin 版本。可以有以下的值

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

### 配置快照仓库
使用 kotlin 版本的快照，需要在 pom 中这样定义：

```pom
<repositories>
	<repository>
		<id>sonatype.oss.snapshots</id>
		<name>Sonatype OSS Snapshot Repository</name>
		<url>http://oss.sonatype.org/content/repositories/snapshots</url>
		<releases>
			<enabled>false</enabled>
		</releases>
		<snapshots>
			<enabled>true</enabled>
		</snapshots>
	</repository>
</repositories>

<pluginRepositories>
	<pluginRepository>
		<id>sonatype.oss.snapshots</id>
		<name>Sonatype OSS Snapshot Repository</name>
		<url>http://oss.sonatype.org/content/repositories/snapshots</url>
	<releases>
		<enabled>false</enabled>
	</releases>
	<snapshots>
		<enabled>true</enabled>
	</snapshots>
	</pluginRepository>
</pluginRepositories>
```

### 依赖
kotlin 有一些扩展标准库可以使用。在 pom 文件中使用如下的依赖：

```pom
<dependencies>
	<dependency>
		<groupId>org.jetbrains.kotlin</groupId>
		<artifactId>kotlin-stdlib</artifactId>
		<version>${kotlin.version}</version>
	</dependency>
</dependencies>
```

### 只编译 kotlin 源码
编译源码需要在源码文件夹打一个标签：

```xml
<sourceDirectory>${project.basedir}/src/main/kotlin</sourceDirectory>
<testSourceDirectory>${project.basedir}/src/test/kotlin</testSourceDirectory>
```

在编译资源是需要引用kotlin Maven Plugin:

```xml
<plugin>
	<artifactId>kotlin-maven-plugin</artifactId>
	<groupId>org.jetbrains.kotlin</groupId>
	<version>${kotlin.version}</version>
	<executions>
		<execution>
			<id>compile</id>
			<phase>compile</phase>
			<goals> <goal>compile</goal> </goals>
		</execution>
		<execution>
			<id>test-compile</id>
			<phase>test-compile</phase>
			<goals> <goal>test-compile</goal> </goals>
		</execution>
	</executions>
</plugin>
```

### 编译 kotlin 和 java 资源
为了编译混合代码的应用，Kotlin 编译器应该在 java 编译器之前先工作。在 maven 中意味着 kotlin-maven-plug 应该在 maven-compiler-plugin 之前。

```xml
<plugin>
	<artifactId>kotlin-maven-plugin</artifactId>
	<groupId>org.jetbrains.kotlin</groupId>
	<version>0.1-SNAPSHOT</version>
	<executions>
		<execution>
			<id>compile</id>
			<phase>process-sources</phase>
			<goals> <goal>compile</goal> </goals>
		</execution>
		<execution>
			<id>test-compile</id>
			<phase>process-test-sources</phase>
			<goals> <goal>test-compile</goal> </goals>
		</execution>
	</executions>
</plugin>
```

### 使用扩展的注解
kotlin 使用扩展的注解解析 java 库的信息。为了明确这些注解，你需要像下面这样：

```xml
<plugin>
	<artifactId>kotlin-maven-plugin</artifactId>
	<groupId>org.jetbrains.kotlin</groupId>
	<version>0.1-SNAPSHOT</version>
	<configuration>
		<annotationPaths>
			<annotationPath>path to annotations root</annotationPath>
		</annotationPaths>
	</configuration>
```

### 例子
你可以在 [Github]() 仓库参考
