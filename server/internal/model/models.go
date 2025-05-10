package model

import "time"

type User struct {
	ID           string `gorm:"primaryKey;autoIncrement"`
	Login        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	Expressions  []Expression
}

type Expression struct {
	ID         string `gorm:"primaryKey"`
	UserID     string `gorm:"not null"`
	Expression string `gorm:"not null"`
	Status     string `gorm:"not null;default:'pending'"`
	Result     float64
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

type Operation struct {
	ID        string `gorm:"primaryKey;autoIncrement"`
	Type      string `gorm:"not null"`
	Duration  int64  `gorm:"not null"` // в миллисекундах
	CreatedAt time.Time
}
