package scikits

import (
	"context"
	"fmt"
	es "github.com/olivere/elastic/v7"
	"log"
	"os"
	"time"
)

func NewEsClient(label string) (*es.Client, error) {
	url := MyViper.GetString(fmt.Sprintf("%s.host", label))
	user := MyViper.GetString(fmt.Sprintf("%s.user", label))
	pass := MyViper.GetString(fmt.Sprintf("%s.pass", label))
	// 创建Client, 连接ES
	client, err := es.NewClient(
		// Go无法连接docker中es，代码设置sniff 为false
		es.SetSniff(false),
		// es 服务地址，多个服务地址使用逗号分隔
		es.SetURL(url),
		// 基于http base auth验证机制的账号和密码
		es.SetBasicAuth(user, pass),
		// 启用gzip压缩
		es.SetGzip(true),
		// 设置监控检查时间间隔
		es.SetHealthcheckInterval(10*time.Second),
		// 设置请求失败最大重试次数
		es.SetMaxRetries(5),
		// 设置错误日志输出
		es.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		es.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))

	return client, err
}

func EsQueryByMatch(label, index, column, text string) *es.SearchResult {
	// label是settings.ini中es的连接配置的标签
	// column是es中要查询的字段名称， text是输入的检索内容
	client, error := NewEsClient(label)
	fmt.Println(error)
	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	matchQuery := es.NewMatchQuery(column, text)

	searchResult, _ := client.Search().
		Index(index).      // 设置索引名
		Query(matchQuery). // 设置查询条件
		//Sort("Created", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(0).      // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(40).     // 设置分页参数 - 每页大小
		Pretty(true). // 查询结果返回可读性较好的JSON格式
		Do(ctx)       // 执行请求

	return searchResult
}
