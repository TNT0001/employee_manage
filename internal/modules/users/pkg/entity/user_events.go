package entity

import (
	"tungnt/emmployee_manage/pkg/share/model"
)

// UserEventsTableName TableName
var UserEventsTableName = "user_events"

type UserEvents struct {
	ID          int    `gorm:"column:id;primaryKey;type:integer(11);not null;autoIncrement;unique" mapstructure:"id"`
	UserID      int    `gorm:"column:user_id;type:integer(11);not null;unique" mapstructure:"user_id"`
	Description string `gorm:"column:description;type:text;not null;unique" mapstructure:"description"`
	TimeLine    string `gorm:"column:time_line;type:date;not null" mapstructure:"time_line"`
	model.BaseModelWithDeleted
}

// TableName func
func (i *UserEvents) TableName() string {
	return UserEventsTableName
}
