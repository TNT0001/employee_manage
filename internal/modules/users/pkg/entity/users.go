package entity

import (
	"tungnt/emmployee_manage/pkg/share/model"
)

// UsersTableName TableName
var UsersTableName = "users"

type Users struct {
	ID             int     `gorm:"column:id;primaryKey;type:integer(11);not null;autoIncrement;unique" mapstructure:"id"`
	KeycloakUserID string  `gorm:"column:keycloak_user_id;type:varchar(255);not null;unique" mapstructure:"keycloak_user_id"`
	UserName       string  `gorm:"column:user_name;type:varchar(255);not null;unique" mapstructure:"user_name"`
	JoinDate       *string `gorm:"column:join_date;type:date;not null" mapstructure:"join_date"`
	TeamID         *int    `gorm:"column:team_id;type:integer(11)" mapstructure:"team_id"`
	model.BaseModelWithDeleted
}

// TableName func
func (i *Users) TableName() string {
	return UsersTableName
}
