
Z 目前只是个设想

#简介

ZxxLang 是一门通用编程语言. 初衷是混合文档和代码于一体并降低 shift 按键使用频率.  ZxxLang 通过忽略不能识别的文本, 减少 shift 按键的使用来实现此目标.

	初衷貌似不靠谱, 事实上 Zxx 的语法更近自然语言.

简便起见, 下文用 Z 替代 ZxxLang 进行描述.

Z 是强类型语言, 虽然可根据上下文进行类型推导, 但最终类型是明确的, 并和相关操作要求匹配.

Z 的语法比较宽松下面的几种输入风格都是合法的:

```
var string s = 'hello word'

var [
	int x
	int y=5
	int z,i=4
]

var int a, string b

var (
	f32 f = 9.0
)

var datetime (
	day = 20160202, now = 20160202T22:48:33
	orz = 20160202T22:48:33Z
)

proc hello string word [
	echo 'hello ' word
]

甚至缩进风格
proc long_function_name =
        var int _one, _two, _three,
	       _four
		doSomeThing()
```

这基本上就是自然语言的风格.

	Z 没有代码风格规范, 书写者自己说了算. 代码规范令人厌恶.
	如果 Z 不能识别某种风格代码, 只说明 Z 不够强壮.

Z 提供了源码格式化工具, 如果你愿意使用, 上面的代码会转换成下面的样子

```
var string s = 'hello word'

var {
	int x
	int y=5
	int z,i=4
}

var int a, string b

var f32 f = 9.0

var datetime {
	day = 20160202
	now = 20160202T22:48:33
	orz = 20160202T22:48:33Z
}

proc hello string word {
	echo 'hello ' word
}

甚至缩进风格
proc long_function_name {
	var int _one, _two, _three,
		_four
	doSomeThing()
}
```

#源文件组织

Z 依靠目录名和文件名组织源代码:

1. 项目     项目这就是一个 git 仓库, 在本地就是一个目录
2. 包名     项目下的子孙目录名
3. 文件     多个文件组成一个包
4. 多包     同目录下的文件可归属于不同的包

源文件目录名和文件名格式:

1. 目录名格式为 `package([-_][a-z0-9]*)?`.
2. 名字格式为 `(package-)?name(_platform)?`.
3. 文件扩展名为 'md' 或者 'zxx'

其中 `package`,`name`,`platform` 为正则 `[a-z]+[a-z0-9]*`.

一些 `package` 具有特别含义.

1. main 可执行
2. test 测试
3. example 样例
4. benchmark 评测
5. hack  只在 main 包中被工具链自动载入.

某些时候急需修改某个依赖包, 那么无需等待维护者更新, 也无需修改依赖包源码, 自己写个外部的 hack `package`  就可以替换掉原包中的代码.

文件名中的包名和子目录重名举例:

```
rep
├── hello/
├── hello-word/
├── use.json
├── hello-word.md
├── rep.md
├── some.md
└── rep-hello.zxx
```

子目录 'hello' , 'hello-word' 和文件 'hello-word.md' 归属同一个包 'hello'.

文件 'rep.md', 'some.md' 归属包 'rep', 也就是它们的上层目录名所表示的包名.

文件 'rep-hello.zxx' 归属包 'rep/rep', 因为它包含了新的包名称.

##use.json

因 'use' 是保留字不能作为包名, 所以用 'use.json' 描述所在目录包细节是无污染的.

同理 'main.json' 用来描述 main 包细节.

	目录的层级关系只是存储的组织形式, 包都是平级的, 没有层级归属.

'use.json' 采用语义化属性名, 通常这无需特别解释, 看到名称就知道作用.


```JSON
{
	"name": "zxx",
	"license": "BSD-2-Clause",
	"version": "0.0.0",
	"repository": {
		"type": "git",
		"url": "https://github.com/zxxLang/zxx.git"
	},
	"keywords": [
		"zxx programming language"
	],
	"author": {
		"name": "YU HengChun",
		"url": "http://achun.github.io/"
	}
}
```

#语素

学术上语素是构成语言的最小单位, Z 有两种语素:

	字面值 源码中的一段文本
	标识符 有格式要求的文本, 是具有明确语义的 Z 实体标记

Z 中它们是互斥的. Z 是这样确定语素的:

1. 是不是符合内建格式的字面值
2. 是不是预声明字面值
3. 是不是预声明标识符
4. 是不是能声明新标识符
5. 判定为 占位字面值

