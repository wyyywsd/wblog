package models

import (
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type Collect struct {
	gorm.Model
	UserId uint
	ArticleID uint
	IsCollect bool
}

func CreateCollect(user_id uint,article_id uint,is_collect bool){
	collect := Collect{
		UserId: user_id,
		ArticleID: article_id,
		IsCollect: is_collect,
	}
	db.W_Db.Create(&collect)
}

func UpdateCollect(collect *Collect,is_collect bool){
	db.W_Db.Model(&collect).Update("is_collect", is_collect)
}

func FindCollectByUserIdAndArticleId(user_id uint,article_id uint)(*Collect,bool,error){
	collect := &Collect{}
	var err error
	err = db.W_Db.Where("deleted_at IS NULL and user_id = ? and article_id = ?",user_id,article_id).First(&collect).Error
	if err != nil{
		if err != gorm.ErrRecordNotFound{
			return nil,false,err
		}
		return nil,false,nil
	}
	return collect,true,err
}
