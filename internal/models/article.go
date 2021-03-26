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
	UserId              uint
	ArticleTitle        string
	ArticleContent      string
	ArticleViews        int `gorm:"default:0"`
	ArticleCommentCount int `gorm:"default:0"`
	ArticleDate         time.Time
	ArticleLikeCount    int        `gorm:"default:0"`
	Comments            []*Comment `gorm:"ForeignKey:ArticleId"`
	Sorts               []*Sort    `gorm:"many2many:artitle_sorts;"`
	Labels              []*Label   `gorm:"many2many:artitle_labels;"`
	IsPublic            bool       `gorm:"default:true"`
}

func CreateArticle(articleTitle string, articleContent string, userId uint, label Label, isPublic bool) {
	//article := models.FindArticleById("1")
	//label := models.FindLabelById("1")
	//db.W_Db.Model(&article).Association("Labels").Append(label)
	article := Article{
		UserId:         userId,
		ArticleTitle:   articleTitle,
		ArticleContent: articleContent,
		ArticleDate:    time.Now(),
	}
	db.W_Db.Create(&article)
	//如果设置了私密 就更新为私密
	if !isPublic {
		db.W_Db.Model(&article).Select("is_public").Updates(map[string]interface{}{"is_public": isPublic})
	}
	db.W_Db.Model(&article).Association("Labels").Append(label)
}

func UpdateArticle(article Article, update_map map[string]interface{}) error {
	err := db.W_Db.Model(&article).Updates(update_map).Error
	fmt.Print("789787978797897879787987978778978797879879787879787", update_map)
	return err
}

func FindArticleById(id string) Article {
	article := Article{}
	db.W_Db.Where("id = ?", id).First(&article)
	return article
}

//获取所有的文章
func AllPublicArticles() ([]*Article, error) {
	return _ListArticle(true)
}

func _ListArticle(isPublic bool) ([]*Article, error) {
	var articles []*Article
	var err error
	if isPublic {
		err = db.W_Db.Where("is_public = ?", true).Find(&articles).Error
	} else {
		err = db.W_Db.Find(&articles).Error
	}
	return articles, err
}

//分页获取所有公开的文章
func PublicArticleLimit(page int, articleCount int) ([]*Article, error) {
	var articles []*Article
	var err error

	err = db.W_Db.Limit(articleCount).Offset((page-1)*articleCount).Where("deleted_at IS NULL and is_public = ?", true).Find(&articles).Error

	return articles, err
}

//获取文章的数量
func ArticleCount() int {
	var count int
	db.W_Db.Table("articles").Where("deleted_at IS NULL and is_public = ?", true).Count(&count)

	fmt.Println("******************************************************************", count)
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
func FindArticlesByLabel(labelId string) []*Article {
	var articles []*Article
	//label := FindLabelById(labelId)
	log.Println("++++++++++")
	//db.W_Db.Model(&label).Association("Labels").Find(&articles)
	//db.W_Db.Table("articles").Select("articles.id, labels.labelId").Joins("left join articles on labels.article_id = articles.id").Where("label.id = ?",labelId).Find(&articles)
	//select articles.*  from articles right  join artitle_labels on artitle_labels.article_id = articles.id;
	db.W_Db.Raw("select articles.*  from articles inner join artitle_labels on "+
		"artitle_labels.article_id = articles.id where artitle_labels.label_id = ? and articles.is_public = ?", labelId, true).Scan(&articles)
	return articles
}

//通过文章发现标签
func (article *Article) FindLabelsByArticle() Label {
	var label Label
	articleId := article.ID
	//label := FindLabelById(label_id)
	log.Println("++++++++++")
	//db.W_Db.Model(&label).Association("Labels").Find(&articles)
	//db.W_Db.Table("articles").Select("articles.id, labels.label_id").Joins("left join articles on labels.articleId = articles.id").Where("label.id = ?",label_id).Find(&articles)
	//select articles.*  from articles right  join artitle_labels on artitle_labels.articleId = articles.id;
	db.W_Db.Raw("select labels.*  from labels inner join artitle_labels on "+
		"artitle_labels.label_id = labels.id where artitle_labels.article_id = ?", articleId).Scan(&label)
	return label

}

//分页获取当前用户的所有文章
func UserArticleLimit(page int, articleCount int, user_id uint) ([]*Article, error) {
	var articles []*Article
	var err error

	err = db.W_Db.Limit(articleCount).Offset((page-1)*articleCount).Where("user_id = ?", user_id).Find(&articles).Error

	return articles, err
}

//获取当前用户文章的数量
func UserArticleCount(userId uint) int {
	var count int
	db.W_Db.Table("articles").Where("deleted_at IS NULL and  user_id = ?", userId).Count(&count)
	fmt.Println("******************************************************************", count)
	return count
}

func (article Article) FindUserIsCollect(user User) bool {
	var userIsCollect bool
	var count int
	db.W_Db.Table("collects").Where("deleted_at IS NULL and  is_collect = ? and article_id = ? and user_id = ?", true, article.ID, user.ID).Count(&count)
	if count > 0 {
		userIsCollect = true
	} else {
		userIsCollect = false
	}
	return userIsCollect
}

//分页显示当前用户收藏的文章
func FindUserCollectArticles(userId uint, articleCount int, page int) ([]*Article, error) {
	var articles []*Article
	var err error
	var collects []*Collect

	err = db.W_Db.Limit(articleCount).Offset((page-1)*articleCount).Where("deleted_at IS NULL and  user_id = ? and is_collect = ?", userId, true).Find(&collects).Error
	for _, collect := range collects {
		articleTemp := FindArticleById(fmt.Sprint(collect.ArticleID))
		articles = append(articles, &articleTemp)
	}
	return articles, err
}

//获取用户收藏的文章的总数量
func UserCollectArticlesCount(userId uint) int {
	var count int
	db.W_Db.Table("collects").Where("deleted_at IS NULL and  user_id = ? and is_collect = ?", userId, true).Count(&count)
	return count
}

//根据文章标题或者内容搜索文章
func FindArticleByKeyWord(keyWord string, page int, articleCount int, isPublic bool, userId uint) ([]*Article, error) {
	var articles []*Article
	var err error
	//or article_content LIKE ?

	err = db.W_Db.Limit(articleCount).Offset((page-1)*articleCount).Where("deleted_at IS NULL").
		Where("  is_public = ? or (user_id =? and is_public = false)", isPublic, userId).
		Where("article_title LIKE ? or article_content LIKE ?", "%"+keyWord+"%", "%"+keyWord+"%").Find(&articles).Error

	return articles, err
}

func KeyWordArticleCount(keyWord string, isPublic bool, userId uint) int {
	var count int
	db.W_Db.Table("articles").Where("deleted_at IS NULL").
		Where("  is_public = ? or (user_id =? and is_public = false)", isPublic, userId).
		Where("article_title LIKE ? or article_content LIKE ?", "%"+keyWord+"%", "%"+keyWord+"%").Count(&count)
	return count
}
