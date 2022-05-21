package entity

import "tungnt/emmployee_manage/pkg/share/model"

// UserRolesTableName TableName
var UserRolesTableName = "user_roles"

type UserRoles struct {
	ID     int `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" mapstructure:"id"`
	UserID int `gorm:"column:user_id;type:integer(11);not null" mapstructure:"user_id"`
	RoleID int `gorm:"column:role_id;type:integer(11);not null" mapstructure:"role_id"`
	model.BaseModel
}

// TableName func
func (i *UserRoles) TableName() string {
	return UserRolesTableName
}
