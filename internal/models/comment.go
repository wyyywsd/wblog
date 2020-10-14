package models

import (
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
	"time"
)

type Comment struct {
	gorm.Model
	UserId uint
	ArticleId uint
	CommentLikeCount int `gorm:"default:0"`
	CommentDate time.Time
	CommentContent string
	ParentCommentId int
}

func CreateComment(article_id uint,user_id uint,comment_content string){
	comment := Comment{
		UserId: user_id,
		ArticleId: article_id,
		CommentContent: comment_content,
		CommentDate: time.Now(),
	}
	db.W_Db.Create(&comment)
}

func FindCommentByArticle(article_id uint) []Comment{
	var comments []Comment
	db.W_Db.Where("article_id = ?", article_id).Find(&comments)
	return comments
}