# Aliyun Log Writer

[![godoc reference](https://godoc.org/github.com/gota33/aliyun-log-writer?status.svg)](https://godoc.org/github.com/gota33/aliyun-log-writer)

此 Writer 用于将通过 Go 1.21 及以上版本 slog 记录的 JSON 格式日志发送到阿里云日志服务.

特点:

- 采用非阻塞设计, 由一个后台线程将日志批量刷到远端日志库.
- 采用轻量级设计, 直接使用 [PutLogs](https://help.aliyun.com/document_detail/29026.html) 接口,
  不依赖于 [阿里云SDK](github.com/aliyun/aliyun-log-go-sdk)
- 除了 slog 也适用于推送其他日志库所记录的 JSON 格式日志

## 安装

`go get -u github.com/gota33/aliyun-log-writer`

## 使用指南

```go
package main

import (
  "log"
  "log/slog"
  "os"
  "time"

  sls "github.com/gota33/aliyun-log-writer"
  "github.com/gota33/aliyun-log-writer/filters"
  "github.com/gota33/aliyun-log-writer/handlers"
  "github.com/gota33/aliyun-log-writer/modifiers"
)

func main() {
  // DEBUG 模式开关
  // sls.SetDebug(true)

  // Config 字段详细说明请看 Config 的字段注释  
  writer, err := sls.New(sls.Config{
    Endpoint:     os.Getenv("APP_ENDPOINT"),                                  // 阿里云日志接入端点, 如: cn-hangzhou-intranet.log.aliyuncs.com
    AccessKey:    os.Getenv("APP_KEY"),                                       // App Key
    AccessSecret: os.Getenv("APP_SECRET"),                                    // App 密钥
    Project:      os.Getenv("APP_PROJECT"),                                   // 日志 Project
    Store:        os.Getenv("APP_STORE"),                                     // 日志 Store
    Topic:        os.Getenv("APP_TOPIC"),                                     // 日志 topic 字段
    Source:       "",                                                         // 日志 source 字段, 默认为 Hostname
    Timeout:      10 * time.Second,                                           // Push 数据超时时间, 默认 1s
    OnError:      func(err error) { log.Printf("[EXAMPLE] error: %s", err) }, // 错误回调
    // 用可选的 Modifier 编辑日志内容
    MessageModifier: modifiers.Chain{
      modifiers.RemapLevelToSysLog(), // 将 slog 的 level 字段重写成 syslog 格式
      modifiers.RenameMessageField(), // 将 slog 的 msg 字段重命名成 message
    },
    // 用可选的 Filter 过滤要发送的日志
    MessageFilter: filters.Chain{
      filters.InfoLevel(), // 只推送 Info 及以上 Level 到服务器
    },
  })
  if err != nil {
    panic(err)
  }

  // 进程结束前把剩余日志推送到阿里云
  defer func() { _ = writer.Close() }()

  // remoteHandler 将日志写入阿里云
  remoteHandler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
    Level:     slog.LevelInfo,
    AddSource: false,
  })

  // localHandler 将日志打印到本地控制台
  localHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level:     slog.LevelDebug,
    AddSource: true,
  })

  // 串联两个 Handler, 使 DEBUG 级别及以上打印到控制台, INFO 级别及以上推送到阿里云
  handler := handlers.MultiHandler(remoteHandler, localHandler)

  // 一些 slog 的常见用法
  logger := slog.New(handler).With("service", "demo")

  logger.Info("Aliyun Log Writer 1111111111111")

  time.Sleep(1 * time.Second)

  logger.Debug("Second Log")

  group := logger.WithGroup("group")
  group.Info("array test", "arr", []int{1, 2, 3})
  group.Info("map test", "map", map[string]string{"field": "value"})
}

```

## 依赖项

```
.
  ├ google.golang.org/protobuf
  └ github.com/pierrec/lz4
```