```
字面值和标识符组成表达式或语句
表达式只能存在于语句中
语句构成源文件
源文件构成 block
```

#字面值

字面值是一段文本. 它的语义类型是在使用时判定的, 有可能被判定为非法值.

##预声明字面值

```
null      无值, 零值
false     假
true      真
NaN       非数
infinite  无穷大
```

行首的空格, tab 缩进以及换行符 CR, LF 其实也是字面值.

##缩进

缩进是行首连续的空格 ' ' 或 Tab '\t', 通常用于排版目的.

Z 不支持混用空格和 Tab 的缩进.

也可能书写者喜欢 Python 风格的缩进代码, Z 的解析器会尝试理解这种格式, 并把它当做分组括号对待.

Z 支持类 Python 风格的缩进代码是因为缩进代码能少输入一个括号, 可降低输入量.
但不保障兼容, 毕竟 Z 不是 Python.

注意: 行首的 '\t' 缩进才是缩进, 行中间出现的 '\t' 缩进被当做尾注释处理

##占位文本

Z 不能识别语义的文本被当做占位文本处理.

#标识符

标识符的正则为 `[a-zA-z_]+[a-zA-Z_0-9]*`.
自定义标识符用于类型名, 常量名, 变量名, 过程名.
预定义标识符见下文.

#符号

##分号

分号 ";" 用来指示语句的结束边界.
源码经过格式美化后分号会被省略. 除非代码 minify 到一行.

换行或者空行会依据上文的完整性推导为分号. 推导规则:

	如果上文至此换行可形成完整的语句, 把换行当做分号.

*注意: 不良的换行可能引起下文非法, 甚至被当做注释.*

##定界符

定界符是一对符号, 表示代码实体的起始和终结. 有三对符号可用

	() 分组定界可用于很多地方, 参见示例
	{} 代码定界用于 block 的开始和结束
	[] 通用定界


方括号 '[]' 更便于输入, 多数情况下可以替代圆括号和花括号. Z 的格式美化工具会替换合适的符号.

##逗号

逗号 ',' 具有多种用途. 可用来分隔声明, 表达式以及语句.

注意逗号和分号的不同.

```
var int x,y = 3+1,6		逗号分隔多个同类型变量
var array[int] = [1,2]	逗号分隔值表达式
var (
	int i;string s		分号分隔一行中多个声明
)

proc each [
	for var int x,y = 0,0; x < 10 ;x++, y = y*x[
		三段式 `for` 语句中第一, 第三段是语句.
		echo x,y
	]
]
```

逗号分组

```
if a and (b or c) [
]
使用逗号分组
if a and, b or c [
]

复杂的分组运算
if a and (x or ((b or c) and (d or e))) [
]
使用逗号分组无法消除所有分组括号
if a and, x or (b or c) and (d or e) [
]

下文将给出进一步消除分组括号的方法
```

##运算符

运算符只能出现在表达式中, 下表按优先级从高到低排列.

| 运算符                   | 解释			    |
|--------------------------|--------------------|
| $                        |模板取值专用		|
| not,!,~                  |一元右结合非,位反	|
| &,\|,xor                 |位与,或,异或		|
| mul,*,/,mod,%,<<,>>      |乘除,整数取模,位移	|
| add,+,-                  |数值加减,字符串连接	|
| to                       |数值至数值			|
| ==,!=,<=,>=,>,<          |比较运算			|
| is,isnot,has             |结果是布尔类型		|
| and                      |与					|
| or                       |或					|

模板取值运算符 '$' 是 Z 语言特征,  不能归类到一元运算符.

```
运算符 mul 等价于四则乘法运算符 '*'.
运算符 add 等价于四则加法运算符 '+'.
加减乘除四则运算要求运算子类型一致且结果类型不变

-3 / 2 == -1 整除
0.1*0.1 isnot 0.01 是浮点数精度造成

取模运算符 mod 的运算子为同类型整型
a mod b 等价于 a - (a/b)*b
分步运算举例

7 mod 4
7 - (7 / 4)*4
7 - 4
3

-7 mod 4
-7 - (-7 / 4)*4
-7 - (-1)*4
-7 - (-4)
-7 + 4
-3

7 mod -4
7 - (7 / -4)*(-4)
7 + (-1)*4
7 + (-4)
7 - 4
3

-7 mod -4
-7 - (-7 / -4)*(-4)
-7 + 1*4
-3

取模运算符 '%' 运算 a % b 时先对取绝对值, 等价 |a| mod |b|.

运算符语法糖可进一步消除分组括号, 对于分组运算例子
if a and (x or b or c) and (d or e) [
]

使用运算符语法糖
if a and x.or b.or c and d.or e [
	语法糖让二元运算变成一元右结合运算, 用伪函数推导一下:
	a and x.or(b.or(c)) and d.or(e)
	a and (x or b.or(c)) and (d or e)
	a and (x or b or c) and (d or e)
]

同理
var int x = i.add 1 mul 5 等同 (i+1)*5
```

