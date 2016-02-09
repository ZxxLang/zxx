
Z 目前只是个设想 

#简介

ZxxLang 是一门文档风格通用编程语言. 初衷是混合文档和代码于一体并降低 shift 按键使用频率.  ZxxLang 通过忽略不能识别的文本, 并减少或替代和 shift 按键有关的符号来实现此目标.

	也许实现目标所用的方法让一些人难于接受, 但这正是 Zxx 的特征.

简便起见, 下文用 Z 替代 ZxxLang 进行描述.

Z 的语法比较宽松下面的几种输入风格都是合法的:

```
var string s = 'hello word'

var [
	int x
	int y=5
	int z,i=4
]

var (
	float f = 9.0
)

var datetime (
	day = 20160202, now = 20160202T22:48:33
	orz = 20160202T22:48:33Z
)

proc hello string word =[
	echo 'hello ' word 
]
```

Z 支持用输入方括号替代圆括号和花括号, 因为方括号不用 shift 键, 更便于输入. 

	Z 没有代码风格规范, 书写者自己说了算. 代码规范令人厌恶.

Z 提供了源码格式美化工具, 如果你愿意使用, 上面的代码会转换成下面的样子

```
var string s = 'hello word'

var {
	int x
	int y=5
	int z,i=4
}

var float f = 9.0

var datetime {
	day = 20160202
	now = 20160202T22:48:33
	orz = 20160202T22:48:33Z
}

proc hello string word ={
	echo 'hello ' word 
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
├── block.json
├── hello-word.md
├── rep.md
├── some.md
└── rep-hello.zxx
```

子目录 'hello' , 'hello-word' 和文件 'hello-word.md' 归属同一个包 'hello'.

文件 'rep.md', 'some.md' 归属包 'rep', 也就是它们的上层目录名所表示的包名.

文件 'rep-hello.zxx' 归属包 'rep/rep', 因为它包含了新的包名称.
 
##block.json

上例中的文件 'block.json' 用来描述该目录下的源码情况.
文件名使用 'block' 而不是 'package' 是因为 Z 中没有 package 概念, 只有块 -- block.

虽然 Z 没有占用 block 关键字, block 就在代码中, 通常被 '{}' 或者 '[]' 包裹着.

显然使用仓库, 项目, 包这些词汇进行描述更友好.

	目录的层级关系只是存储的组织形式, package 或者 block 都是平级的, 没有层级归属. 

'block.json' 采用语义化属性名, 通常这无需特别解释, 看到名称就知道作用.


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

字面值和语素组成表达式或者语句, 表达式只能存在于语句中, 语句构成 block.

#字面值

字面值是就是一段文本. 它的语义类型是在使用时判定的, 有可能被判定为非法值.

##预声明字面值

```
`
null      无值, 零值
false     假
true      真
NaN       非数
infinite  无穷大
`
```

行首的空格, tab 缩进以及换行符 CR, LF 其实也是字面值.

##占位文本

Z 不能识别语义的文本被当做占位文本处理.

#标识符

标识符的正则为 `[a-zA-z_]+[a-zA-Z_0-9]*`. 自定义标识符用于类型名, 常量名, 变量名, 过程名. 预定义标识符下文详述. 

#符号

##分号

分号 ";" 用来指示语句的结束边界. 源码经过格式美化后分号会被省略. 除非代码 minify 到一行. 

换行或者空行会依据上文的完整性推导为分号. 推导规则:

	字符串表达式有独立的规则.
	如果上文至此换行可形成完整的语句, 把换行当做分号.

*注意: 不良的换行可能引起下文非法, 甚至被当做注释.*

##定界符

定界符是一对符号, 表示代码实体的起始和终结. 有三对符号可用

	() 分组定界可用于很多地方, 参见示例
	{} 代码定界用于 block 的开始和结束
	[] 通用定界


方括号 '[]' 更便于输入, 可以替代圆括号和花括号. Z 的格式美化工具会替换合适的符号.

下文中会出现分组方括号的例子.

##逗号

逗号 ',' 具有多种用途. 可用来分隔声明, 表达式以及语句.

