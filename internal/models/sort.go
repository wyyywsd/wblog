package models

import "github.com/jinzhu/gorm"

type Sort struct {
	gorm.Model
	SortName        string
	SortAlias       string
	SortDescription string
	ParentSortId    int
}
