# Zxx

本文采用 ABNF 描述文法, 但不完整, 完整的语法描述参见 [Zxx Definition of ABNFA][].

## 模块

模块是一个路径下多个 Zxx 源文件集合.

模块通过 `use` 引入, 不能被引入的不是模块, 比如编译后的可执行文件或动态链接库.

模块路径命名规则:

```abnf
module-path= name-lite
  (
      1*("." name-lite) 1*("/" name-lite)
    /  *("/" name-lite)
  )

name-lite = ALPHA *(["-"] 1*(ALPHA / DIGIT))
DIGIT     = %x30-39 ; 0-9
ALPHA     = %x41-5A / %x61-7A ; A-Z / a-z
```

最后一段 `name-lite` 作为引入模块的访问名.

路径嵌套是代码组织形式, 不表示依赖关系.

模块的组成:

- use.eson
- 其它文件以及文件夹
- Zxx 源文件
  - 编译指令
    - Fixed
  - 顶级声明
    - 语句
      - 保留字
      - 标识符
      - 分隔符
      - 定界符
      - 注释
    - 注释
  - 注释

### use.eson

模块路径下的 `use.eson` 文件描述模块信息.

    eson 表示 Easy Object Notation, 采用 Zxx 字面值表示法

use.eson 是一个键值为字符串的映射对象.

必有键值:

- `name`        string 合法的 zxx 标识符表示模块名
- `version`     string 最新正式版本号, x.y.z 格式, 参见 [semver.org][]

可选键值:

- `path`        string   约束 `use` 模块路径, 对应拉取到本地的相对路径
- `license`     string   开源许可证 Identifier, 参见 [SPDX License List][]
- `description` string   简短说明
- `homepage`    string   相关主页 URL
- `keywords`    [string] 标签或关键字
- `repository`  [string] 按顺序优先的多个代码托管仓库地址, 必须包含 version 对应的 tag
- `author`      {string} 作者信息

限制:

- 每个字符串不超过 140 个 Unicode 字符
- 版本号 0 <= x <= 99, 0 <= y <= 9999, 0 <= z <= 999

Zxx 自己的版本号规划:

- x 是主版本号 从 0 开始的开发版
- y 是表示两位数的年份+月份, 开始的或发布的
- z 前两位表示日, 第三位从 0 开始
- 第一个版本是 0.1602.020, 表示 开发版 20160202.0

例:

```eson
{
  'name' 'zxx',
  'license' 'BSD-2-Clause',
  'version' '0.1602.020",
  'repository' [
    'git+https://github.com/ZxxLang/zxx.git'
  ],
  'keywords' [
    'Zxx', 'programming', 'language'
  ],
  'author' {
    "name' 'YU HengChun',
    'url' 'https://achun.github.io/''
  }
}
```

## 源文件

源文件采用 UTF-8 编码. 文件命名规则:

```abnf
source-zxx= category *(["-"] (a-z / DIGIT)) *("_" platform) "." (%s"zxx" / %s"md")
category  = 1*a-z
platform  = (a-z / DIGIT) *(["-"] (a-z / DIGIT))
a-z = %x61-7A
```

category:

1. main  拥有 main 函数, 编译生成可执行文件
1. case  用例, 每个都拥有 main 函数, 编译生成可执行文件
1. fail  单元测试, 失败测试
1. fine  单元测试, 成功测试
1. race  性能测试
1. 其它  作为常规模块被 `use` 引入

下划线 `_` 开头的 `platform` 则表示用于特定平台的代码, 顺序无关. 比如:

    linux_darwin_freebsd_64_x86.zxx

上例表示平台为 x86 架构下的 64 位 linux, darwin, freebsd 操作系统.

显然每个 `platform` 值都有特定含义, 需要编译器能够识别和支持.

Zxx 的语法兼容 Markdown 格式, `.md` 文件既是源码也是文档.

    本文 Zxx-spec.md 不符合源文件命名规则, 被排除

### Fixed

Fixed 是预编译指令, 提供不篡改源代码就能修改其它模块的能力:

*警告! Fixed 可能导致灾难性损失!*

1. 沙盒存储, 所有产生的中间文件都使用沙盒存储
1. 文件限制, 只用于 main, case, fail, fine, race 文件
1. 位置限制, 一个文件只有一条, 且在所有声明之前
1. 影响全局, 可以访问, 新增, 覆盖原模块的所有定义并影响全局
1. 强制测试, 生成目标文件前强制通过所有相关模块的单元测试

格式:

```abnf
fixed-directive = %s"Fixed " module-path [exactly] [SP *VCHAR] LF

exactly = "@" version ["-" [version]]
          ; 1.1.1        Same as version is 1.1.1
          ; 0.0.0-       Same as version >= 0.0.0
          ; 0.0.0-1.1.1  Same as 0.0.0 <= version <= 1.1.1
version = 1*2DIGIT "." 1*4DIGIT "." 1*3DIGIT

VCHAR   = %x20-D7FF / %xE000-10FFFF ; Not including U+D800 - U+DFFF
```

文法中用 Unicode code-points 值描述与编码格式无关的字符.
代理对 U+D800 - U+DFFF 经解码后对应 U+10000 - U+10FFFF 的字符.
解码是与编码格式相关的, 所以整个文法描述中不包括代理对.

例:

```zxx
Fixed any module to handling all errors.

fun error-handle(i32 code, string message, any context)
  ; The original method is empty, reserved for Fixed
  handling-all-errors(code, message, context)
```

## 保留字

Zxx 的保留字有:

- 类型 void, self, base, error
- 变量 iota, self, base
- 函数 echo
- 常量 true, false, null, NaN, Infinity
- 运算 un, not, is, is not, in, not in, and, or
- 声明 use, let, def, fun, let pub, def pub, fun pub
- 语句 defer, if, else, elif, for, of, break, continue, throw, catch, out, yield

提示:

    pub 不是关键字, 如 `let pub` 是一个整体
    同理 `is not` 和 `not in` 也是整体

## 标识符

标识符是对模块, 类型, 函数, 参数, 对象, 方法, 字段的命名.

拒绝新标识符使用保留字.

正则:

```js
let
  identifier = RegExp(/([a-zA-Z]|\p{Lo})(-?[a-zA-Z0-9]|\p{Lo})*/,'u');
```

## 下划线

下划线 U+005F `_` 用于:

1. 左值 赋值语句中被丢弃的左值
1. 形参 fun 声明中不被使用的参数
1. 字段 def 声明中不可访问的字段

## 定界符

定界符包括:

- ,  逗号 U+002C 用于界定多个子项
- \n 换行 U+000A 用于续行或表示语句的开始
- () 圆括号用于函数声明和表达式分组
- {} 花括号用于映射类型
- [] 方括号用于列表类型或映射实例, 访问数组下标或映射成员
- '' 字符串字面值定界符
- `` 模板函数定界符
- "" time 字面值定界符

支持 U+000D `CR` 或 U+000DU+000A `CRLF` 风格的换行.

## 分隔符

空格 U+0020 是唯一的分隔符.

## 注释

两种形式的单行注释:

- 顶级注释 行首非顶级声明开始的整行文本
- 单行注释 声明内 "; " 开始至行尾

```abnf
COMMENT = "; " 1*VCHAR
```

## 缩进和续行

缩进是语法, 是声明内, 行开始处多个连续的双空格.

```abnf
indentation = 1*LF 1*"  "
```

续行是一条语句分成多行来写, 通常出于排版需要, 续行不能违反缩进规则.

缩进和续行的原则:

