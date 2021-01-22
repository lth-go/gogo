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

+ 常量fix修正
+ 循环导入校验
+ 代码解耦
+ 全局常量支持布尔值
+ 增加printf函数
+ 增加len,append函数
+ 增加...
+ 增加类型interface{}
+ int 改int64
+ 移除CastExpression
+ 优化包代码
+ 增加map类型
+ 增加struct类型
+ 增加struct方法