二元运算符 'is', 'isnot' 有多种情况.

```
同类型时(null 值会依据右侧算子类型自动推导出同类型值)
a is b 等价 a == b
a isnot b 等价 a != b
int(0) is null 等价 int(0) == 0
'string' is null 等价 'string' == ''

数组, 映射和自定义类型没有被赋值时值为 null, 所以
type empty []
type T [
	int x
]

var array[int] x,
	y = [0]

var empty e,
	T t,
	T t1 = T{x:1}

x is null		结果为真
y is null		结果为假
y isnot null	结果为真
e is null		结果永远为真, 因为无法对 e 的属性赋值
t is null		结果为真
t1 is null		结果为假
```

运算符 'not', 'and', 'or' 的运算结果是 null 或者运算子的值.

#类型

类型名称是个标识符. 预声明类型标识符是保留字.

	保留字 type 是一切类型的根类型.

##布尔

bool 的字面值为 true 或者 false.
在内存中, 布尔类型用整数表示, 0 表示 false, 非 0 表示 true, 具体尺寸由编译器决定.

注意: true, false 只是字面值, 是无类型的.

##整数

预声明固定长度/尺寸整数类型有:

```
u8		无符号的 8 位整数 (0 to 255)
u16		无符号的 16 位整数 (0 to 65535)
u32		无符号的 32 位整数 (0 to 4294967295)
u64		无符号的 64 位整数 (0 to 18446744073709551615)

i8		带符号的 8 位整数 (-128 to 127)
i16		带符号的 16 位整数 (-32768 to 32767)
i32		带符号的 32 位整数 (-2147483648 to 2147483647)
i64		带符号的 64 位整数 (-9223372036854775808 to 9223372036854775807)

byte	和 u8 一样
rune	用 u32 表示的 unicode 码点, 0x0 到 0x10FFFF
```

Z 中用 byte 或 rune 类型可代表单个字符. 可以用字符串对 byte 和 rune 类型进行赋值.

```
var byte b = 'Hello Word' // b 的值为 'h'
var rune r = '世界你好'   // r 的值为 '世'
var rune e = ''           // 空字符串的 rune 或 byte 值为 0.
```

事实上这是由 string 类型的 toByte, toRune 方法进行的转换.

应该注意 rune 在内存中的字节序是由平台决定的.

下列类型的长度/尺寸与运行环境有关:

```
uint     u32 或 u64 的别名
int      i32 或 i64 的别名
```

整数有多种输入格式, 可选前置正负符号位:

1. 0b  开头后跟 [0-1]+ 的二进制表示
2. 0x  开头后跟 [0-9a-f]+ 的十六进制表示
3. 0h  开头后跟十六进制表示, 十六进制字符个数对应固定尺寸类型
4. 0-9 开头后跟 [0-9]* 的十进制表示


##浮点数

预声明的固定长度/尺寸浮点类型有下面这些:

```
f32     IEEE-754 32  位浮点数
f64     IEEE-754 64  位浮点数
f128    IEEE-754 128 位浮点数
```

浮点数有多种输入格式:

1. 可选前置正负符号位后跟十进制数字夹杂一个"." 或 "e"表示的十进制指数
2. 0f 开头后跟 8|16|32 个十六进制字符表示的 IEEE-754 二进制制浮点数

##字符串

string 是一对单引号或者双引号包裹的多行文本.
字符串连接运算符使用 '+' 或者 '-' , 它们是等价的.

