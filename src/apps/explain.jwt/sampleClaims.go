package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JwtSampleNotAfter  = time.Hour
	JwtSampleNotBefore = time.Second * 10
	JwtSampleIssuer    = "MyCooki"
	JwtSampleSubject   = "Пример выпуска JWT токена"
)

var (
	JwtSampleAudience = []string{"all apple", "only limon"}
)

type SampleClaims struct {
	ExpirationTime time.Time
	IssuedAt       time.Time
	NotBefore      time.Time
	Issuer         string
	Subject        string
	Audience       jwt.ClaimStrings
}

func NewSampleClaimDef() (ret *SampleClaims) {

	now := time.Now()

	ret = &SampleClaims{
		ExpirationTime: now.Add(JwtSampleNotBefore),
		IssuedAt:       now,
		NotBefore:      now.Add(JwtSampleNotBefore),
		Issuer:         JwtSampleIssuer,
		Subject:        JwtSampleSubject,
		Audience:       JwtSampleAudience,
	}
	return
}

func (ref *SampleClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: ref.ExpirationTime,
	}, nil
}

func (ref *SampleClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: ref.IssuedAt,
	}, nil
}

func (ref *SampleClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{
		Time: ref.NotBefore,
	}, nil
}

func (ref *SampleClaims) GetIssuer() (string, error) {
	return ref.Issuer, nil
}

func (ref *SampleClaims) GetSubject() (string, error) {
	return ref.Subject, nil
}

func (ref *SampleClaims) GetAudience() (jwt.ClaimStrings, error) {
	return ref.Audience, nil
}
