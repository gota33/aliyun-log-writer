package sls

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gota33/aliyun-log-writer/internal/validator"
)

const (
	DefaultBufferSize = 100
	DefaultTimeout    = 1 * time.Second
	DefaultInterval   = 3 * time.Second
)

type Config struct {
	// 阿里云日志接入地址, 格式: "<region>.log.aliyuncs.com",
	// 例如: "cn-hangzhou-intranet.log.aliyuncs.com",
	// 更多接入点参考: https://help.aliyun.com/document_detail/29008.html?spm=a2c4g.11174283.6.1118.292a1caaVMpfPu
	Endpoint        string
	AccessKey       string          // 密钥对: key
	AccessSecret    string          // 密钥对: secret
	Project         string          // 日志项目名称
	Store           string          // 日志库名称
	Topic           string          // 日志 __topic__ 字段
	Source          string          // 日志 __source__ 字段, 可选, 默认为 hostname
	BufferSize      int             // 本地缓存日志条数, 可选, 默认为 100
	Timeout         time.Duration   // 写缓存最大等待时间, 可选, 默认为 1s
	Interval        time.Duration   // 缓存刷新间隔, 可选, 默认为 3s
	HttpClient      *http.Client    // HTTP 客户端, 可选, 默认为 http.DefaultClient
	MessageModifier MessageModifier // 在发送前编辑日志内容, 可选, 默认为空
	MessageFilter   MessageFilter   // 在发送前过滤日志内容, 可选, 默认为空
	OnError         ErrorListener   // 错误回调, 可选, 默认为空
	UseHttps        bool            // 是否在调用 PutLogs 时使用 Https, 可选, 默认为 false
	uri             *url.URL
}

func (c *Config) validate() (err error) {
	if err := errors.Join(
		validator.Required("Endpoint", c.Endpoint),
		validator.Required("AccessKey", c.AccessKey),
		validator.Required("AccessSecret", c.AccessSecret),
		validator.Required("Project", c.Project),
		validator.Required("Store", c.Store),
		validator.Required("Topic", c.Topic),
	); err != nil {
		return err
	}

	source, _ := os.Hostname()
	c.Source = validator.Coalesce(c.Source, source)
	c.BufferSize = validator.Coalesce(c.BufferSize, DefaultBufferSize)
	c.Timeout = validator.Coalesce(c.Timeout, DefaultTimeout)
	c.Interval = validator.Coalesce(c.Interval, DefaultInterval)

	if c.HttpClient == nil {
		c.HttpClient = http.DefaultClient
	}

	protocol := "http"
	if c.UseHttps {
		protocol = "https"
	}

	c.uri, err = url.Parse(fmt.Sprintf(
		"%s://%s.%s/logstores/%s/shards/lb", protocol, c.Project, c.Endpoint, c.Store))
	if err != nil {
		return validator.IllegalArgument("Endpoint", err.Error())
	}
	return
}
