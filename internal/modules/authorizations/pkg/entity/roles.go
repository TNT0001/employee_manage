package entity

import "tungnt/emmployee_manage/pkg/share/model"

// RolesTableName TableName
var RolesTableName = "roles"

type Roles struct {
	ID               int               `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" mapstructure:"id"`
	Name             string            `gorm:"column:name;type:varchar(25);not null;unique" mapstructure:"name"`
	RolesPermissions []RolePermissions `gorm:"foreignKey:RoleID"`
	model.BaseModel
}

// TableName func
func (i *Roles) TableName() string {
	return RolesTableName
}
