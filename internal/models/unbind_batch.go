package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type UnbindBatch struct {
	gorm.Model
	Status    string
	CarrierId uint
	SimCards  []SimCard
}

func CreateUnbindBatch(carrierId uint, status string) {
	unbindBatch := UnbindBatch{CarrierId: carrierId, Status: status}
	db.W_Db.Create(&unbindBatch)
}

func FindUnbindBatchByPage(batchCount int, page int) ([]*UnbindBatch, error) {
	var unbindBatches []*UnbindBatch
	var err error
	err = db.W_Db.Limit(batchCount).Offset((page - 1) * batchCount).Where("deleted_at IS NULL").Order("created_at desc").Find(&unbindBatches).Error

	return unbindBatches, err
}

func (unbindBatches *UnbindBatch) FindCarrierByUnbindBatch() Carrier {
	var carrier Carrier
	db.W_Db.Where("deleted_at IS NULL and id = ?", unbindBatches.CarrierId).First(&carrier)
	return carrier
}

func (unbindBatches *UnbindBatch) UnbindBatchStatusDisplay() string {
	var status string
	if unbindBatches.Status == "pending" {
		status = "待提交运营商"
	} else if unbindBatches.Status == "processing" {
		status = "已提交运营商"
	} else if unbindBatches.Status == "success" {
		status = "已完成"
	}
	return status

}

func FindUnbindBatchById(id uint) (UnbindBatch, error) {
	var err error
	var unbindBatch UnbindBatch
	err = db.W_Db.Where("deleted_at IS NULL and id = ?", id).First(&unbindBatch).Error
	return unbindBatch, err
}

func UpdateUnbindBatchStatusById(id uint, status string) {
	var unbindBatch UnbindBatch
	unbindBatch, _ = FindUnbindBatchById(id)
	db.W_Db.Model(&unbindBatch).Update("status", status)
}

func DeleteUnbindBatchById(id uint) {
	var unbindBatch UnbindBatch
	db.W_Db.Where("deleted_at IS NULL and id = ?", id).First(&unbindBatch)
	db.W_Db.Delete(&unbindBatch)
}

func UnbindCount() int {
	var count int
	db.W_Db.Table("unbind_batches").Where("deleted_at IS NULL").Count(&count)

	fmt.Println("******************************************************************", count)
	return count
}
