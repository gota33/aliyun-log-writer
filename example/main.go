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
		MessageModifier: modifiers.Chain{
			modifiers.RemapLevelToSysLog(), // 将 slog 的 level 字段重写成 syslog 格式
			modifiers.RenameMessageField(), // 将 slog 的 msg 字段重命名成 message
		},
		MessageFilter: filters.Chain{
			filters.InfoLevel(), // 只推送 Info 及以上 Level 到服务器
		},
	})
	if err != nil {
		panic(err)
	}

	// 进程结束前把存量日志推送到阿里云
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
