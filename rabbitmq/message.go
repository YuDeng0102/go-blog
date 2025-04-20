package rabbitmq

type ESEvent struct {
	ArticleID string `json:"article_id"` // 文章ID
	Field     string `json:"field"`      // 更新字段（views/comments/likes）
	Delta     int    `json:"delta"`      // 变化量（+1/-1）
}
