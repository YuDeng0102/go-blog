package database

import (
	"github.com/gofrs/uuid"
	"server/global"
)

// Comment 评论表
type Comment struct {
	global.MODEL
	ArticleID string    `json:"article_id"` // 文章 ID
	PID       *uint     `json:"p_id"`       // 父评论 ID
	PComment  *Comment  `json:"-" gorm:"foreignKey:PID"`
	Children  []Comment `json:"children" gorm:"foreignKey:PID"`                  // 子评论
	UserUUID  uuid.UUID `json:"user_uuid" gorm:"type:char(36)"`                  // 用户 uuid
	User      User      `json:"user" gorm:"foreignKey:UserUUID;references:UUID"` // 关联的用户
	Content   string    `json:"content"`                                         // 内容
}

// TODO 创建和删除评论时需要更新文章评论数
