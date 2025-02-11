package repository

import (
	"com.mx/crud/internal/models"
	"gorm.io/gorm"
)

type TokenRepository interface {
	Repository[models.Token]
	FindAllValidByUserID(value uint) (*[]models.Token, error)
}

type tokenRepository struct {
	*BaseRepository[models.Token]
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{
		BaseRepository: NewBaseRepository[models.Token](db),
	}
}

func (r *tokenRepository) FindAllValidByUserID(value uint) (*[]models.Token, error) {
	var items []models.Token
	if err := r.DB.Where("user_id = ? and expirated= False and revoked =False ", value).Find(&items).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &items, nil
}

func (r *tokenRepository) RevokeAllUserTokens() {
	// Update with conditions
	//db.Model(&User{}).Where("active = ?", true).Update("name", "hello")
	// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

	// User's ID is `111`:
	//db.Model(&user).Update("name", "hello")

	/*if err := r.DB.Model(&models.Token{}).Where("user_id = ? and expirated= False and revoked =False ", true).Update("name", "hello").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}*/

}
