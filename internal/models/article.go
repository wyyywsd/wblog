package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
	"log"
	"time"
)

type Article struct {
	gorm.Model
	UserId uint
	ArticleTitle string
	ArticleContent string
	ArticleViews int `gorm:"default:0"`
	ArticleCommentCount int `gorm:"default:0"`
	ArticleDate time.Time
	ArticleLikeCount int `gorm:"default:0"`
	Comments      []*Comment `gorm:"ForeignKey:ArticleId"`
	Sorts         []*Sort    `gorm:"many2many:artitle_sorts;"`
	Labels		  []*Label `gorm:"many2many:artitle_labels;"`
	IsPublic bool  `gorm:"default:true"`
}

func CreateArticle(article_title string , article_content string, user_id uint, label Label, is_public bool){
	//article := models.FindArticleById("1")
	//label := models.FindLabelById("1")
	//db.W_Db.Model(&article).Association("Labels").Append(label)
	article := Article{
		UserId: user_id,
		ArticleTitle: article_title,
		ArticleContent: article_content,
		ArticleDate: time.Now(),
	}
	db.W_Db.Create(&article)
	//如果设置了私密 就更新为私密
	if !is_public{
		db.W_Db.Model(&article).Select("is_public").Updates(map[string]interface{}{"is_public": is_public})
	}
	db.W_Db.Model(&article).Association("Labels").Append(label)
}

func UpdateArticle(article Article,update_map map[string]interface{})error{
	err := db.W_Db.Model(&article).Updates(update_map).Error
	fmt.Print("789787978797897879787987978778978797879879787879787",update_map)
	return err
}

func FindArticleById (id string) Article {
	article := Article{}
	db.W_Db.Where("id = ?",id).First(&article)
	return article
}
//获取所有的文章
func AllPublicArticles () ([]*Article,error) {
	return  _ListArticle(true)
}

func _ListArticle(is_public bool)([]*Article,error) {
	var articles []*Article
	var err error
	if is_public {
		err = db.W_Db.Where("is_public = ?", true).Find(&articles).Error
	}else {
		err = db.W_Db.Find(&articles).Error
	}
	return articles,err
}

//分页获取所有公开的文章
func PublicArticleLimit(page int,articleCount int)([]*Article , error){
	var articles []*Article
	var err error

	err = db.W_Db.Limit(articleCount).Offset((page-1)*articleCount).Where("is_public = ?", true).Find(&articles).Error

	return articles,err
}
//获取文章的数量
func ArticleCount()int{
	var count int
	db.W_Db.Table("articles").Where("deleted_at IS NULL and is_public = ?", true).Count(&count)

	fmt.Println("******************************************************************",count)
	return count
}
////获取某个标签文章的数量
//func ArticleLabelCount(label_id string)int{
//	var count int
//	var articles []*Article
//	db.W_Db.Table("articles").Where("is_public = ?", true).Count(&count)
//	db.W_Db.Raw("select articles.*  from articles inner join artitle_labels on " +
//		"artitle_labels.article_id = articles.id where artitle_labels.label_id = ? and articles.is_public = ?",label_id,true).Scan(&articles).Count(&count)
//	fmt.Println("******************************************************************",count)
//	return count
//}

//通过标签寻找文章
func FindArticlesByLabel(label_id string) ([]*Article) {
	var articles []*Article
	//label := FindLabelById(label_id)
	log.Println("++++++++++")
	//db.W_Db.Model(&label).Association("Labels").Find(&articles)
	//db.W_Db.Table("articles").Select("articles.id, labels.label_id").Joins("left join articles on labels.article_id = articles.id").Where("label.id = ?",label_id).Find(&articles)
	//select articles.*  from articles right  join artitle_labels on artitle_labels.article_id = articles.id;
	db.W_Db.Raw("select articles.*  from articles inner join artitle_labels on " +
		"artitle_labels.article_id = articles.id where artitle_labels.label_id = ? and articles.is_public = ?",label_id,true).Scan(&articles)
	return articles
}
//通过文章发现标签
func (article *Article)FindLabelsByArticle() Label {
	var label Label
	article_id := article.ID
	//label := FindLabelById(label_id)
	log.Println("++++++++++")
	//db.W_Db.Model(&label).Association("Labels").Find(&articles)
	//db.W_Db.Table("articles").Select("articles.id, labels.label_id").Joins("left join articles on labels.article_id = articles.id").Where("label.id = ?",label_id).Find(&articles)
	//select articles.*  from articles right  join artitle_labels on artitle_labels.article_id = articles.id;
	db.W_Db.Raw("select labels.*  from labels inner join artitle_labels on " +
		"artitle_labels.label_id = labels.id where artitle_labels.article_id = ?",article_id).Scan(&label)
	return label

}

//分页获取当前用户的所有文章
func UserArticleLimit(page int,articleCount int,user_id uint)([]*Article , error){
	var articles []*Article
	var err error

	err = db.W_Db.Limit(articleCount).Offset((page-1)*articleCount).Where("user_id = ?", user_id).Find(&articles).Error

	return articles,err
}
//获取当前用户文章的数量
func UserArticleCount(user_id uint)int{
	var count int
	db.W_Db.Table("articles").Where("deleted_at IS NULL and  user_id = ?", user_id).Count(&count)
	fmt.Println("******************************************************************",count)
	return count
}