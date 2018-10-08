package common

import (
	"github.com/google/uuid"
	"time"
)

type Video struct {
	JobId     uuid.UUID  `json:"job_id" gorm:"column:job_id"`
	Created   *time.Time `json:"created" gorm:"column:created"`
	Completed *time.Time `json:"completed" gorm:"column:completed"`
	Format    string     `json:"format" gorm:"column:format"`
	Input     string     `json:"input" gorm:"column:input"`
	Output    string     `json:"output" gorm:"column:output"`
	Status    string     `json:"status" gorm:"column:status"`
}
