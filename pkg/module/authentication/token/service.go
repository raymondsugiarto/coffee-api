package token

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/raymondsugiarto/coffee-api/config"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
)

type Service interface {
	GenerateToken(ctx context.Context, userCredentialData e.UserCredentialData) (*entity.LoginDto, error)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) GenerateToken(ctx context.Context, userCredentialData e.UserCredentialData) (*entity.LoginDto, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userCredentialData.ID
	claims["uid"] = userCredentialData.UserID
	if userCredentialData.AdminID != "" {
		claims["aid"] = userCredentialData.AdminID
	}
	if userCredentialData.CompanyID != "" {
		claims["coid"] = userCredentialData.CompanyID
	}
	if userCredentialData.CustomerID != "" {
		claims["cid"] = userCredentialData.CustomerID
	}
	claims["exp"] = time.Now().Add(time.Hour * 720).Unix() // 30 days

	cfg := config.GetConfig()
	t, err := token.SignedString([]byte(cfg.Server.Rest.SecretKey))
	if err != nil {
		return nil, errors.New("errorGeneratetoken")
	}

	// TODO: save to redis
	return &entity.LoginDto{
		Token:   t,
		Expired: strconv.Itoa(int(claims["exp"].(int64))),
	}, nil
}