1. 解析器无需读取下一行就可以确定有续行
1. 多赋值语句的左值必须左对齐
1. 成对的定界符缩进量相同 `()`,`[]`,`{}`,`''`,"``"
1. 兄弟元素缩进量相同

缩进规则:

- 增加缩进时一次增加两个空格
- 子句必须增加缩进

续行规则:

- 拒绝 在匿名类型中续行, 列表, 映射, 匿名函数类型
- 拒绝 在映射的 Key 和 Value 之间有续行
- 拒绝 在 `if` 和 `for` 主句中使用续行
- 允许 在 `=`,`,`,`(`,`[`,`{` 或二元运算符之后, 多行文本之内续行
- 必须 在多行文本内保持缩进
- 最后 不和以上冲突的情况下, 语句中的首个续行必须增加缩进

提示: 方法声明中的点 `.` 不是运算符, 被引号 `''` 包裹重载运算符也不是多行文本

例: 合法的常规缩进

```zxx
let
  _,_,_ = 1,2,3

  _,_=
    [int][
      1,2,3
    ],
    [int][
      1,
      2,
      3
    ]

```

例: 合法的非常规缩进

```zxx
let
  _,
  _,
  _ =
    1,
    2,
    3

  _,_= [int][
    1,2,3
  ], [int][
    1,2,
    3
  ]
```

例: 有问题的缩进

```zxx
Illegal indentation
let
  _,
    _ = 1, 2  ; 左值续行必须左对齐

  _ = [int][1,
  2]          ; 首个续行必须增加缩进

  _ = [int][1,
    2]        ; 同级元素缩进相同, 成对定界符缩进量相同

  _ = [int][1,
    2         ; 同级元素缩进相同
  ]

  _ = [[int]][
      [1,2]   ; 增加缩进时一次增加两个空格
  ]

  _ = [       ; 拒绝在列表类型描述中续行
    int
  ][1,2]

              ; 增加缩进时一次增加两个空格, 注释也不例外
; 过多的减少缩进使得下一行代码被当做注释
  _ = 1 + 2

fun illegal-indentation(int x)
  if x < 0
      echo x  ; 增加缩进时一次增加两个空格

  for i of (1+; 拒绝 在 `def`, `if` 和 `for` 主句中使用续行
    x
  )..10
    echo i
```

例: 合法, 或许正是你需要的

```zxx
let
  _= [string]['
  多行文本内保持缩进
  ']

  _ =
    '首个续行增加缩进, 因为引号不在多行文本内.
     此行前多了一个空格, 两个空格是缩进, 一个空格只是空格.
     剔除换行和缩进后, 最后这两行前的单个空格被保留了.
    '

  _ = [[int]][
    [           ; 第一个续行已经缩进
    1,2,3,4,5,  ; 同级元素缩进相同, 允许子元素不增加缩进
    6
    ]           ; 成对定界符的缩进也对齐了
  ]

  _ = [[int]][[
    1,2,3,
    4,5,6
  ]]
```

## 表达式

表达式产生值且必须被接收处理.

### 运算符和操作符

运算符用于表达式.

| 运算符           | 原始语义                                 |
|------------------|------------------------------------------|
| not,un           | 一元 非真运算                            |
| +,-              | 一元 取正,取负值运算                     |
| ~                | 一元 按位取反                            |
| &,\|,^,<<,>>,>>> | 按位运算 与,或,异或,左移,右移,带符号右移 |
| +,-,*,/,%,**     | 数学运算 加,减,乘,除,取余,幂             |
| <=,<,>=,>        | 比较运算                                 |
| is,is not        | 鉴定运算                                 |
| in,not in        | 成员运算 包含,不包含                     |
| or,and           | 逻辑求值 类型推导, 求值运算              |
| ..               | yield 生成运算符, 习惯上成为 yielder     |
| <-,->            | 输入管道,向输入管道传递参数              |

操作符形成独立的语句.

| 操作符           | 语义                     |
|------------------|--------------------------|
| @=               | 内联赋值                 |
| =                | 赋值                     |

运算优先级: 从低到高

| 运算(操作)                        | 分类                        |
|-----------------------------------|-----------------------------|
| @=,=                              | 内联赋值,赋值               |
| ->,<-                             | 管道                        |
| ,                                 | 逗号                        |
| ..                                | yielder                     |
| or                                | 逻辑求值                    |
| and                               | 逻辑求值                    |
| not x                             | 一元布尔非真                |
| in x,not in x,is,is not,<=,<,>=,> | 元素测试,鉴定,比较运算      |
| |                                 | 位运算                      |
| ^                                 | 位运算                      |
| &                                 | 位运算                      |
| <<,>>,>>>                         | 位移                        |
| +,-                               | 数学运算                    |
| *,/,%                             | 数学运算                    |
| -x,+x,~x                          | 一元取负数,取正,位反        |
| **                                | 数学运算                    |
| un x                              | 一元布尔非真                |
| .,(e..),x(e..),x[e]               | 成员访问,分组,调用,下标访问 |

其中:

- 非运算符: 赋值 `=`, 逗号 `,`, 成员 `,`, 分组 `(e..)`, 调用 `x(e..)`
- `x` 表示运算符重载时的 `self` 位置, 缺省 `self` 在运算符左侧.
- `e` 表示运算符重载时的参数

Zxx 借鉴了 [Python operator][], 不同之处:

- 取消 `//` 因为整数除法 (`/`) 默认行为就是整除
- 取消 `==`,`!=` 使用 `is`, `is not` 替代
- 变更 `not`, `or` 为逻辑求值运算, 支持短路字面值
- 变更 `is`, `is not` 为鉴定运算, 值比较或者类型比较
- 变更 `in`, `not in` 为下标测试
- 变更 `+x` 为取正运算.
- 增加 yielder `..`
- 增加 更高优先级的一元布尔非真运算符 `un`
- 增加 管道

#### 减号

符号 U+002D `-` 被多处使用, 为避免歧义特约定:

    一元负数运算符减号右侧无空格
    二元运算符减号左侧必须有空格, 右侧有空格或换行

#### 取正

一元运算形式 `+x` 是取正运算, 因为类型会发生变化.

```zxx
let
  x = -1
  y = +x
  _ = x is i32 and x is -1
  _ = y is u32 and y is 1
```

编译器为所有的有符号整数类型提供了取正运算, 但无符号整型没有该运算.

#### 取负值

一元运算形式 `-x` 是取负值运算, 类型变成有符号位的类型.
编译器为所有的整数类型和浮点类型提供了取负值运算.
无符号整数类型经过取负值运算, 类型变成有符号整数类型, 且可能发生数据损失.

#### 非真运算

运算符 `not`, `un` 的结果是布尔 `true` 或 `false`.

```zxx
let
  x = 1
  _ = (not x) is false
  _ = (not x is false) is true
  _ = un x is false
  _ = (un x) is false
```

#### 逻辑求值运算

运算符 `and`, `or` 的结果类型是推导出来的. 推导原则:

1. `and` 运算的最后一项决定类型
1. `or` 运算的每一项类型都相同
1. 短路字面值 字面值参与 `and`, `or` 运算时结果是字面值

例:

```zxx
let
  ; 其它语言中的惯用逻辑, 在 Zxx 中有些结果与惯用逻辑不同.
  z = 0
  _ = (z or 0 ) is 0
  _ = (0 and 1 ) is 0
  _ = (0 or 1 ) is 0
  _ = (1 and 0 ) is 1
  _ = (1 or 0 ) is 1
  _ = (false or 1) is false

我们认为上述代码毫无现实意义, 短路字面值能体现实用价值

fun normal(string a, string b, out int)
  ; 常规写法
  if a<b
    out -1
  if a>b
    out 1
  out 0

fun e-g(string a, string b, out int)
  ; 短路字面值写法
  out a is b and 0 or a<b and -1 or 1

  ; 下述写法将编译失败, 因为短路字面值的类型与输出类型不符
  out a is b and 0 or a<b and -1 or '' or 1
```

