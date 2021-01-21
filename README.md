# GoGo

Simple language writen by go.

## Require

```sh
go get -u golang.org/x/tools/cmd/goyacc

# generate parser.go
make parser
```

## Test

```sh
make test
```

## TODO

+ 增加printf函数
+ 增加...
+ 增加类型interface{}
+ int 改int64
+ 增加原生函数int, float, string
+ 移除CastExpression
+ 块重复声明修正,引用全局声明
+ 增加package声明
+ 导入包时,将变量放入静态区
+ 引入包概念, 能使用其他包的函数,变量
