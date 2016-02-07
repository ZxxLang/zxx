
Z 目前只是个设想 

#简介

ZxxLang 是一门文档风格通用编程语言. 初衷是混合文档和代码于一体并降低 shift 按键的使用频率.  ZxxLang 通过忽略不能识别的文本, 并减少或替代和 shift 按键有关的符号来实现此目标.

	减少 shift 按键频率的方法可能会让一些开发者难于接受, 但这是 ZxxLang 固守的特征

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

Z 支持用输入方括号替代圆括号和花括号, 因为方括号不用 shift 键, 更便于输入. Z 提供的格式化工具会进行常规习惯转换. 

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

'block.json' 的采用语义化属性名, 通常这无需特别解释, 看到名称就知道作用.


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

字面值是就是一段无类型的文本. 它的语义类型是在使用时判定的, 有可能被判定为非法值.

##预声明字面值

```
`
null      无值, 零值
false     假
true      真
NaN       浮点类型 非数
infinite  浮点类型 无穷大
`
```

##格式字面值

包括  string, datetime, float, int 等格式的文本, 具体格式见下文类型描述.

##占位文本

Z 不能识别的文本, 而又不确定非法时都会被当做占位文本处理.

#标识符

标识符的正则为 `[a-zA-z_]+[a-zA-z_0-9]*`. 自定义标识符用于类型名, 常量名, 变量名, 过程名. 预定义标识符下文详述. 

#符号

##分号

分号 ";" 用来表示语句的结束边界. 源码经过格式美化后分号会被省略. 除非代码 minify 到一行. 

换行或者空行会依据上文的完整性推导为分号终结符. 推导规则:

	如果上文至此换行可形成完整的语句, 把换行当做分号.

*注意: 不良的换行可能引起下文非法, 甚至被当做注释.*

##定界符

定界符是一对符号, 表示代码实体的起始和终结. 有三对符号可用

	() 分组定界可用于很多地方, 参见示例
	{} 代码定界用于 block 的开始和结束
	[] 通用定界


方括号 '[]' 更便于输入, 输入时可以替代圆括号和花括号. Z 的格式美化工具会替换合适的符号.


##逗号

逗号 ',' 具有多种用途. 可用来分隔声明, 表达式以及语句.

```
var int x,y = 3+1,6     // 分隔声明变量, 分隔表达式
var []int array = [1,2] // 分隔值表达式
var (
	int i,string s      // 分隔声明变量
)

proc loop =[
	for var int x,y = 0,0; x < 10 ;x++, y = y*x[
		// 三段式 `for` 语句中第一, 第三段是语句.
		echo x,y
	]
]
```

逗号分组

```
if a and (b or c) [
]
// 使用逗号分组
if a and, b or c [
]

// 复杂的分组运算
if a and (x or ((b or c) and (d or e))) [
]
// 使用逗号分组无法消除所有分组括号
if a and, x or (b or c) and (d or e) [
]
// 下面的写法 Z 不能识别
if a and, x or, b or c, and, d or e [
]

```

##运算符

运算符只能出现在表达式中, 下表按优先级从高到低排列.

| 运算符                    | 解释           |
|--------------------------|----------------|
| $                        |仅用于模板字符串  |
| &,\|,~,xo                |位与,或,反,异或   |
| <<,>>,mod,mul,%,*,/      |                |
| add,+,-                  |                | 
| to                       |数值至数值       |
| ==,<=,>=,>,<             |比较运算         |
| not,in,is,has            |                |
| and                      |逻辑与           |
| or                       |逻辑或           |


#类型

类型名是个标识符, 预声明类型的标识符是保留字, 所有内置类型实例都是对象, 拥有预定义方法, 详见相关手册.

##布尔

bool 的字面值为 true 或者 false.

##整数

可用带符号的0-9组成的十进制, 0x开头的十六进制, 0b开头的二进制表示.

预声明的固定长度/尺寸数值类型有下面这些:

