package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var global *service

var defaultOption = Option{
	AccessTokenDuration: 24 * time.Hour,
	Issuer:              "SMM",
}

type CustomClaims struct {
	Admin bool `json:"admin"`
	Type  int  `json:"type"`
	jwt.RegisteredClaims
}

type Option struct {
	AccessTokenDuration time.Duration
	Issuer              string
}

type service struct {
	accessTokenDuration time.Duration
	secretKey           string
	issuer              string
}
type Service interface {
	InitGlobal()
	GenerateToken(id primitive.ObjectID, isAdmin bool, tokenType int) (string, error)
	ValidateToken(signedString string) (*CustomClaims, error)
}

func New(secretKey string, opts ...Option) Service {
	var opt Option
	if len(opts) == 0 {
		opt = defaultOption
	} else {
		opt = opts[0]
	}
	return &service{
		secretKey:           secretKey,
		accessTokenDuration: opt.AccessTokenDuration,
		issuer:              opt.Issuer,
	}
}

func (s *service) InitGlobal() {
	global = s
}

func GetGlobal() *service {
	return global
}

func (s *service) GenerateToken(id primitive.ObjectID, isAdmin bool, tokenType int) (string, error) {
	claims := CustomClaims{
		true,
		tokenType,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			ID:        id.Hex(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *service) ValidateToken(signedString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(signedString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	} else {
		return nil, jwt.ErrInvalidType
	}
}
