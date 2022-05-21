package entity

import "tungnt/emmployee_manage/pkg/share/model"

// JobsTableName TableName
var JobsTableName = "jobs"

type Jobs struct {
	ID            int    `gorm:"column:id;primaryKey;type:integer(11);not null;autoIncrement;unique" mapstructure:"id"`
	ProjectName   string `gorm:"column:project_name;type:varchar(255);not null;unique" mapstructure:"project_name"`
	UserID        int    `gorm:"column:user_id;primaryKey;type:integer(11);not null" mapstructure:"user_id"`
	AssignTypeID  *int   `gorm:"column:assign_type_id;type:integer(11)" mapstructure:"assign_type_id"`
	AssignPercent *int   `gorm:"column:assign_percent;type:integer(11)" mapstructure:"assign_percent"`
	model.BaseModel
}

// TableName func
func (i *Jobs) TableName() string {
	return JobsTableName
}