```
'
uint8       无符号的 8 位整数 (0 to 255)
uint16      无符号的 16 位整数 (0 to 65535)
uint32      无符号的 32 位整数 (0 to 4294967295)
uint64      无符号的 64 位整数 (0 to 18446744073709551615)

int8        带符号的 8 位整数 (-128 to 127)
int16       带符号的 16 位整数 (-32768 to 32767)
int32       带符号的 32 位整数 (-2147483648 to 2147483647)
int64       带符号的 64 位整数 (-9223372036854775808 to 9223372036854775807)

byte        和 uint8 一样
rune        和 int32 一样
'
```

Z 中用 byte 或 rune 类型可代表单个字符. 可以用字符串对 byte 和 rune 类型进行赋值. 

```
var byte b = 'Hello Word' // b 的值为 'h'
var rune r = '世界你好'    // r 的值为 '世'
var rune e = ''           // 空字符串赋值详见下文解释
```

下列类型的长度/尺寸与运行环境有关：

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
float32     IEEE-754 32  位浮点数
float64     IEEE-754 64  位浮点数
float128    IEEE-754 128 位浮点数
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

```
proc hello out string =[
	return 'hello'
]

proc word out string =[
	return 'word'
]

var string (
	a = '单引号字符串不支持反斜杠\转义值'
		'直接断行也可以,
		换行会被保留, 前置空白会被剔除.'
		"字符串可以多行混合. 这四行组成一个字符串值"

	b = '字符串还支持嵌入参数 $a $b'.['a'=1,'b'=2]
		'SQL 参数序号风格 $1 $2'.[$a, $b]
		"变量值会被转换为字符串表示,
		这四行组成了一个字符串"

	c = hello() word() // 行内字符串表达式值连接无需运算符

	d = hello()word() // 紧凑书写也可以

	e = hello() -
		' word' // 前行尾部用 '-' 避免括号终结歧义, 否则就成注释了.

	f = '空白行会终止多字符串拼接'

        '上面有个空白行, 变量 c 的值不包括这一行'
)
```

字符串连接操作使用运算符 '-', 因为 '+' 需要 shift 按键, '-' 号更便于输入.

双引号字符串具有简单的模板功能, 双引号字符串也称作模板字符串.

```
echo '{
		"name": "$1",
		"age": $2
	  }'.['tom', 8]
```

输出结果为:

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

输出

```JSON
{"name": "tom","age": 8}
```

##日期和时间

datetime 表示日期和时间, 以 [ISO 8601][] 标准设计. 在源码中使用无 '-' 间隔年, 月, 日的基本格式.

```
var datetime(
	localdate = 20160204
	localtime = 20160204T21:49
	utcdate = 20160204Z
	utctime = 20160204T21:49Z
	zonedate = 20160204T+08 // '+' 号两端是紧凑相连的
	zonetime =  21:49+08
	now = 20160204T214900 // 简化时间格式无 ':' 号
)
```

#表达式

表达式具有以下特征

1. 必定产生运算结果.
2. 必定存在于语句中.

#语句

语句无返回值, 在 Z 源文件中最先出现的语句是注释或声明.

```
proc x out int =[
	return 5
	// Z 中的 return 是个语法糖. 它等同于两条语句
	// result = 5
	// end
]

proc noret =[
]

proc y =[
	x() // 这是非法的.
	// 因为产生了运算结果, 虽然被抛弃了可还是个表达式, 表达式不能独立存在.

	discard x() // 这是一条 discard 语句, Z 不会自动添加 discard.

	noret() // 合法语句, 因为没有产生运算结果
]
```

##注释

注释是对某个标识符或语句体的说明. 注释在 Z 中是语句, Z 有三种注释写法.

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
proc sum int x,y,out int =[
	'sum 返回 x + y 的和' // 此行注释归属于标识符 proc sum
	return x + y // 返回 x+y
]