```
proc hello out string [
	out 'hello'
]

proc word out string [
	out 'word'
]

var string (
	a = '单引号字符串不支持反斜杠\转义值' +
		'直接断行也可以,
		换行会被保留, 前置空白会被剔除.' -
		"连接符使这四行连接为一个字符串值"

	b = '字符串还支持嵌入参数 $a $b'.['a'=1,'b'=2] -
		'SQL 参数序号风格 $1 $2'.[$a, $b]
		"变量值会被转换为字符串表示,
		这四行组成了一个字符串"

	c = hello() + word() 使用 '+' 号连接
	d = hello() - word() 使用 '-' 号连接

	e = hello() - 连接符后的尾注释会被正确处理
		' word'

	f = '空白行会影响字符串拼接' 此处不能用连接符连接下下行

        '上面的空白行产生分号, 赋值语句完结了'

	g = '连续多行字符串, 此行没有使用连接符'
        '这一行只是个注释'

)
```

字符串具有简单的模板功能.

```
echo '{
		"name": "$1",
		"age": $2
	  }'.['tom', 8]
```

输出:

```JSON
{
"name": "tom",
"age": 8
}
```

字符串成员方法 `execute` 会剔除每一行两端的空白符和换行符.

```
echo '{
		"name": "$1",
		"age": $2
	  }'.execute['tom', 8]
```

输出:

```JSON
{"name": "tom","age": 8}
```

##日期和时间

datetime 表示日期和时间, 以 [ISO 8601][] 为标准.

在源码中使用无 '-' 间隔年, 月, 日的基本格式.

```
var datetime(
	localdate = 20160204
	localtime = 20160204T21:49
	utcdate = 20160204Z       尾部的 'Z' 表示时区为 0
	utctime = 20160204T21:49Z
	zonedate = 20160204T+08   加号 '+' 两端是紧凑相连的
	zonetime =  21:49+08
	now = 20160204T214900     可简化时间格式省略 ':' 号
)
```

##数组

array 表示数组类型.

```
var array[int] ai	以 array 开头声明数组

var array[int] (	分组写法
	c
	d = [			赋初值就直接写吧
		1,2,3,
		call()		使用表达式运算结果
	]
)

访问数组元素
var int e = ai[0]	下标访问

多维数组

var array[array[int]] point =[
		[1,2,3]
	]

定长数组

var array[5,int] seq

type Seq {
	array[int]		匿名属性数组
}

var array[5,Seq] seqs
```


##映射

使用 map[keyType,valueType] 声明映射类型.

```
给 map 赋值可以使用 JSON 风格, Z 编译器会检查值的合法性
var map[string,int] m = {
	'age': 13,
	"height": 156		有换行的话可以省略逗号
	"id": 1
}

key 和 value 可以是任意类型, 甚至使用 type
var map[type,string] types = [	万能的方括号
	string: 'string',
	int: 'integer'
]

var map[type,type] assoc = [
	string: types,
	int: f32
]

访问映射
var int e = m['age']	支持直接访问已有的变量
var type t= m[string]
```

Z 是强静态类型语言, 编译后, 所有的类型都是明确的, 可识别的,
在 Z 中 type 的类型就是 type, 可以用于过程的参数类型.

#表达式

表达式具有以下特征

1. 必定产生运算结果.
2. 必定存在于语句中.

表达式运算结果可分为零值和非零值.

#语句

语句无返回值, 语句产生的代码执行顺序和书写顺序一致.

在 Z 源文件中最先出现的语句是注释或声明.

##注释

注释是对某个标识符或语句体的说明.

注释在 Z 中属于语句, 因为注释有明确的格式:

1. Z 中的注释是后置的并紧根语句实体
2. 双斜线尾注释 以 '//' 开始至行尾, '//' 位于行首是这种写法的特例
3. 制表符尾注释 非空白字符之后以 "\t" 开始至行尾.
4. 块注释 以三个以上的减号 '-' 位于一行首个非空白符, 并以此结束.
5. 多字节注释 因标识符不含多字节字符, 这很容易识别.

```
// 这两行是占位文本, 虽然样子像尾注释.
// 纵观整个文本, 无法确定这条注释属于哪个标识符

---
这两行是占位文本, 虽然样子像块注释.
纵观整个文本, 无法确定这条注释属于哪个标识符
---

因为标识符由英文字符和数字组成, 所以多字节注释可以直接写.
同理, 这两行依然是占位

The line is a comment.
上一行英文也是占位, 因为解析器判断 'The' 不是声明保留字.
避开声明保留字, 顶层的英文注释可直接书写.

此行位于 proc sum 上面, 但不是注释.
Z 中的注释是后置的, 这和其它语言不同.
proc sum int x,y,out int [
	// comment for 'proc sum' 非顶层单行英文注释需要双斜线开头
	out x add y		result x+y

	注意 'result x+y' 前面有 '\t', 所以它是尾注释
	而这两行被之前的空行分断被当做占位
]

proc multiByteFriendly out bool [
	使用多字节字符写注释....
	out true 就是这么直接
]
```

