package model

import "time"

type Question struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`

	Answers []Answer `json:"answers,omitempty" gorm:"foreignKey:QuestionID;references:id;constraint:OnDelete:CASCADE"`
}
