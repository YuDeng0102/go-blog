package rabbitmq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/global"
	"server/model/elasticsearch"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"go.uber.org/zap"
)

func StartESUpdateConsumer() {
	conn := global.RmqConn
	ch, err := conn.Channel()
	if err != nil {
		global.Log.Error("RabbitMQ连接失败:", zap.Error(err))
		return
	}
	_, err = ch.QueueDeclare(
		"es_update_queue", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)

	if err != nil {
		global.Log.Error("队列声明失败:", zap.Error(err))
		return
	}

	msgs, err := ConsumeMessages(conn, "es_update_queue")
	if err != nil {
		log.Fatal("消费消息失败:", err)
	}

	// 批量处理配置
	batchSize := 100                 // 每批处理100条
	flushInterval := 5 * time.Second // 最多5秒刷一次
	var batch []ESEvent

	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case msg := <-msgs:
			var event ESEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				// log.Printf("消息解析失败: %v", err)
				global.Log.Error("消息解析失败:", zap.Error(err))
				msg.Nack(false, false) // 丢弃无效消息
				continue
			}
			batch = append(batch, event)
			msg.Ack(false) // 确认消息

			// 达到批量大小立即处理
			if len(batch) >= batchSize {
				processBatch(batch)
				batch = batch[:0]
				ticker.Reset(flushInterval)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				processBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// 批量处理函数
func processBatch(batch []ESEvent) {
	// 按文章ID和字段分组聚合
	updates := make(map[string]map[string]int)
	for _, event := range batch {
		if _, ok := updates[event.ArticleID]; !ok {
			updates[event.ArticleID] = make(map[string]int)
		}
		updates[event.ArticleID][event.Field] += event.Delta
	}

	// 构建批量请求体
	var buf bytes.Buffer
	for articleID, fields := range updates {
		// --- 第一行：操作元数据 ---
		meta := map[string]interface{}{
			"update": map[string]interface{}{
				"_index": elasticsearch.ArticleIndex(), // 替换为你的索引名
				"_id":    articleID,
			},
		}
		metaBytes, _ := json.Marshal(meta)
		buf.Write(metaBytes)
		buf.WriteByte('\n')

		// --- 第二行：脚本和数据 ---
		scriptParts := []string{}
		params := make(map[string]int)
		for field, delta := range fields {
			scriptParts = append(scriptParts, fmt.Sprintf("ctx._source.%s += params.%s", field, field))
			params[field] = delta
		}

		body := map[string]interface{}{
			"script": map[string]interface{}{
				"source": strings.Join(scriptParts, ";"),
				"lang":   "painless",
				"params": params,
			},
		}
		bodyBytes, _ := json.Marshal(body)
		buf.Write(bodyBytes)
		buf.WriteByte('\n')
	}

	// 执行批量操作
	_, err := global.ESClient.Bulk().Raw(bytes.NewReader(buf.Bytes())).Index(elasticsearch.ArticleIndex()).Refresh(refresh.True).Do(context.TODO())
	if err != nil {
		global.Log.Error("批量更新失败:", zap.Error(err))
		return
	}
}