#### 成员访问

符号 U+002E `.` 用于访问模块或对象的下级成员, `.` 和左侧的模块或对象紧凑连接.

#### 下标访问

编译器实现了列表类型的下标访问, 下标赋值. 映射内部由列表实现.

- `x[e]` 中的 `x` 是 `self`, 参数 `e` 是下标或键
- 下标访问 如果下标越界返回缺省值
- 下标赋值 如果下标越界丢弃赋值

#### 元素测试

运算符 `in`, `not in` 用于测试右侧算子是否包含左侧元素.

注意形式 `y in x` 不是 `y not in x` 的语法糖.
形式其中 `x` 是运算符重载方法的 `self`, `y` 是被测试参数.

#### 鉴定运算

运算符 `is`, `is not` 对一个对象进行类型鉴定或值比较, 被鉴定的对象在左侧.

如果两侧的类型不完全一样时, 右侧算子对鉴定算法的影响:

- 函数 取该函数的类型, 由该函数的参数决定
- 接口 验证左侧的类型定义是否包含该右侧的接口名称
- 其它 右侧是否为左侧的 `base` 类型
- 标量 先匹配类型后进行值对比

注意形式 `x is y` 不是 `x is not y` 的语法糖.
形式其中 `x` 是运算符重载方法的 `self`, `y` 是被测试参数且不可能是类型名.

#### 管道

增加管道运算符主要考虑到单词作为方法名语义太具体, 符号具有抽象性.
当为方法命名很纠结时, 可以考虑使用管道.

通过运算符重载实现管道:

1. 输入管道 运算符 `<-` 表示接收一个参数, 拒绝用于表达式, 重载方法必须有一个输出参数
1. 输出管道 运算符 `->` 表示向输入管道传递一个参数, 重载方法必须只有一个输出参数

无需实现输出管道也能使用 `->` 运算符向输入管道传递参数.

注意:

- 对象是否需要打开或者关闭, 比如文件操作
- 运算符 `->` 的级别很低, 需要仔细分析参数是否已经在之前的运算中发生了变化

#### yielder

运算符 `..` 只能用于 `if` 和 `for` 语句, 原语义表示范围 `[x:y)`, 即不包括 `y`.

#### 运算符重载

运算符重载以单引号包裹运算形式作为方法名并实现执行代码.

- 重载方法必须是 `fun pub`, 不能省略 `pub`
- 运算优先级不变
- 语义可能被改变
- 重载方法无嵌套 重载方法中所有的运算都使用原始行为
- yielder 运算符 `..` 的重载方法必须是 `yield` 方法

例: any 模块中的代码, 包含执行体的接口方法即所谓的泛型方法

```zxx
def number
def signed
def unsigned

fun pub signed.'+x'(out self)
  out self < 0 and -self or self

fun pub unsigned.'+x'(out self)
  out self

fun pub number.'..'(self last, yield self)
  for self < last
    yield self
    self += 1

def pub byte byte -- number unsigned
  ; ---------^^^^ 是编译器实现的内部类型
  ; ----^^^^ 是预定义类型

fun pub string.'is not'(self x, out bool)
  out not self.is(x)
fun pub string.'is'(self x, out bool)
  out self.is(x)

fun string.is(self x, out bool)
  if x is null and self is null
    out true
  ; omitted...
```

例: yielder

```zxx
def photo string

fun pub photo.'..'(string, yield photo)
  ; omitted...

fun example()
  for photo of photo('ken')..'me'
    echo photo ; without me
```

例: yielder 实现 step

```zxx
def step u32

fun pub step.'..'(u32 end, yield u32 x)
  for x < end
    yield
    x += base

fun example()
  for x of step(2)..256
    echo x
```

例: 管道

```zxx
def writer

fun pub writer.write(string, out uint)

fun pub writer.'<-'(string)
  _ = self.wirte(string)

fun pub writer.'<-'(string, out uint)
  out self.wirte(string)

def outer
fun pub outer.'->'(out string)
  ; out some string
```

## null

保留字 `null` 表示某个对象无值.

```zxx
let
  _ = 0  is not null
  _ = '' is not null
  _ = any() is null

  _ = bool() is null
  _ = true  is not null
  _ = false is not null

  _ = time() is null
  _ = time.now() is not null
  s = string()
  _ = s is null and s is not ''

  _ = [int][] is not null
  _ = {int}{} is not null

  _ = [int]() is null
  _ = {int}() is null

  _ = noop is fun() and noop is null

  ; illegal
  _ = null          ; null without type
  _ = null is any   ; null without type
  _ = int(null)     ; The int does not accept null

fun nullable(string s, [int] a, interface i, {int} m, time t, any x, fun() f)
  s = null
  a = null
  i = null
  m = null
  t = null
  x = null
  f = null

fun noop()
```

如果函数的参数类型允许 `null` 值, 调用时可以传递 `null`.

## 类型

在 Zxx 中按照类型的组成形式可分为两种:

- 标量类型 没有字段, 没有下标访问
- 复合类型 由标量类型组成, 有字段或者下标访问

在内存中, 标量类型字段存储的是值, 复合类型字段存储的是引用.

标量类型有 void, 整型家族, 浮点家族, 布尔类型以及他们的别名类型, 其它是复合类型.

例:

```zxx
def scalar int

def composite -- iii
  field any-type
```

下文列举预定义类型.

### byte

一个字节是连续的 8 bits 内存, 是基本的内存操作单位.

值范围是正整数的 0 到 255.
支持单个 Unicode `U+0` 到 `U+255` 的字符串字面值.

```zxx
let
  b = byte()
  _ = b is byte and b is 0 and b is not u8
  _ = byte(0) is 0
  _ = byte('a') is 97
  _ = byte('a') is 0x61
  _ = byte('A') is 0x41
```

### 字节序列

字节列表 `[byte]` 是编译器维护的列表类型. 支持字符串转换到 `[byte]`.

```zxx
let
  _ = [byte]('bytes from string literal')
  s = 'hi'
  _ = [byte](`{{s}}`)
```

字节数组 `[capacity,byte]` 是编译器维护的列表类型.

### void

保留字 `void` 表示地址类型, 由编译器实现, 拒绝被继承.

1. 类型地址 一个类型只有一个地址, 永远不变
1. 函数地址 只取该函数的类型地址
1. 模块地址 拒绝取模块对象的地址
1. 对象地址 受多种因素影响该地址可能已被回收
1. 使用限制 可用于鉴定运算, 函数参数, 转换为字符串

例:

```zxx
fun addr(string s, out void)
  out void(s)

fun addr(void v, out void)
  out void(v)

let
  s = 'tom'
  v = void(s)
  x = v
  _ = x is void  and v is void and x is v
  _ = v.string() is string

def type
fun type.noop(any, out void)

def iface
fun iface.noop(any, out void)

fun noop(any, out void)

let
  t = type()
  i = iface()
  _ = void(t.noop) is void(i.noop)
  _ = void(i.noop) is void(noop)
  _ = void(noop)   is void(fun(any,out void))
  _ = noop is not void(noop)

```

### 整数家族

整数家族都是标量类型, 类型有:

