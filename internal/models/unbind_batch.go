package models

import (
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type UnbindBatch struct {
	gorm.Model
	Status string
	CarrierId uint
	SimCards []SimCard

}

func CreateUnbindBatch(carrier_id uint,status string){
	unbind_batch := UnbindBatch{CarrierId: carrier_id, Status: status}
	db.W_Db.Create(&unbind_batch)
}

func FindUnbindBatchByPage(batchCount int,page int)([]*UnbindBatch,error){
	var unbind_batchs []*UnbindBatch
	var err error
	err = db.W_Db.Limit(batchCount).Offset((page-1)*batchCount).Where("deleted_at IS NULL").Order("created_at desc").Find(&unbind_batchs).Error

	return unbind_batchs,err
}

func (unbind_batch *UnbindBatch) FindCarrierByUnbindBatch()Carrier{
	var carrier Carrier
	db.W_Db.Where("deleted_at IS NULL and id = ?", unbind_batch.CarrierId).First(&carrier)
	return carrier
}

func (unbind_batch *UnbindBatch) UnbindBatchStatusDisplay() string{
	var status string
	if (unbind_batch.Status == "pending"){
		status = "待提交运营商"
	}else if (unbind_batch.Status == "processing"){
		status = "已提交运营商"
	}else if (unbind_batch.Status == "success"){
		status = "运营商已完成"
	}
	return status

}

func FindUnbindBatchById(id uint) (UnbindBatch,error){
	var err error
	var unbind_batch UnbindBatch
	err = db.W_Db.Where("deleted_at IS NULL and id = ?",id).First(&unbind_batch).Error
	return unbind_batch,err
}

func UpdateUnbindBatchStatusById(id uint,status string){
	var unbind_batch UnbindBatch
	unbind_batch,_ = FindUnbindBatchById(id)
	db.W_Db.Model(&unbind_batch).Update("status", status)
}

func DeleteUnbindBatchById(id uint){
	var unbind_batch UnbindBatch
	db.W_Db.Where("deleted_at IS NULL and id = ?",id).First(&unbind_batch)
	db.W_Db.Delete(&unbind_batch)
}