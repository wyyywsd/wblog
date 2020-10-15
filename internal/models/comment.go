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
	Floor int
}

func CreateComment(article_id uint,user_id uint,comment_content string,floor int){
	comment := Comment{
		UserId: user_id,
		ArticleId: article_id,
		CommentContent: comment_content,
		Floor: floor,
		CommentDate: time.Now(),
	}
	db.W_Db.Create(&comment)
}

func FindCommentByArticle(article_id string,page int,articleCommentCount int) []Comment{
	var comments []Comment
	db.W_Db.Limit(articleCommentCount).Offset((page-1)*articleCommentCount).Where("article_id = ?", article_id).Find(&comments)
	//db.W_Db.Where("article_id = ?", article_id).Find(&comments)
	return comments
}

func CommentCount(article_id string)int{
	var count int
	db.W_Db.Table("comments").Where("deleted_at IS NULL and article_id = ?", article_id).Count(&count)
	return count
}


