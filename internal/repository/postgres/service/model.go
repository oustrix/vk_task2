package service

type Service struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	UserID   int64 `gorm:"not null"`
	Login    string
	Password string
}
