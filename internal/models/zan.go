package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type Zan struct {
	gorm.Model
	UserId uint
	CommentId uint
	IsZan bool

}

func CreateZan(user_id uint,comment_id uint,is_zan bool){
	zan := Zan{
		UserId: user_id,
		CommentId: comment_id,
		IsZan: is_zan,
	}
	db.W_Db.Create(&zan)
}

func UpdateZan(zan *Zan,is_zan bool){
	fmt.Println("更新获取到的值是",is_zan)
	db.W_Db.Model(&zan).Update("is_zan",is_zan)
}
func FindZanByUserIDAndCommentID(user_id uint,comment_id uint) (*Zan,bool,error){
	zan := &Zan{}
	var err error
	err = db.W_Db.Where("deleted_at IS NULL and user_id = ? and comment_id = ?",user_id,comment_id).First(&zan).Error
	if err != nil{
		//如果不是找不到表的错误，返回
		if err != gorm.ErrRecordNotFound{
			return nil,false,err
		}
		//如果是找不到 就直接返回
		return nil,false,nil
	}
	return zan,true,err

}