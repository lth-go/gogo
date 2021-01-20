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

+ 块重复声明修正,引用全局声明
+ 增加package声明
+ 静态区垃圾回收
+ 操作符左右类型修正
+ 导入包时,将变量放入静态区
+ 引入包概念, 能使用其他包的函数,变量
