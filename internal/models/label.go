package models

import (
	"github.com/jinzhu/gorm"
	"gorm_demo/internal/db"
)

type Label struct {
	gorm.Model
	LabelName        string
	LabelAlias       string
	LabelDescription string
}

func FindLabelById(labelId string) Label {
	label := Label{}
	db.W_Db.Where("id = ?", labelId).First(&label)
	return label
}
func AllLabels() ([]*Label, error) {

	return _ListLabel()
}

func _ListLabel() ([]*Label, error) {
	var labels []*Label
	var err error
	err = db.W_Db.Find(&labels).Error
	return labels, err
}

func CreateLabel(name string) error {
	label := Label{LabelName: name}
	err := db.W_Db.Create(&label).Error
	return err
}
