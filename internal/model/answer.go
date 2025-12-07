package model

import "time"

type Answer struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	QuestionID int       `json:"question_id" gorm:"not null"`
	UserID     string    `json:"user_id" gorm:"type:text;not null"`
	Text       string    `json:"text" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null;default:now()"`
}
