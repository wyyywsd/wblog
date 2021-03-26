package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type Zan struct {
	gorm.Model
	UserId    uint
	CommentId uint
	IsZan     bool
}

func CreateZan(userId uint, commentId uint, isZan bool) {
	zan := Zan{
		UserId:    userId,
		CommentId: commentId,
		IsZan:     isZan,
	}
	db.W_Db.Create(&zan)
}

func UpdateZan(zan *Zan, isZan bool) {
	fmt.Println("更新获取到的值是", isZan)
	db.W_Db.Model(&zan).Update("is_zan", isZan)
}
func FindZanByUserIDAndCommentID(userId uint, commentId uint) (*Zan, bool, error) {
	zan := &Zan{}
	var err error
	err = db.W_Db.Where("deleted_at IS NULL and user_id = ? and comment_id = ?", userId, commentId).First(&zan).Error
	if err != nil {
		//如果不是找不到表的错误，返回
		if err != gorm.ErrRecordNotFound {
			return nil, false, err
		}
		//如果是找不到 就直接返回
		return nil, false, nil
	}
	return zan, true, err

}