- `u8`   无符号  8 位整数 0..255
- `u16`  无符号 16 位整数 0..65535
- `u32`  无符号 32 位整数 0..4294967295
- `u64`  无符号 64 位整数 0..18446744073709551615
- `i8`   有符号  8 位整数 -128..127
- `i16`  有符号 16 位整数 -32768..32767
- `i32`  有符号 32 位整数 -2147483648..2147483647
- `i64`  有符号 64 位整数 -9223372036854775808..9223372036854775807
- `rune` 表示一个 Unicode 码点值, 是 u32 的别名

下列类型的长度/尺寸与运行环境有关:

- `int`  在 32 位架构中是 i32, 在 64 位架构中是 i64
- `uint` 在 32 位架构中是 u32, 在 64 位架构中是 u64

整型字面值缺省类型为 `i32` 或 `u32`. 格式:

```abnf
i32-lit = %x31-39 *DIGIT ["E" 1*DIGIT] ; default i32
u32-lit = "0" (                        ; default u32
      %s"b" 1*64BIT
    / %s"x" 1*16HEXDIG
    / i32-lit
  )
```

提示: 负数中的 `-` 被当做运算符处理

如果字面值超出范围必须显示提升到其它类型, 比如: `u64`,`i64`.

字面值不受损失的情况下可转换为其它数值类型.

```zxx
let
  _ = 0   is i32
  _ = 00  is u32
  _ = 0x01 is u32
  _ = 0b0101010101 is u32
  _ = -0x01 is i32

  u = u8()
  i = i8()
  _ = u is 0 and u is u8 and u is byte
  _ = i is 0 and i is i8 and i is byte
  _ = i16(1) >= 0
  _ = rune('多') is 0x591a
```

### 浮点家族

浮点数是标量类型, 内部遵循 [IEEE_754-2008][], 由符号位 s, 指数域 E, 尾数域 M 组成.

- `f32`    32 bits s(1) + E(8)  + M(23)
- `float`  64 bits s(1) + E(11) + M(52), 缺省类型
- `f128`  128 bits s(1) + E(15) + M(112)

浮点数字面值的格式: 支持 "0f" 开头的十六进制字符串表示的二进制数据

```abnf
float-word = %s"NaN" / %s"Infinity"
float-lit  = 1*DIGIT "." 1*DIGIT ["E" ["+"/ "-"] 1*DIGIT]
float-bin  = %s"0f" 1*2(1*2(8HEXDIG))
              ; 32-bit / 64-bit / 128-bit
```

浮点数运算规则

  NaN 参与比较运算       NaN is not NaN is true, NaN is NaN is true
  NaN 参与数学运算结果为 NaN
  正负零相等             0.0 is -0.0
  零乘除计算符号位       ±0.0 is ±0.0 * 1.0
  Infinity               ±Infinity is 1.0 / ±0.0

### 千位分隔符

十进制的整型和浮点数字面值支持千位分隔符 `-`, 如果使用了千位分隔符, 那么:

1. 小数点前 向左必须每 3 位数字一个 `-`
1. 小数点后 向右必须每 3 位数字一个 `-`

```zxx
let
  _ = 12-345-678   ; 12345678
  _ = 12-345.678-9 ; 12345678.6789
```

为避免歧义, 下述写法会导致解析失败

    1-23
    1-2
    12345.678-9
    12-3.6789
    12-345-67

### bool

布尔是标量类型, 支持三态布尔. 字面值表示:

```abnf
bool-word = %s"true" / %s"false"
```

例:

```zxx
let
  _ = i8(true)  is 1
  _ = i8(false) is 0
  _ = i8(bool()) is -1

  _ = bool() is null
  _ = (un false) is true
  _ = (not true) is false

def bool i8

fun bool.'un'(out bool)
  out base is 1 and true or false

fun bool.'not'(out bool)
  out un self

fun bool.string(out string)
  out base is 0 and 'false' or
  base is 1 and 'true' or 'null'
```

### string

字符串的值是多个 Unicode 字符, 只能整体赋值不可修改, 支持 `null`.

在字符串中使用 U+0000 时应该注意跨语言交换数据的兼容性.

字面值是以单引号包裹的多行文本, 支持反斜线转义字符.

```abnf
string-lit = "'" *string-char "'"

string-char =
    escaped
  / %x20-26 / %x28-D7FF / %xE000-10FFFF
  / 1*LF 1*"  "

escaped =
  %x5C (
      %x78           ; xXX                  U+XX
      2HEXDIG
      ; Otherwise fail
    / %s"u"            ; uXXXX                U+XXXXXX
      (
          "{" 1*6HEXDIG "}"
        / 4HEXDIG
        ; Otherwise fail
      )
    / %x21-7E
    ; The following characters will be escaped
    ; %x62           ; b  backspace         U+0008
    ; %x66           ; f  form feed         U+000C
    ; %x6E           ; n  line feed         U+000A
    ; %x72           ; r  carriage return   U+000D
    ; %x74           ; t  horizontal tab    U+0009

    ; Otherwise fail
  )
```

例:

```zxx
let
  _ = string() is null and '' is not null
  a = '支持多行, 必须保持缩进,
   换行和缩进(续行)缩进
  被剔除'
  _ = a is '支持多行, 必须保持缩进, 换行和缩进(续行)被剔除' is true

  _ = 'hello' + 'word' is 'helloword' is true

  _ = '支持转义\n' is '支持转义\n'
  _ = '支持转义\n
  ' is '支持转义\n'
```

非法缩进的例子:

```zxx
illegal indentation
let
  _ = '
    '
  _ = 'a
    b
  c'
  _ = '
missing-indentation
  '
```

Zxx 支持任意对象转换为 string.

```zxx
use
  io

fun string(any x, out string)
  if any is null
    out 'null'
  ; ...

let
  _ = string() is null
  _ = string() is not ''
  _ = string() is not 'null'
  _ = string('') is ''

  _ = string(string()) is 'null'

  _ = string(false) is 'false'
  _ = string(0) is '0'
  _ = string(0.0) is '0.0'
  _ = string(1.0/3) is '0f3fd5555555555555'
  _ = string(io.file) is 'io.file' ; 显然模块也是一个对象

  illegal = string(bool)
```

### 模板函数

模板函数是一对儿反引号 U+0060 "`" 包裹的多行文本, 可嵌入表达式代码, 不支持反斜线转义.

```abnf
template =
  "`" *(
    1*tpl-char / "{{" *' ' expression "}}" / 1*LF 1*"  "
  ) "`"

tpl-char= ["{"] (%x20-5F / %x61-7A / %x7C-D7FF / %xE000-10FFFF)
```

由 `{{}}` 包裹的是标准的 Zxx 表达式, 值被转换为字符串类型.

例:

```zxx
let
  _ = `` is string
  x = `Hi
   hello {{'word'}}\n` ; 必须保持缩进

  _ = x is 'Hi\n hello word\\n' is true ; 行首缩进被剔除
  _ = `{{ `{{'支持嵌套'}}` }}` is '支持嵌套'

  _ = inc(1) is '1 + 1 is 2'

fun inc(int num, out string)
  out `{{num}} + 1 is {{num+1}}`
```

### any

预定义类型 any 接受任意类型的对象, 拒绝被继承.
事实上 any 和所有接口类型都继承了私有类型 `type`.

```zxx
let
  _ = 0 is not any
  _ = 'any' is not any
  _ = any() is null
  _ = any('tom') is any
  _ = any('tom') is string
  _ = any('tom') is 'tom'
  _ = any([int][0]) is [int]
  _ = any([int][0]) is not [int 128]
  a = my-any()
  _ = a is my-any and a is not any

