package entity

import (
	"tungnt/emmployee_manage/pkg/share/model"
)

// PersonalInfosTableName TableName
var PersonalInfosTableName = "personal_infos"

type PersonalInfos struct {
	ID                    int     `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" mapstructure:"id"`
	UserID                int     `gorm:"column:user_id;type:integer(11);unique;not null" mapstructure:"user_id"`
	Name                  *string `gorm:"column:name;type:varchar(25)" mapstructure:"name"`
	SurName               *string `gorm:"column:sur_name;type:varchar(25)" mapstructure:"sur_name"`
	FullName              string  `gorm:"column:full_name;type:varchar(255);not null;unique" mapstructure:"full_name"`
	Email                 *string `gorm:"column:email;type:varchar(255);unique" mapstructure:"email"`
	FaceBook              *string `gorm:"column:facebook;type:varchar(255);unique" mapstructure:"facebook"`
	Linken                *string `gorm:"column:linken;type:varchar(255);unique" mapstructure:"linken"`
	PhoneNumber           *string `gorm:"column:phone_number;type:string;unique" mapstructure:"phone_number"`
	Address               *string `gorm:"column:address;type:varchar(255)" mapstructure:"address"`
	OnProbationaryPeriod  bool    `gorm:"column:on_probationary_period;type:bool;not null" mapstructure:"on_probationary_period"`
	StartProbationaryDate *string `gorm:"column:start_probationary_date;type:date" mapstructure:"start_probationary_date"`
	EndProbationaryDate   *string `gorm:"column:end_probationary_date;type:date" mapstructure:"end_probationary_date"`
	model.BaseModelWithDeleted
}

// TableName func
func (i *PersonalInfos) TableName() string {
	return PersonalInfosTableName
}