This is an example for comment in English.

```
proc sum int x,y out int [
	------- sum x and y -------
	out x add y		result x + y

	// this is placeholder text follow empty line
]
```

##赋值语句

赋值操作符产生赋值语句.

| 操作符 | 解释   |
|-------|--------|
| ++    |自赋值增一|
| --    |自赋值减一|
| =     |赋值     |


##分支语句

有两个标识符可选, 'if' 和 'switch'.

```
if expr {
	当 expr 的值为真时
	// do something
} else [	也可以使用方括号包裹执行体
]

if expr {
	// do something
} else if expr1 {
}

switch expr {
case a:
	if something {
		break
	}
	// do ...
case b,c:	可选多值匹配
default: {
}

前面的 break 会跳转到这里
```

分支语句中的 break 只能存在于 case 中.
每个 case 语句块结束的位置总是隐含一条 break.


##循环语句

循环语句由保留字 for 开始, 有多种语法.

```
一段式循环条件
for condExpr {
	if expr1 {
		continue
	}

	if expr2 {
		break
	}
}

三段式循环条件
for doSomething; condExpr; doSomethingBeforeNextLoop {
}
常见的
for var i = 0; i < 10; i++{
}

遍历数组
for array as index {
}

for array as index item {
	其中 index 为数组下标, item 为 array[index] 的值
}

遍历 int 数值范围, 从 1 到 10
for 1..10 as index item {
	echo index, item 显然 item 和 index + 1 是相等的
}

遍历立即字符串数组
for ['name', 'nick', 'email'] as index item {
	相当于匿名声明了一个 static [string]
}

遍历已经声明的映射变量 maps
for maps as key {
}

for maps as key val {
}

遍历立即定义的映射需要声明类型
for map[string, int]('a':1,'b':2) as key intVal {
}

遍历类型成员
for int as name typ {
	在 Z 中类型标识符可作为参数
}

甚至通过遍历了解 type
for type as name typ {
}

遍历前文声明的类型 fruit
for fruit as name typ {
}

遍历非映射变量
var fruit f = [color='red']
for f as memberName typ {
	注意 typ 不是 memberName 对应的值, 是 memberName 对应类型
	因为 Z 是强类型语言
}
```

循环, 分支与 is, isnot 运算符的使用:

```
var map[type,string] types ={
	string: 'string',
	int: 'integer'
}

var map[type,type] assoc ={
	string: types,
	int: f32
}

proc fn [
	for assoc as t val {
		if t is int {
			doSomething()
			continue
		}

		if t has string and t[string] isnot null {
			break
		}

		var maps = types(val[string]) 强制类型转换

		switch  {
		case string:

		}
	}
]
```

#声明

声明产生的代码执行顺序由 Z 编译器决定

Z 源代码总是由下列声明开始的.

```
'pub const var static func proc type use'
```

	很明显, 用纯英文写 Z 源码, 避开这些词写顶层注释非常容易.
	如果行首用到这些词, 大写首字母或者随便加个非空白符号就行.


##使用声明

保留字 use 用来声明使用外部包或者声明编译参数.

```
使用单包
use "os"
use path "path"				命名

使用多包分组写法
use (
	"io"
	path "path/filepath"	命名
)

紧凑写法
use ("os", path "path")

声明编译参数

use "-pub"			表示该文件中所有顶层变量,静态,常量,类型,函数都是 'pub' 的.

use "-pubmethod"	表示该文件中所有顶层类型的方法都是 'pub' 的.
```


##导出声明

保留字 static 是对其它声明的修饰.

在声明前加保留字 pub 表示允许外部 block 访问, 否则只允许在 block 内访问.

##静态声明

保留字 static 是声明修饰, 用于修饰变量, 属性和过程.

静态过程详见过程声明.

静态变量或定值属性表示它只被赋值一次.

