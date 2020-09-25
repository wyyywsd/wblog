package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Comment struct {
	gorm.Model
	UserId uint
	ArticleId uint
	CommentLikeCount int
	CommentDate time.Time
	CommentContent string
	ParentCommentId int
}