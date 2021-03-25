package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
	"time"
)

type Comment struct {
	gorm.Model
	UserId           uint
	ArticleId        uint
	CommentLikeCount int `gorm:"default:0"`
	CommentDate      time.Time
	CommentContent   string
	ParentCommentId  int
	Zans             []Zan
	Floor            int
}

func CreateComment(article_id uint, user_id uint, comment_content string, floor int) {
	comment := Comment{
		UserId:         user_id,
		ArticleId:      article_id,
		CommentContent: comment_content,
		Floor:          floor,
		CommentDate:    time.Now(),
	}
	db.W_Db.Create(&comment)
}

func FindCommentByArticle(article_id string, page int, articleCommentCount int) []Comment {
	var comments []Comment
	db.W_Db.Order("created_at asc").Limit(articleCommentCount).Offset((page-1)*articleCommentCount).Where("article_id = ?", article_id).Find(&comments)
	//db.W_Db.Where("article_id = ?", article_id).Find(&comments)
	return comments
}

func CommentCount(article_id string) int {
	var count int
	db.W_Db.Table("comments").Where("deleted_at IS NULL and article_id = ?", article_id).Count(&count)
	return count
}

func FindCommentById(comment_id string) Comment {
	var comment Comment
	db.W_Db.Where("id=?", comment_id).Find(&comment)
	return comment
}

func UpdateComment(comment Comment, update_map map[string]interface{}) {
	db.W_Db.Model(&comment).Updates(update_map)
}

//通过评论发现点赞表
func (comment *Comment) FindZansByComment() []Zan {
	var zans []Zan
	db.W_Db.Model(&comment).Related(&zans)
	fmt.Println(zans)
	return zans
}

//通过评论发现点赞数量
func (comment *Comment) FindTrueZansCountByComment() int {
	var count int
	db.W_Db.Table("zans").Where("deleted_at IS NULL and  is_zan = ? and comment_id = ?", true, comment.ID).Count(&count)
	return count
}

//通过用户判断是否有点赞
func (comment *Comment) FindUserIsZan(user User) bool {
	var user_is_zan bool
	var count int
	db.W_Db.Table("zans").Where("deleted_at IS NULL and  is_zan = ? and comment_id = ? and user_id = ?", true, comment.ID, user.ID).Count(&count)
	if count > 0 {
		user_is_zan = true
	} else {
		user_is_zan = false
	}
	return user_is_zan
}