```
var int x,y = 3+1,6     分隔声明变量, 分隔表达式
var [int] array = [1,2] 分隔值表达式
var (
	int i,string s      分隔声明变量
)

proc loop =[
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

| 运算符                    | 解释			|
|--------------------------|----------------|
| $                        |模板取值是独立运算	|
| ~                        |一元右结合位反		|
| &,\|,xor                 |位与,或,异或		|
| mul,*,/,mod,%,<<,>>      |乘除,整数取模,位移	|
| add,+,-                  |数值加减,字符串连接| 
| to                       |数值至数值		|
| ==,!=,<=,>=,>,<          |比较运算			|
| is,isnot,has             |结果是布尔类型		|
| not                      |一元右结合非		|
| and                      |与				|
| or                       |或				|

模板取值运算符 '$' 是 Z 语言特征,  不能归类到一元运算符.

```
运算符 mul 等价于四则乘法运算符 '*'. 
运算符 add 等价于四则加法运算符 '+'. 
加减乘除四则运算要求运算子类型一致且结果类型不变

-3 / 2 == -1 整除
0.1*0.1 isnot 0.01 浮点数精度造成的

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
if a and (x or ((b or c) and (d or e))) [
]

使用运算符语法糖可该写为
if a and x.or b.or c and d.or e [
	二元 or 运算变成一元右结合运算
	这也意味着保留字 or 不能用作成员名称
]

运算符语法糖混合逗号分组视觉上更清晰
if a and, x or b.or c and d.or e [
]

同理使用加减法运算语法糖
var int x = i.add -1 mul 5 等同 (i-1)*5
```

二元运算符 'is', 'isnot' 有多种情况.

```
同类型时(null 值会依据右侧算子类型自动推导出同类型值)
a is b 等价 a == b
a isnot b 等价 a != b
int(0) is null 等价 int(0) == 0
'string' is null 等价 'string' == '' 

数组, 映射和自定义类型没有被赋值时值为 null, 所以
type empty =[]
type T=[
	int x
]

var [int] x,
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

Z 中运算符 'not', 'and', 'or' 是语法糖, 参见过程和分支语句.

#type

类型名是个标识符. 预声明类型标识符是保留字.

	保留字 type 可用来声明类型, 本身也是个类型, type 的类型就是 type.

##布尔

bool 的字面值为 true 或者 false.

##整数

可前置正负 '+','-' 号的 0-9 的十进制, 或 0x 开头的十六进制, 或 0b 开头的二进制表示.

预声明固定长度/尺寸整数类型有:

```
'
u8       无符号的 8 位整数 (0 to 255)
u16      无符号的 16 位整数 (0 to 65535)
u32      无符号的 32 位整数 (0 to 4294967295)
u64      无符号的 64 位整数 (0 to 18446744073709551615)

i8        带符号的 8 位整数 (-128 to 127)
i16       带符号的 16 位整数 (-32768 to 32767)
i32       带符号的 32 位整数 (-2147483648 to 2147483647)
i64       带符号的 64 位整数 (-9223372036854775808 to 9223372036854775807)

byte        和 u8 一样
rune        和 i32 一样
'
```

Z 中用 byte 或 rune 类型可代表单个字符. 可以用字符串对 byte 和 rune 类型进行赋值. 

```
var byte b = 'Hello Word' // b 的值为 'h'
var rune r = '世界你好'    // r 的值为 '世'
var rune e = ''           // 空字符串值为 0.
```

事实上这是由 string 类型的 toByte, toRune 方法进行的转换.

下列类型的长度/尺寸与运行环境有关:

```
'
uint     32 位或 64 位
int      和 uint 长度一样
'
```

##浮点数

可前置正负 '+','-' 号的十进制数字夹杂一个"." 或 "e"表示的十进制指数.

预声明的固定长度/尺寸浮点类型有下面这些:

```
'
f32     IEEE-754 32  位浮点数
f64     IEEE-754 64  位浮点数
f128    IEEE-754 128 位浮点数
'
```

下列类型的长度/尺寸与运行环境有关：

```
`
float    IEEE-754 32 或 64 位浮点数
`
```

##字符串

string 是一对单引号或者双引号包裹的多行文本.
字符串连接运算符使用 '+' 或者 '-' , 它们是等价的.

