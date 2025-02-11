package models

type Token struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null"`
	Token        string `gorm:"not null"`
	RefreshToken string `gorm:"not null"`
	TokenType    string `gorm:"not null"`
	Expirated    bool   `gorm:"not null"`
	Revoked      bool   `gorm:"not null"`
}