```
static int x = 9				静态变量声明修饰总是省略 'var'
static datetime day = oneday()	直接调用函数赋值
pub static string prefix		导出静态字符串变量

proc oneday int interval, out datetime [
	out datetime.now().add(interval)
]

proc init [
	prefix = 'z' 随便在哪里进行赋值都行
]

proc fn {
	static int x = datetime.time()
	过程内的静态声明
}

proc setName string val [
	static string name

	if name is null {
		用 is null 可判定变量是否已经被赋值.
		这与 name == null 是等价的.
		name = val
	}

	name = val 最简单的写法, Z 会确保只赋值一次
]

proc setOnce string val [
	static string name = val 最简洁的写法
]

声明时直接给属性赋值

type rep {
	static string version = '0.0.0'    定值属性
	string host = 'https://github.com' 初值
}
```

##常量声明

常量是用标识符代表字面值.

```
const (
	CRLF = "\r\n" 常量名使用大写是常见的习惯
	zero = 0      也许你习惯用小写或者混合大小写
	Day  = 20110101T
	ONE  = 1.0
	YES  = true
	code = {
		可用成对的 '[]|{}|()' 包括的文本作为常量值
		其实单引号, 双引号也是成对的.
		虽然此常量命名为 code, 确切类型和用途在用时才显现,
		也许这就是个注释, 谁知道呢
	}

	EASY = [
		1,2,3 随便了, 反正用的时候才有意义
	]
)

const (
	One = iota + 1 // 1
	Two            // 2
	Three          // 3
)
```


##变量声明

保留字 var 是变量声明前缀字符串. 除了语义上的需要, 更实际的作用是让混合文档和源码解析简单些.

```
var (
	byte   b = CRLF 此时常量的值类型显现了
	string s = CRLF
	rune   r = CRLF
)
```

##过程声明

函数, 类型方法统称为过程, 由名字, 参数和过程体组成

```
func toString out string
保留字 func 只声明过程名字和参数类型, 不给参数命名
保留字 out 表示之后的参数被输出
各部分间不使用逗号, 也就是说 func 声明中不允许出现逗号.

func toInt string out int

---
因为 func 声明中不允许出现逗号, 所以参数不能是函数类型.
func callToInt func string out int
这会产生二义性.
---

只能这样声明 callToInt
func callToInt toInt

这意味着函数或过程声明同时也是类型声明.

proc toString out string s [
	保留字 proc 声明带执行体的过程, 必须声明过程名
	除了唯一的 out 参数名外, 其它参数必须声明参数名和类型
	s = 'hello' - s
	end		end 终止过程
]

pub proc walk func callback int out int, out bool [
	过程声明的参数间用逗号分隔, 该过程有两个参数.
	显然过程参数可以是函数, 这是也是 func 参数无逗号的原因.
	唯一的一个 out 参数可以省略参数名, 那表示用 out 替代.
]

使用函数变量
proc fn [
	var func f out string 声明函数变量用 func
	var func c = fn 声明并赋值

	f = proc out string { 匿名过程赋值
		out = 'hello'
	}

	var func fn = proc out int {
		用匿名过程赋值, 直接赋初值可推导参数类型, 否则需要明确参数类型
		out 1
	}

	函数声明只明确参数类型, 无参数名, 无逗号保障了定义函数变量
]
```

如果只用 func 声明了过程而没有实现代码, 那意味着:

```
pub func fn	在外部实现, 或在参数中作为类型约束

func call	来自外部连接库, 或被 block 内的过程输出
```

详细例子

```
proc run {
	proc {
		这是个匿名执行体, 无参数无返回值, 并且直接执行了
		x = 9
	}() 必须执行它, 不然就成了值表达式, 而不是语句

	proc fn {
		具名子过程
	}

	proc noret {
		执行无返回值的具名过程
	}()

	var int x = proc out int {
		执行匿名过程, 返回值被接收
		out 1
	}()

	--- 以下非法
	var fn1 = proc {}	缺少类型声明
	var proc fn2 {}		这是什么鬼
	proc ret out int {
		执行了具名过程, 返回值没被接收, 隐含了一个执行表达式
		表达式不能独立存在
		out 1
	}()
	proc {
		匿名执行体, 未被赋值也未被执行, 隐含了一个过程表达式
		表达式不能独立存在
	}
	---
}
```

过程返回值与表达式和语句的关系

```
proc x out int {
	out 5
}

proc noret {}

proc y {
	x()			非法, 因为产生了运算结果是个表达式, 表达式不能独立存在.

	null x()	用保留字 null 丢弃结果形成语句.

	noret()		合法语句, 因为没有产生运算结果.

	noret[]		输入源码时可以用方括号替代圆括号执行过程
}
```

