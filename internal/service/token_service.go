package service

import (
	"errors"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
)

type TokenService interface {
	IsTokenValid(tipo string, token string) (*models.Token, error)
	UpdateToken(item *models.Token) error
	RevokeAllUserTokens(id uint) error
	SaveToken(item *models.Token) error
}

type tokenService struct {
	tokenRepo repository.TokenRepository
}

func NewTokenService(repo repository.TokenRepository) TokenService {
	return &tokenService{repo}
}

func (s *tokenService) IsTokenValid(tipo string, token string) (*models.Token, error) {

	if tipo != "refreshToken" && tipo != "accessToken" {
		return nil, errors.New("invalid token type")
	} else {
		if tipo == "refreshToken" {
			tipo = "refresh_token"
		}
		if tipo == "accessToken" {
			tipo = "token"
		}
	}

	var tokenEntity = &models.Token{}

	err := s.tokenRepo.FindField(tokenEntity, tipo, token)

	if err != nil {
		return nil, err
	}

	if tokenEntity.ID == 0 {
		return nil, nil
	}
	return tokenEntity, nil
}

func (s *tokenService) UpdateToken(item *models.Token) error {
	return s.tokenRepo.Update(item)
}

func (s *tokenService) RevokeAllUserTokens(id uint) error {
	items, err := s.tokenRepo.FindAllValidByUserID(id)

	if err != nil {
		return err
	}

	if items == nil {
		return nil
	}

	for _, item := range *items {
		item.Revoked = true
		err = s.tokenRepo.Update(&item) //TODO arreglo
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *tokenService) SaveToken(item *models.Token) error {
	return s.tokenRepo.Create(item)
}
