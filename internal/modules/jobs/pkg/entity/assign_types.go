package entity

import "tungnt/emmployee_manage/pkg/share/model"

// AssignTypesTableName TableName
var AssignTypesTableName = "assign_types"

type AssignTypes struct {
	ID   int    `gorm:"column:id;primaryKey;type:integer(11);not null;autoIncrement;unique" mapstructure:"id"`
	Name string `gorm:"column:name;type:varchar(255);not null;unique" mapstructure:"name"`
	model.BaseModel
}

// TableName func
func (i *AssignTypes) TableName() string {
	return AssignTypesTableName
}
