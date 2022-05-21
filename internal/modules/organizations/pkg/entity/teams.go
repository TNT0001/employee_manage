package entity

import "tungnt/emmployee_manage/pkg/share/model"

// TeamsTableName TableName
var TeamsTableName = "teams"

type Teams struct {
	ID           int     `gorm:"column:id;primaryKey;type:bigint;not null;autoIncrement;unique" mapstructure:"id"`
	TeamName     string  `gorm:"column:team_name;type:varchar(255);not null" mapstructure:"team_name"`
	CountryCode  string  `gorm:"column:country_code;type:varchar(2);not null" mapstructure:"country_code"`
	DivisionName *string `gorm:"column:division_name;type:varchar(25)" mapstructure:"division_name"`
	Kind         *string `gorm:"column:kind;type:varchar(25)" mapstructure:"kind"`
	model.BaseModel
}

// TableName func
func (i *Teams) TableName() string {
	return TeamsTableName
}