proc multiByteFriendly out bool =[
	使用多字节字符写注释....
	return true 就是这么直接
]
```

依惯例, 注释最普遍的作用是生成文档. 从这个角度出发, Z 会忽略那无法确定归属的注释. 能向前归属到某个语句的注释被称作注释, 否则被称作备注. 以免描述时产生歧义.

上例中, 字符串 'sum 返回 x + y 的和' 和 `return` 语句尾部的注释可以追溯归属. 其它的注释仅是备注.

Z 中没有注释的注释, 所以注释 'sum 返回 x + y 的和' 尾部的注释仅是备注.

Z 中的注释是后置的, 这和其它语言不同.



无返回值的函数调用可以形成独立的语句.
同理 Z 中匿名函数是匿名函数表达式, 因为它产生了结果.

很明显, 语句可以顺序排列或者嵌套形成语句块.

Z 源文件的代码块总是以下列关键字开始嵌套子代码块.

```
'var static const func type def proc'
```

子代码块可使用的语句更丰富. 

```
type hello is example for Z =[
	string word
]
```

子代码块被包裹在方括号/圆括号中. 此例中' is example for Z ' 被忽略了,
因为在 `type` 语句中 Z 期待定界符 '=', 它之前不能识别的文本被忽略. 此例只是说明 Z 是怎么做的, 不是要推荐这样的写法.

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
proc run ={ // 换个定界符
	var int x
	proc ={
		// 这是个匿名执行体, 无参数无返回值, 并且直接执行了
		x = 9
	}() // 必须执行它, 不然就成了值表达式, 而不是语句

	proc fn ={
		// 具名过程
	}

	proc noret ={
		// 执行无返回值的具名过程 
	}()
	
	var int x = proc out int ={
		// 执行匿名过程, 返回值被接收
		return 1
	}()
	
	// 下面的写法很明显都是非法的
	var fn1 = proc ={} // fn1 前面缺少类型声明
	var proc fn2 ={}  // var proc? 这是什么鬼
	proc ret out int ={
		// 执行了具名过程, 返回值没被接收, 隐含了一个执行表达式
		// 表达式不能独立存在
		return 1
	}()
	proc ={
		// 匿名执行体, 未被赋值也未被执行, 隐含了一个过程表达式
		// 表达式不能独立存在
	}
}
```

Z 中的匿名过程必须被执行, 因为写不出合法的匿名过程赋值语句.

#静态声明

静态标识只进行一次赋值操作.

```
static int x = 9
static datetime day = fn() // 调用函数直接赋值
static string prefix

proc oneday int interval  out datetime =[
	return datetime.now().add(interval)
]

proc init =[
	prefix = 'z' // 随便在哪里进行唯一的一次赋值
]

proc fn ={
	static int x = datetime.time()
	// 过程内的静态声明
}

proc setName string val =[
	static string name

	if name is null {
		// 用 is null 可知道变量是否经过赋值.
		// 这与 name == null 是等价的.
		name = val
	}

	name = val
	// 事实上可以这么用, Z 会确保第一次赋值生效
]

proc setOnce string val =[
	static string name = val
	// 最简洁的写法
]

// 声明时直接赋值

type rep ={
	static string version = '0.0.0'    // 定值静态属性
	string host = 'https://github.com' // 初值
}
```

##常量声明

常量类似静态, 不同的是常量不声明类型, 只声明值, 可用的值类型有字符串, 数值, 布尔值, 时间值.

常量的值是字面值, 因此常量没有明确的类型.

```
const (
	CRLF = "\r\n"
	ZERO = 0
	DAY  = 20110101T
	ONE  = 1.0
	YES  = true
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
	byte   b = CRLF
	string s = CRLF
	rune   r = CRLF
)
```

##类型声明

```
type 
```

#过程声明

函数, 类型方法统称为过程

#表达式

表达式的作用是运算并产生结果, 由标识符和运算符组成.

#语句


##if 语句

```
if expr {
}else  {
}else{
}
```

```
`
add and as asm atomic
break
case const continue
def defer
else end
for from func
go goto
if in iota is
not null
of or out
proc pub
range ref return
shl shr static sub
template to trait type
use
var
xor
yield
`
```

#声明和作用域

#保留字

保留字不能被重新声明. 上文中预声明的值和类型名都是保留字.

以下关键字或类型名被保留. 

```
'
add and as asm atomic
break
case const continue
def defer
else end
for from func
go goto
has
if in is
mul
not null
or out
proc pub
range ref return
static sub
template to trait type
use
var
xor
yield
int128 uint128
'
```



[ISO 8601]: https://en.wikipedia.org/wiki/ISO_8601