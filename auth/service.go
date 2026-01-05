package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	GenerateToken(userID int, ctx context.Context) (string, error)
	ValidateToken(token string, ctx context.Context) (*jwt.Token,error)
}

type jwtService struct {
	Redis *redis.Client
}

func NewService(redisClient *redis.Client) *jwtService  {
	return &jwtService{Redis: redisClient}
}

var jwtSecretToken = []byte(os.Getenv("JWT_SECRETKEY"))

func (s *jwtService) GenerateToken(userID int, ctx context.Context)(string,error) {
	claim:= jwt.MapClaims{
		"id": userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	claim["user_id"] = userID

	token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claim)

	signedToken,err:= token.SignedString(jwtSecretToken)
	if err != nil {
		return signedToken,err
	}

	_,err = s.Redis.SetEx(ctx,signedToken,userID,time.Hour *24 * 30).Result()
	if err != nil {
		return "",err
	}

	return signedToken,nil
}

func (s *jwtService) ValidateToken(tokenString string,ctx context.Context) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecretToken, nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("exp not found")
	}

	if time.Now().Unix() > int64(exp) {
		return nil, errors.New("token expired")
	}

	result,err:= s.Redis.Exists(ctx,tokenString).Result()	
	if err != nil {
		return nil,err
	}

	if result == 0 {
		return nil,errors.New("Unauthorized")
	}

	return jwtToken, nil
}