```
proc hello out string =[
	return 'hello'
]

proc word out string =[
	return 'word'
]

var string (
	a = '单引号字符串不支持反斜杠\转义值' -
		'直接断行也可以,
		换行会被保留, 前置空白会被剔除.' -
		"字符串可以多行混合. 这四行组成一个字符串值"

	b = '字符串还支持嵌入参数 $a $b'.['a'=1,'b'=2] -
		'SQL 参数序号风格 $1 $2'.[$a, $b]
		"变量值会被转换为字符串表示,
		这四行组成了一个字符串"

	c = hello() + word() 字符串连接可以使用 '+' 号

	d = hello() - word() 字符串连接也可以使用 '-' 号

	e = hello() - 尾注释会被正确处理
		' word' 前行尾部用 '-' 避免括号终结歧义, 否则就成注释了.

	f = '空白行会终止多字符串拼接' 此处不能用连接符连接下下行

        '上面的空白行产生分号, 赋值语句完结了'

	g = '连续多行字符串, 此行没有使用连接符'
        '这一行只是个注释'

)
```

字符 '+' 或者 '-' 都是字符串连接符, '-' 号更便于输入.
同一行内结果为字符串的表达式, 或引号包裹的连续多行字符串连接可省略 '+' , '-' 号.

双引号字符串具有简单的模板功能, 双引号字符串也称作模板字符串.

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

使用字符串预定义方法 `execute` 会剔除每一行两端的空白符和换行符.

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

datetime 表示日期和时间, 以 [ISO 8601][] 标准设计. 在源码中使用无 '-' 间隔年, 月, 日的基本格式.

```
var datetime(
	localdate = 20160204
	localtime = 20160204T21:49
	utcdate = 20160204Z // 尾部的 'Z' 表示时区为 0
	utctime = 20160204T21:49Z
	zonedate = 20160204T+08 // '+' 号两端是紧凑相连的
	zonetime =  21:49+08
	now = 20160204T214900 // 简化时间格式无 ':' 号
)
```

##数组

数组是预定义类型, 用前缀的方括号表示, 没有单独的保留字标识符.

```
var [int] a 看上去非常规, 录入却很自然 
var [int] ( 类型分组声明写法
	c
	d = [ 赋初值就直接写吧
		1,2,3,
		call() 使用表达式运算结果
	]
)

访问数组元素
var int e = d[0] 直接访问已有的变量

多维数组嵌套中括号即可

var [[int]] point =[
		[1,2,3]
	]
```

##映射

映射具有 key-value 结构, key 和 value 的具体类型有使用者定以.

```
给 map 赋值可以使用 JSON 风格, Z 编译器会检查值的合法性
var map[string]int m = {
	'age': 13, 
	"height": 156 有换行的话可以省略逗号
	"id": 1 
}

key 和 value 可以是任意类型
var map[type]string types = [ 万能的方括号
	string: 'string',
	int: 'integer'
]

甚至可以使用 type
var map[type]type assoc = [
	string: types,
	int: float
]

访问映射
var int e = m['age'] 支持直接访问已有的变量
var type t= m[string]
``` 

Z 是强静态类型语言, 编译后, 所有的类型都是明确的, 可识别的,
在 Z 中 type 的类型就是 type, 可以用于过程的参数类型.

#表达式

表达式具有以下特征

1. 必定产生运算结果.
2. 必定存在于语句中.

#语句

语句无返回值, 语句产生的代码执行顺序和书写顺序一致.

在 Z 源文件中最先出现的语句是注释或声明.

##注释

注释是对某个标识符或语句体的说明. 注释在 Z 中属于语句, Z 有三种注释写法.

1. 尾注释以 '//' 开始, 直到行尾
2. 块注释以三个以上的减号 '-' 位于一行非空白符之首开始, 并以此结束.
3. 字符串值注释, 孤立的字符串值也会被当做注释.

