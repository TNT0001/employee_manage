package entity

import "tungnt/emmployee_manage/pkg/share/model"

// PermissionsTableName TableName
var PermissionsTableName = "permissions"

type Permissions struct {
	ID   int    `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" mapstructure:"id"`
	Name string `gorm:"column:name;type:varchar(25);not null;unique" mapstructure:"name"`
	model.BaseModel
}

// TableName func
func (i *Permissions) TableName() string {
	return PermissionsTableName
}