def pub my-any
fun pub my-any.string(out string)
```

any 类型的 `string` 方法返回一个精简的字符串描述所接受的对象.

```zxx
fun pub any.string(out string)
```

例:

```zxx
let
  _ = any('string').string() is '\'string\''
  _ = any(1).string() is '1'
  _ = any(01).string() is '01'
  _ = any(noop).string() is 'fun()'
  _ = any([int][1,2,3]).string() is '[int][]'

fun noop()
```

### time

time 类型表示日期和时间, 字面值由双引号包裹, 参考了 [ISO 8601][] 格式.

字面值(包括非闰秒)中的 `LSC` 表示在运算中加入[闰秒][Leap_second]造成的影响.

必须有 `LSC` 标记闰秒才能被识别和验证, 否则 `60` 秒是非法的.

IERS 提前发布的 [leapseconds][] 被用来验证闰秒的具体时间,
影响及算法依据参见 [leap-seconds.list][].

```abnf
time =
  DQUOTE (
      ymd [%s"T"] [hms [ns] [LSC]] [TZ]
    / [%s"T"] hms [ns] [LSC] [TZ]
    / ns
    / TZ
    ; Otherwise fail
  )
  DQUOTE

ymd = 4DIGIT ["-"] 2DIGIT ["-"] 2DIGIT
  ; year-month-day
hms = 2DIGIT [":"] 2DIGIT [":"] 2DIGIT
  ; hour:minute:second
ns  = "." 1*9DIGIT
  ; nanosecond
LSC = %s"LSC"
  ; leap second effects in time calculations
TZ  = %s"Z" / offset
  ; UTC Timezone or UTC offset hours [minutes]
offset = ("+" / "-") 1*2(2DIGIT)
```

例:

```zxx
let
  localdate   = "2016-02-04"
  localtime   = "2016-02-04T21:49:00"
  utcdate     = "20160204Z"
  utctime     = "20160204T21:49:01Z"
  zonedate    = "20160204T+08"
  since       = "20160204T21:49:00.123456789+08"
  hhmmss      = "21:49:00"
  thhmmss     = "T21:49:00"
  nanoseconds = ".999999"

  _ = since is time

  _ = "20160204+8000"
  _ = "214900+8000"
  _ = ".9"
  _ = "+8000"
  _ = "-70"
  _ = "Z"
  _ = "20160204214900.123456789+08"
```

提示: 使用者应该知道如何使用数据不全的时间, 否则可能造成非预期的结果

### 列表

列表是元素序列, 下标是元素在序列中的序号.

- 元素值可修改
- 下标从 0 开始
- 列表类型可被继承
- 列表对象被整体赋值时必须带上类型描述

```abnf
list-type  = "[" base-type [SP capacity] "]"
capacity   = ["+"] 1*DIGIT

list       = list-type elements
elements   =
  "["
    *(expression "," *SP [COMMENT] [LF 1*"  "])
    [expression]
  "]"

external   = identifier "." identifier
base-type  = list-type / map-type / anonymous / external / identifier
```

扩容前缀 `+` 表示 `capacity` 是每次扩容的长度, 当容量不够时.

根据容量不同列表可分为两种:

- 切片 容量可变, 未声明 `capacity` 或者带 `+`
- 数组 容量固定, 声明了 `capacity` 且不带 `+`

数组使用上和切片唯一的区别就是容量固定.

如果声明切片时容量为 0 也由编译器决定扩容方案.

如果声明数组时容量为 0, 首次初始化元素的个数就是容量.
当然对于变量, 必须在声明时一次性赋初值, 对于字段则可以单独赋值.

声明字段类型或函数的参数类型时可用 `[T 0]` 表示数组类型, `[T +0]` 表示切片类型.

切片缺省的扩容行为由编译器决定, 扩容失败将产生错误.

例:

```zxx
let
  _ = [any]() is null
  _ = [any][] is not null
  _ = [int 3]() is not null

  slice = [int][1,2,3]
  _ = slice is [int]
  _ = slice.length is 3
  _ = slice[0] is 1
  _ = slice[100] is 0
  _ = slice is [int]
  _ = slice is [int +0]

  array = [int 3][]
  _ = array is not [int]
  _ = array is not [int +0]
  _ = array is [int 3]
  _ = array is [int 0]
  _ = array[0] is 0
  _ = array[100] is 0

  multi = [multi-type][
    1
    'Supports multiple element types'
    2
    'Comma can be omitted before Newline-Indentation'
    [3,4]
  ]
  _ = multi is not [int]
  _ = multi is not [any]
  _ = multi is multi-type

  _ = multi[0] is any
  _ = multi[0] is int
  _ = multi[1] is any
  _ = multi[1] is string
  _ = multi[2] is any
  _ = multi[2] is [int]
  _ = multi[100] is null

def multi-type int string [int]

允许做为基类型
def list [string]
```

#### 属性和方法

编译器为列表实现了:

- `length`      uint 拥有的元素个数
- `capacity`    uint 当前分配的容量, >= length
- `size`        uint 容量直接占用的字节量, 与 capacity 线性相关
- `slice(i,j)`  切片 返回 `[x:y)` 的新切片, x,y 类型为 uint
- `remove(i)`   如果下标 i 存在移除该元素, 直接影响 length

数组的 `slice` 方法也返回切片类型.

`slice` 返回的新切片的元素值变化不直接影响(可能间接影响)原数组的元素值.

```zxx
let
  array = [int 0][
    1,
    2,
    3
  ]
  _ = array.length is 3

  slice = array.slice(1,3)
  _ = slice.length is 2
  _ = slice[0] is 2 and slice[1] is 3
  _ = slice is [int +0]
  _ = slice is not [int]
```

### 映射

映射是有序 Key-Value 元素集合, 键类型为 `any` 且值唯一, 允许修改属性值.

- 映射总是可扩容的
- 映射类型可被继承
- 映射对象被整体赋值时必须带上类型描述
- 映射的键值不能是字面的 `null`, 因为字面值 `null` 不是 `any` 类型

```abnf
map-type   = "{" base-type "}"
map        = map-type propertys

propertys  =
  "{"
    *(property "," *SP [COMMENT] [LF 1*"  "])
    [property]
  "}"

property   = expression SP expression
```

映射使用两个切片实现, 源码在 `any` 模块中.

例:

```zxx
let
  _ = {any}{} is null
  _ = {int string}{} is null
  nil = any()

  map = {bool int string}{
    'true'  true,
    'false' false,
    nil     null,
    string  'string'
    int     'int'
    1+2     `Sum is {{1+2}}`
    fun()   'typeable'
  }

  _ = map[nil] is null
  _ = map[3] is 'Sum is 3'
  _ = map[string] is 'string'
  _ = map[fun()] is 'typeable'

  config = {string}{
    'ZBIN'  '~/zxx/bin'
    'ZPATH' '~/zxx/src', 'ZGETFIRST' 'github.com'
  }

fun asiic(out {byte} map)
  if map is null
    map = {byte}{}

  map[n] = byte(n) for n of 0..128 ; Syntactic sugar

  m = map(0..1)
  m[0] = 100
  echo map[0] is 100 ; true

  x = map('string'..1)
  ; if original-map is null or map.exist('string') is false
  echo x is {byte}  ; true
  echo x is null    ; true


允许做为基类型
def list {string}
```

## 声明

每个声明都有独立的语法.

### use

保留字 `use` 声明引入模块, 引入的模块仅在文件内可用.

一个文件只能有一个 `use` 声明且在其它声明之前.

```abnf
use-decl= %s"use" LF
  1*("  " module-path [exactly] [SP global] [COMMENT] LF)

