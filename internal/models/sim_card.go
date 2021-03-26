package models

import (
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type SimCard struct {
	gorm.Model
	Iccid          string
	Msisdn         string
	UnbindBatchID  uint
	AgentName      string
	ReplaceReason  string
	EquipmentPhoto string
}

func FindSimCardsByUnbindBatch(unbindBatch UnbindBatch) ([]*SimCard, error) {
	var simCards []*SimCard
	var err error
	err = db.W_Db.Model(&unbindBatch).Order("agent_name asc").Related(&simCards).Error
	return simCards, err
}

func CreateSimCards(agentName string, iccid string, msisdn string, unbindBatchId uint, replaceReason string, equipmentPhoto string) {
	sim_card := SimCard{
		AgentName:      agentName,
		Iccid:          iccid,
		Msisdn:         msisdn,
		UnbindBatchID:  unbindBatchId,
		ReplaceReason:  replaceReason,
		EquipmentPhoto: equipmentPhoto,
	}
	db.W_Db.Create(&sim_card)
}

func DeleteSimCardById(id uint) {
	var simCard SimCard
	db.W_Db.Where("deleted_at IS NULL and id = ?", id).First(&simCard)
	db.W_Db.Delete(simCard)
}

func (simCard *SimCard) ReplaceReasonDisplay() string {
	var reason string
	if simCard.ReplaceReason == "equipment_damage" {
		reason = "设备损坏"
	} else if simCard.ReplaceReason == "test_equipment" {
		reason = "测试设备"
	} else if simCard.ReplaceReason == "misoperation" {
		reason = "误操作"
	} else {
		reason = "数据库没有记录此项原因"
	}

	return reason
}
