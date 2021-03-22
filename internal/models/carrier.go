package models

import (
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type Carrier struct {
	gorm.Model
	Name string
	UnbindBatchs  []UnbindBatch
}

func AllCarriers ()([]*Carrier,error){
	var err error
	var carriers []*Carrier
	err = db.W_Db.Find(&carriers).Error

	return carriers,err
}