global  = identifier / "_"
```

例:

```zxx
use
  os               ; 引入 ZPATH 下的 os 模块, 缺省命名 os
  w3c.org/fetch    ; 引入 w3c.org/fetch 模块, 缺省命名 fetch
  path/to/dir pkg  ; 引入 path/to/dir 模块, 自定义命名 pkg
  ./pkg-path-file  ; 支持向深层的相对路径, 缺省命名 pkg-path-file
  some/pkg _       ; 引入一个不使用的模块
```

### let

保留字 `let` 声明顶级对象并赋初值, 保留字 `let pub` 声明的顶级对象提供外部访问.

```abnf
let-decl= %s"let" [%s" pub"] LF
  1*("  " global *("," *SP global) *SP "=" *SP
  expression *("," *SP expression) [COMMENT] LF)
```

对象是某个类型的实例, 顶级对象在模块内共享, 不能与其它标识符冲突.

顶级对象:

- 拒绝被整体赋值
- 成员具有可读写性
- 不会产生变量提升 能被整体赋值的是局部变量

依照命名惯例:

1. 无成员的顶级对象可称为常量
1. 有成员的顶级对象可称为变量, 注意拒绝被整体赋值

编译器的工作顺序:

1. 确认标识符没有冲突
1. 根据右侧表达式确认对象类型
1. 为对象创建标识符, 即创建对象并分配空间
1. 表达式求值并赋值

例:

```zxx
let
  a = true
  b = 'global'
  config = init()

fun init(out {string} c)
  a = 1                ; 局部变量
  c = {string}{0 b}    ; 这个 b 是顶级常量
  b = 'local'          ; 这个 b 是局部的, 并被后续代码使用

  echo config is {string} ; true
  echo config is null     ; true
  echo c is {string}      ; true
  echo c is null          ; true

  config[0] = 'member'
  config[1] = config
  c[config] = b
  out c

let
  _ = config[config] is 'local'
  _ = config[0] is 'global'
  _ = config[1] is config
```

### def

保留字 `def` 声明新类型, `def pub` 声明的类型允许外部访问.

允许覆盖预定义类型名(事实上是继承), 但不能覆盖两次.

```abnf
def-decl   =
  %s"def " [%s"pub "] identifier
  *(SP base-type) [SP "--" 1*(SP implement)] [*SP COMMENT] LF
  *("  " field-decl)

implement  = external / identifier
field-decl =
  (identifier / "_") SP base-type [*SP COMMENT] LF
```

不同组合产生的类型区别:

- 标量类型 无字段, 单继承标量类型
- 复合类型 有字段, 来自单继承的或新增的
- 别名类型 单继承, 无新增字段, 可能是标量类型或复合类型
- 接口类型 无继承, 无字段, 有方法, 来自新增的或实现接口
- 枚举常量 单继承整数类型, 声明字段时只有字段名
- 类型约束 多继承表示类型限制, 无字段, 无方法, 无实现接口

字段名可以为 `_` 但不可访问. 通常用于占位填充.

别名类型和基类型是两个不同的类型, 必须进行显示转换

例:

```zxx
def num int
def Num
  _ int

let
  i = int()
  n = num()
  N = Num()
  _ = i is int and n is int N is not int
  _ = i is 0 and n is 0 and N is null

非法的空接口
def illegal-empty-interface

非法继承 any
def illegal-based-any any

非法实现 any
def illegal-interface -- any

接口
def interface
fun interface.method()

结构体
def fruit
  string      name
  string      brand, color
  [string]    flavors

别名, integer 是 int 的别名, 可以给它声明新的方法
def integer int

因为声明了新的字段, INT 不是 int 的别名
def INT int
  int x

别名类型和原类型是两个不同类型, 但可以进行显示转换
let
  _ = integer() is 0
  _ = integer() is not int
  _ = int(integer()) is int
  _ = integer(int()) is integer
  _ = typeof(int()) is 'int'
  _ = typeof(integer()) is 'integer'
  _ = typeof(true) is null

fun typeof(any x, out string)
  out any is int and 'int' or
    any is integer and 'integer' or null
```

#### 枚举和约束

类型约束的基类型是 `any` 且约定了可接受的类型, 因此拒绝单继承类型约束.
多继承类型约束时得到的还是继承约束.

枚举常量的基类型是它继承的整数类型, 允许实现接口.
编译器为枚举常量实现了 `string` 方法, 缺省返回对应的字段名.
转换到枚举类型之前必须进行越界检查, 形如: `if x in enum-type`.

例:

```zxx
def types int string fun(int) enum

def enum u8
  _
  one   ; one is u8 and one is 1
  two   ; two is u8 and two is 2

def flags enum
  three ; three is u8 and three is 3

let
  _ = u8(0) not in enum
  _ = 1 in enum and u8(2) in enum
  _ = 1 is not enum and u8(2) is not enum
  _ = enum.one is 1 and enum.one is u8(1)
  x = enum.one
  _ = x is enum and x is u8 and x is 1
  _ = x.string() is 'one'
```

用例: 编译器负责检查语句中的枚举约束

```zxx
; 如果分支枚举了所有的可能, 那不会例外情况, 即不应该使用 else

fun example1(types x)
  if y, x of int
    z = u8(y)
    ; example2(flags(z)) 将编译失败, z 可能越界
    if z in flags
      ; 编译器会实现枚举常量的 in 运算, 除非重载
      ; 编译通过, 因为已经用 in 进行了越界检查, z 不可能越界
      ; 如果没有使用 in 进行了越界检查, 编译失败
      ; 间接通过函数检查会被判断为未进行越界检查
      example2(flags(z))
  of string
    dosomething()
  of byte
    dosomething()
  of enum                ; 注意 of 子句的顺序
    dosomething()
  of u8
    dosomething()
  ; omitted...

fun example2(flags x)
  ; x is u8
  if x is flags.one
    dosomething()
  elif x is flags.two
    dosomething()
  elif x is flags.three
    dosomething()
  ; 穷举了枚举值, 无需判断例外
```

### 继承与可访问性

外部可访问性(读写)受两个因素影响:

- 对象是否被暴漏, `pub` 修饰的或者被 `out`
- 对象的类型是否是 `pub` 修饰的

显然被暴漏的对象总是可以被整体访问, 它的字段, 下标, 方法的可访问性:

- `pub` 修饰的方法可被访问
- `pub` 修饰的类型, 它的字段和下标可被访问

例:

```zxx
def interface

fun pub interface.pub-method()
fun interface.method()

def pub base
  bField int

def type base  -- interface
  tField int

fun pub example(out type)
  ; omitted....

def pub class base  -- interface
  cField int

fun pub example(out class)
  ; omitted....
```

在外部得到 `type` 类型的对象后, 可访问:

- bField
- pub-method()

但外部无法使用类型名 `type`.

在外部得到 `class` 类型的对象后, 可访问:

- bField
- cField
- pub-method()

## 语句

语句位于函数或方法声明中, 无返回值. 分为三类:

1. 基本语句 由语句保留字开头
1. 赋值语句 用赋值符号分配表达式的值到左侧标识符
1. 调用语句 无返回值的函数或方法调用

### 赋值语句

赋值语句左侧为值接收对象, 右侧是表达式, 赋值操作符是 `=`, 支持复合赋值.

1. 赋值语句 左值和右值各有一个
1. 复合赋值 是语法糖. 形式 `x += y` 总是转换为 `x = x + y`
1. 多值赋值 不是语法糖, 包括 `iota` 自增, 不等价拆分成多条赋值语句

```abnf
assignment =
  receiver *SP
  (
      compound-assignment *SP expression
    / *("," *SP receiver) *SP ["@"] "=" *SP expression *(comma expression)
  )

receiver = identifier *("." identifier / "[" expression "]")