```
// 这是一条尾注释, 可以分成连续的多行书写. 
// 纵观整个文本, 无法确定这条注释属于哪个标识符

---
这是个块注释, 可以多行书写.
纵观整个文本, 无法确定这条注释属于哪个标识符
---

因为标识符由英文字符和数字组成
所以多字节字符开头的注释可以直接写

The line is a comment.
上一行英文也是注释, 原因见下文关于声明的描述.

此行注释位于 proc sum 上面, 但不属于 sum 的注释
proc sum int x,y,out int =[
	'sum 返回 x + y 的和' // 此行前部注释归属于标识符 proc sum
	return x + y // 返回 x+y
]

proc multiByteFriendly out bool =[
	使用多字节字符写注释....
	return true 就是这么直接
]
```

依惯例, 注释最普遍的作用是生成文档. 从这个角度出发, Z 会忽略那无法确定归属的注释. 能向前归属到某个语句的注释被称作注释, 否则被称作备注. 以免描述时产生歧义.

上例中, 字符串 'sum 返回 x + y 的和' 和 `return` 语句尾部的注释可以追溯归属. 其它的注释仅是备注.

Z 中的注释是后置的, 这和其它语言不同.

Z 中没有注释的注释, 所以注释 'sum 返回 x + y 的和' 尾部的注释仅是备注.

##赋值语句

赋值语句不存在优先级.

| 操作符 | 解释   |
|-------|--------|
| ++    |自赋值增一|
| --    |自赋值减一|
| =     |赋值     |

前例 `proc` 声明中的 '=' 也是赋值操作符, 值就是过程(函数)执行体本身.

	在 Z 中过程执行体是个值表达式, 它的运算符是 '执行'.
	只不过 '执行' 本身如果有返回值, 那么就是执行表达式, 否则就是执行语句.


```
proc run ={ 换个定界符
	var int x
	proc ={
		这是个匿名执行体, 无参数无返回值, 并且直接执行了
		x = 9
	}() 必须执行它, 不然就成了值表达式, 而不是语句

	proc fn ={
		具名过程
	}

	proc noret ={
		执行无返回值的具名过程 
	}()
	
	var int x = proc out int ={
		执行匿名过程, 返回值被接收
		return 1
	}()
	
	下面的写法很明显都是非法的
	var fn1 = proc ={} // fn1 前面缺少类型声明
	var proc fn2 ={}  // var proc? 这是什么鬼
	proc ret out int ={
		执行了具名过程, 返回值没被接收, 隐含了一个执行表达式
		表达式不能独立存在
		return 1
	}()
	proc ={
		匿名执行体, 未被赋值也未被执行, 隐含了一个过程表达式
		表达式不能独立存在
	}
}
```

Z 中的匿名过程必须被执行, 因为写不出合法的匿名过程赋值语句.

```
proc x out int =[
	return 5
	// Z 中的 return 是个语法糖. 它等同于两条语句
	// out = 5
	// end
]

proc noret =[
]

proc y =[
	x() 非法, 因为产生了运算结果是个表达式, 表达式不能独立存在.
	
	discard x() 可使用保留字 discard 丢弃结果形成语句.

	noret() 合法语句, 因为没有产生运算结果.
]
```

##分支语句

有两个标识符可选, 'if' 和 'switch'.

```
if expr {
	// do something
} else [ 也可以使用方括号包裹执行体
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
case b,c: 可选多值匹配
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
for map[string]int{'a':1,'b':2} as key intVal {
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
var map[type]string types ={
	string: 'string',
	int: 'integer'
}

var map[type]type assoc ={
	string: types,
	int: float
}

proc fn =[
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
'pub const var static func proc type'
```

很明显, 用纯英文写 Z 源码, 避开这些词写顶层注释非常容易.
如果行首用到这些词, 大写首字母或者随便加个非空白符号就行.

##导出声明

声明默认都是非导出的, 如果允许外部 block 访问, 可在声明前加保留字 pub 进行修饰.

##静态声明

静态声明标识符 'static' 可对变量, 'func', 'proc' 以及属性声明进行修饰.

