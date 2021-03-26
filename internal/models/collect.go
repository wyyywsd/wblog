package models

import (
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type Collect struct {
	gorm.Model
	UserId    uint
	ArticleID uint
	IsCollect bool
}

func CreateCollect(userId uint, articleId uint, isCollect bool) {
	collect := Collect{
		UserId:    userId,
		ArticleID: articleId,
		IsCollect: isCollect,
	}
	db.W_Db.Create(&collect)
}

func UpdateCollect(collect *Collect, isCollect bool) {
	db.W_Db.Model(&collect).Update("is_collect", isCollect)
}

func FindCollectByUserIdAndArticleId(userId uint, articleId uint) (*Collect, bool, error) {
	collect := &Collect{}
	var err error
	err = db.W_Db.Where("deleted_at IS NULL and user_id = ? and article_id = ?", userId, articleId).First(&collect).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, false, err
		}
		return nil, false, nil
	}
	return collect, true, err
}
