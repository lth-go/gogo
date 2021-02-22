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

+ 类型转换修正
+ for range
+ 代码解耦
+ int改int64
+ 增加struct方法
+ 增加指针

函数调用时栈

```
9 -- sp
8 func index
7 arg3
6 arg2
5 arg1
4 r2
3 r1
```

函数调用后的栈

```
13 -- sp
12 local2
11 local2
10 local1
9  local1
8  callInfo -- base
7  arg3
6  arg2
5  arg1
4  r2
3  r1
```