```
static int x = 9 修饰变量总是省略 var
static datetime day = fn() 调用函数直接赋值
pub static string prefix 这个是导出的

proc oneday int interval  out datetime =[
	return datetime.now().add(interval)
]

proc init =[
	prefix = 'z' 随便在哪里进行赋值都行
]

proc fn ={
	static int x = datetime.time()
	过程内的静态声明
}

proc setName string val =[
	static string name

	if name is null {
		用 is null 可判定变量是否已经被赋值.
		这与 name == null 是等价的.
		name = val
	}

	name = val 最简单的写法, Z 会确保只赋值一次
]

proc setOnce string val =[
	static string name = val 最简洁的写法
]

声明时直接给属性赋值

type rep ={
	static string version = '0.0.0'    定值静态属性
	string host = 'https://github.com' 初值
}
```

##常量声明

常量只声明值, 常见的值类型有字符串, 数值, 布尔值, 时间值.

常量的值是字面值, 常量在使用时才能显现出类型以及合法性.

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
保留字 func 只声明名字和参数类型
并且, 参数中只写类型, 不能带参数名
多个参数间不使用逗号分隔
保留字 out 表示之后的参数被输出

func toInt string out int

proc toString out string s =[
	保留字 proc 声明带执行体的过程
	s = 'hello' - s
	end // end 指示终止过程
]

pub proc walk func callback int out int, out bool =[
	这个过程有两个参数, 第一个参数是个函数类型
	显然 func 参数无逗号的设计显现出作用
]

使用函数变量
proc fn =[
	var func f out string 声明函数变量用 func
	var func c = fn 声明并赋值
	
	f = proc out string ={ 匿名过程赋值
		out = 'hello'
	}
]
```

如果只用 func 声明了过程而没有实现代码, 那意味着:

```
pub func fn	在外部实现, 或在参数中作为类型约束

func call	来自外部连接库, 或被 block 内的过程输出
```

关于过程的参数:

1. 单纯的输入参数必须具有类型和名字
2. 输出参数只有一个的话, 可以省略参数名 
3. 参数都有零值或缺省值,  null 是通用的零值.
4. 调用过程不必传递所有参数.
5. out 后续的参数具有输入输出双向性.
6. 保留字 return 是个语法糖, 等同先设置输出参数再 end

```
proc one out int =[ 此例省略了唯一的输出参数名
	out = 1 用 out 替代唯一的输出参数
	return 2 这只是个语法糖, 等同
	out = 2
	end
]

proc fn out int x,y =[
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
	end 保留字 end 指示结束过程, 可以省略它

	多输出参数也可以使用 return 语法糖
	return {
		x:1,
		y:2
	}
]
```


##类型声明

使用保留字 type 声明新类型.
Z 中没有别名, 声明类型就是定义类型成员.

```
type empty =[] 无属性类型

type fruit ={ 带属性的水果类型
	bool clean true 属性初值, 水果都是无污染产品
	string brand,color
	flavors []string 多种味道
}

type apple ={
	use fruit[color='red'] 类型复合并且可以设置初值
}

成员方法是个过程, 这种过程可以叫做 '实例方法'
proc fruit.isSweet out bool =[
	return self.flavors isnot null and
		self.flavors.has('sweet')
]

匿名复合, 下例中的 path 导出了类型 string 所有的导出成员.
type path =[
	pub use string
]

给类型定义静态方法, 外部使用时直接用 path.fn('string')
这种过程可以叫做 '类型方法'
pub static proc path.last string dir =[
	return dir.split('/').slice(-1)
]

proc path.name =[
	实例方法调用类型方法也可以用 self
	if self == '/' {
		return '/'
	}else{
		return self.last(self)
		因为 path 只符合了一种类型, 可以直接使用 self
		否则需要使用 self.last(self.string)
	}
]
```

Z 中没有 class, 只是简单的类型复合.
在上例中保留字 self 代指 proc isSweet 的接收类型 fruit.
Z 支持几个特别的属性名可省略 self 进行访问, 当然这首先需要使用者声明这些属性

```
type node =[
	pub static name string 属性可以是静态的, 表示只赋值一次
	node parent
	[]node child
	root this
]

type root =[
	use node
]

proc root.count out int =[
	return child.length()
]

proc node.fn =[
	if this.count() [
		// do something
	]
]
```

注意: Z 只是对 this, parent, child 属性提供了便捷访问方式, 不关心, 也不明白它们与 self 的真正关系.


[ISO 8601]: https://en.wikipedia.org/wiki/ISO_8601