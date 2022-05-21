package entity

import "tungnt/emmployee_manage/pkg/share/model"

// UserPermissionsTableName TableName
var UserPermissionsTableName = "user_permissions"

type UserPermissions struct {
	ID           int    `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" mapstructure:"id"`
	UserID       string `gorm:"column:user_id;type:integer(11);not null" mapstructure:"user_id"`
	PermissionID string `gorm:"column:permission_id;type:integer(11);not null" mapstructure:"permission_id"`
	model.BaseModel
}

// TableName func
func (i *UserPermissions) TableName() string {
	return UserPermissionsTableName
}
