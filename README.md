# Zxx

WIP

Zxx 是静态类型编程语言, 双空格缩进, 单行注释, 没有 `HTAB`.

    初衷是降低书写疲劳

[语法规格][Zxx-spec]

相关工具:

- [ZxxSublime][] for SublimeText 3

演示:

```zxx
使用外部模块
use
  os
  path/module   mod ; 注释以 "; " 开始至行尾

顶级对象, 类型推导
let
  name = 'Zxx'                ; string
  bytes= [byte]('Lang')       ; list
  since= "20160202T22:48:33"  ; time

  one  = 1                    ; i32
  two  = 02                   ; u32
  Π    = 3.1415               ; float

  list = [int][               ; 续行
    1,2,
    3
  ]

  map  = {string}{
    list  'list',
    map   'map',
  }

  _ = list is [int] and
    list.length is 0 and list is not null

  ; 多变量赋值
  min, max = 0, 1000

  ; iota 计数
  _,KB,MB,GB,TB,PB,EB,ZB,YB = filesize(1 << 10*iota)

别名
def filesize u64

结构体
def file
  name string
  size filesize
  func fun(int)

接口
def humanString

方法
fun humanString.human(out string)

fun filesize.string(out string, self dd)
  if self >= YB
    out 'YB+'

函数
fun sum(
  int x,
  int y,
  out int
)
  out x + y

斐波那契数列 F0=0, F1=1 ... Fn=Fn-1 + Fn-2

短路字面值递归法 Fn
fun fib(u32 n, out u32)
  out n or 0 or
    n is 1 and 1 or
    self(n - 1) + self(n - 2)

循环法 Fn
fun fibonacci(u32 n, out u32 Fn)
  Fn-1 = 00
  for _ of 0..n ; 0 <= _ < n
    Fn, Fn-1 = Fn + Fn-1, Fn
  fun.name

yield 斐波那契数列
fun fibSeries(u32 n, yield u32 Fn)
  Fn-1 = 00
  for _ of 0..n
    yield
    Fn, Fn-1 = Fn + Fn-1, Fn
```

## LICENSE

[BSD 2-Clause License](https://github.com/ZxxLang/zxx/blob/master/LICENSE)

Copyright (c) 2018 The Zxx Authors All rights reserved.

[ZxxSublime]: https://github.com/ZxxLang/ZxxSublime
[Zxx-spec]: https://github.com/ZxxLang/zxx/blob/master/Zxx-spec.md