compound-assignment =
  (
      "**"
    / "+" / "-" / "*" / "/" / "%"
    / "|" / "^" / "&" / "<<" / ">>>" / ">>"
  ) "="
```

例:

```zxx

fun example()
  x, y, z = 1, 2 , x + y
  echo x is 1 and y is 2 and z is 0

  x = 1
  y = 2
  z = x+y
  echo x is 1 and y is 2 and z is 3

  y, z = z + 1, y
  echo y is 4 and z is 2

  x += y + z
  echo y is 7
```

注意 `self` 是个变量, 对 `self` 赋值不会影响外部.

```zxx
def pipe

fun pipe.write(string)
  ; dosomething
```

### 内联赋值

内联的语义是内部(局部)对象关联, 不是其它语言的内联函数.
内联赋值 `@=` 右侧表达式中的函数调用会改变传递的局部对象, 拒绝传递顶级对象.

```zxx
fun inc-and-sum(int x,int y, out int)
  x+=1
  out x + y

fun call()
  x = 1
  z @= inc-and-sum(x,2) + inc-and-sum(x,2)
  echo x, z ; 3, 7
```

所以内联的语义是: 内部对象关联, 显然不能传递顶级对象

### echo

保留字 `echo` 用于调试输出信息, 由编译器实现.

```abnf
echo = %s"echo" 1*SP expression *(comma expression)
```

### if

保留字
if 语句有三种形式.

if-else 选择, `elif` 和 `else if` 等价.

  if condition
    [statements]
  [else if conditionN
    [alternate statements]]
  [elif conditionN
    [alternate statements]]
  [else
    [alternate statements]]

if-of-type 断言, 当 `[instance,] any` 与首个相同时可省略, 编译器会补全, 下同.

  if [instance,] any of type
    [statements]
  [else if [instance,] any of typeN
    [statements]
  [elif [instance,] any of typeN
    [statements]
  [of typeN
    [statements]
  [else
    [alternate statements]]

if-of-yeilder 循环, 缺省 yield 缓冲长度 1, 自定义值范围为 0..255

循环结束条件:

1. 所有 yielder 结束
1. 使用 `break` 结束

if-of-yeilder 中可使用 `continue` 语句.

  if yield-out[,yield-out...] of yielder()
    [statements]
  [else if yield-out[,yield-out...] of yielderN()
    [statements]
  [elif yield-out[,yield-out...] of yielderN()
    [statements]
  [of yielderN()
    [statements]
  [else
    [alternate statements]]

  if yield-out[,yield-out...] of yield(u8, yielder())
    [statements]
  [of yield(u8, yielderN())
    [statements]
  [else
    [alternate statements]]

```zxx
fun example(
  any x,
  fun(yield int, int) yielder,
  fun(yield string, int) worker,
  fun(yield int, int, int) another
  )

  if x is null
    echo null
  elif x is string
    echo 'string'
  else
    echo 'other'

  if v, x of string
    echo 'x is string: ' + v
  of int
    echo 'x is int: ', v + 0
  else
    echo 'other: ', v

  ; if-of-yeilder
  if x, y of yielder()
    echo 'yield int, int:', x, y
  of worker()
    echo 'yield string, int:', x, y
  else
    if something()
      break
    if some-thing()
      continue
    echo '...'


  ; 可以自定义缓冲长度
  if x, y of yield(2, yielder())
    echo 'yield int, int:', x, y
  of yield(3, worker())
    echo 'yield string, int:', x, y

  ; 如果几个 yielder 的输出参数个数不同, 可利用 else 解决
  if x, y of yielder()
    echo 'yield int, int:', x, y
  elif x,y,z of another()
    echo 'yield int, int, int:', x, y, x
```

### for

for 语句有多种形式, 执行体中可以包含 break, continue.

for-else: 当条件从未成立时 else 执行体被执行.

  for [condition]
    [statements]
  [else
    [alternate statements]]

for-of-range:

  for [key,] value of x...y
    [statements]

for-of-yeilder: 缺省 yield 缓冲长度 0, 自定义值范围为 0..255

  for yield-out[,yield-out...] of yielder()
    [statements]

  for yield-out[,yield-out...] of yield(u8, yielder())
    [statements]

```zxx
fun example(u8 x, fun(yield int, int) yielder)

  for
    echo x
    if x is 0
      break
    x-=1

  for y of 0...3
    echo y

  ; 缺省 yield 缓冲长度为 0
  for i, j of yielder()
    echo i, j
    if i < 0
      break

  ; 自定义 yield 缓冲长度
  for i, j of yield(2, yielder())
    echo i, j
    if i < 0
      break
```

#### 创建实例

Zxx 中没有 `new` 或者构造函数概念, 依照惯例使用该词仅为方便描述.

可以直接使用 `type()` 创建实例. 显然代码中不能创建与类型同名的无参数函数.

```zxx
let
  _ = fruit() is null ; fruit 的定义在前文中
  _ = newFruit() is null
  _ = fruit('apple', 'lucky', 'red') is not null
  _ = flavors('sweet', 'sour') is not null

允许同名函数, 但必须有参数
fun fruit(string name, string brand,string color, out self)
  ; self 在参数声明中表示类型 fruit, 执行体中表示实例
  ; 显然返回的实例是不可执行的
  self.name, self.brand, self.color = name, brand, color

可以使用变参, 相关约束见函数声明部分
fun flavors(string flavors...,out fruit)
  fruit.flavors = flavors

fun newFruit(out fruit)
  ; 允许无执行体
```

#### 局部结构体类型

允许在函数或者方法中声明局部结构体类型, 但不能声明新方法.

```zxx
def coord
  float x, y

fun example(out coord)
  def coord
    int x
    int y

  out coord(int(1), int(2))
```

### fun

保留字 `fun` 声明一个函数或方法. 格式举例:

```zxx
fun name()        ; 无参数
fun name(int)     ; 一个输入参数, int 类型, 无名
fun name(any x)   ; 一个输入参数, any 类型, 名子为 x
fun name(out int) ; 一个输出参数, int 类型, 无名

fun vals(int...) ; 变参必须是最后一个输入参数, 允许无形参名

各种组合
fun func(int, bool x, string y..., out type, bool ok, self)
```

语义:

  声明一个函数类型 func(...)
  该类型有一个命名为 name 的实现

函数的特征

1. 函数的类型由参数类型决定, 参数类型相同的函数类型相同
1. 缺省形参名与类型名称相同, 显然该名称在函数中不再是类型
1. 输入参数支持变参, 变参必须是最后一个参数
1. 非变参函数支持函数重载, 参数类型必须有区别
1. 执行体中能被整体赋值的变量都是局部变量
1. `self`,`base` 在参数声明中可作为类型, 在执行体中是变量
1. `out` 用于输出参数声明只能出现一次.

### out

保留字 `out` 用于执行体中表示结束执行体. 格式:

    out [arguments...]

注意 out 语句不是系列写法的语法糖.

  params0 = ,..., paramsN = arguments0, ..., argumentsN
  out

```zxx
fun fn(out string s)
  ; 调用时传入的输出变量值被改变了.
  s = 'hello' + s
  ; out 语句不是必须的, 当然也可以写成
  ; out 'hello' + s

fun call()
  fn('world') is 'helloworld'
```

### base

保留字 `base` 用于方法:

- 位于参数声明时 表示方法所属类型的基类型
- 位于方法执行体 代表 `self` 的基类型值

*在方法执行体中 `base` 和 `self` 具有联动性*

### self

保留字 `self` 用于函数或方法:

- 位于参数声明时 表示方法所属类型
- 位于函数执行体 表示该函数本身
- 位于方法执行体 表示所属类型实例

```zxx
use
  os

