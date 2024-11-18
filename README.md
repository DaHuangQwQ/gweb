# gweb
an easy go web framework

```shell
go get github.com/DaHuangQwQ/gweb
```
## 示例
```golang
package main

import (
	"github.com/DaHuangQwQ/gweb"
	"github.com/DaHuangQwQ/gweb/context"
)

func main() {
	server := gweb.Default()

	server.Get("/", func(ctx *context.Context) {
		_ = ctx.RespJSONOK("hello world")
		return
	})

	err := server.Start(":8081")
	if err != nil {
		panic(err)
	}
}
```

## route tree
- 静态匹配
- 通配符匹配
- 参数路径

## context
- 处理输出
- 处理输入

## file
- 文件上传
- 文件下载

## session
- 基于内存实现
- 基于 redis 实现
- 基于 cookie 实现

## template
Go template

## AOP
1. access log
2. err handler
3. opentelemetry
4. prometheus
5. recover
