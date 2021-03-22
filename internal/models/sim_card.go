package models

import (
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type SimCard struct {
	gorm.Model
	Iccid string
	Msisdn string
	UnbindBatchID uint
	AgentName string
	ReplaceReason string
	EquipmentPhoto string
}

func FindSimCardsByUnbindBatch(unbind_batch UnbindBatch)([]*SimCard,error){
	var sim_cards []*SimCard
	var err error
	err = db.W_Db.Model(&unbind_batch).Order("agent_name asc").Related(&sim_cards).Error
	return sim_cards,err
}

func CreateSimCards(agent_name string,iccid string,msisdn string,unbind_batch_id uint,replace_reason string,equipment_photo string){
	sim_card := SimCard{
		AgentName: agent_name,
		Iccid: iccid,
		Msisdn: msisdn,
		UnbindBatchID: unbind_batch_id,
		ReplaceReason: replace_reason,
		EquipmentPhoto: equipment_photo,
	}
	db.W_Db.Create(&sim_card)
}

func DeleteSimCardById(id uint){
	var sim_card SimCard
	db.W_Db.Where("deleted_at IS NULL and id = ?", id).First(&sim_card)
	db.W_Db.Delete(sim_card)
}

func (sim_card *SimCard) ReplaceReasonDisplay()string{
	var reason string
	if (sim_card.ReplaceReason == "equipment_damage"){
		reason = "设备损坏"
	}else if (sim_card.ReplaceReason == "test_equipment"){
		reason = "测试设备"
	}else if (sim_card.ReplaceReason == "misoperation"){
		reason = "误操作"
	}else {
		reason = "数据库没有记录此项原因"
	}

	return reason
}