fun recursive(int n)
  if n>0
    self(n - 1)
  ; ^^^^---- recursive

fun callback(string, os.fileInfo, error, out bool deepin)

fun walk(string path, callback)

  for path of os.file-list(path)
    info = os.file-info(path)
    if callback(path, info, null) and info.isdir()
      self(path, callback)

  catch err
    callback(path, null, err)
```

执行体中的 self 是个可赋值变量, 但不会改变所属类型实例.

```zxx
let
  obj = type(1)
  _ = obj.example() is 0
  _ = obj.x is 1

def type
  int x

fun type.example(out int)
  ; self 在执行体内就是个变量
  ; Zxx 中没有 this, 可以把 self 看做 self = this
  self = self(0) ; self = type(0)
  out self.x
```

### 参数传递

两种固定传递方式

  值浅拷贝 输入参数, 左值接收参数.
  引用传递 传入输出参数.

示例

```zxx
def obj
  x int

fun beFalse(obj a, out obj b)
  echo void(a) is void(b)
  ; false, 因为 a 是个值拷贝对象

fun isSame(out obj a, obj b)
  echo void(a) is void(b)

fun test()
  a = obj()
  isSame(a, a) ; true 因为是引用传递
  b = a
  isSame(a, b) ; false

def box
  obj o
  int y

fun setx(box b)
  b.y = 1
  b.o.x = 2 ; 字段是结构体的话, 指向同一块内存
fun example()
  c = box(obj(0), 0)
  setx(c)
  echo c.o.x  ; 2
  echo c.y    ; 0
```

### defer

defer 声明的语句在 out|throw 语句之后执行. 多条 defer 以 `后进先出` 顺序执行.

格式:

  defer
    statements

显然 defer 和 out 语句很可能发生嵌套.

示例:

```zxx
fun good-example()
  file = open('somefile')

  defer file.close() ; defer 不会立即执行, 被推延到 out|throw 语句之后.

  processing(file)

  catch
    ; 忽略错误
    out ; 执行 file.close()

  dosomething()

  if true
    throw errno, message ; 执行 file.close()
  defer
    echo 'never'
```

### yield

如前文所示 yield 可用于

1. 参数声明   与 out 互斥
1. of yield   定制 0..255 缓冲长度, 非 0 时以协程或线程方式工作.
1. 独立语句   表示 yield 输出

以协程或线程方式工作时, 接收到的数据是缓冲区的副本, 副本算法:

  copy1, copyN = yield-argument1, yield-argumentN

注意如果 yield-argument 实现了 `=` 运算符重载, 其行为会被执行.

```zxx
fun fibSeries(u32 n, yield u32 Fn)
  Fn-1 = 01
  for 0..n
    yield
    Fn, Fn-1 = Fn + Fn-1, Fn

fun fib-echo(u32 n)
  ; 无缓冲非协程
  for val of fibSeries(n)
    echo val

  ; 带缓冲多协程
  for val of yield(1, fibSeries(n))
    echo val
```

### 局部函数

局部函数在函数或方法执行体中声明. 局部函数中会发生变量提升, 但不会提升至顶级.
拒绝在局部函数中再声明类型或函数.

示例:

```zxx
fun loop-echo(uint n)
  x = uint(0)
  fun update(out bool)
    echo x
    x+=1
    out x < n

  for update()
```

局部函数被暴漏到外部函数时可能产生闭包.

示例:

```zxx
fun loop(uint n)
  for cond(n)()

fun cond(uint n, out fun(out bool) fn)
  x = uint(0)
  fun update(out bool)
    echo x
    x+=1
    out x < n ; x, n 的生命周期取 cond 执行体与返回函数的生命周期的大值
  out update
```

### 错误处理

产生的错误都应该被抛出并捕获, 推荐总是在 `main` 中包含错误处理代码.
未被捕获的错误会导致主进程结束.

保留字 `error` 用于产生错误:

```abnf
error = %s"error(" *SP error-arguments *SP ")"

error-arguments =
  code-expression *SP "," *SP
  message-expression
  [*SP "," *SP context-expression]
```

产生错误时必须附加 `i32` 类型的错误代码和 `string` 类型的错误信息,
以及可选的 `any` 类型上下文附加数据.

保留字 `throw` 直接抛出所有的错误参数或者一个 `error` 对象, 并构建堆栈信息:

堆栈信息包括文件路径和关键语句所在的行列位置.

```abnf
throw = %s"throw" 1*SP ( error-arguments / identifier )
```

保留字 `catch` 捕获错误, 使用变量名接收所有的参数和可选的堆栈信息.

```abnf
catch =
  %s"catch" 1*SP
  code-identifier *SP "," *SP
  message-identifier *SP "," *SP
  [
    context-identifier
    [*SP "," *SP stack-identifier]
  ]
  [*SP COMMENT] LF

  1*"  " statements
```

例:

```zxx
let
  _ = open('filename')

fun open(string filename, out os.file)
  if not os.exist(filename)
    throw 1, 'file name is empty', filename

fun processing(os.file d)
  if isempty(file)
    throw 2, 'file is empty'

代码块对 catch 的影响
fun example(string filename)
  processing(open('config.yaml'))
  if filename
    _ = open(filename)
    catch code, message, context, stack
      ; 不能捕获 processing(open('config.yaml')) 的错误
      ; 可捕获 open(filename) 的错误
      dosomething()

  catch code, message, context, stack
    ; 可捕获 processing 的错误.
    if code is 0
      echo 'eof'
    else
      throw code, message, context ; 继续弹出

fun example2(string filename)
  processing(open('config.yaml'))

  if filename and open(filename)
    catch code, message, context
      dosomething()

  catch code, message, context
    ; 可捕获 processing 和两个 open 的错误
    if code is 0
      echo 'eof'
    else
      dosomething()
```

替换成 try...catch 结构很容易理解代码块对 catch 的影响.

```js
function example(filename) {
  try {
    processing(open('config.yaml'));
    if (filename) {
      try {
        _ = open(filename);
      } catch (error) {
        dosomething();
      }
    }
  } catch (error) {
    if (error === null)
      echo('eof');
    else
      dosomething();
  }
}

function example2(filename) {
  try {
    processing(open('config.yaml'));
    if (filename && open(filename)) {
      try {
        // nothing
      } catch (error) {
        dosomething();
      }
    }
  } catch (error) {
    if (error === null)
      echo('eof');
    else
      dosomething();
  }
}
```

## LICENSE

[BSD 2-Clause License](https://github.com/ZxxLang/zxx/blob/master/LICENSE)

Copyright (c) 2018 The Zxx Authors All rights reserved.

[IEEE_754-2008]: https://en.wikipedia.org/wiki/IEEE_754-2008
[syntax-across-languages]: http://rigaux.org/language-study/syntax-across-languages.html
[scripting]: http://hyperpolyglot.org/scripting
[Zxx Definition of ABNFA]: https://github.com/ZxxLang/abnfa/blob/master/grammar/zxx.abnf
[Unicode]: https://en.wikipedia.org/wiki/Unicode
[Python operator]: https://docs.python.org/3/reference/expressions.html#operator-precedence
[SPDX License List]: https://spdx.org/licenses/
[semver.org]: https://semver.org/
[ISO 8601]: https://en.wikipedia.org/wiki/ISO_8601
[leap-seconds.list]: https://www.ietf.org/timezones/data/leap-seconds.list
[leapseconds]: https://www.ietf.org/timezones/data/leapseconds
[Leap_second]: https://zh.wikipedia.org/zh-hans/%E9%97%B0%E7%A7%92