package entity

import "tungnt/emmployee_manage/pkg/share/model"

// RolePermissionsTableName TableName
var RolePermissionsTableName = "role_permissions"

type RolePermissions struct {
	ID           int         `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" mapstructure:"id"`
	RoleID       int         `gorm:"column:role_id;type:integer(11);not null" mapstructure:"role_id"`
	PermissionID int         `gorm:"column:permission_id;type:integer(11);not null" mapstructure:"permission_id"`
	Role         Roles       `gorm:"foreignKey:RoleID"`
	Permission   Permissions `gorm:"foreignKey:PermissionID"`
	model.BaseModel
}

// TableName func
func (i *RolePermissions) TableName() string {
	return RolePermissionsTableName
}