关于过程:

1. 参数声明中 out 前的参数必须具有类型和名字
2. 参数声明中 out 后的参数具有输入输出双向性
3. 只有一个 out 参数时, 可以省略参数名
4. 参数都有零值或缺省值,  null 是通用的零值.
5. 调用过程不必传递所有参数.
6. 过程体中 out 既可以替代输出参数, 也可以用作返回语句
7. 函数或者过程声明也是类型声明.

```
proc one out int [		省略唯一的输出参数名
	if out is null [
		out = 1			赋值语句
	]
	out++
	out out				合法但不简洁, 其它语言写作 return out
]

proc fn out int x,y [	多个输出参数必须具名
	out.x = 1
	out.y = 2
	可把 out 当做对象用
	out =[
		x = 1
		y = 2
	]
	也可以用 JSON 风格
	out ={
		x: 1,
		y: 2
	}
	直接输出参数也可以
	out {
		x:1,
		y:2
	}
]
```


##类型声明

使用保留字 type 声明自定义类型.

```
type T {		声明类型 T
	string name	定义属性.
}

type A int		别名声明

type math {}	无属性类型, 详见下文
```

别名声明的类型不能扩展方法, 最常见的是为将来增加类型属性提供了方便.
在 Z 中别名和原类型是完全等价的, 别名为代码维护提供了方便.

定义成员的格式依然是: 类型在前, 属性名在后

```
type fruit {						带属性的水果类型
	pub static string name			定值属性只被赋值一次
	bool              clean=true	赋初值为真
	string            brand, color	同类型属性列表
	[string]          flavors		多种口味并存

	pub func          calculate		属性类型可以是个函数
}

匿名复合使用前缀保留字 use, apple 复合了 fruit 的所有成员
type apple {
	use fruit =[			类型复合并设置初值
		name='apple',
		color='red'
	]
}

匿名复合并导出类型 string 的所有成员.
pub type path [
	pub use string
]

成员方法, 是对 '实例的方法' 声明
pub proc fruit.isSweet out bool [
	out self.flavors isnot null and
		self.flavors.has('sweet')
]

静态方法, 是对 '类型的方法' 声明, 使用时用 path.last()
pub static proc path.last out string [
	out out.split('/').slice(-1)
]

proc path.name out string [
	类型推导 支持 self 替代第一个匿名复合类型
	if self == '/' { 等价使用 self.string
		out '/'
	}
	out self.last(self) 等价于
	out self.last(self.string)

	更简单的写法
	out self is '/' and self or self.last(self)
]

无属性类型示例
pub {
	type math {}

	static {
		导出静态类型方法, 使用上类似命名空间
		proc math.pow int x, uint pow out int {
			out = 1
			for ;pow; pow-- {
				out *=x
			}
		}

		proc math.pow int x, int pow out f64 {
			var [
				int step = pow>0 and 1 or -1
				f64 y = pow > 0 and x.f64 or 1/x.f64
			]

			out = 1
			for  ;pow; pow -= step {
				out *=y
			}
		}
	}
}
```

Z 中没有 class, 只是简单的类型复合. 保留字 self 代指类型实例或类型本身.
Z 支持几个特别的属性名可省略 self 进行访问, 当然这首先需要使用者声明这些属性

```
type node [
	node   parent
	[node] child
	root   this
]

type root [
	use node
]

proc root.count out int [
	out child.length()
]

proc node.fn [
	if this.count() [
		// do something
	]
]
```

注意: Z 只是对 this, parent, child 属性提供了便捷访问方式, 不关心, 也不明白它们与 self 的真正关系.

#类型约束

相对于其它语言中的接口概念, Z 中使用保留字 'as' 进行修饰.

```
pub type integer [] 类似接口表示所有的整型.

事实上预定义整型实现了下列方法
pub func integer {
	add as integer out as integer
	mul as integer out as integer
	div as integer out as integer
	i8 out i8
	i16 out i16
	i32 out i32
	i64 out i64
	u8 out u8
	u16 out u16
	u32 out u32
	u64 out u64
}

显然上面使用了分组, 指的是一组实例方法

使用约束
pub proc opAdd as integer i, u, out i64 [
	out i.add(u).i64()
]
```


[ISO 8601]: https://en.wikipedia.org/wiki/ISO_8601